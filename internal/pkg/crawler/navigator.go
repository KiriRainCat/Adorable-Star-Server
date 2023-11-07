package crawler

import (
	"errors"
	"strings"
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
func Login(page *rod.Page, name string, pwd string, uid int) error {
	err := rod.Try(func() {
		// Enter basic school info
		page.Timeout(time.Second * 2).MustElement("#text_school1").MustInput("Georgia School Ningbo")
		page.Timeout(time.Second * 2).MustElement("#text_city1").MustInput("Ningbo")
		page.Timeout(time.Second * 2).MustElement("#showcity > div.menuspace").MustClick()
		WaitStable(page)
		page.Timeout(time.Second * 2).MustElement("#menulist_region1 > div[val='xx_xx']").MustClick()

		// Enter user account for login
		page.Timeout(time.Second * 2).MustElement("#text_studid1").MustInput(name)
		page.Timeout(time.Second * 2).MustElement("#text_password1").MustInput(pwd)

		page.Timeout(time.Second * 2).MustElement("#loginbtn").MustClick()
		WaitStable(page, 800)
	})
	if err != nil {
		return err
	}

	// Check if request blocked by Cloudflare
	if strings.Contains(page.MustElement("body").MustText(), "malicious") {
		return errors.New("crawler request blocked by cloudflare")
	}

	// Select the newest school year
	err = rod.Try(func() {
		page.Timeout(time.Second * 2).MustElement("#schoolyeartab").MustClick()
		page.Timeout(time.Second * 2).MustElement("#schoolyearlist > div:nth-child(1)").MustClick()
	})
	if err != nil {
		return err
	}

	WaitStable(page)
	return nil
}

// Get all options from the nav bar
func NavGetOptions(page *rod.Page) (opts rod.Elements, courses rod.Elements, err error) {
	WaitStable(page)
	err = rod.Try(func() {
		page.Timeout(time.Second * 2).MustElement("#touchnavbtn").MustClick()
		opts = page.Timeout(time.Second * 2).MustElements("#sidebar > div[val]")
		courses = page.Timeout(time.Second * 2).MustElements("#sidebar > div[click*='grades']")
	})
	return
}

// Navigate to designated target on the nav bar
func NavNavigate(page *rod.Page, target *rod.Element) error {
	WaitStable(page, 800)
	err := rod.Try(func() {
		target.Timeout(time.Second * 2).MustClick()
	})
	WaitStable(page)
	return err
}
