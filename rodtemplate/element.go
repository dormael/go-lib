package rodtemplate

import (
	"github.com/go-rod/rod"
)

type ElementTemplate struct {
	*rod.Element
}

func (e ElementTemplate) Has(selector string) bool {
	has, _, err := e.Element.Has(selector)
	if err != nil {
		panic(err)
	}

	return has
}

func (e ElementTemplate) El(selector string) *ElementTemplate {
	return &ElementTemplate{Element: e.MustElement(selector)}
}

func (e ElementTemplate) Els(selector string) ElementsTemplate {
	return toElementsTemplate(e.MustElements(selector))
}

func toElementsTemplate(elements rod.Elements) ElementsTemplate {
	est := make([]*ElementTemplate, 0)
	for idx := range elements {
		est = append(est, &ElementTemplate{Element: elements[idx]})
	}

	return est
}

func (e ElementTemplate) ElE(selector string) (*rod.Element, error) {
	return e.Element.Element(selector)
}

func (e ElementTemplate) ElementAttribute(selector string, name string) *string {
	return e.El(selector).MustAttribute(name)
}

func (e ElementTemplate) Height() float64 {
	quad := e.MustShape().Quads[0]

	return quad[7] - quad[1]
}

type ElementsTemplate []*ElementTemplate
