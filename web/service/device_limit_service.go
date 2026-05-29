package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mhsanaei/3x-ui/v3/database"
	"github.com/mhsanaei/3x-ui/v3/database/model"
	"github.com/mhsanaei/3x-ui/v3/logger"
	"gorm.io/gorm"
)

type DeviceLimitService struct{}

type ClientDevice struct {
	DeviceID   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
	LastSeen   int64  `json:"lastSeen"`
}

const (
	// Device stale cutoff: 30 days
	DeviceStaleCutoffDays = 30
)

// CheckDeviceLimit checks if a client has exceeded device limit
// Returns (allowed, error)
func (s *DeviceLimitService) CheckDeviceLimit(clientEmail string, limit int, newDeviceID string) (bool, error) {
	if limit <= 0 {
		return true, nil // No limit
	}

	db := database.GetDB()
	var record model.InboundClientDevices

	err := db.Where("client_email = ?", clientEmail).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil // No devices recorded yet
		}
		return false, err
	}

	var devices []ClientDevice
	err = json.Unmarshal([]byte(record.Devices), &devices)
	if err != nil {
		return false, err
	}

	// Count unique devices from last 30 days
	now := time.Now().Unix()
	staleCutoff := now - (int64(DeviceStaleCutoffDays) * 24 * 60 * 60)

	activeDevices := make(map[string]struct{})
	for _, dev := range devices {
		if dev.LastSeen > staleCutoff {
			activeDevices[dev.DeviceID] = struct{}{}
		}
	}

	// Add new device if not seen before
	if _, exists := activeDevices[newDeviceID]; !exists {
		activeDevices[newDeviceID] = struct{}{}
	}

	// Check if exceeded limit
	if len(activeDevices) > limit {
		return false, nil // Limit exceeded
	}

	return true, nil
}

// RecordDeviceAccess records or updates device information
func (s *DeviceLimitService) RecordDeviceAccess(clientEmail, deviceID, deviceName string) error {
	db := database.GetDB()
	now := time.Now().Unix()

	var record model.InboundClientDevices
	err := db.Where("client_email = ?", clientEmail).First(&record).Error

	var devices []ClientDevice
	if err == nil {
		json.Unmarshal([]byte(record.Devices), &devices)
	}

	// Update existing device or add new
	found := false
	for i := range devices {
		if devices[i].DeviceID == deviceID {
			devices[i].LastSeen = now
			if deviceName != "" {
				devices[i].DeviceName = deviceName
			}
			found = true
			break
		}
	}

	if !found {
		devices = append(devices, ClientDevice{
			DeviceID:   deviceID,
			DeviceName: deviceName,
			LastSeen:   now,
		})
	}

	devicesJSON, err := json.Marshal(devices)
	if err != nil {
		return err
	}

	if err == nil && record.Id > 0 {
		// Update existing
		record.Devices = string(devicesJSON)
		return db.Save(&record).Error
	}

	// Create new
	return db.Create(&model.InboundClientDevices{
		ClientEmail: clientEmail,
		Devices:     string(devicesJSON),
	}).Error
}

// GetClientDevices returns all devices for a client
func (s *DeviceLimitService) GetClientDevices(clientEmail string) ([]ClientDevice, error) {
	db := database.GetDB()
	var record model.InboundClientDevices

	err := db.Where("client_email = ?", clientEmail).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []ClientDevice{}, nil
		}
		return nil, err
	}

	var devices []ClientDevice
	err = json.Unmarshal([]byte(record.Devices), &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// ClearOldDevices removes devices not seen beyond stale cutoff
func (s *DeviceLimitService) ClearOldDevices(clientEmail string) error {
	db := database.GetDB()
	now := time.Now().Unix()
	staleCutoff := now - (int64(DeviceStaleCutoffDays) * 24 * 60 * 60)

	var record model.InboundClientDevices
	err := db.Where("client_email = ?", clientEmail).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	var devices []ClientDevice
	json.Unmarshal([]byte(record.Devices), &devices)

	filtered := make([]ClientDevice, 0)
	for _, dev := range devices {
		if dev.LastSeen > staleCutoff {
			filtered = append(filtered, dev)
		}
	}

	if len(filtered) == 0 {
		return db.Delete(&record).Error
	}

	devicesJSON, _ := json.Marshal(filtered)
	record.Devices = string(devicesJSON)
	return db.Save(&record).Error
}

// RemoveDevice removes a specific device from client's device list
func (s *DeviceLimitService) RemoveDevice(clientEmail, deviceID string) error {
	db := database.GetDB()
	var record model.InboundClientDevices

	err := db.Where("client_email = ?", clientEmail).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	var devices []ClientDevice
	json.Unmarshal([]byte(record.Devices), &devices)

	filtered := make([]ClientDevice, 0)
	for _, dev := range devices {
		if dev.DeviceID != deviceID {
			filtered = append(filtered, dev)
		}
	}

	if len(filtered) == 0 {
		return db.Delete(&record).Error
	}

	devicesJSON, _ := json.Marshal(filtered)
	record.Devices = string(devicesJSON)
	return db.Save(&record).Error
}
