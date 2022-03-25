package sfdc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type QueryResponse[T any] struct {
	Done           bool   `json:"done" sfdc:"-"`
	NextRecordsURL string `json:"nextRecordsUrl" sfdc:"-"`
	Records        []T    `json:"records" sfdc:"-"`
	TotalSize      int    `json:"totalSize" sfdc:"-"`
}

func errorForResponse(r io.Reader) error {
	var errorMessage []struct {
		Message   string `json:"message"`
		ErrorCode string `json:"errorCode"`
	}
	if err := json.NewDecoder(r).Decode(&errorMessage); err != nil {
		return errors.New("got bad result")
	}
	return fmt.Errorf("%s (%s)", errorMessage[0].Message, errorMessage[0].ErrorCode)
}
