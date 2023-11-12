package crawler

import (
	"errors"
	"time"

	"github.com/go-rod/rod"
)

// Wait for loading to disappear
func WaitStable(page *rod.Page) {
	page.MustWaitElementsMoreThan("#busy[class*='hide']", 0)
}

// Use user provided account info to log into Jupiter
func Login(page *rod.Page, name string, pwd string) error {
	err := rod.Try(func() {
		page.Timeout(time.Second*30).MustWaitElementsMoreThan("#text_school1", 0)

		// Enter basic school info
		page.Timeout(time.Second * 2).MustElement("#text_school1").MustInput("Georgia School Ningbo")
		page.Timeout(time.Second * 2).MustElement("#text_city1").MustInput("Ningbo")
		page.Timeout(time.Second * 2).MustElement("#showcity > div.menuspace").MustClick()
		page.WaitStable(time.Millisecond * 100)
		page.Timeout(time.Second * 2).MustElement("#menulist_region1 > div[val='xx_xx']").MustClick()

		// Enter user account for login
		page.Timeout(time.Second * 2).MustElement("#text_studid1").MustInput(name)
		page.Timeout(time.Second * 2).MustElement("#text_password1").MustInput(pwd)

		page.Timeout(time.Second * 2).MustElement("#loginbtn").MustClick()
	})
	if err != nil {
		return err
	}

	// Select the newest school year
	WaitStable(page)
	err = rod.Try(func() {
		page.Timeout(time.Second * 2).MustElement("#schoolyeartab").MustClick()
		page.Timeout(time.Second * 2).MustElement("#schoolyearlist > div:nth-child(1)").MustClick()
	})
	if err != nil {
		return errors.New("invalidJupiterAccount")
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
	WaitStable(page)
	err := rod.Try(func() {
		target.Timeout(time.Second * 2).MustClick()
	})
	WaitStable(page)
	return err
}
