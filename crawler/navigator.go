package crawler

import (
	"errors"
	"time"

	"github.com/go-rod/rod"
)

// Wait page stable for [ms_optional] ms, if not defined, default [200] ms
func WaitStable(page *rod.Page, ms_optional ...int) {
	ms := 100
	if len(ms_optional) > 0 {
		ms = ms_optional[0]
	}
	page.WaitStable(time.Millisecond * time.Duration(ms))
}

// Use user provided account info to log into Jupiter
func Login(page *rod.Page, name string, pwd string) error {
	err := rod.Try(func() {
		// Enter basic school info
		page.MustElement("#text_school1").MustInput("Georgia School Ningbo")
		page.MustElement("#text_city1").MustInput("Ningbo")
		page.MustElement("#showcity > div.menuspace").MustClick()
		WaitStable(page)
		page.MustElement("#menulist_region1 > div[val='xx_xx']").MustClick()

		// Enter user account for login
		page.MustElement("#text_studid1").MustInput(name)
		page.MustElement("#text_password1").MustInput(pwd)

		page.MustElement("#loginbtn").MustClick()
	})
	if err != nil {
		return errors.New("登录 Jupiter 时发生未知异常: " + err.Error())
	}

	// Select the newest school year
	WaitStable(page, 800)
	page.MustElement("#schoolyeartab").MustClick()
	page.MustElement("#schoolyearlist > div:nth-child(1)").MustClick()

	WaitStable(page)
	return nil
}

// Get all options from the nav bar
func NavGetOptions(page *rod.Page) (opts rod.Elements, courses rod.Elements) {
	WaitStable(page)
	page.MustElement("#touchnavbtn").MustClick()
	opts, courses = page.MustElements("#sidebar > div[val]"), page.MustElements("#sidebar > div[click*='grades']")
	return
}

// Navigate to designated target on the nav bar
func NavNavigate(page *rod.Page, target *rod.Element) {
	WaitStable(page, 800)
	target.MustClick()
	WaitStable(page)
}
