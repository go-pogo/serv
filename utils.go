// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"encoding/json"
	"github.com/go-pogo/errors"
	"net/http"
)

// WriteJSON encodes v to JSON and writes it to w.
func WriteJSON(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return errors.WithStack(json.NewEncoder(w).Encode(v))
}
