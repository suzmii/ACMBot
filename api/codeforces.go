package api

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/errorx/usererr"
	"github.com/suzmii/ACMBot/util"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"

	"math/rand"
)

var limiter = rate.NewLimiter(4, 1)

type CodeforcesResponse[T any] struct {
	/*
		codeforces响应数据的基本格式:
			Result->期望的数据
			Comment->失败时返回的提示信息
	*/
	Status  string `json:"status"`
	Result  T      `json:"result"`
	Comment string `json:"comment"`
}

type CodeforcesAPIError struct {
	Api      string
	Args     map[string]any
	Response CodeforcesResponse[any]
}

func (e CodeforcesAPIError) Error() string {
	return fmt.Sprintf("Codeforces API Error: api: %s, args: %v, response: %v", e.Api, e.Args, e.Response)
}

func (e CodeforcesAPIError) TransToUserErr() error {
	if strings.Contains(e.Response.Comment, "not found") {
		return usererr.ErrUserNotFound()
	}
	// Other types of errors can be handled here as needed

	return usererr.ErrCodeforcesAPIFailure
}

func fetchCodeforcesAPI[T any](apiMethod string, args map[string]any, cfg *subconfig.API) (T, error) {
	logger.Tracef("fetchCodeforcesAPI-finished, %v, %v", apiMethod, args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	err = limiter.Wait(ctx)
	if err != nil {
		return util.Zero[T](), fmt.Errorf("failed to wait for rate limiter: %w; usererr: %w", err, usererr.ErrTooManyRequests)
	}

	apiURL := "https://codeforces.com/api/"

	args["apiKey"] = cfg.CodeforcesKey
	args["time"] = strconv.Itoa(int(time.Now().Unix()))

	var sortedArgs []string
	for k, v := range args {
		sortedArgs = append(sortedArgs, fmt.Sprintf("%v=%v", k, v))
	}
	sort.Strings(sortedArgs)

	randStr := strconv.Itoa(rand.Intn(900000) + 100000)
	hashSource := randStr + "/" + apiMethod + "?" + strings.Join(sortedArgs, "&") + "#" + cfg.CodeforcesSecret

	h := sha512.New()
	h.Write([]byte(hashSource))
	hashSig := hex.EncodeToString(h.Sum(nil))

	apiFullURL := apiURL + apiMethod + "?"
	for _, arg := range sortedArgs {
		apiFullURL += arg + "&"
	}
	apiFullURL += "apiSig=" + randStr + hashSig

	logger.Debug("calling: ", apiFullURL)

	resp, err := http.Get(apiFullURL)
	if err != nil {
		return util.Zero[T](), fmt.Errorf("failed to call codeforces api: %w; usererr: %w", err, usererr.ErrCodeforcesAPIFailure)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.Zero[T](), fmt.Errorf("failed to read response body while call codeforces api: %w; usererr: %w", err, usererr.ErrCodeforcesAPIFailure)
	}

	// status code 不为 200 时, 也会正常返回 body, 故无需处理该情况
	// if resp.StatusCode != http.StatusOK {
	// 	return util.Zero[T](), fmt.Errorf("codeforces return bad status code: %s; usererr: %w", resp.Status, usererr.ErrCodeforcesAPIFailure)
	// }

	var res CodeforcesResponse[T]
	if err := json.Unmarshal(body, &res); err != nil {
		return util.Zero[T](), fmt.Errorf("failed to unmarshal codeforces response: %w; usererr: %w", err, usererr.ErrCodeforcesAPIFailure)
	}
	if res.Status != "OK" {
		cferr := CodeforcesAPIError{
			Api:  apiMethod,
			Args: args,
			Response: CodeforcesResponse[any]{
				Status:  res.Status,
				Result:  res.Result,
				Comment: res.Comment,
			},
		}
		return util.Zero[T](), fmt.Errorf("codeforces api call successfully but status is not ok: %w; usererr: %w", cferr, cferr.TransToUserErr())
	}
	logger.Tracef("fetchCodeforcesAPI-finished, %v, %v, %v", apiMethod, args, res.Result)
	return res.Result, nil
}

type User struct {
	Username     string `json:"handle"`
	Avatar       string `json:"titlePhoto"`
	Rating       int    `json:"rating"`
	MaxRating    int    `json:"maxRating"`
	FriendCount  int    `json:"friendOfCount"`
	Organization string `json:"organization"`
	CreatedAt    int64  `json:"registrationTimeSeconds"`
	Rank         string `json:"rank"`    // 称号
	MaxRank      string `json:"maxRank"` // 最高称号
}

func (api *API) FetchCodeforcesUsersInfo(handles []string, checkHistoricHandles bool) ([]User, error) {
	users, err := fetchCodeforcesAPI[[]User]("user.info", map[string]any{
		"handles":              strings.Join(handles, ";"),
		"checkHistoricHandles": checkHistoricHandles,
	}, &api.cfg)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return nil, usererr.ErrUserNotFound(strings.Join(handles, ";"))
		}
		return nil, err
	}
	return users, nil
}

func (api *API) FetchCodeforcesUserInfo(handle string, checkHistoricHandles bool) (*User, error) {
	users, err := api.FetchCodeforcesUsersInfo([]string{handle}, checkHistoricHandles)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, usererr.ErrUserNotFound(handle)
	}
	return &users[0], nil
}

type Problem struct {
	ContestID      int      `json:"contestId"`
	ProblemSetName string   `json:"problemsetName"`
	Index          string   `json:"index"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Rating         int      `json:"rating"`
	Tags           []string `json:"tags"`
}

type Submission struct {
	ID                  uint    `json:"id"`
	ContestID           int     `json:"contestId"`
	At                  int64   `json:"creationTimeSeconds"`
	RelativeTimeSeconds int64   `json:"relativeTimeSeconds"` // 相对于比赛开始的时间（秒）
	Status              string  `json:"verdict"`
	Problem             Problem `json:"problem"`
	Language            string  `json:"programmingLanguage"`
}

func (api *API) FetchCodeforcesSubmissions(handle string, from, count int) ([]Submission, error) {
	return fetchCodeforcesAPI[[]Submission]("user.status", map[string]any{
		"handle": handle,
		"from":   from,
		"count":  count,
	}, &api.cfg)
}

func (api *API) FetchCodeforcesSubmissionsAfter(handle string, after time.Time) ([]Submission, error) {
	perFetch := 500
	if after.IsZero() {
		perFetch = 10000
	}
	allSubmissions := make([]Submission, 0, perFetch)
	count := 1
	for {
		correct, err := api.FetchCodeforcesSubmissions(handle, count, perFetch)
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
		} else if correctEnd.Before(after) { // 有部分submission早于期望时间
			for _, submission := range correct {
				// 早于或等于的都不要
				if !time.Unix(submission.At, 0).After(after) {
					break
				}
				allSubmissions = append(allSubmissions, submission)
			}
			break
		} else { // 全部submission都在期望时间之后
			allSubmissions = append(allSubmissions, correct...)
		}
		if len(correct) < perFetch {
			break
		}
		count += perFetch
	}
	return allSubmissions, nil
}

type RatingRecord struct {
	At          int64  `json:"ratingUpdateTimeSeconds"`
	Rank        int    `json:"rank"`
	Username    string `json:"handle"`
	ContestID   int    `json:"contestId"`
	ContestName string `json:"contestName"`
	OldRating   int    `json:"oldRating"`
	NewRating   int    `json:"newRating"`
}

func (api *API) FetchCodeforcesRatingRecords(handle string) ([]RatingRecord, error) {
	return fetchCodeforcesAPI[[]RatingRecord]("user.rating", map[string]any{
		"handle": handle,
	}, &api.cfg)
}

type CodeforcesRace struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int64  `json:"durationSeconds"`
	StartTimeSeconds    int64  `json:"startTimeSeconds"`
	RelativeTimeSeconds int64  `json:"relativeTimeSeconds"`
}

func (api *API) FetchCodeforcesContestList(gym bool) ([]CodeforcesRace, error) {
	return fetchCodeforcesAPI[[]CodeforcesRace]("contest.list", map[string]any{
		"gym": gym,
	}, &api.cfg)
}

type ProblemSample struct {
	Input  string
	Output string
}

type ProblemStatement struct {
	URL         string
	Title       string
	TimeLimit   string
	MemoryLimit string
	Statement   []string
	Input       []string
	Output      []string
	Samples     []ProblemSample
}

func BuildCodeforcesProblemURL(problem Problem) string {
	if problem.ContestID > 0 {
		return fmt.Sprintf("https://codeforces.com/problemset/problem/%d/%s", problem.ContestID, problem.Index)
	}
	if problem.ProblemSetName != "" {
		return fmt.Sprintf("https://codeforces.com/problemset/problem/%s/%s", problem.ProblemSetName, problem.Index)
	}
	return ""
}

func (api *API) FetchCodeforcesProblemStatement(problem Problem) (*ProblemStatement, error) {
	problemURL := BuildCodeforcesProblemURL(problem)
	if problemURL == "" {
		return nil, fmt.Errorf("invalid problem identity: contestId=%d, problemsetName=%s, index=%s", problem.ContestID, problem.ProblemSetName, problem.Index)
	}

	req, err := http.NewRequest(http.MethodGet, problemURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for problem page: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch problem page: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch problem page, status=%s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse problem page: %w", err)
	}

	ps := doc.Find(".problem-statement").First()
	if ps.Length() == 0 {
		return nil, fmt.Errorf("failed to parse statement: .problem-statement not found")
	}

	statement := &ProblemStatement{
		URL:         problemURL,
		Title:       normalizeText(ps.Find(".title").First().Text()),
		TimeLimit:   normalizeText(ps.Find(".time-limit").First().Text()),
		MemoryLimit: normalizeText(ps.Find(".memory-limit").First().Text()),
	}

	if statement.Title == "" {
		statement.Title = fmt.Sprintf("%s %s", problem.ProblemSetName, problem.Index)
	}

	ps.Children().Each(func(_ int, child *goquery.Selection) {
		class, _ := child.Attr("class")
		switch class {
		case "header", "input-specification", "output-specification", "sample-test", "note":
			return
		}
		child.Find("p").Each(func(_ int, p *goquery.Selection) {
			if text := normalizeText(p.Text()); text != "" {
				statement.Statement = append(statement.Statement, text)
			}
		})
	})
	if len(statement.Statement) == 0 {
		ps.Find("p").Each(func(_ int, p *goquery.Selection) {
			if text := normalizeText(p.Text()); text != "" {
				statement.Statement = append(statement.Statement, text)
			}
		})
	}

	ps.Find(".input-specification p").Each(func(_ int, p *goquery.Selection) {
		if text := normalizeText(p.Text()); text != "" {
			statement.Input = append(statement.Input, text)
		}
	})
	ps.Find(".output-specification p").Each(func(_ int, p *goquery.Selection) {
		if text := normalizeText(p.Text()); text != "" {
			statement.Output = append(statement.Output, text)
		}
	})

	inputSamples := ps.Find(".sample-test .input pre")
	outputSamples := ps.Find(".sample-test .output pre")
	count := inputSamples.Length()
	if outputSamples.Length() < count {
		count = outputSamples.Length()
	}
	for i := 0; i < count; i++ {
		statement.Samples = append(statement.Samples, ProblemSample{
			Input:  extractPreText(inputSamples.Eq(i)),
			Output: extractPreText(outputSamples.Eq(i)),
		})
	}

	return statement, nil
}

var preTagRegex = regexp.MustCompile(`<[^>]+>`)

func extractPreText(sel *goquery.Selection) string {
	if sel == nil || sel.Length() == 0 {
		return ""
	}
	raw, err := sel.Html()
	if err != nil {
		return normalizeText(sel.Text())
	}
	raw = strings.ReplaceAll(raw, "<br/>", "\n")
	raw = strings.ReplaceAll(raw, "<br />", "\n")
	raw = strings.ReplaceAll(raw, "<br>", "\n")
	raw = preTagRegex.ReplaceAllString(raw, "")
	raw = html.UnescapeString(raw)
	return strings.TrimSpace(raw)
}

func normalizeText(s string) string {
	s = html.UnescapeString(s)
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	parts := strings.Fields(s)
	return strings.Join(parts, " ")
}

func EscapeCodeforcesProblemTag(tag string) string {
	return url.QueryEscape(strings.TrimSpace(tag))
}
