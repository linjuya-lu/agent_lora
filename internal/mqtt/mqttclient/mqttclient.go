package mqttclient

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/linjuya-lu/device_lora_agent_go/internal/config"
)

// 根据broker URL 和 clientID 创建并连接 MQTT 客户端
func NewClient(brokerURL, clientID string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		// 设置自动重连，心跳，超时等
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetKeepAlive(60 * time.Second).
		SetPingTimeout(10 * time.Second)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if ok := token.WaitTimeout(10 * time.Second); !ok {
		return nil, fmt.Errorf("MQTT 连接超时")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("MQTT 连接失败: %w", err)
	}
	return client, nil
}

// EdgexMessage 是 EdgeX MessageBus 的通用消息格式
type EdgexMessage struct {
	ApiVersion    string      `json:"apiVersion"`
	ReceivedTopic string      `json:"receivedTopic,omitempty"`
	CorrelationID string      `json:"correlationID"`
	RequestID     string      `json:"requestID"`
	ErrorCode     int         `json:"errorCode"`
	Payload       interface{} `json:"payload,omitempty"`
	ContentType   string      `json:"contentType"`
}

// SinkPayload 表示 edgex/service/data/sink 消息中的 payload 部分
type SinkPayload struct {
	Type      string `json:"Type"`      // sink: 网关自身参数，sensor: 传感器数据
	Eid       string `json:"Eid"`       // 模块 EID，作为设备标识
	Timestamp uint64 `json:"Timestamp"` // 世纪秒时间戳（网关本地时间）
	Datalen   int    `json:"Datalen"`   // 原始数据长度
	Data      string `json:"Data"`      // 原始数据（16 进制 HEX 字符串）
}

// SinkCommandPayload 表示 edgex/core/command/sink 消息中的 payload 部分
type SinkCommandPayload struct {
	Eid       string `json:"Eid"`       // 网关模块 EID
	Timestamp uint64 `json:"Timestamp"` // 世纪秒时间戳（边代本地时间）
	Cmd       string `json:"Cmd"`       // 网关监测数据查询命令，如 "getdata"
}

// 内部用：把 Payload 换成 RawMessage，避免二次编码
type edgexEnvelope struct {
	ApiVersion    string          `json:"apiVersion"`
	ReceivedTopic string          `json:"receivedTopic,omitempty"`
	CorrelationID string          `json:"correlationID"`
	RequestID     string          `json:"requestID"`
	ErrorCode     int             `json:"errorCode"`
	Payload       json.RawMessage `json:"payload,omitempty"`
	ContentType   string          `json:"contentType"`
}

// SubscribeSinkData 订阅指定 topic（如 "edgex/service/data/sink"），解析后把 Data 放入 SinkHexDataCh
func SubscribeSinkData(cli mqtt.Client, topic string, qos byte) error {
	log.Printf("🔔 Subscribing sink data: %s", topic)
	tok := cli.Subscribe(topic, qos, sinkMsgHandler)
	tok.Wait()
	return tok.Error()
}

func sinkMsgHandler(_ mqtt.Client, msg mqtt.Message) {
	// 1) 解外层
	var env edgexEnvelope
	if err := json.Unmarshal(msg.Payload(), &env); err != nil {
		log.Printf("❌ 解析 EdgexMessage 失败: %v; payload=%s", err, string(msg.Payload()))
		return
	}
	if len(env.Payload) == 0 {
		log.Printf("❌ payload 为空")
		return
	}

	// 2) 解 payload → SinkPayload
	var sp SinkPayload
	if err := json.Unmarshal(env.Payload, &sp); err != nil {
		log.Printf("❌ 解析 SinkPayload 失败: %v; payload=%s", err, string(env.Payload))
		return
	}

	// 3) 可选校验：长度与 HEX 是否匹配
	if len(sp.Data)%2 != 0 {
		log.Printf("⚠ HEX 长度不是偶数: len=%d", len(sp.Data))
	}
	if sp.Datalen >= 0 && sp.Datalen != len(sp.Data)/2 {
		log.Printf("⚠ Datalen(%d) 与 HEX 实际字节数(%d) 不一致", sp.Datalen, len(sp.Data)/2)
	}

	// 4) 投递到全局通道（HEX 文本）
	select {
	case config.SinkHexDataCh <- sp.Data:
		// OK
	default:
		log.Printf("⚠ SinkHexDataCh 已满，丢弃一帧")
	}

	// 如果你要直接得到原始字节，可同时解码后投递：
	if raw, err := hex.DecodeString(sp.Data); err == nil {
		// select { case SinkRawDataCh <- raw: default: }
		_ = raw
	}
}

//--------------------------------用法------------------------------
// // 业务侧异步读取：
// go func() {
// 	for hexStr := range mq.SinkHexDataCh {
// 		// TODO: 你的处理
// 		log.Printf("收到 HEX: %s", hexStr)
// 	}
// }()

// PublishSinkCommand 发布 edgex/core/command/sink 消息
// 参数：
//   - client: MQTT 客户端
//   - topic: 要发布的 MQTT 主题，比如 "edgex/core/command/sink"
//   - eid:   网关模块 EID
//   - cmd:   命令字符串，比如 "getdata"
func PublishSinkCommand(client mqtt.Client, topic, eid, cmd string) error {
	// 1. 内层 payload
	payload := SinkCommandPayload{
		Eid:       eid,
		Timestamp: uint64(time.Now().Unix()), // 世纪秒，可以替换为本地特定时间源
		Cmd:       cmd,
	}

	// 2. 外层通用消息
	msg := EdgexMessage{
		ApiVersion:    "v3",
		ReceivedTopic: "",
		CorrelationID: uuid.NewString(),
		RequestID:     uuid.NewString(),
		ErrorCode:     0,
		Payload:       payload,
		ContentType:   "application/json",
	}

	// 3. 序列化 JSON
	body, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("❌ JSON Marshal error: %v\n", err)
		return err
	}

	// 发布前打印主题和消息体
	fmt.Printf("⮉ Publishing MQTT topic=%s, message=%s\n", topic, string(body))

	// 4. 发布并等待完成
	tok := client.Publish(topic, 0, false, body)
	tok.Wait()
	if err := tok.Error(); err != nil {
		fmt.Printf("❌ Publish error: %v\n", err)
	} else {
		fmt.Printf("✅ Publish succeeded: topic=%s\n", topic)
	}
	return tok.Error()
}

// 使用例子
// err := PublishSinkCommand(mqttClient, "edgex/core/command/sink", "238A08411011", "getdata")
// if err != nil {
// 	fmt.Println("发布失败:", err)
// }
