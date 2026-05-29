package database

import (
	"fmt"
	"log"

	"github.com/mhsanaei/3x-ui/v3/database/model"
	"gorm.io/gorm"
)

// MigrateDeviceLimit creates the inbound_client_devices table
func MigrateDeviceLimit(db *gorm.DB) error {
	log.Println("[DB] Migrating device limit table...")

	type InboundClientDevices struct {
		Id          int    `gorm:"primaryKey;autoIncrement"`
		ClientEmail string `gorm:"unique"`
		Devices     string `gorm:"type:text"`
		CreatedAt   int64
		UpdatedAt   int64
	}

	if !db.Migrator().HasTable(&InboundClientDevices{}) {
		if err := db.Migrator().CreateTable(&InboundClientDevices{}); err != nil {
			return fmt.Errorf("failed to create inbound_client_devices table: %w", err)
		}
		log.Println("[DB] Created inbound_client_devices table")
	}

	return nil
}
