/**
 * Building, Floor, Apartment API endpoints
 */

import apiClient from './client.js';

const buildingAPI = {
  // ========== Building Management ==========

  /**
   * Get all buildings
   */
  async getBuildings() {
    return apiClient.get('/buildings');
  },

  /**
   * Get building by ID
   */
  async getBuildingByID(id) {
    return apiClient.get(`/buildings/${id}`);
  },

  /**
   * Create building (Admin/Manager only)
   */
  async createBuilding(data) {
    return apiClient.post('/buildings', data);
  },

  /**
   * Update building (Admin/Manager only)
   */
  async updateBuilding(id, data) {
    return apiClient.put(`/buildings/${id}`, data);
  },

  /**
   * Delete building (Admin only)
   */
  async deleteBuilding(id) {
    return apiClient.delete(`/buildings/${id}`);
  },

  // ========== Floor Management ==========

  /**
   * Get all floors for a building
   */
  async getFloorsByBuilding(buildingID) {
    return apiClient.get(`/buildings/${buildingID}/floors`);
  },

  /**
   * Create floor (Admin/Manager only)
   */
  async createFloor(data) {
    return apiClient.post('/floors', data);
  },

  // ========== Apartment Management ==========

  /**
   * Get all apartments
   */
  async getApartments() {
    return apiClient.get('/apartments');
  },

  /**
   * Get apartment by ID
   */
  async getApartmentByID(id) {
    return apiClient.get(`/apartments/${id}`);
  },

  /**
   * Create apartment (Admin/Manager only)
   */
  async createApartment(data) {
    return apiClient.post('/apartments', data);
  },

  /**
   * Update apartment status (Admin/Manager only)
   */
  async updateApartmentStatus(id, data) {
    return apiClient.put(`/apartments/${id}/status`, data);
  },
};

export default buildingAPI;
