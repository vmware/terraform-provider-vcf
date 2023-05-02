/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/vmware/terraform-provider-vcf/internal/provider"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return provider.Provider()
		},
	}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "vmware/vcf"
	}

	plugin.Serve(opts)
}
