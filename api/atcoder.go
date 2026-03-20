package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/errorx/usererr"
)

const (
	baseURL     = "https://atcoder.jp"
	apiBaseURL  = "https://kenkoooo.com/atcoder"
	rateLimit   = 50 * time.Millisecond
	httpTimeout = 30 * time.Second
)

// 全局限流器，所有 API 请求共享
var apiLimiter = rate.NewLimiter(rate.Every(rateLimit), 1)

type AtcoderUser struct {
	Username         string // Atcoder用户名
	Avatar           string // 头像URL
	Rank             string
	Rating           int
	IsProvisional    bool   // 分数是否与水平相符
	Dan              string // 段位
	PromotionMessage string // 升段信息
	HighestRating    int
	RatedMatches     int
	LastCompeted     string
}

type AtcoderSubmission struct {
	SubmissionId   uint    `json:"id"`
	SubmissionTime int64   `json:"epoch_second"` // 提交时间（UNIX时间戳）
	ProblemId      string  `json:"problem_id"`
	ContestId      string  `json:"contest_id"`
	Usernames      string  `json:"user_id"`
	Language       string  `json:"language"`
	Point          float32 `json:"point"`
	Length         int     `json:"length"`
	Status         string  `json:"result"` // 提交状态
	ExecutionTime  int     `json:"execution_time"`
}

type AtcoderProblem struct {
	Id    string  `json:"id"`
	Point float64 `json:"point"`
}

type AtcoderContest struct {
	Id             string `json:"id"`                 // 比赛ID
	StartTime      int64  `json:"start_epoch_second"` // 开始时间（UNIX时间戳）
	DurationSecond int    `json:"duration_second"`    // 持续时间（秒）
	Title          string `json:"title"`              // 完整标题
	RateChange     string `json:"rate_change"`        // Rated分数范围
}

// API 相关函数
func fetchAtcoderAPI[T any](suffix string, args map[string]any) (*T, error) {
	requestURL := apiBaseURL + "/" + suffix + "?"
	for k, v := range args {
		requestURL += k + "=" + fmt.Sprint(v) + "&"
	}
	requestURL = requestURL[:len(requestURL)-1]

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	// 使用全局限流器
	if err := apiLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	response, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res T
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (api *API) FetchAtcoderSubmissionListAfter(username string, after time.Time) (*[]AtcoderSubmission, error) {
	return fetchAtcoderAPI[[]AtcoderSubmission]("atcapi-api/v3/user/submissions", map[string]any{
		"user":        username,
		"from_second": max(0, after.Unix()),
	})
}

func (api *API) FetchAtcoderUser(username string) (*AtcoderUser, error) {
	logger := logger.WithFields(logrus.Fields{
		"handle": username,
		"action": "fetch_user",
	})

	logger.Info("Starting fetch user")
	user := &AtcoderUser{Username: username}
	var e error

	d := api.atcoderClient.Clone()

	d.OnHTML("tr", func(h *colly.HTMLElement) {
		var err error
		switch h.ChildText("th") {
		case "Rank":
			user.Rank = h.ChildText("td")
		case "Rating":
			td := h.DOM.Find("td")
			user.Rating, err = strconv.Atoi(td.Find(`span[class^="user-"]`).Text())
			if err != nil {
				logger.Infof("Failed to convert rating to int: %v", err)
				user = nil
				e = fmt.Errorf("parse error: %v", err)
			}
			user.IsProvisional = td.Find("span").HasClass("bold small")
		case "Highest Rating":
			h.DOM.Find("td").Find("span").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					user.HighestRating, err = strconv.Atoi(s.Text())
					if err != nil {
						logger.Infof("Failed to convert highest rating to int: %v", err)
						user = nil
						e = fmt.Errorf("parse error: %v", err)
					}
				case 2:
					user.Dan = s.Text()
				case 3:
					user.PromotionMessage = s.Text()
				}
			})
		case "Rated Matches":
			user.RatedMatches, err = strconv.Atoi(h.ChildText("td"))
			if err != nil {
				logger.Infof("Failed to convert rated matches to int: %v", err)
				user = nil
				e = fmt.Errorf("parse error: %v", err)
			}
		case "Last Competed":
			user.LastCompeted = h.ChildText("td")
		}
	})

	d.OnHTML("img.avatar", func(h *colly.HTMLElement) {
		user.Avatar = h.Attr("src")
	})

	d.OnError(func(r *colly.Response, err error) {
		if r == nil {
			return
		}

		switch r.StatusCode {
		case http.StatusNotFound:
			e = usererr.ErrUserNotFound(username)
		default:
			e = fmt.Errorf("failed to fetch Atcoder user: %v", err)
		}
		logger.WithFields(logrus.Fields{
			"handle": username,
			"error":  e,
		}).Info("Failed to fetch Atcoder user")
	})

	url := baseURL + "/users/" + username
	logger.Infof("Visiting: %v", url)
	_ = d.Visit(url)

	logger.Info("Completed fetch user")
	return user, e
}
