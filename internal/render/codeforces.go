package render

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/internal/render/core"
	"github.com/suzmii/ACMBot/internal/render/rendermodel"
)

func RatingDetail(records rendermodel.CodeforcesRatingRecords) ([]byte, error) {
	var buffer bytes.Buffer
	records.EchartsJS = core.ResourceEcharts
	if err := core.GetTemplate(core.TemplateCodeforcesRatingRecords).Execute(&buffer, records); err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}
	return core.HtmlAutoSize(buffer)
}

func ProfileV2(user rendermodel.CodeforcesUserProfile) ([]byte, error) {
	logrus.Debug(user)
	user.TailwindJS = core.ResourceTailwind
	user.FontCSS = core.ResourceZsft184
	var buffer bytes.Buffer
	if err := core.GetTemplate(core.TemplateCodeforcesProfileV2).Execute(&buffer, user); err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}
	logrus.Tracef("Rendered profile:\n %s", buffer.String())
	return core.HtmlAutoSize(buffer)
}
