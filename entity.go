package sfdc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
	"unicode/utf8"
)

type Entity[T any] struct {
	instance     *Instance
	name         string
	allFields    string
	taggedFields string
}

func NewEntity[T any](instance *Instance) *Entity[T] {
	var result Entity[T]
	var example T
	typ := reflect.TypeOf(example)
	name := typ.Name()
	result.instance = instance
	result.name = name
	result.taggedFields = fieldsForType(typ)
	return &result
}

func (e *Entity[T]) SetName(name string) {
	e.name = name
}

func (e *Entity[T]) TaggedFields() string {
	return e.taggedFields
}

func (e *Entity[T]) BuildQuery(fields string, constraints string) string {
	query := fmt.Sprintf("SELECT %v FROM %s", fields, e.name)
	if utf8.RuneCountInString(constraints) > 0 {
		query = fmt.Sprintf("%v WHERE %v", query, constraints)
	}
	return query
}

func (e *Entity[T]) Query(ctx context.Context, query string) ([]T, error) {
	var r QueryResponse[T]
	results := []T{}
	uri, err := e.instance.QueryAllURL()
	if err != nil {
		return nil, err
	}
	q := uri.Query()
	q.Set("q", query)
	uri.RawQuery = q.Encode()
	reqURI := uri.String()
	for !r.Done {
		if r.NextRecordsURL != "" {
			reqURI = fmt.Sprintf("%v%v", e.instance.url, r.NextRecordsURL)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURI, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		res, err := e.instance.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode >= 400 {
			return nil, errorForResponse(res.Body)
		}
		r.Records = nil
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return nil, err
		}
		results = append(results, r.Records...)
	}

	return results, nil
}

// QueryAsync returns a channel that T are written to.
// The channel is closed when all records have been written.
// Errors are written to the returned error channel.
// The query aborts when an error is encountered.
func (e *Entity[T]) QueryAsync(ctx context.Context, query string) (<-chan []T, <-chan error) {
	result := make(chan []T)
	errs := make(chan error, 1)
	uri, err := e.instance.QueryAllURL()
	if err != nil {
		errs <- err
		close(result)
		return result, errs
	}
	go func() {
		var r QueryResponse[T]

		q := uri.Query()
		q.Set("q", query)
		uri.RawQuery = q.Encode()
		reqURI := uri.String()
		for !r.Done {
			if r.NextRecordsURL != "" {
				reqURI = fmt.Sprintf("%v%v", e.instance.url, r.NextRecordsURL)
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURI, nil)
			if err != nil {
				errs <- err
				break
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			res, err := e.instance.client.Do(req)
			if err != nil {
				errs <- err
				break
			}
			defer res.Body.Close()
			if res.StatusCode >= 400 {
				errs <- errorForResponse(res.Body)
				break
			}
			r.Records = nil
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				errs <- err
				break
			}
			result <- r.Records
		}
		close(result)
	}()

	return result, errs
}

// List finds all T objects.
func (e *Entity[T]) List(ctx context.Context) (<-chan []T, <-chan error) {
	return e.QueryAsync(ctx, "")
}

// ListModifiedSince finds all T objects modified since some point in time.
func (e *Entity[T]) ListModifiedSince(ctx context.Context, since time.Time) (<-chan []T, <-chan error) {
	return e.QueryAsync(ctx, fmt.Sprintf("LastModifiedDate > %s", since.Format(DateTimeLayout)))
}
