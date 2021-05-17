// +build helmchart

// Copyright 2017-2021 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package helmchart

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilder_GetHash(t *testing.T) {
	r := require.New(t)
	dir, err := os.Getwd()
	r.NoError(err)

	path := dir + "/../../../karta-infra/terraform/modules/apps/rentals/charts/atlant-rentals-web"
	chart, err := New(path, map[string]interface{}{})
	r.NoError(err)
	// fmt.Println("chart.Hash", chart.Hash)
	r.Equal(len(chart.Hash), 64)
}
