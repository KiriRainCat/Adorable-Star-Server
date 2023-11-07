package crawler

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/model"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var d = dao.Jupiter
var browser *rod.Browser
var pagePool rod.PagePool
var pageCreate func() *rod.Page

var TaskPool []func(...any) any

// Initialize crawler with page pool to execute tasks asynchronously
func Init() {
	// Launch differently in armbian (linux-arm-hf) and windows (dev-env)
	if gin.Mode() == gin.ReleaseMode {
		u := launcher.New().Bin("/bin/chromium-browser").
			Set("--no-sandbox").
			Set("--headless").
			Set("--remote-debugging-port", "7999").
			Set("--disable-gpu").
			Set("--disable-dev-shm-usage").
			Set("--disable-setuid-sandbox").
			Set("--no-first-run").
			Set("--no-zygote").
			Set("--single-process").
			MustLaunch()

		browser = rod.New().ControlURL(u).MustConnect()
	} else {
		browser = rod.New().NoDefaultDevice().MustConnect()
	}

	// Create page pool for multithreading
	pagePool = rod.NewPagePool(10)
	pageCreate = func() *rod.Page {
		return browser.MustIncognito().MustPage()
	}
	pagePool.Put(pagePool.Get(pageCreate).MustNavigate("https://login.jupitered.com/login/"))
}

// Composite function for fetch and store data from Jupiter
func CrawlerJob(uid int) {
	// Fetch all data
	courseList, assignmentsList, gpa, err := FetchData(uid)
	if err != nil {
		return
	}

	// Store fetched data to database
	StoreData(uid, gpa, courseList, assignmentsList)
}

// Open a webpage for Jupiter
func OpenJupiterPage() *rod.Page {
	// Navigate to Jupiter Ed login page
	page := pagePool.Get(pageCreate).MustSetCookies().MustNavigate("https://login.jupitered.com/login/")

	// Bypass Cloudflare detection for crawler
	page.MustEvalOnNewDocument("const newProto = navigator.__proto__;delete newProto.webdriver;navigator.__proto__ = newProto;")

	return page
}

// Fetch all Jupiter data for a student
func FetchData(uid int) (courseList []*model.Course, assignmentsList [][]*model.Assignment, gpa string, err error) {
	// Get a page to access Jupiter
	page := OpenJupiterPage()
	defer pagePool.Put(page)

	// Find user's Jupiter account info
	data, _ := d.GetDataByUID(uid)

	// Login
	if err = Login(page, data.Account, data.Password); err != nil {
		return
	}

	// Get courses from nav bar
	_, courses, err := NavGetOptions(page)
	if err != nil {
		return
	}
	if rod.Try(func() { page.MustElement("#touchnavbtn").MustClick() }) != nil {
		return
	}

	// Navigate through all of the courses to fetch course data
	for idx := range courses {
		// Nav bar navigation to course page
		_, courses, err := NavGetOptions(page)
		if err != nil {
			continue
		}
		var courseName string
		if rod.Try(func() { courseName = courses[idx].Timeout(time.Second * 2).MustText() }) != nil {
			continue
		}
		if NavNavigate(page, courses[idx]) != nil {
			continue
		}

		// Get grade
		courseList = append(courseList, GetCourseGrade(page, courseName, uid))

		// Get all assignments
		assignmentsList = append(assignmentsList, GetCourseAssignments(page, courseName, uid))
	}

	// Fetch GPA and report card image
	gpa = FetchReportAndGPA(page)

	return
}

// Store all fetched data to database
func StoreData(uid int, gpa string, courseList []*model.Course, assignmentsList [][]*model.Assignment) {
	var count = 0
	wg := &sync.WaitGroup{}

	// Get stored course list
	storedCourses, err := d.GetCoursesByUID(uid)
	if err != nil {
		return
	}

	// Check if course already exist in stored course list
	for idx, course := range courseList {
		var same = false
		var old *model.Course
		for _, storedCourse := range storedCourses {
			// Use update instead of create new course when found same course
			if storedCourse.Title == course.Title {
				// When both courses are completely equivalent
				if storedCourse.LetterGrade == course.LetterGrade && storedCourse.PercentGrade == course.PercentGrade {
					same = true
					break
				}

				old = storedCourse
				course.CopyFromOther(storedCourse)
				break
			}
		}

		// Insert or update course
		if !same {
			count++
			if course.ID == 0 {
				d.InsertCourse(course)
			} else {
				d.UpdateCourse(old, course)
			}
		}

		// Asynchronously store assignments data
		wg.Add(1)
		courseTitle := course.Title
		assignments := assignmentsList[idx]
		go func() {
			count += StoreAssignmentsData(uid, courseTitle, assignments)
			wg.Done()
		}()
	}
	wg.Wait()
	if count != 0 {
		log.Printf("(Jupiter Data Crawling): [%v] DB Actions Done for User [%v]\n", count, uid)
	}

	// Update fetch time and GPA
	d.UpdateFetchTimeAndGPA(uid, gpa)
}

// Store all fetched assignments data to database
func StoreAssignmentsData(uid int, courseTitle string, assignments []*model.Assignment) int {
	var count = 0
	wg := &sync.WaitGroup{}

	// Get stored assignment list
	storedAssignments, _ := d.GetAssignmentsByCourseAndUID(uid, courseTitle)

	// List for new assignments
	var newAssignments []*model.Assignment

	// Check if assignment already exist in stored assignment list
	for _, assignment := range assignments {
		var same = false
		var old *model.Assignment
		for _, storedAssignment := range storedAssignments {
			// Use update instead of create new assignment when found same assignment
			if storedAssignment.Title == assignment.Title && storedAssignment.Due == assignment.Due && storedAssignment.From == assignment.From {
				// When both courses are completely equivalent
				if storedAssignment.Desc == assignment.Desc && storedAssignment.Score == assignment.Score && storedAssignment.Status == assignment.Status {
					same = true
				}

				old = storedAssignment
				assignment.CopyFromOther(storedAssignment)
				break
			}
		}

		// Insert or update assignment
		if !same {
			count++
			if assignment.ID == 0 {
				// If assignment is new put it into tmp list
				newAssignments = append(newAssignments, assignment)
			} else {
				d.UpdateAssignment(old, assignment)
			}
		}
	}

	// Too much new assignments, store them directly without description
	if len(newAssignments) > 5 { // TODO: 数字暂时的，确定负载后再改
		for _, assignment := range newAssignments {
			d.InsertAssignment(assignment)
		}
	} else {
		// Asynchronously fetch descriptions and store new assignments
		for _, assignment := range newAssignments {
			wg.Add(1)
			assignmentC := assignment
			go func() {
				assignmentWithDesc := FetchAssignmentDesc(assignmentC)
				d.InsertAssignment(assignmentWithDesc)
				wg.Done()
			}()
		}
		wg.Wait()
	}

	return count
}

// Fetch a student's GPA and report card image
func FetchReportAndGPA(page *rod.Page) string {
	// Navigate to report card page
	opts, _, err := NavGetOptions(page)
	if err != nil {
		return ""
	}
	if NavNavigate(page, opts[5]) != nil {
		return ""
	}

	// Get GPA and report card image
	return GetReportCardAndGPA(page, 1)
}

// Fetch assignment description
func FetchAssignmentDesc(assignment *model.Assignment) *model.Assignment {
	// TODO: Logic Implement
	return assignment
}
