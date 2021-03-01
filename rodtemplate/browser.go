package rodtemplate

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
)

type Credential struct {
	LoginGateURL          string
	LoginAfterURL         string
	LoginLinkSelector     string
	LoginInputSelector    string
	PasswordInputSelector string

	LoginURL string

	LoginSuccessSelector string

	ID       string
	Password string
	EnvID    string

	EnvPassword string

	CaptchaHandler         func(pt *PageTemplate) error
	LoginLinkHandler       func(pt *PageTemplate) error
	LoginPostSubmitHandler func(pt *PageTemplate) error
	LoginSuccessCheckHandler func(pt *PageTemplate) (bool, error)
}

type BrowserTemplate struct {
	*rod.Browser
}

func (b *BrowserTemplate) Login(c Credential) (*PageTemplate, error) {
	var pt *PageTemplate

	page := b.MustPage(c.LoginGateURL)
	if err := page.WaitLoad(); err != nil {
		if false == cdp.ErrCtxDestroyed.Is(err) {
			panic(err)
		}
		log.Println(err.Error(), "occurred occasionally but has no problem")
	}

	pt = &PageTemplate{P: page}
	pt.MaximizeToWindowBounds()

	pages, err := b.Browser.Pages()
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		if p.FrameID == pt.P.FrameID {
			continue
		}
		p.MustClose()
	}

	pt.WaitLoad()

	if c.LoginURL != "" {
		if err := pt.Navigate(c.LoginURL); err != nil {
			return nil, err
		}
	} else if c.LoginLinkHandler != nil {
		if err := c.LoginLinkHandler(pt); err != nil {
			return nil, err
		}
	} else {
		pt.WaitLoadAndIdle()
		pt.ClickWhenAvailable(c.LoginLinkSelector)
	}

	login := &Login{PageTemplate: pt, Credential: c}

	if err := login.Validate(); err != nil {
		return nil, err
	}

	if err := login.Submit(); err != nil {
		return nil, err
	}

	return pt, nil
}

func NewBrowserTemplate(b *rod.Browser) *BrowserTemplate {
	return &BrowserTemplate{Browser: b}
}
