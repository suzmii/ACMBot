package core

import (
	_ "embed"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type Template int

var templates = make(map[Template]*template.Template)

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
	//go:embed templates/race_calendar.gohtml
	TemplateContentRaceCalendar string
)

var templateContents = map[Template]*string{
	TemplateCodeforcesProfileV1:     &TemplateContentCodeforcesProfileV1,
	TemplateCodeforcesProfileV2:     &TemplateContentCodeforcesProfileV2,
	TemplateCodeforcesRatingRecords: &TemplateContentCodeforcesRatingRecords,
	TemplateAtcoderProfile:          &TemplateContentAtcoderProfile,
	TemplateQQGroupRank:             &TemplateContentQQGroupRank,
	TemplateRaceCalendar:            &TemplateContentRaceCalendar,
}

// GetTemplate returns a compiled template by name.
func GetTemplate(name Template) *template.Template {
	tmpl, ok := templates[name]
	if !ok {
		log.Fatalf("Template %v not found", name)
	}

	return tmpl
}
