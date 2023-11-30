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
	for _, el := range data {
		assignment := &model.Assignment{UID: uid, From: courseName}

		rod.Try(func() {
			if date := FormatJupiterDueDate(el.Timeout(time.Second * 2).MustElement("td:nth-child(2)").MustText()); date != "" {
				due, _ := time.Parse("2006-01-02", date)
				assignment.Due = due
			}
		})
		rod.Try(func() {
			assignment.Title = el.Timeout(time.Second * 2).MustElement("td:nth-child(3)").MustText()
		})
		rod.Try(func() {
			assignment.Score = el.MustElement("td:nth-child(4)").MustText()
		})

		assignments = append(assignments, assignment)
	}

	return
}

// Get description for an assignment
func GetAssignmentDesc(page *rod.Page) string {
	WaitStable(page)

	var desc string
	err := rod.Try(func() {
		desc = page.Timeout(time.Second*2).MustElementR("#mainpage > div", "/Directions/").MustElement("div").MustHTML()
	})
	if err != nil {
		return ""
	}

	desc = strings.ReplaceAll(desc, " style=\"padding:0px 20px; max-width:472px;\"", "")
	desc = strings.ReplaceAll(desc, "<b>Directions</b><br>", "")
	return desc
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
	page.MustElement("table.bord > tbody").MustScreenshot(util.GetCwd() + "/storage/img/report/" + strconv.Itoa(uid) + ".png")

	return
}
