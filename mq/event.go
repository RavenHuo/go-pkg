package mq

import "time"

type EventInfo struct {
	MsgID     string      `json:"msg_id,omitempty"` // 全局唯一消息ID，去重使用
	Type      EventType   `json:"type,omitempty"`   // 事件类型ID
	Timestamp *time.Time  `json:"ts,omitempty"`     // 事件发生时间戳
	Info      interface{} `json:"info,omitempty"`   // 事件详情
	Seq       int64
}
