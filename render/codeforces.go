package render

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/suzmii/ACMBot/render/internal"
	"github.com/suzmii/ACMBot/util/logx"
)

type CodeforcesRatingChange struct {
	At        int64 `json:"at"`
	NewRating int   `json:"newRating"`
}

type CodeforcesRatingRecords struct {
	Data   []CodeforcesRatingChange
	Handle string

	EchartsJS template.JS
}

type CodeforcesRatingLevel string

const (
	CodeforcesRatingLevelNewbie                   CodeforcesRatingLevel = "Newbie"
	CodeforcesRatingLevelPupil                    CodeforcesRatingLevel = "Pupil"
	CodeforcesRatingLevelSpecialist               CodeforcesRatingLevel = "Specialist"
	CodeforcesRatingLevelExpert                   CodeforcesRatingLevel = "Expert"
	CodeforcesRatingLevelCandidateMaster          CodeforcesRatingLevel = "CM"
	CodeforcesRatingLevelMaster                   CodeforcesRatingLevel = "Master"
	CodeforcesRatingLevelInternationalMaster      CodeforcesRatingLevel = "IM"
	CodeforcesRatingLevelGrandmaster              CodeforcesRatingLevel = "GM"
	CodeforcesRatingLevelInternationalGrandmaster CodeforcesRatingLevel = "IGM"
	CodeforcesRatingLevelLegendaryGrandmaster     CodeforcesRatingLevel = "LGM"
	CodeforcesRatingLevelTourist                  CodeforcesRatingLevel = "Tourist"
)

func Rating2Level(rating int) CodeforcesRatingLevel {
	switch {
	case rating < 1200:
		return CodeforcesRatingLevelNewbie
	case rating < 1400:
		return CodeforcesRatingLevelPupil
	case rating < 1600:
		return CodeforcesRatingLevelSpecialist
	case rating < 1900:
		return CodeforcesRatingLevelExpert
	case rating < 2100:
		return CodeforcesRatingLevelCandidateMaster
	case rating < 2300:
		return CodeforcesRatingLevelMaster
	case rating < 2400:
		return CodeforcesRatingLevelInternationalMaster
	case rating < 2600:
		return CodeforcesRatingLevelGrandmaster
	case rating < 3000:
		return CodeforcesRatingLevelInternationalGrandmaster
	case rating < 4000:
		return CodeforcesRatingLevelLegendaryGrandmaster
	default:
		return CodeforcesRatingLevelTourist
	}
}

type CodeforcesUserSolvedData struct {
	Range   string
	Percent float32
}

type CodeforcesUserProfile struct {
	Avatar    string
	Handle    string
	MaxRating int
	FriendOf  int
	Rating    int
	Level     CodeforcesRatingLevel
	Solved    int

	SolvedData []CodeforcesUserSolvedData

	TailwindJS template.JS
	FontCSS    template.CSS
}

// executeTemplate 是通用的模板渲染函数
// templateName: 模板名称
// data: 模板数据
// setupFunc: 数据设置函数，用于在渲染前设置数据（如注入资源）
func (r *Render) executeTemplate(ctx context.Context, templateName internal.Template, data interface{}, setupFunc func(interface{})) ([]byte, error) {
	defer logx.TraceWall(logger, "executeTemplate")
	var buffer bytes.Buffer

	// 应用数据设置函数
	if setupFunc != nil {
		setupFunc(data)
	}

	if err := r.GetTemplate(templateName).Execute(&buffer, data); err != nil {
		return nil, fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return r.RenderWithAutoSize(ctx, buffer)
}

func (r *Render) RatingDetail(ctx context.Context, records CodeforcesRatingRecords) ([]byte, error) {
	logger.Debug("rendering rating detail", records)
	return r.executeTemplate(
		ctx,
		internal.TemplateCodeforcesRatingRecords,
		&records,
		func(data interface{}) {
			if d, ok := data.(*CodeforcesRatingRecords); ok {
				d.EchartsJS = internal.ResourceEcharts
			}
		},
	)
}

func (r *Render) ProfileV2(ctx context.Context, user CodeforcesUserProfile) ([]byte, error) {
	defer logx.TraceWall(logger, "ProfileV2")()
	logger.Trace("rendering codeforces profile v2 ", user)
	return r.executeTemplate(
		ctx,
		internal.TemplateCodeforcesProfileV2,
		&user,
		func(data interface{}) {
			if d, ok := data.(*CodeforcesUserProfile); ok {
				d.TailwindJS = internal.ResourceTailwind
				d.FontCSS = internal.ResourceZsft184
			}
		},
	)
}
