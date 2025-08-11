// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver provides an implementation of a ProtocolDriver interface.
package driver

import (
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/device-sdk-go/v4/pkg/interfaces"
	dsModels "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/linjuya-lu/device_lora_agent_go/internal/config"
	"github.com/linjuya-lu/device_lora_agent_go/internal/mqtt/mqttclient"
)

type LoraAgentDriver struct {
	lc         logger.LoggingClient
	asyncCh    chan<- *dsModels.AsyncValues
	locker     sync.Mutex
	sdk        interfaces.DeviceServiceSDK
	mqttClient mqtt.Client
}

var once sync.Once
var driver *LoraAgentDriver

func NewUartDeviceDriver() interfaces.ProtocolDriver {
	once.Do(func() {
		driver = new(LoraAgentDriver)
	})
	return driver
}

func (d *LoraAgentDriver) Initialize(sdk interfaces.DeviceServiceSDK) error {
	d.sdk = sdk
	d.lc = sdk.LoggingClient()
	d.asyncCh = sdk.AsyncValuesChannel()

	// —— 1. 初始化 MQTT 客户端 —— //
	brokerURL := "tcp://localhost:1883"
	clientID := "lora-agent-client"

	client, err := mqttclient.NewClient(brokerURL, clientID)
	if err != nil {
		return fmt.Errorf("初始化 MQTT 客户端失败: %w", err)
	}
	d.mqttClient = client

	return nil
}

func (d *LoraAgentDriver) Start() error {
	// 配置文件和串口参数
	devicesYAML := "../cmd/res/devices/devices.yaml"
	profilesDir := "../cmd/res/profiles"

	// 初始化资源
	if err := config.InitDeviceResources(devicesYAML, profilesDir); err != nil {
		return fmt.Errorf("初始化设备资源失败: %w", err)
	}
	d.lc.Infof("lora 代理已启动")
	return nil
}

func (d *LoraAgentDriver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	d.locker.Lock()
	defer d.locker.Unlock()
	d.lc.Infof("HandleReadCommands 调用: 设备=%s, 请求资源数=%d", deviceName, len(reqs))
	// 取缓存
	values, ok := config.GetDeviceValues(deviceName)
	if !ok {
		return nil, fmt.Errorf("设备 %s 未找到或无可用值", deviceName)
	}
	for _, req := range reqs {
		resName := req.DeviceResourceName

		//  静态资源从 config 缓存读取
		val, exists := values[resName]
		if !exists {
			return nil, fmt.Errorf("设备 %s 上未找到资源 %s 的值", deviceName, resName)
		}
		cv := &dsModels.CommandValue{
			DeviceResourceName: resName,
			Type:               req.Type,
			Value:              val,
			Origin:             time.Now().UnixNano(),
			Tags:               map[string]string{},
		}
		d.lc.Infof("读取值: %s.%s = %v", deviceName, resName, val)
		res = append(res, cv)
	}
	return res, nil
}

func (d *LoraAgentDriver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {
	d.locker.Lock()
	defer d.locker.Unlock()

	d.lc.Infof("HandleWriteCommands 调用: 设备=%s, 写入请求数=%d", deviceName, len(reqs))
	// 请求数与参数数必须一致
	if len(reqs) != len(params) {
		d.lc.Errorf("请求数与参数数不匹配: %d vs %d", len(reqs), len(params))
		return fmt.Errorf("请求数与参数数不匹配")
	}
	for i, req := range reqs {
		resName := req.DeviceResourceName
		cv := params[i]
		// 命令类型转换
		v, _ := cv.Int8Value()
		d.lc.Infof("Int8Value = %d", v)
		// 如果是时间参数查询且值为 1
		if resName == "Time_Parameter_Query" && v == 1 {
			if err := d.handleTimeParameterQuery(deviceName); err != nil {
				return err
			}
		}
		// 如果是时间参数设置且值为 1
		if resName == "Time_Parameter_Set" && v == 1 {
			if err := d.handleTimeParameterSet(deviceName); err != nil {
				return err
			}
		}
		// 如果是复位命令且值为 1
		if resName == "Reset_Set" && v == 1 {
			if err := d.handleResetCommand(deviceName); err != nil {
				return err
			}
		}
		// 如果是ID查询命令且值为 1
		if resName == "ID_Query" && v == 1 {
			if err := d.handleIdQuery(deviceName); err != nil {
				return err
			}
		}
		// 如果是所有通用参数查询命令且值为 1
		if resName == "General_Parameter_Query" && v == 1 {
			if err := d.handleGeneParaQuery(deviceName); err != nil {
				return err
			}
		}
		// 如果是所有告警数据查询命令且值为 1
		if resName == "Alarm_Parameter_Query" && v == 1 {
			if err := d.handleIdAlarmParaQuery(deviceName); err != nil {
				return err
			}
		}
		// 如果是所有检测参数查询命令且值为 1
		if resName == "Monitoring_Data_Query" && v == 1 {
			if err := d.handleIdMoniDataQuery(deviceName); err != nil {
				return err
			}
		}
		// 如果是网络拓扑查询命令且值为 1
		if resName == "Router_Parameter_Query" && v == 1 {

		}
	}
	return nil
}

func (d *LoraAgentDriver) Stop(force bool) error {
	d.lc.Info("LoraAgentDriver.Stop: device-lora-agent driver is stopping...")

	return nil
}

func (d *LoraAgentDriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debugf("新设备已添加: %s", deviceName)
	// 获取 Device 对象
	dev, err := d.sdk.GetDeviceByName(deviceName)
	if err != nil {
		return fmt.Errorf("获取设备 %s 失败: %w", deviceName, err)
	}
	// 从 Device 中取出 Profile 名称
	profileName := dev.ProfileName
	// 获取对应的 DeviceProfile
	prof, err := d.sdk.GetProfileByName(profileName)
	if err != nil {
		return fmt.Errorf("获取设备配置文件 %s 失败: %w", profileName, err)
	}
	// 针对每个资源执行初始化，传递默认值和类型
	for _, dr := range prof.DeviceResources {
		resName := dr.Name
		defaultValue := dr.Properties.DefaultValue
		valueType := dr.Properties.ValueType
		if err := config.DeviceInit(deviceName, resName, defaultValue, valueType); err != nil {
			return fmt.Errorf("初始化设备 %s 资源 %s 失败：%v", deviceName, resName, err)
		}
		d.lc.Infof("已将设备 %s 的资源 %s 初始化为默认值: %s (类型: %s)", deviceName, resName, defaultValue, valueType)
	}
	return nil
}

func (d *LoraAgentDriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debugf("Device %s is updated", deviceName)

	dev, err := d.sdk.GetDeviceByName(deviceName)
	if err != nil {
		return fmt.Errorf("获取设备 %s 失败: %w", deviceName, err)
	}
	profileName := dev.ProfileName
	prof, err := d.sdk.GetProfileByName(profileName)
	if err != nil {
		return fmt.Errorf("获取设备配置文件 %s 失败: %w", profileName, err)
	}
	// 针对每个资源重新初始化，传递默认值和类型
	for _, dr := range prof.DeviceResources {
		resName := dr.Name
		defaultValue := dr.Properties.DefaultValue
		valueType := dr.Properties.ValueType
		if err := config.DeviceInit(deviceName, resName, defaultValue, valueType); err != nil {
			return fmt.Errorf("更新设备 %s 资源 %s 失败：%v", deviceName, resName, err)
		}
		d.lc.Infof("已将设备 %s 的资源 %s 重新初始化为默认值: %s (类型: %s)", deviceName, resName, defaultValue, valueType)
	}

	d.lc.Infof("已刷新设备 %s 的资源值为最新默认配置", deviceName)
	return nil
}

func (d *LoraAgentDriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.lc.Debugf("Device %s is removed", deviceName)

	// 删除运行时值表
	if err := config.DeleteDeviceValues(deviceName); err != nil {
		d.lc.Errorf("删除设备 %s 的运行时值失败: %v", deviceName, err)
		return fmt.Errorf("删除设备 %s 的运行时值失败: %w", deviceName, err)
	}
	// 删除 sensorID 到 deviceName 的所有映射
	if err := config.DeleteSensorIDMappingsByDevice(deviceName); err != nil {
		d.lc.Errorf("删除设备 %s 的传感器映射失败: %v", deviceName, err)
		return fmt.Errorf("删除设备 %s 的传感器映射失败: %w", deviceName, err)
	}
	d.lc.Infof("已移除设备 %s 的所有运行时数据和映射", deviceName)
	return nil
}

func (d *LoraAgentDriver) Discover() error {
	return fmt.Errorf("driver's Discover function isn't implemented")
}

func (d *LoraAgentDriver) ValidateDevice(device models.Device) error {
	d.lc.Debug("Driver's ValidateDevice function isn't implemented")
	return nil
}
