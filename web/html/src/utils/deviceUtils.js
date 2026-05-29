/**
 * Device utility functions for device-based access control
 */

/**
 * Generate or retrieve unique device ID from localStorage
 */
export function getDeviceID() {
  const DEVICE_ID_KEY = 'x_ui_device_id';
  let deviceID = localStorage.getItem(DEVICE_ID_KEY);
  
  if (!deviceID) {
    deviceID = generateUUID();
    localStorage.setItem(DEVICE_ID_KEY, deviceID);
  }
  
  return deviceID;
}

/**
 * Get device name based on user agent
 */
export function getDeviceName() {
  const ua = navigator.userAgent;
  
  if (ua.includes('Windows')) return 'Windows PC';
  if (ua.includes('Mac')) return 'Mac';
  if (ua.includes('iPhone')) return 'iPhone';
  if (ua.includes('iPad')) return 'iPad';
  if (ua.includes('Android')) return 'Android Device';
  if (ua.includes('Linux')) return 'Linux';
  
  return 'Unknown Device';
}

/**
 * Generate UUID v4
 */
function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

/**
 * Get client device info
 */
export function getClientDeviceInfo() {
  return {
    deviceId: getDeviceID(),
    deviceName: getDeviceName(),
  };
}
