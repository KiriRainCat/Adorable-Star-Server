package crawler

import (
	"adorable-star/model"
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
func GetCourseGrade(page *rod.Page, courseName string) *model.Course {
	WaitStable(page)
	el, err := page.Timeout(time.Millisecond * 100).Element("table > tbody > tr.baseline.botline.printblue")
	if err != nil {
		return &model.Course{Title: courseName}
	}

	return &model.Course{
		Title:        courseName,
		LetterGrade:  el.MustElement(":nth-child(2)").MustText(),
		PercentGrade: el.MustElement(":nth-child(3)").MustText(),
	}
}

// Use the current page of course to crawl all assignments
func GetCourseAssignments(page *rod.Page, courseName string) (assignments []*model.Assignment) {
	// Get course assignments
	data := page.MustElements("table > tbody[click*='goassign'] > tr:nth-child(2)")
	for _, assignment := range data {
		due, _ := time.Parse("2006-01-02", FormatJupiterDueDate(assignment.MustElement(":nth-child(2)").MustText()))
		assignments = append(assignments,
			&model.Assignment{
				From:  courseName,
				Due:   due,
				Title: assignment.MustElement(":nth-child(3)").MustText(),
				Score: assignment.MustElement(":nth-child(4)").MustText(),
			},
		)
	}

	return
}

// Get description for an assignment
func GetAssignmentDesc(page *rod.Page) string {
	return ""
}

// Get all grades and data from report card page
func GetReportCard() {}
