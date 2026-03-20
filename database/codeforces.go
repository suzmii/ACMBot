package database

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/suzmii/ACMBot/database/sqlc"
	"github.com/suzmii/ACMBot/util/jsonx"
)

type CodeforcesRatingRecord struct {
	Rating int       `json:"new_rating"`
	At     time.Time `json:"at"`
}

type CodeforcesRatingRecords struct {
	MaxRating  int                      `json:"max_rating"`
	CurrRating int                      `json:"rating"`
	Records    []CodeforcesRatingRecord `json:"rating_changes"`

	UpdatedAt time.Time `json:"updated_at"`
}

func (db *dbStore) UpdateCodeforcesRatingRecords(ctx context.Context, userId int, originRecords []CodeforcesRatingRecord) (*CodeforcesRatingRecords, error) {
	var records2insert CodeforcesRatingRecords
	records2insert.Records = originRecords
	if len(originRecords) > 0 {
		records2insert.CurrRating = originRecords[len(originRecords)-1].Rating
		records2insert.MaxRating = originRecords[0].Rating
		for _, record := range originRecords {
			if record.Rating > records2insert.MaxRating {
				records2insert.MaxRating = record.Rating
			}
		}
	}
	records2insert.UpdatedAt = time.Now()

	jsonb, err := jsonx.Marshal(records2insert)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal records: %v", err)
	}

	err = db.Querier.UpdateCodeforcesRatingRecordsRaw(ctx, sqlc.UpdateCodeforcesRatingRecordsRawParams{
		ID:            int64(userId),
		RatingRecords: jsonb,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update codeforces rating records: %v", err)
	}

	return &records2insert, nil
}

type CodeforcesProblemForQuery struct {
	ID     string   `json:"id"`
	Rating int      `json:"rating"`
	Tags   []string `json:"tags"`
}

type CodeforcesSubmissionStatistics struct {
	TotalCount int `json:"total_count"` // 提交总数

	Ac []CodeforcesProblemForQuery `json:"ac"` // 通过的题目列表, 需要去重

	LastSubmissionAt time.Time `json:"last_id_of_submission_in_statistics"` // 统计信息中最晚的提交
	UpdatedAt        time.Time `json:"updated_at"`                          // 最后更新时间
	Version          int       `json:"version"`                             // 版本号, 如果后续提交统计结构变化, 据此进行自动化升级
}

func (db *dbStore) UpdateCodeforcesSubmissionStatistics(ctx context.Context, userId int) (*CodeforcesSubmissionStatistics, error) {
	// FIXME: 使用事务; 插入前需要对submission进行排序
	// 获取用户信息, 读取当前的提交统计
	user, err := db.GetCodeforcesUserByID(ctx, int64(userId))
	if err != nil {
		return nil, fmt.Errorf("failed to get codeforces user: %v", err)
	}

	// 最后一个提交的时间
	lastSubmissionAt := user.SubmissionStatistics.LastSubmissionAt
	newSubmissions, err := db.Querier.GetCodeforcesSubmissionsAfter(ctx, sqlc.GetCodeforcesSubmissionsAfterParams{
		UserID: user.ID,
		At:     pgtype.Timestamptz{Time: lastSubmissionAt, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get submissions after id: %v", err)
	}

	// 3. 计算新的统计数据
	// 用于去重的AC题目集合
	acProblemMap := make(map[string]CodeforcesProblemForQuery)

	// 首先将已有的AC题目添加到map中
	for _, acProblem := range user.SubmissionStatistics.Ac {
		acProblemMap[acProblem.ID] = acProblem
	}

	// 处理新的提交
	newTotalCount := user.SubmissionStatistics.TotalCount

	for _, submission := range newSubmissions {
		newTotalCount++

		if submission.At.Time.Sub(lastSubmissionAt) > 0 {
			lastSubmissionAt = submission.At.Time
		}

		// 解析problem信息
		var problem CodeforcesProblem
		if err := jsonx.Unmarshal(submission.Problem, &problem); err != nil {
			continue // 跳过解析失败的problem
		}

		// 如果是AC状态, 添加到AC列表(去重), 后续的提交会覆盖前面的
		if submission.Status == "OK" {
			problemID := problem.ID()
			acProblemMap[problemID] = CodeforcesProblemForQuery{
				ID:     problemID,
				Rating: problem.Rating,
				Tags:   problem.Tags,
			}
		}
	}

	// 将map转换为slice
	acList := make([]CodeforcesProblemForQuery, 0, len(acProblemMap))
	for _, problem := range acProblemMap {
		acList = append(acList, problem)
	}

	// 构建新的统计信息
	newStats := CodeforcesSubmissionStatistics{
		TotalCount:       newTotalCount,
		Ac:               acList,
		LastSubmissionAt: lastSubmissionAt,
		UpdatedAt:        time.Now(),
		Version:          user.SubmissionStatistics.Version, // TODO: Version 管理
	}

	// 更新数据库
	jsonb, err := jsonx.Marshal(newStats)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submission statistics: %v", err)
	}

	err = db.Querier.UpdateCodeforcesSubmissionStatisticsRaw(ctx, sqlc.UpdateCodeforcesSubmissionStatisticsRawParams{
		ID:                   user.ID,
		SubmissionStatistics: jsonb,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update submission statistics: %v", err)
	}

	return &newStats, nil
}

type CodeforcesProblem struct {
	ContestID      int      `json:"contestId"`
	ProblemSetName string   `json:"problemsetName"`
	Index          string   `json:"index"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Rating         int      `json:"rating"`
	Tags           []string `json:"tags"`
}

func (p CodeforcesProblem) ID() string {
	if p.ContestID == 0 {
		return p.ProblemSetName + p.Index
	}
	return fmt.Sprintf("%d%s", p.ContestID, p.Index)
}

type CodeforcesSubmission struct {
	UserID  int64
	Problem CodeforcesProblem
	Status  string
	At      time.Time
}

func (db *dbStore) CreateCodeforcesSubmissions(ctx context.Context, submissions []CodeforcesSubmission) error {
	if len(submissions) == 0 {
		return nil
	}

	userIDs := make([]int64, len(submissions))
	problems := make([][]byte, len(submissions))
	statuses := make([]string, len(submissions))
	ats := make([]pgtype.Timestamptz, len(submissions))

	slices.SortFunc(submissions, func(a, b CodeforcesSubmission) int {
		return int(a.At.Sub(b.At))
	})

	for i, submission := range submissions {
		userIDs[i] = submission.UserID
		statuses[i] = submission.Status
		ats[i] = pgtype.Timestamptz{Time: submission.At, Valid: true}

		problemJSON, err := jsonx.Marshal(submission.Problem)
		if err != nil {
			return fmt.Errorf("failed to marshal problem: %v", err)
		}
		problems[i] = problemJSON
	}

	err := db.Querier.CreateCodeforcesSubmissionsRaw(ctx, sqlc.CreateCodeforcesSubmissionsRawParams{
		Column1: userIDs,
		Column2: problems,
		Column3: statuses,
		Column4: ats,
	})
	if err != nil {
		return fmt.Errorf("failed to create codeforces submissions: %v", err)
	}

	return nil
}

// CodeforcesUserWithRecords 包含用户信息及解析后的rating记录和提交统计
type CodeforcesUserWithRecords struct {
	sqlc.CodeforcesUser
	RatingRecords        *CodeforcesRatingRecords
	SubmissionStatistics *CodeforcesSubmissionStatistics
}

type CreateCodeforcesUserParams struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	FriendNum int32  `json:"friend_num"`
}

func (db *dbStore) CreateCodeforcesUser(ctx context.Context, params *CreateCodeforcesUserParams) (*CodeforcesUserWithRecords, error) {
	result := &CodeforcesUserWithRecords{
		RatingRecords: &CodeforcesRatingRecords{
			MaxRating:  0,
			CurrRating: 0,
			Records:    []CodeforcesRatingRecord{},
			UpdatedAt:  time.Time{},
		},
		SubmissionStatistics: &CodeforcesSubmissionStatistics{
			TotalCount:       0,
			Ac:               []CodeforcesProblemForQuery{},
			LastSubmissionAt: time.Time{},
			UpdatedAt:        time.Time{},
			Version:          1,
		},
	}

	rtjsonb, err := jsonx.Marshal(result.RatingRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rating records: %w", err)
	}

	ssjsonb, err := jsonx.Marshal(result.SubmissionStatistics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submission statistics: %w", err)
	}

	user, err := db.CreateCodeforcesUserRaw(ctx, sqlc.CreateCodeforcesUserRawParams{
		Username:             params.Username,
		AvatarUrl:            params.AvatarUrl,
		FriendNum:            params.FriendNum,
		RatingRecords:        rtjsonb,
		SubmissionStatistics: ssjsonb,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	result.CodeforcesUser = user
	return result, nil
}

// parseCodeforcesUser 将sqlc.CodeforcesUser解析为CodeforcesUserWithRecords
func (db *dbStore) parseCodeforcesUser(user sqlc.CodeforcesUser) (*CodeforcesUserWithRecords, error) {
	result := &CodeforcesUserWithRecords{
		CodeforcesUser: user,
	}

	// 解析rating记录
	if len(user.RatingRecords) > 0 {
		var records CodeforcesRatingRecords
		if err := jsonx.Unmarshal(user.RatingRecords, &records); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rating records: %w", err)
		}
		result.RatingRecords = &records
	}

	// 解析提交统计
	if len(user.SubmissionStatistics) > 0 {
		var stats CodeforcesSubmissionStatistics
		if err := jsonx.Unmarshal(user.SubmissionStatistics, &stats); err != nil {
			return nil, fmt.Errorf("failed to unmarshal submission statistics: %w", err)
		}
		result.SubmissionStatistics = &stats
	}

	return result, nil
}

// GetCodeforcesUserByID 按ID获取Codeforces用户及其解析后的rating记录和提交统计
func (db *dbStore) GetCodeforcesUserByID(ctx context.Context, userID int64) (*CodeforcesUserWithRecords, error) {
	user, err := db.Querier.GetCodeforcesUserByIDRaw(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codeforces user by id: %w", err)
	}
	return db.parseCodeforcesUser(user)
}

// GetCodeforcesUserByUsername 按用户名获取Codeforces用户及其解析后的rating记录和提交统计
func (db *dbStore) GetCodeforcesUserByUsername(ctx context.Context, username string) (*CodeforcesUserWithRecords, error) {
	user, err := db.Querier.GetCodeforcesUserByUsernameRaw(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get codeforces user by username: %w", err)
	}
	return db.parseCodeforcesUser(user)
}
