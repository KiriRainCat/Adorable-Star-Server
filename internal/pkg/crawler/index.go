package crawler

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/config"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
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
		browser = rod.New().ControlURL(config.Config.Crawler.BrowserSocketUrl).MustConnect()
	} else {
		browser = rod.New().MustConnect()
	}

	// Create page pool for multithreading
	pagePool = rod.NewPagePool(config.Config.Crawler.MaxParallel)
	pageCreate = func() *rod.Page {
		return browser.MustIncognito().MustPage()
	}
	pagePool.Put(pagePool.Get(pageCreate))

	// Start scheduled crawler job
	go CrawlerJob()
	t := time.NewTicker(time.Minute * time.Duration(config.Config.Crawler.FetchInterval))

	// For every 30 minutes, fetch jupiter data for user
	go func() {
		for range t.C {
			go CrawlerJob()
		}
	}()
}

// Composite function for fetch and store data from Jupiter for each user
func CrawlerJob(uid ...int) {
	// When single user called api
	if len(uid) > 0 {
		// Start job for single user
		go func() {
			// Fetch all data
			courseList, assignmentsList, gpa, err := FetchData(uid[0])
			if err != nil {
				return
			}

			// Store fetched data to database
			StoreData(uid[0], gpa, courseList, assignmentsList)
		}()
	}

	// Get all active users
	users, err := dao.User.GetActiveUsers()
	if err != nil {
		return
	}

	// Loop through and start crawler job for users
	for _, user := range users {
		uid := user.ID
		go func() {
			// Fetch all data
			courseList, assignmentsList, gpa, err := FetchData(uid)
			if err != nil {
				return
			}

			// Store fetched data to database
			StoreData(uid, gpa, courseList, assignmentsList)
		}()
	}
}

// Open a webpage for Jupiter
func OpenJupiterPage(uid int, notPool ...bool) (page *rod.Page, err error) {
	// Whether using page pool
	if len(notPool) > 0 && notPool[0] {
		page = browser.MustIncognito().MustPage().MustSetCookies()
	} else {
		page = pagePool.Get(pageCreate).MustSetCookies()
	}

	// Navigate to Jupiter Ed login page
	err = page.Navigate("https://login.jupitered.com/login/")
	if err != nil {
		return
	}

	// Bypass Cloudflare detection for crawler
	page.MustEvalOnNewDocument("const newProto = navigator.__proto__;delete newProto.webdriver;navigator.__proto__ = newProto;")

	// TODO: 记得以后再开启 cloudflare 拦截检测
	// // Check if request blocked by Cloudflare
	// if strings.Contains(page.MustElement("body").MustText(), "malicious") {
	// 	dao.Message.Insert(&model.Message{
	// 		UID:  uid,
	// 		Type: -1,
	// 		Msg:  "cfToken",
	// 	})
	// }

	return
}

// Verify if a Jupiter Ed account is valid
func VerifyAccount(uid int, account string, pwd string) error {
	// Open a Jupiter page that's not affected by page pool size
	page, err := OpenJupiterPage(uid, true)
	defer page.MustClose()
	if err != nil {
		return err
	}

	// Check account
	if err := Login(page, account, pwd); err != nil {
		return err
	}

	return nil
}

// Fetch all Jupiter data for a student
func FetchData(uid int) (courseList []*model.Course, assignmentsList [][]*model.Assignment, gpa string, err error) {
	// Find user's Jupiter account info
	data, err := d.GetDataByUID(uid)
	if err != nil {
		return
	}

	// Get a page to access Jupiter
	page, err := OpenJupiterPage(uid)
	defer pagePool.Put(page)
	if err != nil {
		return
	}

	// Login
	if err = Login(page, data.Account, data.Password); err != nil {
		return
	}

	// Get courses from nav bar
	_, courses, err := NavGetOptions(page)
	if err != nil {
		return
	}
	if rod.Try(func() { page.Timeout(time.Second * 2).MustElement("#touchnavbtn").MustClick() }) != nil {
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
				if strings.Contains(course.LetterGrade, storedCourse.LetterGrade) && strings.Contains(course.PercentGrade, storedCourse.PercentGrade) {
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
		log.Printf("INFO [%v] DB Actions Done for User [%v]\n", count, uid)
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
				assignmentWithDesc := FetchAssignmentDesc(uid, assignmentC)
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
	opts[5].MustClick()

	// Get GPA and report card image
	return GetReportCardAndGPA(page, 1)
}

// Fetch assignment description
func FetchAssignmentDesc(uid int, assignment *model.Assignment) *model.Assignment {
	// Get a page to access Jupiter
	page, err := OpenJupiterPage(uid)
	defer pagePool.Put(page)
	if err != nil {
		return assignment
	}

	// Find user's Jupiter account info
	data, _ := d.GetDataByUID(uid)

	// Login
	if err := Login(page, data.Account, data.Password); err != nil {
		return assignment
	}

	// Get courses from nav bar
	_, courses, err := NavGetOptions(page)
	if err != nil {
		return assignment
	}

	// Navigate to the course the assignment is located
	for _, course := range courses {
		if course.MustText() == assignment.From {
			err = NavNavigate(page, course)
			if err != nil {
				return assignment
			}
			break
		}
	}

	// Get course assignments
	WaitStable(page)
	var elements rod.Elements
	err = rod.Try(func() {
		elements = page.Timeout(time.Second * 2).MustElements("table > tbody[click*='goassign'] > tr:nth-child(2)")
	})
	if err != nil {
		return assignment
	}

	// Find and click the targeted assignment
	for _, el := range elements {
		due := strconv.Itoa(int(assignment.Due.Month())) + "/" + strconv.Itoa(assignment.Due.Day())
		if strings.Contains(el.MustElement(":nth-child(3)").MustText(), assignment.Title) && (assignment.Due.Year() == 1 || strings.Contains(el.MustElement(":nth-child(2)").MustText(), due)) {
			el.MustClick()
			break
		}
	}

	// Get assignment desc
	desc := GetAssignmentDesc(page)
	assignment.Desc = desc

	return assignment
}
