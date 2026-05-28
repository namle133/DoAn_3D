/**
 * Constants for Vinhomes Property Management System
 */

// API Configuration
const API_CONFIG = {
  BASE_URL: 'http://localhost:8080/api/v1',
  TIMEOUT: 10000,
};

// User Roles
const ROLES = {
  ADMIN: 'admin',
  MANAGER: 'manager',
  STAFF: 'staff',
  RESIDENT: 'resident',
};

// Apartment Status
const APARTMENT_STATUS = {
  EMPTY: 'empty',
  RENTED: 'rented',
  MAINTENANCE: 'maintenance',
  RESERVED: 'reserved',
};

// Apartment Status Colors for GIS 3D Model (UC 2.7.8)
const APARTMENT_COLORS = {
  [APARTMENT_STATUS.EMPTY]: 0x4CAF50,       // Xanh lá — Trống
  [APARTMENT_STATUS.RENTED]: 0x2196F3,      // Xanh dương — Đang thuê
  [APARTMENT_STATUS.MAINTENANCE]: 0xF44336, // Đỏ — Bảo trì
  [APARTMENT_STATUS.RESERVED]: 0xFF9800,    // Cam — Đặt trước
};

// GIS display status for expiring contracts
const GIS_DISPLAY_STATUS = {
  EXPIRING: 'expiring', // Vàng — Sắp hết HĐ
};

// Contract Status
const CONTRACT_STATUS = {
  PENDING: 'pending',
  ACTIVE: 'active',
  EXTENDED: 'extended',
  TERMINATED: 'terminated',
  EXPIRED: 'expired',
};

// Invoice Status
const INVOICE_STATUS = {
  DRAFT: 'draft',
  ISSUED: 'issued',
  OVERDUE: 'overdue',
  PAID: 'paid',
  CANCELLED: 'cancelled',
};

// Maintenance Priority
const MAINTENANCE_PRIORITY = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
  URGENT: 'urgent',
};

// Maintenance Status
const MAINTENANCE_STATUS = {
  NEW: 'new',
  ASSIGNED: 'assigned',
  IN_PROGRESS: 'in_progress',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
};

// Notification Types
const NOTIFICATION_TYPES = {
  CONTRACT_EXPIRING: 'contract_expiring',
  INVOICE_OVERDUE: 'invoice_overdue',
  MAINTENANCE_UPDATE: 'maintenance_update',
  PAYMENT_RECEIVED: 'payment_received',
  CONTRACT_APPROVED: 'contract_approved',
  CONTRACT_REJECTED: 'contract_rejected',
  GENERAL: 'general',
};

// Local Storage Keys
const STORAGE_KEYS = {
  AUTH_TOKEN: 'vinhomes_auth_token',
  REFRESH_TOKEN: 'vinhomes_refresh_token',
  CURRENT_USER: 'vinhomes_current_user',
  USER_ROLE: 'vinhomes_user_role',
};

// Session Storage Keys
const SESSION_KEYS = {
  CURRENT_BUILDING: 'current_building',
  CURRENT_FLOOR: 'current_floor',
  CURRENT_APARTMENT: 'current_apartment',
};

export {
  API_CONFIG,
  ROLES,
  APARTMENT_STATUS,
  APARTMENT_COLORS,
  CONTRACT_STATUS,
  INVOICE_STATUS,
  MAINTENANCE_PRIORITY,
  MAINTENANCE_STATUS,
  NOTIFICATION_TYPES,
  STORAGE_KEYS,
  SESSION_KEYS,
};
