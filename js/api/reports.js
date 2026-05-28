/**
 * Reports and Dashboard API endpoints
 */

import apiClient from './client.js';

const reportsAPI = {
  /**
   * Get dashboard KPIs
   */
  async getDashboardKPIs() {
    return apiClient.get('/dashboard/kpis');
  },

  /**
   * Get audit logs (Admin only)
   */
  async getAuditLogs() {
    return apiClient.get('/audit-logs');
  },

  /**
   * Get financial report
   */
  async getFinancialReport(params) {
    return apiClient.post('/reports/financial', params);
  },
};

export default reportsAPI;
