/**
 * Invoice API endpoints
 */

import apiClient from './client.js';

const invoiceAPI = {
  /**
   * Get all invoices
   */
  async getInvoices() {
    return apiClient.get('/invoices');
  },

  /**
   * Get invoice by ID
   */
  async getInvoiceByID(id) {
    return apiClient.get(`/invoices/${id}`);
  },

  /**
   * Create invoice (Admin/Manager only)
   */
  async createInvoice(data) {
    return apiClient.post('/invoices', data);
  },

  /**
   * Update invoice (Admin/Manager only)
   */
  async updateInvoice(id, data) {
    return apiClient.put(`/invoices/${id}`, data);
  },

  /**
   * Record payment (Staff/Manager only)
   */
  async recordPayment(invoiceID, data) {
    return apiClient.post(`/invoices/${invoiceID}/payments`, data);
  },

  /**
   * Get overdue invoices
   */
  async getOverdueInvoices() {
    return apiClient.get('/invoices/overdue');
  },

  /**
   * Get financial report
   */
  async getFinancialReport(data) {
    return apiClient.post('/reports/financial', data);
  },
};

export default invoiceAPI;
