package internal

import (
	_ "embed"
	"html/template"
)

type Template int

var templates = make(map[Template]*template.Template)

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
	case TemplateQQGroupRank:
		return "TemplateQQGroupRank"
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
	TemplateQQGroupRank
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
	//go:embed templates/qq_group_rank.gohtml
	TemplateContentQQGroupRank string
)

var templateContents = map[Template]*string{
	TemplateCodeforcesProfileV1:     &TemplateContentCodeforcesProfileV1,
	TemplateCodeforcesProfileV2:     &TemplateContentCodeforcesProfileV2,
	TemplateCodeforcesRatingRecords: &TemplateContentCodeforcesRatingRecords,
	TemplateAtcoderProfile:          &TemplateContentAtcoderProfile,
	TemplateQQGroupRank:             &TemplateContentQQGroupRank,
}

// GetTemplate returns a compiled template by name.
func (r *Render) GetTemplate(name Template) *template.Template {
	tmpl, ok := templates[name]
	if !ok {
		logger.Warnf("Template %#v not found", name)
		return nil
	}
	return tmpl
}
