package model

type Result struct {
	Data *CEPResponse
	Fail *ErrResult
}

type ErrResult struct {
	Err        error
	StatusCode int
}
