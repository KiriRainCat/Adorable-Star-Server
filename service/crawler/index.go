package crawler

import (
	"adorable-star/dao"
	"adorable-star/model"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var browser *rod.Browser
var pagePool rod.PagePool
var pageCreate func() *rod.Page
var d *dao.JupiterDAO

var TaskPool []func(...any) any

// Initialize crawler with page pool to execute tasks asynchronously
func Init(jupiterDao *dao.JupiterDAO) {
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

	// Get dao deps
	d = jupiterDao
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
func FetchData(uid int) (courseList []model.Course, assignmentsList [][]*model.Assignment, err error) {
	// Get a page to access Jupiter
	page := OpenJupiterPage()
	defer pagePool.Put(page)

	// Find user's Jupiter account info
	data, _ := d.GetDataByUID(uid)

	// Login
	err = Login(page, data.Account, data.Password)
	if err != nil {
		return
	}

	// Get courses from nav bar
	_, courses := NavGetOptions(page)
	page.MustElement("#touchnavbtn").MustClick()

	// Navigate through all of the courses to fetch course data
	for idx := range courses {
		// Nav bar navigation to course page
		_, courses := NavGetOptions(page)
		courseName := courses[idx].MustText()
		NavNavigate(page, courses[idx])

		// Get grade
		courseList = append(courseList, *GetCourseGrade(page, courseName))

		// Get all assignments
		assignmentsList = append(assignmentsList, GetCourseAssignments(page, courseName))
	}

	// Fetch GPA and report card image
	FetchReportAndGPA(page)

	return
}

// Fetch multiple assignments' description
func FetchAssignmentsDesc(page *rod.Page, ids []int) error {
	return nil
}

// Fetch a student's GPA and report card image
func FetchReportAndGPA(page *rod.Page) error {
	// Navigate to report card page
	opts, _ := NavGetOptions(page)
	NavNavigate(page, opts[5])

	// Get GPA and report card image
	GetReportCardAndGPA(page, 1)

	return nil
}
