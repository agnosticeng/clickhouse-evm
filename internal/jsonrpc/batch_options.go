package jsonrpc

type BatchOptionsFunc func(*BatchOptions)

type BatchOptions struct {
	chunkSize        int
	concurrencyLimit int
}

func NewBatchOptions(optFuncs ...BatchOptionsFunc) *BatchOptions {
	var res BatchOptions

	for _, optFunc := range optFuncs {
		optFunc(&res)
	}

	return &res
}

func WithChunkSize(value int) BatchOptionsFunc {
	return func(opts *BatchOptions) {
		opts.chunkSize = value
	}
}

func WithConcurrencyLimit(value int) BatchOptionsFunc {
	return func(opts *BatchOptions) {
		opts.concurrencyLimit = value
	}
}
