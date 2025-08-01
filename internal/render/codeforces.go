package render

import (
	"bytes"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/internal/render/rendermodel"
)

func RatingDetail(records rendermodel.CodeforcesRatingRecords) ([]byte, error) {
	var buffer bytes.Buffer
	records.EchartsJS = ResourceEcharts
	if err := GetTemplate(TemplateCodeforcesRatingRecords).Execute(&buffer, records); err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}
	return Html(
		&playwright.BrowserNewPageOptions{
			DeviceScaleFactor: &[]float64{2.0}[0],
			Viewport: &playwright.Size{
				Width:  1000,
				Height: 500,
			},
		}, buffer,
	)
}

func ProfileV2(user rendermodel.CodeforcesUserProfile) ([]byte, error) {
	logrus.Debug(user)
	user.TailwindJS = ResourceTailwind
	user.FontCSS = ResourceZsft184
	var buffer bytes.Buffer
	if err := GetTemplate(TemplateCodeforcesProfileV2).Execute(&buffer, user); err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}
	logrus.Tracef("Rendered profile:\n %s", buffer.String())
	return Html(
		&playwright.BrowserNewPageOptions{
			DeviceScaleFactor: &[]float64{2.0}[0],
			Viewport: &playwright.Size{
				Width:  300,
				Height: 400,
			},
		}, buffer,
	)
}
