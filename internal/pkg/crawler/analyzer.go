package crawler

import (
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/util"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Parse raw date "9/2" into something like "2023-09-02"
func FormatJupiterDueDate(raw string) string {
	if raw == "" {
		return "0001-01-01"
	}
	if strings.Contains(raw, "-") {
		return raw
	}

	parts := strings.Split(raw, "/")
	if len(parts[0]) < 2 {
		parts[0] = "0" + parts[0]
	}
	if len(parts[1]) < 2 {
		parts[1] = "0" + parts[1]
	}

	// If the date is back in time for 4 month, change the year to next year
	if month, _ := strconv.Atoi(parts[0]); int(time.Now().Month())-4 > month {
		return strconv.Itoa(time.Now().Year()+1) + "-" + parts[0] + "-" + parts[1]
	}

	return strconv.Itoa(time.Now().Year()) + "-" + parts[0] + "-" + parts[1]
}

// Use the current page of course to crawl course grade
func GetCourseGrade(page *rod.Page, courseName string, uid int) *model.Course {
	WaitStable(page)

	course := &model.Course{Title: courseName, UID: uid}

	el, err := page.Timeout(time.Second).Element("table > tbody > tr.baseline.botline.printblue")
	if err != nil {
		return course
	}

	rod.Try(func() {
		course.LetterGrade = el.Timeout(time.Second * 2).MustElement("td:nth-child(2)").MustText()
	})
	rod.Try(func() {
		course.PercentGrade = el.Timeout(time.Second * 2).MustElement("td:nth-child(3)").MustText()
	})

	return course
}

// Use the current page of course to crawl all assignments
func GetCourseAssignments(page *rod.Page, courseName string, uid int) (assignments []*model.Assignment) {
	WaitStable(page)

	// Get course assignments
	var data rod.Elements
	err := rod.Try(func() {
		data = page.Timeout(time.Second * 2).MustElements("table > tbody[click*='goassign'] > tr:nth-child(2)")
	})
	if err != nil {
		return
	}

	// Get information about each assignment
	for idx, el := range data {
		assignment := &model.Assignment{UID: uid, From: courseName}

		err := rod.Try(func() {
			due, _ := time.Parse("2006-01-02", FormatJupiterDueDate(el.Timeout(time.Second*2).MustElement("td:nth-child(2)").MustText()))
			assignment.Due = due
			assignment.Title = el.Timeout(time.Second * 2).MustElement("td:nth-child(3)").MustText()
			assignment.Score = el.MustElement("td:nth-child(4)").MustText()
		})

		// Prevent element temporary nil pointer resolving
		if err != nil {
			rod.Try(func() {
				el = page.Timeout(time.Second * 2).MustElement("table > tbody[click*='goassign']:nth-child(" + strconv.Itoa(idx+2) + ") > tr:nth-child(2)")
				due, _ := time.Parse("2006-01-02", FormatJupiterDueDate(el.Timeout(time.Second*2).MustElement("td:nth-child(2)").MustText()))
				assignment.Due = due
				assignment.Title = el.Timeout(time.Second * 2).MustElement("td:nth-child(3)").MustText()
				assignment.Score = el.MustElement("td:nth-child(4)").MustText()
			})
		}

		assignments = append(assignments, assignment)
	}

	return
}

// Get description for an assignment
func GetAssignmentDesc(page *rod.Page) (desc string) {
	WaitStable(page)

	err := rod.Try(func() {
		desc = page.Timeout(time.Second * 2).MustElement("#mainpage > div[class*='selectable wrap']").MustHTML()
	})
	if err != nil {
		return ""
	}

	desc = strings.ReplaceAll(desc, " style=\"display:block; max-width:472px; padding:12px 0px 12px 0px\"", "")
	desc = strings.ReplaceAll(desc, " style=\"padding:0px 20px; max-width:472px;\"", "")
	desc = strings.ReplaceAll(desc, " style=\"max-width:472px;\"", "")
	desc = strings.ReplaceAll(desc, "<b>Directions</b><br>", "")
	return desc
}

// Get teacher feedbacks for an assignment
func GetTeacherFeedback(page *rod.Page, uid int, id int) (feedback string) {
	WaitStable(page)

	err := rod.Try(func() {
		feedback = page.Timeout(time.Second * 2).MustElement("div:nth-child(3) > div > div:nth-child(8)").MustText()
	})
	if err != nil {
		return ""
	}

	rod.Try(func() {
		page.Timeout(time.Second).MustElement("div.momentum").
			MustScreenshot(util.GetCwd() + "/storage/" + strconv.Itoa(uid) + "/feedback/" + strconv.Itoa(id) + ".png")
	})

	return
}

// Check whether the assignment has turn in button
func HasTurnIn(page *rod.Page) int {
	WaitStable(page)

	err := rod.Try(func() {
		page.Timeout(time.Millisecond*200).MustElementR("div.btn", "/^Turn In/")
	})

	if err == nil {
		return 1
	}

	return -1
}

// Get turn in able types
func GetTurnInTypes(page *rod.Page) (list []string) {
	WaitStable(page)

	err := rod.Try(func() {
		page.Timeout(time.Millisecond*200).MustElementR("div.btn", "/^Turn In/").MustClick()
	})
	if err != nil {
		return
	}

	err = rod.Try(func() {
		page.Timeout(time.Millisecond * 200).MustElement("input[onchange*='uploadfiles(this)']")
	})
	if err == nil {
		list = append(list, "Files")
	}

	err = rod.Try(func() {
		page.Timeout(time.Millisecond * 200).MustElement("tr[click*='picknewtext()']")
	})
	if err == nil {
		list = append(list, "Juno Doc")
	}

	return
}

// Get the turn inned list
func GetTurnInnedList(page *rod.Page) (list []string) {
	WaitStable(page)

	var elList rod.Elements
	err := rod.Try(func() {
		elList = page.Timeout(time.Millisecond * 200).MustElements("div[dblclick*='downloadopen(clickval)'] > table > tbody > tr > td:nth-child(2)")
	})
	if err != nil {
		return
	}

	for _, el := range elList {
		list = append(list, el.MustText())
	}

	return
}

// Use the current page of report card to crawl GPA and report card image
func GetReportCardAndGPA(page *rod.Page, uid int) (gpa string) {
	WaitStable(page)

	// Get newest GPA
	err := rod.Try(func() {
		gpa = page.Timeout(time.Second * 2).MustElement("tr.blue.topbotline > td:last-child").MustText()
	})
	if err != nil {
		return ""
	}

	// Take a screenshot of the report card section
	page.MustElement("table.bord > tbody").MustScreenshot(util.GetCwd() + "/storage/" + strconv.Itoa(uid) + "/report.png")

	return
}
