package handler

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/errorx"
	"github.com/suzmii/ACMBot/errorx/usererr"	
	"github.com/suzmii/ACMBot/render"
	"github.com/suzmii/ACMBot/util"
)

var problemIDRe   = regexp.MustCompile(`^(\d{1,6})([A-Za-z][0-9A-Za-z]*)$`)
var ratingRangeRe = regexp.MustCompile(`^(\d{3,4})\s*[~～-]\s*(\d{3,4})$`)

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


func parseRandomProblemFilter(expr string) (database.CodeforcesProblemFilter, error) {
	filter := database.CodeforcesProblemFilter{}
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return filter, nil
	}

	tokens := strings.Fields(expr)
	hasRating := false
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		if m := ratingRangeRe.FindStringSubmatch(token); len(m) == 3 {
			if hasRating {
				return filter, usererr.ErrInvalidRatingRange
			}
			minV, _ := strconv.Atoi(m[1])
			maxV, _ := strconv.Atoi(m[2])
			if minV > maxV {
				minV, maxV = maxV, minV
			}
			if err := util.ValidateRating(minV); err != nil {
				return filter, err
			}
			if err := util.ValidateRating(maxV); err != nil {
				return filter, err
			}
			filter.MinRating = &minV
			filter.MaxRating = &maxV
			hasRating = true
			continue
		}

		if util.IsNumber(token) {
			if hasRating {
				return filter, errorx.NewUserError("只能设置一种 rating 筛选（单值或范围）")
			}
			v, _ := strconv.Atoi(token)
			if err := util.ValidateRating(v); err != nil {
				return filter, err
			}
			filter.MinRating = &v
			filter.MaxRating = &v
			hasRating = true
			continue
		}

		filter.Tags = append(filter.Tags, strings.ToLower(token))
	}

	return filter, nil
}


func (h *Handler) GetCodeforcesRandomProblemImage(ctx context.Context, filterExpr string) ([]byte, error) {
	filter, err := parseRandomProblemFilter(filterExpr)
	if err != nil {
		return nil, err
	}

	problem, err := h.store.GetRandomCodeforcesProblem(ctx, filter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usererr.ErrNoRandomProblemMatched
		}
		return nil, fmt.Errorf("failed to fetch random problem from database: %w", err)
	}

	statement, err := h.api.FetchCodeforcesProblemStatement(api.Problem{
		ContestID:      problem.ContestID,
		ProblemSetName: problem.ProblemSetName,
		Index:          problem.Index,
		Name:           problem.Name,
		Type:           problem.Type,
		Rating:         problem.Rating,
		Tags:           problem.Tags,
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
		Rating:      problem.Rating,
		Tags:        problem.Tags,
		TimeLimit:   statement.TimeLimit,
		MemoryLimit: statement.MemoryLimit,
		Statement:   statement.Statement,
		Input:       statement.Input,
		Output:      statement.Output,
		Samples:     samples,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render random problem image: %w", err)
	}
	return image, nil
}
