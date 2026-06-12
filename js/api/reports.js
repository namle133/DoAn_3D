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
   * Get financial report (UC 2.7.39)
   * @param {{ start_date: string, end_date: string, status?: string }} params — dates as YYYY-MM-DD
   */
  async getFinancialReport(params) {
    return apiClient.post('/reports/financial', params);
  },
};

export default reportsAPI;
