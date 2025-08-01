package client

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/internal/api/model"
	"github.com/suzmii/ACMBot/internal/errs"
	"github.com/suzmii/ACMBot/internal/util"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"math/rand"

	"github.com/sirupsen/logrus"
)

var cfLock sync.Mutex

var (
	key    string
	secret string
)

func init() {
	cfg := config.LoadConfig().API
	key = cfg.CodeforcesKey
	secret = cfg.CodeforcesSecret

	if key == "" {
		logrus.Fatal("CODEFORCES_KEY environment variable not set")
	}
	if secret == "" {
		logrus.Fatal("CODEFORCES_SECRET environment variable not set")
	}
}

func fetchCodeforcesAPI[T any](apiMethod string, args map[string]any) (T, error) {
	cfLock.Lock()
	defer cfLock.Unlock()
	time.Sleep(500 * time.Millisecond)
	type codeforcesResponse[T any] struct {
		/*
			codeforces响应数据的基本格式:
				Result->期望的数据
				Comment->失败时返回的提示信息
		*/
		Status  string `json:"status"`
		Result  T      `json:"result"`
		Comment string `json:"comment"`
	}

	apiURL := "https://codeforces.com/api/"

	args["apiKey"] = key
	args["time"] = strconv.Itoa(int(time.Now().Unix()))

	var sortedArgs []string
	for k, v := range args {
		sortedArgs = append(sortedArgs, fmt.Sprintf("%v=%v", k, v))
	}
	sort.Strings(sortedArgs)

	randStr := strconv.Itoa(rand.Intn(900000) + 100000)
	hashSource := randStr + "/" + apiMethod + "?" + strings.Join(sortedArgs, "&") + "#" + secret

	h := sha512.New()
	h.Write([]byte(hashSource))
	hashSig := hex.EncodeToString(h.Sum(nil))

	apiFullURL := apiURL + apiMethod + "?"
	for _, arg := range sortedArgs {
		apiFullURL += arg + "&"
	}
	apiFullURL += "apiSig=" + randStr + hashSig

	logrus.Debug("calling: ", apiFullURL)

	resp, err := http.Get(apiFullURL)
	if err != nil {
		return util.Zero[T](), err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.Zero[T](), err
	}

	var res codeforcesResponse[T]
	if err := json.Unmarshal(body, &res); err != nil {
		return util.Zero[T](), err
	}
	if res.Status != "OK" {
		return util.Zero[T](), fmt.Errorf(res.Comment)
	}

	return res.Result, nil
}

func FetchCodeforcesUsersInfo(handles []string, checkHistoricHandles bool) ([]model.User, error) {
	users, err := fetchCodeforcesAPI[[]model.User]("user.info", map[string]any{
		"handles":              strings.Join(handles, ";"),
		"checkHistoricHandles": checkHistoricHandles,
	})
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return nil, errs.ErrUserNotFound(strings.Join(handles, ";"))
		}
		return nil, err
	}
	return users, nil
}

func FetchCodeforcesSubmissions(handle string, from, count int) ([]model.Submission, error) {
	return fetchCodeforcesAPI[[]model.Submission]("user.status", map[string]any{
		"handle": handle,
		"from":   from,
		"count":  count,
	})
}

func FetchCodeforcesRatingRecords(handle string) ([]model.RatingChange, error) {
	return fetchCodeforcesAPI[[]model.RatingChange]("user.rating", map[string]any{
		"handle": handle,
	})
}

func FetchCodeforcesContestList(gym bool) ([]model.CodeforcesRace, error) {
	return fetchCodeforcesAPI[[]model.CodeforcesRace]("contest.list", map[string]any{
		"gym": gym,
	})
}

func FetchCodeforcesUserInfo(handle string, checkHistoricHandles bool) (*model.User, error) {
	users, err := FetchCodeforcesUsersInfo([]string{handle}, checkHistoricHandles)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, errs.ErrUserNotFound(handle)
	}
	return &users[0], nil
}

func FetchCodeforcesRatingRecordsAfter(handle string, after time.Time) ([]model.RatingChange, error) {
	allChanges, err := FetchCodeforcesRatingRecords(handle)
	if err != nil {
		return nil, err
	}
	if after.Unix() == 0 {
		return allChanges, nil
	}
	for i, change := range allChanges {
		if time.Unix(change.At, 0).After(after) {
			return allChanges[i:], nil
		}
	}
	return []model.RatingChange{}, nil
}

func FetchCodeforcesSubmissionsAfter(handle string, after time.Time) ([]model.Submission, error) {
	perFetch := 500
	if after.IsZero() {
		perFetch = 10000
	}
	allSubmissions := make([]model.Submission, 0, perFetch)
	count := 1
	for {
		correct, err := FetchCodeforcesSubmissions(handle, count, perFetch)
		if err != nil {
			return nil, err
		}
		if len(correct) == 0 {
			break
		}
		correctStart := time.Unix(correct[0].At, 0)            // 最晚的submission
		correctEnd := time.Unix(correct[len(correct)-1].At, 0) // 最早的submission
		// 所有submission都早于期望时间
		if correctStart.Before(after) {
			break
			// 有部分submission早于期望时间
		} else if correctEnd.Before(after) {
			for _, submission := range correct {
				// 早于或等于的都不要
				if !time.Unix(submission.At, 0).After(after) {
					break
				}
				allSubmissions = append(allSubmissions, submission)
			}
			break
			// 全部submission都在期望时间之后
		} else {
			allSubmissions = append(allSubmissions, correct...)
		}
		if len(correct) < perFetch {
			break
		}
		count += perFetch
	}
	return allSubmissions, nil
}
