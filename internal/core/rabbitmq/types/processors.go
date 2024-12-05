package types

type Identifiable interface {
	Id() int64
}

type RequestProcessor[Request Identifiable] interface {
	// ProcessRequest processes the request and returns whether the request can be reprocessed
	ProcessRequest(request Request) (reprocessable bool, err error)
	// ReprocessFailedCallback is a callback to be called if the reprocessing fails
	ReprocessFailedCallback(request Request) error
}

type BatchProcessor[Request Identifiable] interface {
	// ProcessBatch processes the batch and returns whether the batch should be reprocessed
	ProcessBatch(batch []Request) (reprocessable bool, err error)
	// ReprocessFailedCallback is a callback to be called if the reprocessing fails
	ReprocessFailedCallback(batch []Request) error
}
