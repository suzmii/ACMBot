package fetcher

import (
	"github.com/suzmii/ACMBot/internal/errs"
	"testing"
	"time"
)

var (
	existedHandles   = []string{"tourist", "rng_58", "chokudai", "gori", "petr", "jiangly", "keep_running"}
	unexistedHandles = []string{"dongxuelian", "daijiangqi", "aminuosi"}
)

func TestAtcoderUserSubmissionsFromEpoch(t *testing.T) {
	for _, handle := range existedHandles {
		subs, err := FetchAtcoderSubmissionListFrom(handle, time.Unix(0, 0).Unix())
		if err != nil {
			t.Error(err)
		} else if len(*subs) == 0 {
			t.Logf("Existed user %v doesn't has any submission", handle)
		}
	}

	for _, handle := range unexistedHandles {
		subs, err := FetchAtcoderSubmissionListFrom(handle, time.Unix(0, 0).Unix())
		if err != nil {
			t.Error(err)
		} else if len(*subs) != 0 {
			t.Errorf("Unexisted user %v submission length should 0 from epoch!", handle)
		}
	}
}

func TestAtcoderUserSubmissionsFromNow(t *testing.T) {
	for _, handle := range existedHandles {
		subs, err := FetchAtcoderSubmissionListFrom(handle, time.Now().Unix())
		if err != nil {
			t.Error(err)
		} else if len(*subs) != 0 {
			t.Error("User submission length should 0 from now!")
		}
	}

	for _, handle := range unexistedHandles {
		subs, err := FetchAtcoderSubmissionListFrom(handle, time.Now().Unix())
		if err != nil {
			t.Error(err)
		} else if len(*subs) != 0 {
			t.Error("User submission length should 0 from now!")
		}
	}
}

// func TestAtcoderUserSubumissions(t *testing.T) {
// 	const maxCap = 500
// 	for _, handle := range existedHandles {
// 		subsFromEpoch, err := FetchAtcoderSubmissionListFrom(handle, time.Unix(0, 0).Unix())
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		subs, err := FetchAtcoderSubmissionList(handle)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if len(*subs) <= maxCap && len(*subs) != len(*subsFromEpoch) {
// 			t.Error("User submission length(Less than maxCap)should be equal to the length from epoch!")
// 		}
// 	}
// }

func TestAtcoderProblems(t *testing.T) {
	problems, err := FetchAtcoderProblemList()
	if err != nil {
		t.Error(err)
	} else if len(*problems) == 0 {
		t.Error("Problem list should not be empty!")
	}
}

func TestAtcoderContests(t *testing.T) {
	contests, err := FetchAtcoderContestList()
	if err != nil {
		t.Error(err)
	} else if len(*contests) == 0 {
		t.Error("Contest list should not be empty!")
	}
}

func TestAtcoderUser(t *testing.T) {
	// 动态生成测试用例
	var tests []struct {
		name    string
		handle  string
		wantErr error
	}

	// 添加有效用户测试
	for _, handle := range existedHandles {
		tests = append(tests, struct {
			name    string
			handle  string
			wantErr error
		}{
			name:    "valid user - " + handle,
			handle:  handle,
			wantErr: nil,
		})
	}

	// 添加无效用户测试
	for _, handle := range unexistedHandles {
		tests = append(tests, struct {
			name    string
			handle  string
			wantErr error
		}{
			name:    "invalid user - " + handle,
			handle:  handle,
			wantErr: &errs.ErrHandleNotFound{Handle: handle},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := FetchAtcoderUser(tt.handle)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("FetchUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if user == nil {
					t.Error("FetchUser() returned nil user without error")
				}
			} else {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("FetchUser() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
