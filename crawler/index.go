package crawler

import (
	"adorable-star/model"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

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

// Open a webpage for Jupiter
func OpenJupiterPage() *rod.Page {
	// Navigate to Jupiter Ed login page
	page := pagePool.Get(pageCreate).MustSetCookies().MustNavigate("https://login.jupitered.com/login/")

	// Bypass Cloudflare detection for crawler
	page.MustEvalOnNewDocument("const newProto = navigator.__proto__;delete newProto.webdriver;navigator.__proto__ = newProto;")

	return page
}

// Fetch all Jupiter data for a student
func FetchData(name string, pwd string) (courseList []model.Course, assignmentsList [][]model.Assignment, err error) {
	// Get a page to access Jupiter
	page := OpenJupiterPage()
	defer pagePool.Put(page)

	// Login
	err = Login(page, name, pwd)
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

	return
}

func FetchAssignmentsDesc(ids []int) error {
	return nil
}
