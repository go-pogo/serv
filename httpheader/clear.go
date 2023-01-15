// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpheader

const (
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Clear-Site-Data
	ClearSiteData = "Clear-Site-Data"

	ClearSiteData_All     ClearSiteDataDirective = "*"
	ClearSiteData_Cache   ClearSiteDataDirective = "cache"
	ClearSiteData_Cookies ClearSiteDataDirective = "cookies"
	ClearSiteData_Storage ClearSiteDataDirective = "storage"
)

type ClearSiteDataDirective string

func (csd ClearSiteDataDirective) String() string { return string(csd) }

func SetClearSiteData(h Header, directives ...ClearSiteDataDirective) {
	h.Del(ClearSiteData)
	for _, d := range directives {
		h.Add(ClearSiteData, d.String())
	}
}

func SetClearSiteDataAll(h Header) { h.Set(ClearSiteData, "*") }
