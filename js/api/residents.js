/**
 * Resident API endpoints
 */

import apiClient from './client.js';

const residentAPI = {
  /**
   * Get all residents
   */
  async getResidents() {
    return apiClient.get('/residents');
  },

  /**
   * Get resident by ID
   */
  async getResidentByID(id) {
    return apiClient.get(`/residents/${id}`);
  },

  /**
   * Create resident (Admin/Manager/Staff only)
   */
  async createResident(data) {
    return apiClient.post('/residents', data);
  },

  /**
   * Update resident profile
   */
  async updateResidentProfile(id, data) {
    return apiClient.put(`/residents/${id}`, data);
  },

  // ========== Dependent Management ==========

  /**
   * Add dependent to resident
   */
  async addDependent(residentID, data) {
    return apiClient.post(`/residents/${residentID}/dependents`, data);
  },

  /**
   * Remove dependent
   */
  async removeDependent(dependentID) {
    return apiClient.delete(`/dependents/${dependentID}`);
  },

  // ========== Vehicle Management ==========

  /**
   * Register vehicle for resident
   */
  async registerVehicle(residentID, data) {
    return apiClient.post(`/residents/${residentID}/vehicles`, data);
  },

  /**
   * Get resident vehicles
   */
  async getResidentVehicles(residentID) {
    return apiClient.get(`/residents/${residentID}/vehicles`);
  },

  /**
   * Unregister vehicle
   */
  async unregisterVehicle(vehicleID) {
    return apiClient.delete(`/vehicles/${vehicleID}`);
  },
};

export default residentAPI;
