package rodtemplate

import (
	"errors"
	"fmt"
	"os"
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

				loginPt = &PageTemplate{p: iPage}
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

	loginPt.Input(c.LoginInputSelector, c.ID)
	loginPt.Input(c.PasswordInputSelector, c.Password)
	loginPt.PressKey(input.Enter)

	pt.WaitLoadAndIdle()

	for i := 0; i < 100; i++ {
		currentPageURL := pt.URL()
		if loginPageURL == currentPageURL {
			pt.WaitLoadAndIdle()
		} else if pt.El(c.LoginSuccessSelector).MustText() != "" {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	pt.WaitLoadAndIdle()
	currentPageURL := pt.URL()

	if loginPageURL == currentPageURL {
		return errors.New("failed to login")
	}

	return nil
}