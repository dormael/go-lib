package rodtemplate

import (
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

type PageTemplate struct {
	P *rod.Page
}

func (p *PageTemplate) Navigate(url string) error {
	if p.P == nil {
		return errors.New("page is nil")
	}

	p.P.MustWaitLoad()
	p.P.MustWaitIdle()

	p.P.MustNavigate(url)

	p.P.MustWaitLoad()
	p.P.MustWaitIdle()

	return nil
}

func (p *PageTemplate) ClickElement(selector string) {
	p.P.MustWaitIdle()

	el := p.P.MustElement(selector)
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
		p.P.Mouse.MustMove(point.X, point.Y)
	} else {
		if cErr, ok := err.(*cdp.Error); ok {
			log.Println("failed to get element shape", cErr)
		} else {
			panic(err)
		}
	}
}

func (p *PageTemplate) URL() string {
	return p.P.MustInfo().URL
}

func (p *PageTemplate) Input(selector string, value string) {
	for i := 0; i < 100; i++ {
		if true == p.P.MustHas(selector) {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	if false == p.P.MustHas(selector) {
		log.Fatalf("failed to find input having selector %s\n", selector)
	}

	el := p.P.MustElement(selector)
	el.MustClick().MustSelectAllText().MustInput(value)
}

func (p *PageTemplate) PressKey(keyCode int32) {
	p.P.Keyboard.MustPress(keyCode)
}

func (p *PageTemplate) WaitLoadAndIdle() {
	p.P.MustWaitNavigation()
	p.WaitLoad()
	p.WaitIdle()
}

func (p *PageTemplate) Has(selector string) bool {
	has, _, err := p.P.Has(selector)
	if err != nil {
		panic(err)
	}

	return has
}

func (p *PageTemplate) El(selector string) *ElementTemplate {
	return &ElementTemplate{Element: p.P.MustElement(selector)}
}

func (p *PageTemplate) Els(selector string) ElementsTemplate {
	return toElementsTemplate(p.P.MustElements(selector))
}

func (p *PageTemplate) Reload() {
	p.P.MustReload()
}

func (p *PageTemplate) FrameID() proto.PageFrameID {
	return p.P.FrameID
}

func (p *PageTemplate) WaitIdle() {
	p.P.MustWaitIdle()
}

func (p *PageTemplate) WaitLoad() {
	if err := p.P.WaitLoad(); err != nil {
		if cErr, ok := err.(*cdp.Error); ok {
			log.Println("failed to wait", cErr)
		} else {
			panic(err)
		}
	}
}

func (p *PageTemplate) WaitRepaint() {
	if err := p.P.WaitRepaint(); err != nil {
		log.Println("failed to wait", err)
	}
}

func (p *PageTemplate) ScrollTop() {
	p.P.Keyboard.MustPress(input.Home)
}

func (p *PageTemplate) ScrollTo(e *ElementTemplate) {
	quad := e.MustShape().Quads[0]
	ybottom := quad[7]
	if err := p.P.Mouse.Scroll(0.0, ybottom, 1); err != nil {
		log.Println("failed to scroll mouse", err)
	}
}

func (p *PageTemplate) Body() string {
	return p.El("body").MustHTML()
}

func (p *PageTemplate) HTML() string {
	return p.El("html").MustHTML()
}

func (p *PageTemplate) Event() <-chan *rod.Message {
	return p.P.Event()
}

func (p *PageTemplate) MaximizeToWindowBounds() {
	bounds := p.P.MustGetWindow()
	p.SetViewport(bounds.Width, bounds.Height)
}

func (p *PageTemplate) SetViewport(width, height int) {
	p.P.MustSetViewport(width, height, 0, false)
}

func (p *PageTemplate) Dump(dumpPath string) []byte {
	return p.P.MustScreenshotFullPage(dumpPath)
}

func (p *PageTemplate) ScreenShot(el *ElementTemplate, dumpPath string, yDelta float64) []byte {
	err := el.ScrollIntoView()
	if err != nil {
		panic(err)
	}

	bounds := p.P.MustGetWindow()
	quad := el.MustShape().Quads[0]

	req := &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      quad[1] + yDelta,
			Width:  float64(bounds.Width),
			Height: quad[7] - quad[1],
			Scale:  1,
		},
	}

	byteArr, errScreenShot := p.P.Screenshot(false, req)
	if errScreenShot != nil {
		panic(errScreenShot)
	}

	errWrite := ioutil.WriteFile(dumpPath, byteArr, 0644)
	if errWrite != nil {
		panic(errWrite)
	}

	return byteArr
}

func (e ElementTemplate) ScreenShotElement(p *PageTemplate, dumpPath string, yDelta float64) []byte {
	err := e.ScrollIntoView()
	if err != nil {
		panic(err)
	}

	quad := e.MustShape().Quads[0]

	req := &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      quad[0],
			Y:      quad[1] + yDelta,
			Width:  quad[2] - quad[0],
			Height: quad[7] - quad[1],
			Scale:  1,
		},
	}

	byteArr, errScreenShot := p.P.Screenshot(false, req)
	if errScreenShot != nil {
		panic(errScreenShot)
	}

	errWrite := ioutil.WriteFile(dumpPath, byteArr, 0644)
	if errWrite != nil {
		panic(errWrite)
	}

	return byteArr
}

func NewPageTemplate(p *rod.Page) *PageTemplate {
	return &PageTemplate{P: p}
}
