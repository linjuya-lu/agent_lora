// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/v4/pkg/startup"
	"github.com/linjuya-lu/device_lora_agent_go/internal/driver"
	"github.com/linjuya-lu/device_lora_agent_go/internal/version"
)

const (
	serviceName string = "device-lora-agent"
)

func main() {
	d := driver.NewUartDeviceDriver()
	startup.Bootstrap(serviceName, version.Version, d)
}
