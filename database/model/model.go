package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// InboundClientDevices stores device information associated with inbound clients for device-based access control
type InboundClientDevices struct {
	Id          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ClientEmail string    `json:"clientEmail" form:"clientEmail" gorm:"unique"`
	Devices     string    `json:"devices" form:"devices" gorm:"type:text"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// ClientDevice represents a single device accessing the inbound
type ClientDevice struct {
	DeviceID   string `json:"deviceId"`   // Unique device identifier
	DeviceName string `json:"deviceName"` // Device name/model
	LastSeen   int64  `json:"lastSeen"`   // Last connection timestamp
}

// TableName specifies the table name for InboundClientDevices
func (InboundClientDevices) TableName() string {
	return "inbound_client_devices"
}

// MarshalJSON emits the Devices column as a real JSON array instead of an escaped string
func (ic InboundClientDevices) MarshalJSON() ([]byte, error) {
	type alias InboundClientDevices
	return json.Marshal(struct {
		*alias
		Devices json.RawMessage `json:"devices"`
	}{
		alias:   (*alias)(&ic),
		Devices: jsonStringFieldToRaw(ic.Devices),
	})
}

// UnmarshalJSON accepts devices as either a JSON array or a JSON-encoded string
func (ic *InboundClientDevices) UnmarshalJSON(data []byte) error {
	type alias InboundClientDevices
	aux := struct {
		*alias
		Devices json.RawMessage `json:"devices"`
	}{
		alias: (*alias)(ic),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ic.Devices = jsonStringFieldFromRaw(aux.Devices)
	return nil
}

// jsonStringFieldToRaw converts a JSON-text string to a raw JSON message
func jsonStringFieldToRaw(s string) json.RawMessage {
	if s == "" {
		return nil
	}
	var parsed any
	if err := json.Unmarshal([]byte(s), &parsed); err != nil {
		return nil
	}
	raw, _ := json.Marshal(parsed)
	return raw
}

// jsonStringFieldFromRaw converts a raw JSON message to a JSON-text string
func jsonStringFieldFromRaw(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "[]"
	}
	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "[]"
	}
	b, _ := json.Marshal(parsed)
	return string(b)
}
