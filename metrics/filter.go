// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metrics

import (
	"net/http"
)

type ResponseStatus int

const (
	ResponseStatusSuccess     ResponseStatus = 1 << iota // 200-299
	ResponseStatusRedirect                               // 300-399
	ResponseStatusClientError                            // 400-499
	ResponseStatusServerError                            // 500-599

	ResponseStatusNone   ResponseStatus = 0
	ResponseStatusAll                   = ResponseStatusSuccess | ResponseStatusRedirect | ResponseStatusClientError | ResponseStatusServerError
	ResponseStatusErrors                = ResponseStatusClientError | ResponseStatusServerError
)

func (l ResponseStatus) InRange(statusCode int) bool {
	return (l&ResponseStatusSuccess != 0 && statusCode >= 200 && statusCode <= 299) ||
		(l&ResponseStatusRedirect != 0 && statusCode >= 300 && statusCode <= 399) ||
		(l&ResponseStatusClientError != 0 && statusCode >= 400 && statusCode <= 499) ||
		(l&ResponseStatusServerError != 0 && statusCode >= 500 && statusCode <= 599)
}

func LimitCodes(stat ResponseStatus, next Recorder) Recorder {
	switch stat {
	case ResponseStatusNone:
		return nil
	case ResponseStatusAll:
		return next
	}
	return RecorderFunc(func(met Metrics, req *http.Request) {
		if stat.InRange(met.Code) {
			next.Record(met, req)
		}
	})
}
