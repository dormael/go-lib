package rodtemplate

import (
	"errors"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

type PageTemplate struct {
	p *rod.Page
}

func (p *PageTemplate) Navigate(url string) error {
	if p.p == nil {
		return errors.New("page is nil")
	}

	p.p.MustWaitLoad()
	p.p.MustWaitIdle()

	p.p.MustNavigate(url)

	p.p.MustWaitLoad()
	p.p.MustWaitIdle()

	return nil
}

func (p *PageTemplate) ClickElement(selector string) {
	p.p.MustWaitIdle()

	el := p.p.MustElement(selector)
	p.MoveMouseTo(el)

	el.MustClick()
}

func (p *PageTemplate) ClickWhenAvailable(selector string) bool {
	for i := 0; i < 1000; i++ {
		if true == p.Has(selector) {
			el := p.El(selector)
			if true == el.MustVisible() {
				el.MustFocus()
				el.MustScrollIntoView()
				el.MustClick()
				return true
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	return false
}

func (p *PageTemplate) FocusWhenAvailable(selector string) bool {
	for i := 0; i < 1000; i++ {
		if true == p.Has(selector) {
			el := p.El(selector)
			el.MustFocus()
			return true
		}
		time.Sleep(time.Millisecond * 100)
	}
	return false
}

func (p *PageTemplate) MoveMouseTo(el *rod.Element) {
	shape, err := el.Shape()
	if err == nil {
		point := shape.OnePointInside()
		p.p.Mouse.MustMove(point.X, point.Y)
	} else {
		if cErr, ok := err.(*cdp.Error); ok {
			log.Println("failed to get element shape", cErr)
		} else {
			panic(err)
		}
	}
}

func (p *PageTemplate) URL() string {
	return p.p.MustInfo().URL
}

func (p *PageTemplate) Input(selector string, value string) {
	for i := 0; i < 100; i++ {
		if true == p.p.MustHas(selector) {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	if false == p.p.MustHas(selector) {
		log.Fatalf("failed to find input having selector %s\n", selector)
	}

	el := p.p.MustElement(selector)
	el.MustClick().MustSelectAllText().MustInput(value)
}

func (p *PageTemplate) PressKey(keyCode int32) {
	p.p.Keyboard.MustPress(keyCode)
}

func (p *PageTemplate) WaitLoadAndIdle() {
	p.p.MustWaitNavigation()
	p.WaitLoad()
	p.WaitIdle()
}

func (p *PageTemplate) Has(selector string) bool {
	has, _, err := p.p.Has(selector)
	if err != nil {
		panic(err)
	}

	return has
}

func (p *PageTemplate) El(selector string) *ElementTemplate {
	return &ElementTemplate{Element: p.p.MustElement(selector)}
}

func (p *PageTemplate) Els(selector string) ElementsTemplate {
	return toElementsTemplate(p.p.MustElements(selector))
}

func (p *PageTemplate) Reload() {
	p.p.Reload()
}

func (p *PageTemplate) FrameID() proto.PageFrameID {
	return p.p.FrameID
}

func (p *PageTemplate) WaitIdle() {
	p.p.MustWaitIdle()
}

func (p *PageTemplate) WaitLoad() {
	if err := p.p.WaitLoad(); err != nil {
		if cErr, ok := err.(*cdp.Error); ok {
			log.Println("failed to wait", cErr)
		} else {
			panic(err)
		}
	}
}

func (p *PageTemplate) ScrollTop() {
	p.p.Keyboard.MustPress(input.Home)
}

func (p *PageTemplate) Body() string {
	return p.El("body").MustHTML()
}

func (p *PageTemplate) Event() <-chan *rod.Message {
	return p.p.Event()
}

func NewPageTemplate(p *rod.Page) *PageTemplate {
	return &PageTemplate{p: p}
}
