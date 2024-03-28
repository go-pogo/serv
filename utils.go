// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"encoding/json"
	"github.com/go-pogo/errors"
	"net/http"
)

const contentTypeJSON = "application/json"

// WriteJSON encodes v to JSON and writes it to w.
func WriteJSON(w http.ResponseWriter, v any) error {
	if v == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return errors.WithStack(err)
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	return nil
}

func WriteJSONError(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}

	w.WriteHeader(errors.GetStatusCodeOr(err, http.StatusInternalServerError))
	return WriteJSON(w, struct {
		Error string `json:"error"`
	}{Error: err.Error()})
}
