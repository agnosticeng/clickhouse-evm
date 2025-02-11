package jsonrpc

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
)

type CallOptionsFunc func(*CallOptions)

type CallOptions struct {
	maxBatchSize          int
	maxConcurrentRequests int
	disableBatch          bool

	failOnError             bool
	failOnRetryableError    bool
	failOnNull              bool
	retryableErrorPredicate func(error) bool
	retryableStatusCodes    []int

	retryInitialInterval     time.Duration
	retryRandomizationFactor float64
	retryMultiplier          float64
	retryMaxInterval         time.Duration
	retryMaxElapsedTime      time.Duration
	retryMaxTries            uint
}

func WithMatchBatchSize(value int) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.maxBatchSize = value
	}
}

func WithMaxConcurrentRequests(value int) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.maxConcurrentRequests = value
	}
}

func WithDisableBatch(b bool) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.disableBatch = b
	}
}

func WithFailOnError(b bool) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.failOnError = b
	}
}

func WithFailOnNull(b bool) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.failOnNull = b
	}
}

func WithRetryableErrorPredicate(f func(error) bool) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryableErrorPredicate = f
	}
}

func WithRetryableStatusCodes(statuses []int) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryableStatusCodes = statuses
	}
}

func WithRetryInitialInterval(d time.Duration) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryInitialInterval = d
	}
}

func WithRetryRandomizationFactor(f float64) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryRandomizationFactor = f
	}
}

func WithRetryMultiplier(m float64) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryMultiplier = m
	}
}

func WithRetryMaxInterval(d time.Duration) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryMaxInterval = d
	}
}

func WithRetryMaxElapsedTime(d time.Duration) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryMaxElapsedTime = d
	}
}

func WithRetryMaxTries(i uint) CallOptionsFunc {
	return func(opts *CallOptions) {
		opts.retryMaxTries = i
	}
}

func (opts *CallOptions) ParseFromEndpoint(endpoint string) error {
	u, err := url.Parse(endpoint)

	if err != nil {
		return err
	}

	frag, err := url.ParseQuery(u.Fragment)

	if err != nil {
		return err
	}

	if s := frag.Get("max-batch-size"); len(s) > 0 {
		if i, err := strconv.ParseInt(s, 10, 64); err != nil {
			return err
		} else {
			opts.maxBatchSize = int(i)
		}
	}

	if s := frag.Get("max-concurrent-requests"); len(s) > 0 {
		if i, err := strconv.ParseInt(s, 10, 64); err != nil {
			return err
		} else {
			opts.maxConcurrentRequests = int(i)
		}
	}

	if s := frag.Get("disable-batch"); len(s) > 0 {
		if b, err := strconv.ParseBool(s); err != nil {
			return err
		} else {
			opts.disableBatch = b
		}
	}

	if s := frag.Get("fail-on-error"); len(s) > 0 {
		if b, err := strconv.ParseBool(s); err != nil {
			return err
		} else {
			opts.failOnError = b
		}
	}

	if s := frag.Get("fail-on-retryable-error"); len(s) > 0 {
		if b, err := strconv.ParseBool(s); err != nil {
			return err
		} else {
			opts.failOnRetryableError = b
		}
	}

	if s := frag.Get("fail-on-null"); len(s) > 0 {
		if b, err := strconv.ParseBool(s); err != nil {
			return err
		} else {
			opts.failOnNull = b
		}
	}

	if s := frag.Get("retryable-status-codes"); len(s) > 0 {
		clear(opts.retryableStatusCodes)

		for _, v := range strings.Split(s, ",") {
			if len(v) == 0 {
				continue
			}

			if i, err := strconv.ParseInt(v, 10, 64); err != nil {
				return err
			} else {
				opts.retryableStatusCodes = append(opts.retryableStatusCodes, int(i))
			}
		}
	}

	if s := frag.Get("retry-initial-interval"); len(s) > 0 {
		if d, err := time.ParseDuration(s); err != nil {
			return err
		} else {
			opts.retryInitialInterval = d
		}
	}

	if s := frag.Get("retry-randomization-factor"); len(s) > 0 {
		if f, err := strconv.ParseFloat(s, 64); err != nil {
			return err
		} else {
			opts.retryRandomizationFactor = f
		}
	}

	if s := frag.Get("retry-multiplier"); len(s) > 0 {
		if f, err := strconv.ParseFloat(s, 64); err != nil {
			return err
		} else {
			opts.retryMultiplier = f
		}
	}

	if s := frag.Get("retry-max-interval"); len(s) > 0 {
		if d, err := time.ParseDuration(s); err != nil {
			return err
		} else {
			opts.retryMaxInterval = d
		}
	}

	if s := frag.Get("retry-max-elapsed-time"); len(s) > 0 {
		if d, err := time.ParseDuration(s); err != nil {
			return err
		} else {
			opts.retryMaxElapsedTime = d
		}
	}

	if s := frag.Get("retry-max-tries"); len(s) > 0 {
		if i, err := strconv.ParseUint(s, 10, 64); err != nil {
			return err
		} else {
			opts.retryMaxTries = uint(i)
		}
	}

	return nil
}

func (opts *CallOptions) ToExponentialBackoff() *backoff.ExponentialBackOff {
	var bo = backoff.NewExponentialBackOff()

	if opts.retryInitialInterval > 0 {
		bo.InitialInterval = opts.retryInitialInterval
	}

	if opts.retryRandomizationFactor > 0 {
		bo.RandomizationFactor = opts.retryRandomizationFactor
	}

	if opts.retryMultiplier > 0 {
		bo.Multiplier = opts.retryMultiplier
	}

	if opts.retryMaxInterval > 0 {
		bo.MaxInterval = opts.retryMaxInterval
	}

	return bo
}

func (opts *CallOptions) GetRetryMaxElapsedTimeOrDefault() time.Duration {
	if opts.retryMaxElapsedTime > 0 {
		return opts.retryMaxElapsedTime
	} else {
		return time.Minute * 15
	}
}

func (opts *CallOptions) GetRetryMaxTriesOrDefault() uint {
	if opts.retryMaxTries > 0 {
		return opts.retryMaxTries
	} else {
		return 20
	}
}

func NewCallOptions(optFuncs ...CallOptionsFunc) *CallOptions {
	var res CallOptions

	for _, optFunc := range optFuncs {
		optFunc(&res)
	}

	return &res
}
