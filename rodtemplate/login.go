package rodtemplate

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/go-rod/rod/lib/input"
)

type Login struct {
	*PageTemplate
	Credential Credential
}

func (l *Login) Validate() error {
	if l.Credential.ID == "" {
		l.Credential.ID = os.Getenv(l.Credential.EnvID)
	}

	if l.Credential.Password == "" {
		l.Credential.Password = os.Getenv(l.Credential.EnvPassword)
	}

	if l.Credential.ID == "" || l.Credential.Password == "" {
		return fmt.Errorf("id, password is required as parameter or os environment variables with names(%s, %s)", l.Credential.EnvID, l.Credential.EnvPassword)
	}

	return nil
}

func (l *Login) Submit() error {
	c := l.Credential
	pt := l.PageTemplate

	pt.WaitLoad()
	loginPageURL := pt.URL()

	pt.WaitLoad()

	var loginPt *PageTemplate

	for i := 0; i < 100; i++ {
		if true == pt.Has(c.LoginInputSelector) {
			loginPt = pt
		} else if true == pt.Has("iframe") {
			for _, e := range pt.Els("iframe") {
				iPage, err := e.Frame()
				if err != nil {
					continue
				}

				loginPt = &PageTemplate{P: iPage}
				if true == loginPt.Has(c.LoginInputSelector) {
					break
				}
			}
		}

		if loginPt != nil {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	if loginPt == nil {
		return errors.New("failed to find login input selector " + c.LoginInputSelector)
	}

	if c.CaptchaHandler != nil {
		errCaptcha := c.CaptchaHandler(loginPt)
		if errCaptcha != nil {
			return errCaptcha
		}
	}

	loginPt.Input(c.LoginInputSelector, c.ID)
	loginPt.Input(c.PasswordInputSelector, c.Password)
	loginPt.PressKey(input.Enter)

	pt.WaitLoadAndIdle()

	if c.LoginSuccessClickSelector != "" {
		for i := 0; i < 100; i++ {
			if false == pt.Has(c.LoginSuccessClickSelector) {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			pt.El(c.LoginSuccessClickSelector).MustClick()
			break
		}

		pt.WaitLoadAndIdle()
	}

	if c.LoginPostSubmitHandler != nil {
		errPostSubmit := c.LoginPostSubmitHandler(pt)
		if errPostSubmit != nil {
			return errPostSubmit
		}
	}

	for i := 0; i < 100; i++ {
		currentPageURL := pt.URL()
		if loginPageURL == currentPageURL {
			pt.WaitLoadAndIdle()
		} else if pt.Has(c.LoginSuccessSelector) && pt.El(c.LoginSuccessSelector).MustText() != "" {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	pt.WaitLoadAndIdle()
	currentPageURL := pt.URL()

	if loginPageURL == currentPageURL {
		screenshotPath := os.Getenv("PATH_SCREENSHOT")
		screenshotLoginFailed := os.Getenv("SCREENSHOT_LOGIN_FAILED")
		if screenshotLoginFailed == "1" && screenshotPath != "" {
			if _, err := os.Stat(screenshotPath); err != nil {
				if true == os.IsNotExist(err) {
					if err = os.MkdirAll(screenshotPath, 0755); err != nil {
						log.Println(err)
						return errors.New("failed to login")
					}
				} else {
					log.Println(err)
					return errors.New("failed to login")
				}
			}

			screenshotFile := fmt.Sprintf("loginfailed.%s.png", time.Now().Format("20060102150405"))
			pt.ScreenShot(pt.El("html"), path.Join(screenshotPath, screenshotFile), 0)
		}
		return errors.New("failed to login")
	}

	return nil
}
