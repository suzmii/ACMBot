package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/errorx/usererr"
	"github.com/suzmii/ACMBot/render"
	"github.com/suzmii/ACMBot/util/transer"
)

func normalizeAtcoderAvatarURL(avatar string) string {
	switch {
	case strings.HasPrefix(avatar, "//"):
		return "https:" + avatar
	case strings.HasPrefix(avatar, "/"):
		return "https://atcoder.jp" + avatar
	default:
		return avatar
	}
}

func (h *Handler) initAtcoderUser(ctx context.Context, username string) (*database.AtcoderUserWithRecords, error) {
	userInfo, err := h.api.FetchAtcoderUser(username)
	if err != nil {
		if errors.Is(err, usererr.ErrUserNotFound(username)) {
			return nil, usererr.ErrUserNotFound(username)
		}
		return nil, fmt.Errorf("failed to fetch user info from atcoder api: %w", err)
	}

	if userInfo == nil {
		return nil, fmt.Errorf("failed to fetch user info from atcoder api: empty response")
	}

	user, err := h.store.CreateAtcoderUser(ctx, &database.CreateAtcoderUserParams{
		Username:         username,
		AvatarUrl:        normalizeAtcoderAvatarURL(userInfo.Avatar),
		Rank:             userInfo.Rank,
		Rating:           int32(userInfo.Rating),
		HighestRating:    int32(userInfo.HighestRating),
		PromotionMessage: userInfo.PromotionMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create atcoder user in database: %w", err)
	}

	return user, nil
}

func (h *Handler) getAndInitAtcoderUser(ctx context.Context, username string) (*database.AtcoderUserWithRecords, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	user, err := h.store.GetAtcoderUserByUsername(ctx, username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get atcoder user from database: %w", err)
		}

		user, err = h.initAtcoderUser(ctx, username)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return user, nil
}

func (h *Handler) updateAtcoderSubmissionStatistics(ctx context.Context, user *database.AtcoderUserWithRecords) (*database.AtcoderUserWithRecords, error) {
	fetchedSubmissions, err := h.api.FetchAtcoderSubmissionListAfter(user.Username, user.SubmissionStatistics.LastSubmissionAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch atcoder submissions: %w", err)
	}

	if fetchedSubmissions != nil {
		err = h.store.CreateAtcoderSubmissions(ctx, transer.AtcoderSubmissionApi2DB(user.ID, *fetchedSubmissions))
		if err != nil {
			return nil, fmt.Errorf("failed to create atcoder submissions: %w", err)
		}
	}

	user.SubmissionStatistics, err = h.store.UpdateAtcoderSubmissionStatistics(ctx, int(user.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to update atcoder submission statistics: %w", err)
	}

	return user, nil
}

func calculateAtcoderSolvedData(acProblems []database.AtcoderProblemForQuery) []render.AtcoderUserSolvedData {
	ranges := []struct {
		name string
		min  float64
		max  float64
	}{
		{"<200", 0, 200},
		{"200+", 200, 400},
		{"400+", 400, 800},
		{"800+", 800, math.MaxFloat64},
	}

	total := len(acProblems)
	if total == 0 {
		return make([]render.AtcoderUserSolvedData, 0)
	}

	result := make([]render.AtcoderUserSolvedData, 0, len(ranges))
	for _, r := range ranges {
		count := 0
		for _, problem := range acProblems {
			if problem.Point >= r.min && problem.Point < r.max {
				count++
			}
		}

		percent := float32(count) / float32(total) * 100
		result = append(result, render.AtcoderUserSolvedData{
			Range:   r.name,
			Percent: percent,
		})
	}

	return result
}

func (h *Handler) GetAtcoderUserProfileImage(ctx context.Context, username string) ([]byte, error) {
	user, err := h.getAndInitAtcoderUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get atcoder user: %w", err)
	}

	if user.SubmissionStatistics == nil {
		user.SubmissionStatistics = &database.AtcoderSubmissionStatistics{
			TotalCount:       0,
			Ac:               []database.AtcoderProblemForQuery{},
			LastSubmissionAt: time.Time{},
			UpdatedAt:        time.Time{},
			Version:          1,
		}
	}

	if user.SubmissionStatistics.UpdatedAt.IsZero() || time.Since(user.SubmissionStatistics.UpdatedAt) > 4*time.Hour {
		user, err = h.updateAtcoderSubmissionStatistics(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to update atcoder submission statistics: %w", err)
		}
	}

	renderProfile := render.AtcoderUserProfile{
		Avatar:           user.AvatarUrl,
		Handle:           user.Username,
		MaxRating:        int(user.HighestRating),
		PromotionMessage: user.PromotionMessage,
		Rating:           int(user.Rating),
		Level:            render.AtcoderRating2Level(int(user.Rating), user.Rank),
		Solved:           len(user.SubmissionStatistics.Ac),
		SolvedData:       calculateAtcoderSolvedData(user.SubmissionStatistics.Ac),
	}

	imageData, err := h.render.AtcoderProfile(ctx, renderProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to render atcoder profile image: %w", err)
	}

	return imageData, nil
}
