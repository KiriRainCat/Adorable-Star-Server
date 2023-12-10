package crawler

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/util"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var d = dao.Jupiter
var browser *rod.Browser
var browserWithProxy *rod.Browser
var browserWithoutProxy *rod.Browser
var pagePool rod.PagePool
var pageCreate func() *rod.Page

var PagePoolLoad int
var PendingTaskCount int

// Initialize crawler with page pool to execute tasks asynchronously
func Init() {
	// Find browser executable path
	bin, _ := launcher.LookPath()

	// With proxy
	browserWithProxy = rod.New().ControlURL(
		launcher.
			New().
			Bin(bin).
			Proxy("localhost:" + strconv.Itoa(config.Config.Crawler.ProxyPort)).
			MustLaunch()).
		MustConnect()

	// Without proxy
	browserWithoutProxy = rod.New().ControlURL(
		launcher.
			New().
			Bin(bin).
			MustLaunch()).
		MustConnect()

	// Set browser
	browser = browserWithProxy

	// Create page pool for multithreading
	pagePool = rod.NewPagePool(config.Config.Crawler.MaxParallel)
	pageCreate = func() *rod.Page {
		return browser.MustIncognito().MustPage()
	}

	// Immediately start one crawler job
	{
		users, err := dao.User.GetActiveUsers()
		if err != nil {
			return
		}

		var ids []int
		for _, user := range users {
			if user.Status >= 100 {
				ids = append(ids, user.ID)
			}
		}

		CrawlerJob(ids...)
		CrawlerJob()
	}

	// Start scheduled crawler job
	t2 := time.NewTicker(time.Minute * time.Duration(config.Config.Crawler.FetchInterval-5))
	go func() {
		for range t2.C {
			users, err := dao.User.GetActiveUsers()
			if err != nil {
				return
			}

			var ids []int
			for _, user := range users {
				if user.Status >= 100 {
					ids = append(ids, user.ID)
				}
			}

			CrawlerJob(ids...)
			CrawlerJob()
			t2.Reset(time.Minute * time.Duration(config.Config.Crawler.FetchInterval-5))
		}
	}()
}

func SwitchBrowser(id int) {
	switch id {
	case 0:
		browser = browserWithProxy
	case 1:
		browser = browserWithoutProxy
	default:
		browser = browserWithoutProxy
	}
}

// Composite function for fetch and store data from Jupiter for each user
func CrawlerJob(uid ...int) {
	// When specific user ids are passed
	if len(uid) == 1 {
		// Start job for single user
		startedAt := time.Now()

		// Fetch all data
		courseList, assignmentsList, gpa, err := FetchData(uid[0])
		if err != nil {
			return
		}

		// Store fetched data to database
		StoreData(uid[0], gpa, courseList, assignmentsList, &startedAt)
		return
	} else if len(uid) > 0 {
		for _, id := range uid {
			// Start job for single user
			startedAt := time.Now()

			// Fetch all data
			courseList, assignmentsList, gpa, err := FetchData(id)
			if err != nil {
				return
			}

			// Store fetched data to database
			StoreData(id, gpa, courseList, assignmentsList, &startedAt)
		}
		return
	}

	// Get all active users
	users, err := dao.User.GetActiveUsers()
	if err != nil {
		return
	}

	// Loop through and start crawler job for users
	wg := sync.WaitGroup{}
	for _, user := range users {
		if user.Status >= 100 {
			continue
		}

		wg.Add(1)
		uid := user.ID
		go func() {
			startedAt := time.Now()

			// Fetch all data
			courseList, assignmentsList, gpa, err := FetchData(uid)
			if err != nil {
				return
			}

			// Store fetched data to database
			StoreData(uid, gpa, courseList, assignmentsList, &startedAt)
			wg.Done()
		}()
	}
	wg.Wait()
}

// Open a webpage for Jupiter
func OpenJupiterPage(uid int, notPool ...bool) (page *rod.Page, err error) {
	// Whether using page pool
	if len(notPool) > 0 && notPool[0] {
		page = browser.MustIncognito().MustPage().MustSetCookies()
	} else {
		time.Sleep(time.Minute / 2 * time.Duration(rand.Float32()))
		PendingTaskCount++
		page = pagePool.Get(pageCreate).MustSetCookies()
		PendingTaskCount--
		PagePoolLoad++
	}

	// Bypass Cloudflare detection for crawler
	page.MustEvalOnNewDocument("const newProto = navigator.__proto__;delete newProto.webdriver;navigator.__proto__ = newProto;")

	// Try to bypass cloudflare
	dataList, _ := d.GetNewestCfbp()
	for i := 0; i < len(dataList); i++ {
		// Add cloudflare bypass Cookie
		data := dataList[i]
		page.MustSetExtraHeaders("Cookie", "cfbp="+data.Cfbp)

		// Navigate to Jupiter Ed login page
		err = page.Navigate("https://login.jupitered.com/login/")
		if err != nil {
			// Notify developer
			dao.Message.Insert(&model.Message{
				UID:  1,
				Type: -1,
				Msg:  "browserProxyErr",
			})

			// Switch to browser without proxy
			page.MustClose()
			browser = browserWithoutProxy
			return OpenJupiterPage(uid)
		}

		// Check if request blocked by Cloudflare
		text := ""
		rod.Try(func() {
			text = page.Timeout(time.Second).MustElement("body > div > div").MustText()
		})
		if !strings.Contains(text, "malicious") {
			dao.User.UpdateCfbp(data.UID, data.Cfbp)
			return
		}
	}

	// Do this when cfbp list have no element
	err = page.Navigate("https://login.jupitered.com/login/")
	if err == nil {
		return
	}

	dao.Message.Insert(&model.Message{
		UID:  uid,
		Type: -1,
		Msg:  "cfToken",
	})
	return page, errors.New("requestBlocked")
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
	defer func() {
		PagePoolLoad--
		page.MustNavigate("about:blank")
		pagePool.Put(page)
	}()
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
	for idx := 0; idx < len(courses); idx++ {
		// Nav bar navigation to course page
		_, courses, err := NavGetOptions(page)
		if err != nil {
			idx--
			continue
		}
		var courseName string
		if rod.Try(func() { courseName = courses[idx].Timeout(time.Second * 2).MustText() }) != nil {
			idx--
			continue
		}
		if NavNavigate(page, courses[idx]) != nil {
			idx--
			continue
		}

		// Get grade
		courseList = append(courseList, GetCourseGrade(page, courseName, uid))

		// Get all assignments
		assignmentsList = append(assignmentsList, GetCourseAssignments(page, courseName, uid))
	}

	// Fetch GPA and report card image
	gpa = FetchReportAndGPA(page, uid)

	return
}

// Store all fetched data to database
func StoreData(uid int, gpa string, courseList []*model.Course, assignmentsList [][]*model.Assignment, startedAt *time.Time) {
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
				if strings.Contains(storedCourse.LetterGrade, course.LetterGrade) && strings.Contains(storedCourse.PercentGrade, course.PercentGrade) {
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
		now := time.Now()
		diff := (now.Minute()*60 + now.Second()) - (startedAt.Minute()*60 + startedAt.Second())
		util.Log("crawler", "INFO [%v] DB Actions for User [%v] (%vs)", count, uid, diff)
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
	var newAssignments [][]*model.Assignment

	// Check if assignment already exist in stored assignment list
	for _, assignment := range assignments {
		var same = false
		var old *model.Assignment
		for idx, storedAssignment := range storedAssignments {
			if storedAssignment == nil {
				continue
			}

			// Use update instead of create new assignment when found same assignment
			if storedAssignment.Title == assignment.Title && storedAssignment.Due.YearDay() == assignment.Due.YearDay() {
				old = storedAssignment
				assignment.CopyFromOther(storedAssignment)

				// When both courses are completely equivalent
				if storedAssignment.Desc == assignment.Desc && storedAssignment.Score == assignment.Score && storedAssignment.Status == assignment.Status {
					same = true
				}
				storedAssignments[idx] = nil
				break
			}
		}

		// Insert or update assignment
		if !same && assignment.Title != "" {
			count++
			if assignment.ID == 0 {
				// If assignment is new put it into tmp list
				if assignment != nil {
					newAssignments = append(newAssignments, []*model.Assignment{old, assignment})
				}
			} else {
				d.UpdateAssignment(old, assignment)
			}
		}
	}

	// Delete nonexisting assignments from database
	for _, assignment := range storedAssignments {
		if assignment != nil {
			count++
			d.DeleteAssignment(assignment.ID)
		}
	}

	// Too much new assignments, store them directly without description
	if len(newAssignments) > 5 {
		for _, assignment := range newAssignments {
			d.InsertAssignment(assignment[0], assignment[1])
		}
	} else {
		// Asynchronously fetch descriptions and store new assignments
		for _, assignment := range newAssignments {
			wg.Add(1)
			tmp := assignment
			go func() {
				assignmentWithDesc := FetchAssignmentDesc(uid, tmp[1])
				d.InsertAssignment(tmp[0], assignmentWithDesc)
				wg.Done()
			}()
		}
		wg.Wait()
	}

	return count
}

// Fetch a student's GPA and report card image
func FetchReportAndGPA(page *rod.Page, uid int) string {
	// Navigate to report card page
	opts, _, err := NavGetOptions(page)
	if err != nil {
		return ""
	}
	opts[5].MustClick()

	// Get GPA and report card image
	return GetReportCardAndGPA(page, uid)
}

// Fetch assignment description
func FetchAssignmentDesc(uid int, assignment *model.Assignment) *model.Assignment {
	// Check if there is existing descriptions that's within expiration time
	storedAssignment, err := dao.Jupiter.GetAssignmentByInfo(assignment.Title, &assignment.Due, assignment.From)
	if err == nil && time.Now().Unix()-storedAssignment.DescFetchedAt.Unix() < 3600 {
		assignment.Desc = storedAssignment.Desc
		return assignment
	}

	// Get a page to access Jupiter
	page, err := OpenJupiterPage(uid)
	defer func() {
		PagePoolLoad--
		page.MustNavigate("about:blank")
		pagePool.Put(page)
	}()
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
	assignment.DescFetchedAt = time.Now()

	return assignment
}
