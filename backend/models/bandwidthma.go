// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Model for BandwidthMA
	BandwidthMA tracks moving average for download and upload bandwidth
*/

package models

import "github.com/mxmCherry/movavg"

type BandwidthMA struct {
	Download movavg.MA
	Upload   movavg.MA
}
