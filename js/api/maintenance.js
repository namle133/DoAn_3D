/**
 * Maintenance API endpoints
 */

import apiClient from './client.js';

const maintenanceAPI = {
  /**
   * Get all maintenance requests
   */
  async getMaintenanceRequests() {
    return apiClient.get('/maintenance-requests');
  },

  /**
   * Get maintenance request by ID
   */
  async getMaintenanceRequestByID(id) {
    return apiClient.get(`/maintenance-requests/${id}`);
  },

  /**
   * Create maintenance request
   */
  async createMaintenanceRequest(data) {
    return apiClient.post('/maintenance-requests', data);
  },

  /**
   * Assign maintenance request (Admin/Manager only)
   */
  async assignMaintenanceRequest(id, data) {
    return apiClient.post(`/maintenance-requests/${id}/assign`, data);
  },

  /**
   * Update maintenance request status
   */
  async updateMaintenanceStatus(id, status) {
    return apiClient.put(`/maintenance-requests/${id}/status`, { status });
  },

  /**
   * Complete maintenance request
   */
  async completeMaintenanceRequest(id, data = {}) {
    return apiClient.post(`/maintenance-requests/${id}/complete`, data);
  },

  // ========== Equipment Management ==========

  /**
   * Register equipment (Admin/Manager only)
   */
  async registerEquipment(data) {
    return apiClient.post('/equipment', data);
  },

  /**
   * Get equipment by building
   */
  async getEquipmentByBuilding(buildingID) {
    return apiClient.get(`/buildings/${buildingID}/equipment`);
  },

  /**
   * Get equipment due for maintenance
   */
  async getEquipmentDueForMaintenance() {
    return apiClient.get('/equipment/due-for-maintenance');
  },

  /**
   * Update equipment maintenance record (Staff/Manager only)
   */
  async updateEquipmentMaintenanceRecord(equipmentID, data) {
    return apiClient.put(`/equipment/${equipmentID}/maintenance`, data);
  },
};

export default maintenanceAPI;
