package config

// 全局写入通道，传入待发送的帧数据
var WriteChan = make(chan []byte, 100)

// ===== 全局通道：后续业务从这里读 HEX 字符串（或你也可以改成 []byte） =====
var SinkHexDataCh = make(chan string, 1024) // HEX 文本
