package handler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/errorx/usererr"
	"github.com/suzmii/ACMBot/render"
)

var problemIDRe = regexp.MustCompile(`^(\d{1,6})([A-Za-z][0-9A-Za-z]*)$`)

func parseCodeforcesProblemID(problemID string) (contestID int, index string, err error) {
	problemID = strings.TrimSpace(problemID)
	m := problemIDRe.FindStringSubmatch(problemID)
	if len(m) != 3 {
		return 0, "", usererr.ErrInvalidProblemID
	}

	contestID, err = strconv.Atoi(m[1])
	if err != nil {
		return 0, "", usererr.ErrInvalidProblemID
	}

	index = strings.ToUpper(m[2])
	return contestID, index, nil
}

func (h *Handler) GetCodeforcesProblemStatementImage(ctx context.Context, problemID string) ([]byte, error) {
	contestID, index, err := parseCodeforcesProblemID(problemID)
	if err != nil {
		return nil, err
	}

	statement, err := h.api.FetchCodeforcesProblemStatement(api.Problem{
		ContestID: contestID,
		Index:     index,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch problem statement: %w", err)
	}

	samples := make([]render.CodeforcesProblemSample, 0, len(statement.Samples))
	for _, s := range statement.Samples {
		samples = append(samples, render.CodeforcesProblemSample{
			Input:  s.Input,
			Output: s.Output,
		})
	}

	image, err := h.render.RandomProblem(ctx, render.CodeforcesRandomProblem{
		Title:       statement.Title,
		URL:         statement.URL,
		Rating:      0,
		Tags:        nil,
		TimeLimit:   statement.TimeLimit,
		MemoryLimit: statement.MemoryLimit,
		Statement:   statement.Statement,
		Input:       statement.Input,
		Output:      statement.Output,
		Samples:     samples,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render problem statement image: %w", err)
	}
	return image, nil
}
