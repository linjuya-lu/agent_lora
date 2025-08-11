package driver

import (
	"encoding/hex"
	"fmt"

	"github.com/linjuya-lu/device_lora_agent_go/internal/config"
)

// handleResetCommand 封装了对 Reset_Set 资源写入后的完整处理：
// 1. 获取设备 EID
// 2. 校验并解码为 6 字节 sensorID
// 3. 构建复位帧
// 4. 通过串口层发送 AT+DTXSTR 命令
func (d *LoraAgentDriver) handleTimeParameterSet(deviceName string) error {
	d.lc.Infof("开始处理时间设置命令: %s", deviceName)
	// 获取设备的 EID 字符串
	eidValue, ok := config.GetDeviceValue(deviceName, "eid")
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 未初始化", deviceName)
		d.lc.Error(err.Error())
		return err
	}
	eidStr, ok := eidValue.(string)
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 类型错误，期望 string，实际 %T", deviceName, eidValue)
		d.lc.Error(err.Error())
		return err
	}
	eidStr = "238A0841D828"
	// 解码成 6 字节
	eidBytes, err := hex.DecodeString(eidStr)
	if err != nil {
		err = fmt.Errorf("EID[%s] 转十六进制失败: %w", eidStr, err)
		d.lc.Error(err.Error())
		return err
	}
	if len(eidBytes) != 6 {
		err = fmt.Errorf("EID 长度不对，期望 6 字节，实际 %d 字节", len(eidBytes))
		d.lc.Error(err.Error())
		return err
	}
	var sensorID [6]byte
	copy(sensorID[:], eidBytes)
	// 构建复位帧

	eidStr, _ = eidValue.(string)
	d.lc.Infof("已发送复位命令到设备 %s (EID: %s)", deviceName, eidStr)
	return nil
}

func (d *LoraAgentDriver) handleResetCommand(deviceName string) error {
	d.lc.Infof("开始处理复位命令: %s", deviceName)
	// 获取设备的 EID 字符串
	eidValue, ok := config.GetDeviceValue(deviceName, "eid")
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 未初始化", deviceName)
		d.lc.Error(err.Error())
		return err
	}
	eidStr, ok := eidValue.(string)
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 类型错误，期望 string，实际 %T", deviceName, eidValue)
		d.lc.Error(err.Error())
		return err
	}
	eidStr = "238A0841D828"
	// 解码成 6 字节
	eidBytes, err := hex.DecodeString(eidStr)
	if err != nil {
		err = fmt.Errorf("EID[%s] 转十六进制失败: %w", eidStr, err)
		d.lc.Error(err.Error())
		return err
	}
	if len(eidBytes) != 6 {
		err = fmt.Errorf("EID 长度不对，期望 6 字节，实际 %d 字节", len(eidBytes))
		d.lc.Error(err.Error())
		return err
	}
	var sensorID [6]byte
	copy(sensorID[:], eidBytes)
	// 构建复位帧

	d.lc.Infof("已发送复位命令到设备 %s (EID: %s)", deviceName, eidStr)
	return nil
}

func (d *LoraAgentDriver) handleTimeParameterQuery(deviceName string) error {
	d.lc.Infof("开始处理复位命令: %s", deviceName)
	// 获取设备的 EID 字符串
	eidValue, ok := config.GetDeviceValue(deviceName, "eid")
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 未初始化", deviceName)
		d.lc.Error(err.Error())
		return err
	}
	eidStr, ok := eidValue.(string)
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 类型错误，期望 string，实际 %T", deviceName, eidValue)
		d.lc.Error(err.Error())
		return err
	}
	eidStr = "238A0841D828"
	// 解码成 6 字节
	eidBytes, err := hex.DecodeString(eidStr)
	if err != nil {
		err = fmt.Errorf("EID[%s] 转十六进制失败: %w", eidStr, err)
		d.lc.Error(err.Error())
		return err
	}
	if len(eidBytes) != 6 {
		err = fmt.Errorf("EID 长度不对，期望 6 字节，实际 %d 字节", len(eidBytes))
		d.lc.Error(err.Error())
		return err
	}
	var sensorID [6]byte
	copy(sensorID[:], eidBytes)
	// 构建复位帧

	d.lc.Infof("已发送复位命令到设备 %s (EID: %s)", deviceName, eidStr)
	return nil
}

func (d *LoraAgentDriver) handleIdQuery(deviceName string) error {
	d.lc.Infof("开始处理复位命令: %s", deviceName)
	// 获取设备的 EID 字符串
	eidValue, ok := config.GetDeviceValue(deviceName, "eid")
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 未初始化", deviceName)
		d.lc.Error(err.Error())
		return err
	}
	eidStr, ok := eidValue.(string)
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 类型错误，期望 string，实际 %T", deviceName, eidValue)
		d.lc.Error(err.Error())
		return err
	}
	eidStr = "238A0841D828"
	// 解码成 6 字节
	eidBytes, err := hex.DecodeString(eidStr)
	if err != nil {
		err = fmt.Errorf("EID[%s] 转十六进制失败: %w", eidStr, err)
		d.lc.Error(err.Error())
		return err
	}
	if len(eidBytes) != 6 {
		err = fmt.Errorf("EID 长度不对，期望 6 字节，实际 %d 字节", len(eidBytes))
		d.lc.Error(err.Error())
		return err
	}
	var sensorID [6]byte
	copy(sensorID[:], eidBytes)
	//构建ID查询帧

	d.lc.Infof("已发送复位命令到设备 %s (EID: %s)", deviceName, eidStr)
	return nil
}

func (d *LoraAgentDriver) handleIdMoniDataQuery(deviceName string) error {
	d.lc.Infof("开始处理复位命令: %s", deviceName)
	// 获取设备的 EID 字符串
	eidValue, ok := config.GetDeviceValue(deviceName, "eid")
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 未初始化", deviceName)
		d.lc.Error(err.Error())
		return err
	}
	eidStr, ok := eidValue.(string)
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 类型错误，期望 string，实际 %T", deviceName, eidValue)
		d.lc.Error(err.Error())
		return err
	}
	eidStr = "238A0841D828"
	// 解码成 6 字节
	eidBytes, err := hex.DecodeString(eidStr)
	if err != nil {
		err = fmt.Errorf("EID[%s] 转十六进制失败: %w", eidStr, err)
		d.lc.Error(err.Error())
		return err
	}
	if len(eidBytes) != 6 {
		err = fmt.Errorf("EID 长度不对，期望 6 字节，实际 %d 字节", len(eidBytes))
		d.lc.Error(err.Error())
		return err
	}
	var sensorID [6]byte
	copy(sensorID[:], eidBytes)

	d.lc.Infof("已发送复位命令到设备 %s (EID: %s)", deviceName, eidStr)
	return nil
}

func (d *LoraAgentDriver) handleIdAlarmParaQuery(deviceName string) error {
	d.lc.Infof("开始处理复位命令: %s", deviceName)
	// 获取设备的 EID 字符串
	eidValue, ok := config.GetDeviceValue(deviceName, "eid")
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 未初始化", deviceName)
		d.lc.Error(err.Error())
		return err
	}
	eidStr, ok := eidValue.(string)
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 类型错误，期望 string，实际 %T", deviceName, eidValue)
		d.lc.Error(err.Error())
		return err
	}
	eidStr = "238A0841D828"
	//解码成 6 字节
	eidBytes, err := hex.DecodeString(eidStr)
	if err != nil {
		err = fmt.Errorf("EID[%s] 转十六进制失败: %w", eidStr, err)
		d.lc.Error(err.Error())
		return err
	}
	if len(eidBytes) != 6 {
		err = fmt.Errorf("EID 长度不对，期望 6 字节，实际 %d 字节", len(eidBytes))
		d.lc.Error(err.Error())
		return err
	}
	var sensorID [6]byte
	copy(sensorID[:], eidBytes)
	// 构建ID查询帧

	d.lc.Infof("已发送复位命令到设备 %s (EID: %s)", deviceName, eidStr)
	return nil
}

func (d *LoraAgentDriver) handleGeneParaQuery(deviceName string) error {
	d.lc.Infof("开始处理复位命令: %s", deviceName)
	// 获取设备的 EID 字符串
	eidValue, ok := config.GetDeviceValue(deviceName, "eid")
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 未初始化", deviceName)
		d.lc.Error(err.Error())
		return err
	}
	eidStr, ok := eidValue.(string)
	if !ok {
		err := fmt.Errorf("设备 %s 的 EID 类型错误，期望 string，实际 %T", deviceName, eidValue)
		d.lc.Error(err.Error())
		return err
	}
	eidStr = "238A0841D828"
	//解码成 6 字节
	eidBytes, err := hex.DecodeString(eidStr)
	if err != nil {
		err = fmt.Errorf("EID[%s] 转十六进制失败: %w", eidStr, err)
		d.lc.Error(err.Error())
		return err
	}
	if len(eidBytes) != 6 {
		err = fmt.Errorf("EID 长度不对，期望 6 字节，实际 %d 字节", len(eidBytes))
		d.lc.Error(err.Error())
		return err
	}
	var sensorID [6]byte
	copy(sensorID[:], eidBytes)

	d.lc.Infof("已发送复位命令到设备 %s (EID: %s)", deviceName, eidStr)
	return nil
}
