package param

import (
	"github.com/suzmii/ACMBot/internal/errs"
	"strings"
)

func AsCodeforcesUsername(p []string) (string, error) {
	if len(p) == 0 {
		return "", errs.ErrNoUsername
	}
	return strings.ToLower(p[0]), nil
}
