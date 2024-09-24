package jsonrpc

type CallOptionsFunc func(*CallOptions)

type CallOptions struct {
	maxBatchSize          int
	maxConcurrentRequests int
	disableBatch          bool
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
