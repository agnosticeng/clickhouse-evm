package jsonrpc

import (
	"net/url"
	"strconv"
)

type CallOptionsFunc func(*CallOptions)

type CallOptions struct {
	maxBatchSize            int
	maxConcurrentRequests   int
	disableBatch            bool
	failOnError             bool
	failOnRetryableError    bool
	failOnNull              bool
	retryableErrorPredicate func(error) bool
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

	return nil
}

func NewCallOptions(optFuncs ...CallOptionsFunc) *CallOptions {
	var res CallOptions

	for _, optFunc := range optFuncs {
		optFunc(&res)
	}

	return &res
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
