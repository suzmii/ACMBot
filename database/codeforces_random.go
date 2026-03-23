package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (db *dbStore) GetRandomCodeforcesProblem(ctx context.Context, filter CodeforcesProblemFilter) (*CodeforcesProblem, error) {
	if db.pool == nil {
		return nil, fmt.Errorf("random problem query requires root store")
	}

	var (
		conds []string
		args  []any
	)

	conds = append(conds, "status = 'OK'")
	conds = append(conds, "problem->>'index' <> ''")
	conds = append(conds, "problem->>'name' <> ''")

	if filter.MinRating != nil {
		args = append(args, *filter.MinRating)
		conds = append(conds, fmt.Sprintf("COALESCE((problem->>'rating')::int, 0) >= $%d", len(args)))
	}
	if filter.MaxRating != nil {
		args = append(args, *filter.MaxRating)
		conds = append(conds, fmt.Sprintf("COALESCE((problem->>'rating')::int, 0) <= $%d", len(args)))
	}

	if len(filter.Tags) > 0 {
		normalizedTags := make([]string, 0, len(filter.Tags))
		for _, tag := range filter.Tags {
			tag = strings.ToLower(strings.TrimSpace(tag))
			if tag != "" {
				normalizedTags = append(normalizedTags, tag)
			}
		}
		if len(normalizedTags) > 0 {
			args = append(args, normalizedTags)
			conds = append(conds, fmt.Sprintf(
				`EXISTS (
					SELECT 1
					FROM jsonb_array_elements_text(COALESCE(problem->'tags', '[]'::jsonb)) AS t
					WHERE lower(t) = ANY($%d::text[])
				)`, len(args)))
		}
	}

	query := `
		WITH filtered AS (
			SELECT DISTINCT ON (
				COALESCE(problem->>'contestId', ''),
				COALESCE(problem->>'problemsetName', ''),
				COALESCE(problem->>'index', '')
			) problem
			FROM codeforces_submissions
			WHERE ` + strings.Join(conds, " AND ") + `
			ORDER BY
				COALESCE(problem->>'contestId', ''),
				COALESCE(problem->>'problemsetName', ''),
				COALESCE(problem->>'index', ''),
				at DESC
		)
		SELECT problem
		FROM filtered
		ORDER BY random()
		LIMIT 1
	`

	var raw []byte
	err := db.pool.QueryRow(ctx, query, args...).Scan(&raw)
	if err != nil {
		return nil, fmt.Errorf("failed to query random codeforces problem: %w", err)
	}

	var p CodeforcesProblem
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal random codeforces problem: %w", err)
	}
	return &p, nil
}
