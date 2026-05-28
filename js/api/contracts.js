/**
 * Contract API endpoints
 */

import apiClient from './client.js';

const contractAPI = {
  /**
   * Get all contracts
   */
  async getContracts() {
    return apiClient.get('/contracts');
  },

  /**
   * Get contract by ID
   */
  async getContractByID(id) {
    return apiClient.get(`/contracts/${id}`);
  },

  /**
   * Create contract (Admin/Manager/Staff only)
   */
  async createContract(data) {
    return apiClient.post('/contracts', data);
  },

  /**
   * Update contract
   */
  async updateContract(id, data) {
    return apiClient.put(`/contracts/${id}`, data);
  },

  /**
   * Approve contract (Admin/Manager only)
   */
  async approveContract(id, data = {}) {
    return apiClient.post(`/contracts/${id}/approve`, data);
  },

  /**
   * Extend contract (Admin/Manager only)
   */
  async extendContract(id, data) {
    return apiClient.post(`/contracts/${id}/extend`, data);
  },

  /**
   * Terminate contract (Admin/Manager only)
   */
  async terminateContract(id, data = {}) {
    return apiClient.post(`/contracts/${id}/terminate`, data);
  },
};

export default contractAPI;
