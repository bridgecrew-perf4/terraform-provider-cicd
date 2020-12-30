// Copyright 2017-2020 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"

	"github.com/AtlantPlatform/infra/terraform-provider-cicd/cicd"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cicd.Provider,
	})
}
