package crawler

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/util"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

var d = dao.Jupiter
var browser *rod.Browser
var browserWithProxy *rod.Browser
var browserWithoutProxy *rod.Browser
var pagePool rod.PagePool
var pageCreate func() *rod.Page

var PagePoolLoad int
var PendingTaskCount int

var FetchDataRateLimiter = []int{}

var GptPage *rod.Page
var GptAccessToken string

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

	// Create a page for Pandora Next(GPT)
	if config.Config.GPT.Enable {
		go GetGptAccessToken()
	}

	// Set browser
	browser = browserWithProxy

	// Create page pool for multithreading
	pagePool = rod.NewPagePool(config.Config.Crawler.MaxParallel)
	pageCreate = func() *rod.Page {
		return stealth.MustPage(browser.MustIncognito())
	}

	// Immediately start one crawler job
	func() {
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
	}()

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
		browser = browserWithoutProxy
	case 1:
		browser = browserWithProxy
	default:
		return
	}
}

func GetGptAccessToken() {
	if GptPage == nil {
		GptPage = browserWithoutProxy.MustIncognito().MustPage()
		GptPage.Navigate(config.Config.GPT.URL)
	} else {
		GptPage.Close()
		time.Sleep(time.Millisecond * 250)
		GptPage = browserWithoutProxy.MustIncognito().MustPage()
		GptPage.Navigate(config.Config.GPT.URL)
	}

	go GptPage.HijackRequests().MustAdd("*/session", func(ctx *rod.Hijack) {
		res := make(map[string]string)
		ctx.MustLoadResponse()
		json.Unmarshal(ctx.Response.Payload().Body, &res)
		GptAccessToken = "Bearer " + res["accessToken"]
	}).Run()

	for i := 0; i < 3; i++ {
		GptPage.MustWaitLoad()

		rod.Try(func() {
			GptPage.Timeout(time.Second).MustElement("#username").Input(config.Config.GPT.Username)
			GptPage.Timeout(time.Second).MustElement("#password").Input(config.Config.GPT.Password)
		})

		time.Sleep(time.Second * 2)
		GptPage.Keyboard.Type(input.Key(13))

		if GptPage.Timeout(time.Second*30).WaitElementsMoreThan("img[alt*='User']", 0) == nil {
			break
		}

		GptPage.MustReload()
		util.Log("crawler", "INFO GPT Proxy Err")
	}
}

// Composite function for fetch and store data from Jupiter for each user
func CrawlerJob(uid ...int) {
	// When specific user ids are passed
	if len(uid) == 1 {
		// Rate limiter
		if util.IfExistInSlice(FetchDataRateLimiter, uid[0]) {
			return
		}
		FetchDataRateLimiter = util.Append(FetchDataRateLimiter, uid[0])
		defer func() {
			FetchDataRateLimiter = util.RemoveFromSlice(FetchDataRateLimiter, uid[0])
		}()

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
			func() {
				// Rate limiter
				if util.IfExistInSlice(FetchDataRateLimiter, id) {
					return
				}
				FetchDataRateLimiter = util.Append(FetchDataRateLimiter, id)
				defer func() {
					FetchDataRateLimiter = util.RemoveFromSlice(FetchDataRateLimiter, id)
				}()

				// Start job for single user
				startedAt := time.Now()

				// Fetch all data
				courseList, assignmentsList, gpa, err := FetchData(id)
				if err != nil {
					return
				}

				// Store fetched data to database
				StoreData(id, gpa, courseList, assignmentsList, &startedAt)
			}()
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
			// Rate limiter
			if util.IfExistInSlice(FetchDataRateLimiter, uid) {
				return
			}
			FetchDataRateLimiter = util.Append(FetchDataRateLimiter, uid)
			defer func() {
				FetchDataRateLimiter = util.RemoveFromSlice(FetchDataRateLimiter, uid)
			}()

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
		page = stealth.MustPage(browser).MustSetCookies()
	} else {
		time.Sleep(time.Minute / 2 * time.Duration(rand.Float32()))
		PendingTaskCount++
		page = pagePool.Get(pageCreate).MustSetCookies()
		PagePoolLoad++
		PendingTaskCount--
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
			PagePoolLoad--
			pagePool.Put(page)
			pages, _ := browser.Pages()
			for _, page := range pages {
				page.MustClose()
			}
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

// Get a page to that's logged inned already
func OpenPageAndLogin(uid int) (page *rod.Page, err error) {
	// Get a page to access Jupiter
	page, err = OpenJupiterPage(uid)
	if err != nil {
		return
	}

	// Find user's Jupiter account info
	data, _ := d.GetDataByUID(uid)

	// Login
	if err := Login(page, data.Account, data.Password); err != nil {
		return page, err
	}

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
	// Get a page to access Jupiter
	page, err := OpenPageAndLogin(uid)
	defer func() {
		PagePoolLoad--
		page.MustNavigate("about:blank")
		pagePool.Put(page)
	}()
	if err != nil {
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
		if ClickTarget(page, courses[idx]) != nil {
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

/*
conds[0] states isSingle for only one assignment being passed.

conds[1] states forceTurnInnedArrEmpty for copying method.
*/
func StoreAssignmentsData(uid int, courseTitle string, assignments []*model.Assignment, conds ...bool) int {
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
				if len(conds) > 1 && conds[1] {
					assignment.CopyFromOther(storedAssignment, true)
				} else {
					assignment.CopyFromOther(storedAssignment, false)
				}

				// When both courses are completely equivalent
				if storedAssignment.Desc == assignment.Desc &&
					storedAssignment.Score == assignment.Score &&
					storedAssignment.Status == assignment.Status &&
					storedAssignment.Feedback == assignment.Feedback &&
					storedAssignment.TurnInAble == assignment.TurnInAble &&
					reflect.DeepEqual(storedAssignment.TurnInnedList, assignment.TurnInnedList) {
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

				// Need an extra fetch for assignment that got new score for teacher feedbacks
				if old.Score != assignment.Score {
					assignmentC := assignment
					go func() {
						StoreAssignmentsData(uid, assignmentC.From, []*model.Assignment{FetchAssignmentDetail(uid, assignmentC, true)}, true)
					}()
				}
			}
		}
	}

	// Delete nonexisting assignments from database
	if len(conds) == 0 {
		for _, assignment := range storedAssignments {
			if assignment != nil {
				count++
				// Only allow deletion when the new data is retrieved in the same quarter
				if assignment.Quarter == assignments[0].Quarter {
					d.DeleteAssignment(assignment.ID)
				}
			}
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
				assignmentWithDetail := FetchAssignmentDetail(uid, tmp[1])
				d.InsertAssignment(tmp[0], assignmentWithDetail)
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
func FetchAssignmentDetail(uid int, assignment *model.Assignment, force ...bool) *model.Assignment {
	// Check if there is existing descriptions that's within expiration time
	storedAssignment, err := dao.Jupiter.GetAssignmentByInfo(assignment.Title, &assignment.Due, assignment.From)
	if err == nil && time.Now().Unix()-storedAssignment.DescFetchedAt.Unix() < 1800 && (len(force) == 0 || !force[0]) {
		if storedAssignment.Feedback == "" || assignment.Feedback != "" {
			assignment.Desc = storedAssignment.Desc
			assignment.TurnInAble = storedAssignment.TurnInAble
			assignment.TurnInTypes = storedAssignment.TurnInTypes
			return assignment
		}
	}

	// Get a page to access Jupiter
	page, err := OpenPageAndLogin(uid)
	defer func() {
		PagePoolLoad--
		page.MustNavigate("about:blank")
		pagePool.Put(page)
	}()
	if err != nil {
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
			err = ClickTarget(page, course)
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

		var title, dueDate string
		err := rod.Try(func() {
			title = el.MustElement(":nth-child(3)").MustText()
			dueDate = el.MustElement(":nth-child(2)").MustText()
		})
		if err != nil {
			return assignment
		}

		if strings.Contains(title, assignment.Title) && (assignment.Due.Year() == 1 || strings.Contains(dueDate, due)) {
			err := rod.Try(func() { el.MustClick() })
			if err != nil {
				return assignment
			}
			break
		}
	}

	// Get assignment details
	assignment.Desc = GetAssignmentDesc(page)
	assignment.Feedback = GetTeacherFeedback(page, uid, assignment.ID)
	assignment.TurnInAble = HasTurnIn(page)
	page.WaitStable(time.Millisecond * 100)
	if assignment.TurnInAble == 1 {
		assignment.TurnInTypes = GetTurnInTypes(page)
		assignment.TurnInnedList = GetTurnInnedList(page)
	}
	assignment.DescFetchedAt = time.Now()

	return assignment
}

// Turn in JunoDoc/Files for assignment
func TurnIn(uid int, id int, turnInType string, files ...string) error {
	// Get a page to access Jupiter
	page, err := OpenPageAndLogin(uid)
	defer func() {
		PagePoolLoad--
		page.MustNavigate("about:blank")
		pagePool.Put(page)
	}()
	if err != nil {
		return err
	}

	// Get courses from nav bar
	_, courses, err := NavGetOptions(page)
	if err != nil {
		return err
	}

	// Get assignment info from database
	assignment, err := d.GetAssignmentByID(id)
	if err != nil {
		return err
	}

	// Navigate to the course the assignment is located
	for _, course := range courses {
		if course.MustText() == assignment.From {
			err = ClickTarget(page, course)
			if err != nil {
				return err
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
		return err
	}

	// Find and click the targeted assignment
	for _, el := range elements {
		due := strconv.Itoa(int(assignment.Due.Month())) + "/" + strconv.Itoa(assignment.Due.Day())

		var title, dueDate string
		err := rod.Try(func() {
			title = el.MustElement(":nth-child(3)").MustText()
			dueDate = el.MustElement(":nth-child(2)").MustText()
		})
		if err != nil {
			return err
		}

		if strings.Contains(title, assignment.Title) && (assignment.Due.Year() == 1 || strings.Contains(dueDate, due)) {
			err := rod.Try(func() { el.MustClick() })
			if err != nil {
				return err
			}
			break
		}
	}

	// Starting turn in process
	if turnInType == "JunoDoc" {
		// Click `Turn In`
		WaitStable(page)
		err := rod.Try(func() {
			page.Timeout(time.Second*2).MustElementR("div.btn", "/^Turn In/").MustClick()
		})
		if err != nil {
			return err
		}

		// Click `New Juno Doc`
		page.WaitStable(time.Millisecond * 100)
		err = rod.Try(func() {
			page.Timeout(time.Millisecond * 200).MustElement("tr[click*='picknewtext()']").MustClick()
		})
		if err != nil {
			return err
		}
		WaitStable(page)

		// Enter the title and text
		text := strings.Split(files[0], "|")
		err = rod.Try(func() {
			page.Timeout(time.Millisecond * 200).MustElement("#text_title").MustInput(text[0])
			page.Timeout(time.Millisecond * 200).MustElement("#text_writetext").MustInput(text[1])
		})
		if err != nil {
			return err
		}

		// Click turn in
		err = rod.Try(func() {
			page.Timeout(time.Millisecond*200).MustElementR("div.btn", "/^Turn In/").MustClick()
		})
		if err != nil {
			return err
		}
	} else {
		for _, path := range files {
			// Click `Turn In`
			WaitStable(page)
			err := rod.Try(func() {
				page.Timeout(time.Second*2).MustElementR("div.btn", "/^Turn In/").MustClick()
			})
			if err != nil {
				return err
			}

			// Upload file
			page.WaitStable(time.Millisecond * 100)
			err = rod.Try(func() {
				page.Timeout(time.Millisecond * 200).MustElement("input[onchange*='uploadfiles(this)']").SetFiles([]string{path})
			})
			if err != nil {
				return err
			}
		}
	}

	// Update turn inned list
	assignment.TurnInnedList = GetTurnInnedList(page)
	StoreAssignmentsData(uid, assignment.From, []*model.Assignment{assignment}, true)

	return nil
}

// Un-submit JunoDoc/Files for assignment
func UnSubmit(uid int, id int, name string) error {
	// Get a page to access Jupiter
	page, err := OpenPageAndLogin(uid)
	defer func() {
		PagePoolLoad--
		page.MustNavigate("about:blank")
		pagePool.Put(page)
	}()
	if err != nil {
		return err
	}

	// Get nav options
	opts, _, err := NavGetOptions(page)
	if err != nil {
		return err
	}

	// Go to `My Files`
	err = ClickTarget(page, opts[2])
	if err != nil {
		return err
	}

	// Show all submitted works
	err = rod.Try(func() {
		page.Timeout(time.Millisecond*200).MustElementR("div.btnl", "/^Show All/").MustClick()
	})
	if err != nil {
		return err
	}

	// Un-submit the targeted work
	WaitStable(page)
	assignment, err := d.GetAssignmentByID(id)
	if err != nil {
		return err
	}

	var elList rod.Elements
	err = rod.Try(func() {
		elList = page.Timeout(time.Millisecond * 200).MustElements("tr[val]")
	})
	if err != nil {
		return err
	}

	// Remove files from `Not Turned In` section
	elList = slices.DeleteFunc(elList, func(el *rod.Element) bool {
		return strings.Contains(el.MustHTML(), "mb</td>")
	})

	var target *rod.Element
	for _, el := range elList {
		if el.MustElement("td:nth-child(2)").MustText() == name {
			target = el
			break
		}
	}
	if target == nil {
		return errors.New("targetNotFound")
	}

	// Click `Delete`
	err = rod.Try(func() {
		target.MustClick()

		page.WaitStable(time.Millisecond * 100)
		page.Timeout(time.Millisecond * 200).MustElement("#deletebtn").MustClick()

		page.WaitStable(time.Millisecond * 100)
		page.Timeout(time.Millisecond*100).MustElementR("#promptdelete > div > div > div.btn", "/^Unsubmit/").MustClick()
	})
	if err != nil {
		return err
	}

	//* -------------------- For updating assignment detail -------------------- *//
	// Get courses from nav bar
	_, courses, err := NavGetOptions(page)
	if err != nil {
		return err
	}

	// Navigate to the course the assignment is located
	for _, course := range courses {
		if course.MustText() == assignment.From {
			err = ClickTarget(page, course)
			if err != nil {
				return err
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
		return err
	}

	// Find and click the targeted assignment
	for _, el := range elements {
		due := strconv.Itoa(int(assignment.Due.Month())) + "/" + strconv.Itoa(assignment.Due.Day())

		var title, dueDate string
		err := rod.Try(func() {
			title = el.MustElement(":nth-child(3)").MustText()
			dueDate = el.MustElement(":nth-child(2)").MustText()
		})
		if err != nil {
			return err
		}

		if strings.Contains(title, assignment.Title) && (assignment.Due.Year() == 1 || strings.Contains(dueDate, due)) {
			err := rod.Try(func() { el.MustClick() })
			if err != nil {
				return err
			}
			break
		}
	}

	// Update turn inned list
	assignment.TurnInnedList = GetTurnInnedList(page)
	StoreAssignmentsData(uid, assignment.From, []*model.Assignment{assignment}, true, true)

	return nil
}
