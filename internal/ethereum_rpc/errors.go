package ethereum_rpc

import (
	"regexp"
)

var (
	retryableErrorPatterns = []*regexp.Regexp{
		regexp.MustCompile(`intrinsic gas too low`),
	}
)

func RetryableErrorPredicate(err error) bool {
	for _, re := range retryableErrorPatterns {
		if re.MatchString(err.Error()) {
			return true
		}
	}

	return false
}
