package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/suzmii/ACMBot/database"
	"github.com/suzmii/ACMBot/errorx/usererr"
	"github.com/suzmii/ACMBot/render"
	"github.com/suzmii/ACMBot/util/logx"
	"github.com/suzmii/ACMBot/util/transer"
)

// initCodeforcesUser 拉取并创建用户
func (h *Handler) initCodeforcesUser(ctx context.Context, username string) (*database.CodeforcesUserWithRecords, error) {
	userInfo, err := h.api.FetchCodeforcesUserInfo(username, false)
	if err != nil {
		if errors.Is(err, usererr.ErrUserNotFound(username)) {
			return nil, usererr.ErrUserNotFound(username)
		}
		return nil, fmt.Errorf("failed to fetch user info from API: %w", err)
	}

	user, err := h.store.CreateCodeforcesUser(ctx, &database.CreateCodeforcesUserParams{
		Username:  username,
		AvatarUrl: userInfo.Avatar,
		FriendNum: int32(userInfo.FriendCount),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}
	return user, nil
}

// getAndInitUser 获取用户，如果不存在则尝试拉取用户信息
func (h *Handler) getAndInitUser(ctx context.Context, username string) (*database.CodeforcesUserWithRecords, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty") // TODO: change to usererr
	}

	user, err := h.store.GetCodeforcesUserByUsername(ctx, username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get user from database: %w", err)
		}
		user, err = h.initCodeforcesUser(ctx, username)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return user, err
}

func (h *Handler) updateUserRatingRecords(ctx context.Context, user *database.CodeforcesUserWithRecords) (*database.CodeforcesUserWithRecords, error) {
	fetchedRecords, err := h.api.FetchCodeforcesRatingRecords(user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch codeforces rating records: %w", err)
	}

	user.RatingRecords, err = h.store.UpdateCodeforcesRatingRecords(ctx, int(user.ID), transer.CodeforcesRatingRecordApi2DB(fetchedRecords))
	if err != nil {
		return nil, fmt.Errorf("failed to update user  rating records: %w", err)
	}

	return user, err
}

// GetCodeforcesRatingImage 获取Codeforces用户的rating记录图片
// 输入: username - Codeforces用户名
// 输出: []byte - 图片数据, error - 错误信息
func (h *Handler) GetCodeforcesRatingImage(ctx context.Context, username string) ([]byte, error) {
	user, err := h.getAndInitUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 检查更新rating
	if time.Since(user.RatingRecords.UpdatedAt) > 4*time.Hour {
		user, err = h.updateUserRatingRecords(ctx, user)
	}

	// 检查是否有rating记录
	if len(user.RatingRecords.Records) == 0 {
		return nil, usererr.ErrNoRatingRecords
	}

	// 渲染图片
	imageData, err := h.render.RatingDetail(ctx, render.CodeforcesRatingRecords{
		Handle: username,
		Data:   transer.CodeforcesRatingRecordDB2Render(user.RatingRecords.Records),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render rating image: %w", err)
	}

	return imageData, nil
}

// TODO: 统计函数从数据库挪过来这里？或者再开一个算法模块
// func (h *Handler) calculateCodeforcesSubmissionStatistics(ctx context.Context, user *database.CodeforcesUserWithRecords) (*database.CodeforcesUserWithRecords, error) {

// }

func (h *Handler) updateCodeforcesSubmissionStatistics(ctx context.Context, user *database.CodeforcesUserWithRecords) (*database.CodeforcesUserWithRecords, error) {
	fetchedSubmissions, err := h.api.FetchCodeforcesSubmissionsAfter(user.Username, user.SubmissionStatistics.LastSubmissionAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch codeforces rating records: %w", err)
	}

	err = h.store.CreateCodeforcesSubmissions(ctx, transer.CodeforcesSubmissionApi2DB(user.ID, fetchedSubmissions))
	if err != nil {
		return nil, fmt.Errorf("failed to create codeforces submissions: %w", err)
	}

	user.SubmissionStatistics, err = h.store.UpdateCodeforcesSubmissionStatistics(ctx, int(user.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to update codeforces submission statistics: %w", err)
	}

	return user, err
}

// 统计解题数区间
func calculateCodeforcesSolvedData(acProblems []database.CodeforcesProblemForQuery) []render.CodeforcesUserSolvedData {
	defer logx.TraceWall(logger, "calculateCodeforcesSolvedData")()
	// TODO: 使用配置文件？
	ranges := []struct {
		name string
		min  int
		max  int
	}{
		{"800+", 800, 1400},
		{"1400+", 1400, 2000},
		{"2000+", 2000, 2600},
		{"2600+", 2600, math.MaxInt},
	}

	total := len(acProblems)
	if total == 0 {
		return make([]render.CodeforcesUserSolvedData, 0)
	}

	result := make([]render.CodeforcesUserSolvedData, 0, len(ranges))
	for _, r := range ranges {
		count := 0
		for _, problem := range acProblems {
			if problem.Rating >= r.min && problem.Rating < r.max {
				count++
			}
		}
		percent := float32(count) / float32(total) * 100
		result = append(result, render.CodeforcesUserSolvedData{
			Range:   r.name,
			Percent: percent,
		})
	}

	return result
}

// GetCodeforcesUserProfileImage 获取Codeforces用户资料图片
// 输入: username - Codeforces用户名
// 输出: []byte - 图片数据, error - 错误信息
func (h *Handler) GetCodeforcesUserProfileImage(ctx context.Context, username string) ([]byte, error) {
	defer logx.TraceWall(logger, "GetCodeforcesUserProfileImage")()
	logger.Tracef("GetCodeforcesUserProfileImage: username=%s", username)

	user, err := h.getAndInitUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	logger.Tracef("GetCodeforcesUserProfileImage: user.ID=%d, user.Username=%s", user.ID, user.Username)
	logger.Tracef("GetCodeforcesUserProfileImage: RatingRecords=%v, SubmissionStatistics=%v",
		user.RatingRecords, user.SubmissionStatistics)

	// 检查更新rating
	if user.RatingRecords == nil || user.RatingRecords.UpdatedAt.IsZero() || time.Since(user.RatingRecords.UpdatedAt) > 4*time.Hour {
		logger.Tracef("GetCodeforcesUserProfileImage: updating rating records")
		user, err = h.updateUserRatingRecords(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to update rating records: %w", err)
		}
	}

	// 检查更新submission
	if user.SubmissionStatistics == nil || user.SubmissionStatistics.UpdatedAt.IsZero() || time.Since(user.SubmissionStatistics.UpdatedAt) > 4*time.Hour {
		logger.Tracef("GetCodeforcesUserProfileImage: updating submission statistics")
		user, err = h.updateCodeforcesSubmissionStatistics(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to update submission statistics: %w", err)
		}
	}

	solvedData := calculateCodeforcesSolvedData(user.SubmissionStatistics.Ac)

	logger.Tracef("GetCodeforcesUserProfileImage: total solved=%d, solvedData=%v", len(user.SubmissionStatistics.Ac), solvedData)

	// 转换为render需要的格式
	var maxRating, currRating int
	if user.RatingRecords != nil {
		maxRating = user.RatingRecords.MaxRating
		currRating = user.RatingRecords.CurrRating
	}

	logger.Tracef("GetCodeforcesUserProfileImage: maxRating=%d, currRating=%d", maxRating, currRating)

	renderProfile := render.CodeforcesUserProfile{
		Avatar:     user.AvatarUrl,
		Handle:     user.Username,
		MaxRating:  maxRating,
		FriendOf:   int(user.FriendNum),
		Rating:     currRating,
		Level:      render.Rating2Level(currRating),
		Solved:     len(user.SubmissionStatistics.Ac),
		SolvedData: solvedData,
	}

	// 渲染图片
	imageData, err := h.render.ProfileV2(ctx, renderProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to render profile image: %w", err)
	}

	return imageData, nil
}
