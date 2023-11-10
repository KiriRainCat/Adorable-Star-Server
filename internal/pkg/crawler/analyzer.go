package crawler

import (
	"adorable-star/internal/global"
	"adorable-star/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Parse raw date "9/2" into something like "2023-09-02"
func FormatJupiterDueDate(raw string) string {
	if raw == "" {
		return ""
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
	return strconv.Itoa(time.Now().Year()) + "-" + parts[0] + "-" + parts[1]
}

// Use the current page of course to crawl course grade
func GetCourseGrade(page *rod.Page, courseName string, uid int) *model.Course {
	WaitStable(page)
	el, err := page.Timeout(time.Millisecond * 100).Element("table > tbody > tr.baseline.botline.printblue")
	if err != nil {
		return &model.Course{Title: courseName, UID: uid}
	}

	return &model.Course{
		UID:          uid,
		Title:        courseName,
		LetterGrade:  el.MustElement(":nth-child(2)").MustText(),
		PercentGrade: el.MustElement(":nth-child(3)").MustText(),
	}
}

// Use the current page of course to crawl all assignments
func GetCourseAssignments(page *rod.Page, courseName string, uid int) (assignments []*model.Assignment) {
	// Get course assignments
	var data rod.Elements
	err := rod.Try(func() {
		data = page.Timeout(time.Second * 2).MustElements("table > tbody[click*='goassign'] > tr:nth-child(2)")
	})
	if err != nil {
		return
	}

	// Get information about each assignment
	for _, el := range data {
		var assignment *model.Assignment

		// Prevent dead lock
		err := rod.Try(func() {
			due, _ := time.Parse("2006-01-02", FormatJupiterDueDate(el.Timeout(time.Second*2).MustElement(":nth-child(2)").MustText()))
			assignment = &model.Assignment{
				UID:   uid,
				From:  courseName,
				Due:   due,
				Title: el.Timeout(time.Second * 2).MustElement(":nth-child(3)").MustText(),
				Score: el.Timeout(time.Second * 2).MustElement(":nth-child(4)").MustText(),
			}
		})

		// Only add assignment when no error occurred
		if err == nil {
			assignments = append(assignments, assignment)
		}
	}

	return
}

// Get description for an assignment
func GetAssignmentDesc(page *rod.Page) string {
	WaitStable(page)

	var desc string
	err := rod.Try(func() {
		desc = page.Timeout(time.Second * 2).MustElement("#mainpage >div:nth-child(6) > div").MustText()
	})
	if err != nil {
		return ""
	}

	return strings.Replace(desc, "Directions\n", "", 1)
}

// Use the current page of report card to crawl GPA and report card image
func GetReportCardAndGPA(page *rod.Page, uid int) (gpa string) {
	WaitStable(page, 800)

	// Get newest GPA
	err := rod.Try(func() {
		gpa = page.Timeout(time.Second * 2).MustElement("tr.blue.topbotline > td:last-child").MustText()
	})
	if err != nil {
		return ""
	}

	// Take a screenshot of the report card section
	page.MustElement("table.bord > tbody").MustScreenshot(global.GetCwd() + "/storage/img/report/" + strconv.Itoa(uid) + ".png")

	return
}
