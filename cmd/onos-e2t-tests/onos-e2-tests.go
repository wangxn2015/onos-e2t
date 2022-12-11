// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/wangxn2015/helmit/pkg/registry"
	"github.com/wangxn2015/helmit/pkg/test"
	"github.com/wangxn2015/onos-e2t/test/e2"
)

func main() {
	registry.RegisterTestSuite("e2", &e2.TestSuite{})
	test.Main()
}
