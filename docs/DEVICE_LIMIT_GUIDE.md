# Device Limit Feature Guide

## Overview

The Device Limit feature replaces IP-based access control with device-based access control. This allows clients to connect from multiple devices (up to a configured limit) while maintaining security.

## Architecture

### Database Schema

```sql
CREATE TABLE inbound_client_devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_email TEXT UNIQUE NOT NULL,
    devices TEXT NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Device Information Structure

Each device is stored as a JSON object:

```json
{
  "deviceId": "unique-device-identifier",
  "deviceName": "Device Model/Name",
  "lastSeen": 1234567890
}
```

## Components

### 1. DeviceLimitService (`web/service/device_limit_service.go`)

Provides core device limit functionality:

- **CheckDeviceLimit(email, limit, newDeviceID)** - Validates if a new device exceeds the limit
- **RecordDeviceAccess(email, deviceID, deviceName)** - Records device access
- **GetClientDevices(email)** - Retrieves all devices for a client
- **RemoveDevice(email, deviceID)** - Removes a specific device
- **ClearOldDevices(email)** - Removes devices not seen in 30 days

### 2. Database Model (`database/model/model.go`)

Defines the `InboundClientDevices` model with proper JSON marshaling.

### 3. API Endpoints (`web/controller/device_limit_api.go`)

#### GET `/api/client/devices/:email`
Returns all devices for a client.

```bash
curl http://localhost:2053/api/client/devices/user@example.com
```

Response:
```json
{
  "msg": "success",
  "devices": [
    {
      "deviceId": "uuid-1",
      "deviceName": "Windows PC",
      "lastSeen": 1234567890
    }
  ]
}
```

#### DELETE `/api/client/devices/:email/:deviceId`
Removes a specific device.

```bash
curl -X DELETE http://localhost:2053/api/client/devices/user@example.com/uuid-1
```

#### DELETE `/api/client/devices/:email`
Clears all devices for a client.

```bash
curl -X DELETE http://localhost:2053/api/client/devices/user@example.com
```

### 4. Frontend Utilities (`web/html/src/utils/deviceUtils.js`)

- **getDeviceID()** - Returns unique device ID
- **getDeviceName()** - Returns device name
- **getClientDeviceInfo()** - Returns complete device info

## Usage Example

### In Client Code

```go
deviceSvc := &service.DeviceLimitService{}

// Check if device is allowed
allowed, err := deviceSvc.CheckDeviceLimit(
    "user@example.com",
    5, // limit
    "device-uuid-123", // new device ID
)

if !allowed {
    // Device limit exceeded
    return errors.New("device limit exceeded")
}

// Record device access
err = deviceSvc.RecordDeviceAccess(
    "user@example.com",
    "device-uuid-123",
    "Windows PC",
)
```

### In Frontend (Vue.js)

```vue
<script setup>
import { getClientDeviceInfo } from '@/utils/deviceUtils'

const deviceInfo = getClientDeviceInfo()
console.log(deviceInfo)
// Output: { deviceId: 'uuid-xxx', deviceName: 'Windows PC' }
</script>
```

## Configuration

### Client Model Field

In `database/model/client.go`, add:

```go
type Client struct {
    // ... other fields ...
    LimitDevice int `json:"limitDevice"` // Number of allowed devices (0 = unlimited)
}
```

### Stale Device Cutoff

Devices are considered stale after 30 days of inactivity. Modify in `device_limit_service.go`:

```go
const DeviceStaleCutoffDays = 30
```

## Migration from IP Limit

1. Run database migration to create `inbound_client_devices` table
2. Update client model to use `limitDevice` instead of `limitIp`
3. Update client validation logic
4. Update frontend UI components

## API Routes

Register these routes in your router setup:

```go
// Device management endpoints
router.GET("/api/client/devices/:email", a.GetClientDevices)
router.DELETE("/api/client/devices/:email/:deviceId", a.ClearClientDevice)
router.DELETE("/api/client/devices/:email", a.ClearAllClientDevices)
```

## Testing

### Test Device Limit Validation

```bash
# Check device limit
curl -X GET http://localhost:2053/api/client/devices/test@example.com

# Remove a device
curl -X DELETE http://localhost:2053/api/client/devices/test@example.com/device-id
```

## Security Considerations

1. **Device ID Uniqueness** - Ensure device IDs are truly unique per browser/device
2. **Timestamp Accuracy** - Keep server time synchronized
3. **Stale Device Cleanup** - Automatically remove old devices
4. **Rate Limiting** - Implement rate limiting on device registration

## Future Enhancements

- Device fingerprinting for more robust device identification
- Device trust levels (trusted/untrusted devices)
- Device-specific access policies
- Device activity logs
- Push notifications on new device login
