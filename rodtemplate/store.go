package rodtemplate

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/proto"
)

type Credential struct {
	LoginGateURL          string
	LoginLinkSelector     string
	LoginInputSelector    string
	PasswordInputSelector string
	LoginURL              string

	LoginSuccessSelector string
	ID                   string

	Password    string
	EnvID       string
	EnvPassword string
}

type StoreTemplate struct {
	*rod.Browser
}

func (s *StoreTemplate) Login(c Credential) (*PageTemplate, error) {
	var pt *PageTemplate

	page := s.MustPage(c.LoginGateURL)
	if err := page.WaitLoad(); err != nil {
		if false == cdp.ErrCtxDestroyed.Is(err) {
			panic(err)
		}
		log.Println(err.Error(), "occurred occasionally but has no problem")
	}

	pt = &PageTemplate{p: page}

	bounds := pt.p.MustGetWindow()
	pt.p.SetViewport(&proto.EmulationSetDeviceMetricsOverride{Width: bounds.Width, Height: bounds.Height})

	pages, err := s.Browser.Pages()
	if err != nil {
		return nil, err
	}

	for _, p := range pages {
		if p.FrameID == pt.p.FrameID {
			continue
		}
		p.MustClose()
	}

	pt.WaitLoad()

	if c.LoginURL != "" {
		if err := pt.Navigate(c.LoginURL); err != nil {
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

func NewStoreTemplate(b *rod.Browser) *StoreTemplate {
	return &StoreTemplate{Browser: b}
}
