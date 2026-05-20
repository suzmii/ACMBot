package render

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

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
func (r *Render) executeTemplate(ctx context.Context, templateName internal.Template, data interface{}) ([]byte, error) {
	defer logx.TraceWall(logger, "executeTemplate")
	var buffer bytes.Buffer

	if err := r.GetTemplate(templateName).Execute(&buffer, data); err != nil {
		return nil, fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return r.RenderWithAutoSize(ctx, buffer.String())
}

func (r *Render) RatingDetail(ctx context.Context, records CodeforcesRatingRecords) ([]byte, error) {
	logger.Debug("rendering rating detail", records)
	return r.executeTemplate(
		ctx,
		internal.TemplateCodeforcesRatingRecords,
		&CodeforcesRatingRecords{
			Data:      records.Data,
			Handle:    records.Handle,
			EchartsJS: internal.ResourceEcharts,
		},
	)
}

func (r *Render) ProfileV2(ctx context.Context, user CodeforcesUserProfile) ([]byte, error) {
	defer logx.TraceWall(logger, "ProfileV2")()
	logger.Trace("rendering codeforces profile v2 ", user)
	return r.executeTemplate(
		ctx,
		internal.TemplateCodeforcesProfileV2,
		&CodeforcesUserProfile{
			Avatar:     user.Avatar,
			Handle:     user.Handle,
			MaxRating:  user.MaxRating,
			FriendOf:   user.FriendOf,
			Rating:     user.Rating,
			Level:      user.Level,
			Solved:     user.Solved,
			SolvedData: user.SolvedData,
			TailwindJS: internal.ResourceTailwind,
			FontCSS:    internal.ResourceZsft184,
		},
	)
}

type RankUser struct {
	Rank   int
	Handle string
	Avatar string
	Rating int
	Level  string
}

type RankUsers struct {
	Users      []RankUser
	Time       time.Time
	TailwindJS template.JS
	FontCSS    template.CSS
}

func (r *Render) Rank(ctx context.Context, users []RankUser) ([]byte, error) {
	defer logx.TraceWall(logger, "Rank")()
	logger.Trace("rendering codeforces Rank ", users)
	return r.executeTemplate(
		ctx,
		internal.TemplattCodeforcesRank,
		&RankUsers{
			Time:       time.Now(),
			Users:      users,
			TailwindJS: internal.ResourceTailwind,
			FontCSS:    internal.ResourceZsft184,
		},
	)
}
