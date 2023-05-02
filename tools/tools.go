//go:build tools
// +build tools

/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package tools

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
