package internal

import (
	_ "embed"
	"html/template"
)

type Template int

// String 返回模板的名称
func (t Template) String() string {
	switch t {
	case TemplateCodeforcesProfileV1:
		return "TemplateCodeforcesProfileV1"
	case TemplateCodeforcesProfileV2:
		return "TemplateCodeforcesProfileV2"
	case TemplateCodeforcesRatingRecords:
		return "TemplateCodeforcesRatingRecords"
	case TemplateAtcoderProfile:
		return "TemplateAtcoderProfile"
	case TemplattCodeforcesRank:
		return "TemplattCodeforcesRank"
	case TemplateRaceCalendar:
		return "TemplateRaceCalendar"
	default:
		return "Unknown"
	}
}

const (
	TemplateCodeforcesProfileV1 Template = iota + 1
	TemplateCodeforcesProfileV2
	TemplateCodeforcesRatingRecords
	TemplateAtcoderProfile
	TemplattCodeforcesRank
	TemplateRaceCalendar
)

var (
	//go:embed templates/codeforces_profile_v1.gohtml
	TemplateContentCodeforcesProfileV1 string
	//go:embed templates/codeforces_profile_v2.gohtml
	TemplateContentCodeforcesProfileV2 string
	//go:embed templates/codeforces_rating_change.gohtml
	TemplateContentCodeforcesRatingRecords string
	//go:embed templates/atcoder_profile.gohtml
	TemplateContentAtcoderProfile string
	//go:embed templates/codeforces_rating_rank.gohtml
	TemplateContentCodeforcesRank string
)

// templateContents 是 embed 出来的源码，只读，可在多 Render 实例间安全共享。
var templateContents = map[Template]*string{
	TemplateCodeforcesProfileV1:     &TemplateContentCodeforcesProfileV1,
	TemplateCodeforcesProfileV2:     &TemplateContentCodeforcesProfileV2,
	TemplateCodeforcesRatingRecords: &TemplateContentCodeforcesRatingRecords,
	TemplateAtcoderProfile:          &TemplateContentAtcoderProfile,
	TemplattCodeforcesRank:          &TemplateContentCodeforcesRank,
}

// GetTemplate 返回该 Render 实例下编译好的模板。
func (r *Render) GetTemplate(name Template) *template.Template {
	tmpl, ok := r.templates[name]
	if !ok {
		logger.Warnf("Template %#v not found", name)
		return nil
	}
	return tmpl
}
