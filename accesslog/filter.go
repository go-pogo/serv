// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
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

func LimitCodes(rs ResponseStatus, next Logger) Logger {
	if next == nil {
		return nil
	}

	switch rs {
	case ResponseStatusNone:
		return nil
	case ResponseStatusAll:
		return next
	}
	return loggerFunc(func(ctx context.Context, det Details, req *http.Request) {
		if rs.InRange(det.StatusCode) {
			next.Log(ctx, det, req)
		}
	})
}
