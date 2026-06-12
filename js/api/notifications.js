/**
 * Notification API endpoints
 */

import apiClient from './client.js';

const notificationAPI = {
  /**
   * Get all user notifications
   */
  async getUserNotifications(scope) {
    const qs = scope === 'all' ? '?scope=all' : '';
    return apiClient.get('/notifications' + qs);
  },

  /**
   * Get unread notifications
   */
  async getUnreadNotifications() {
    return apiClient.get('/notifications/unread');
  },

  /**
   * Create notification (Admin/Manager only)
   */
  async createNotification(data) {
    return apiClient.post('/notifications', data);
  },

  /**
   * Mark notification as read
   */
  async markNotificationAsRead(id) {
    return apiClient.put(`/notifications/${id}/read`, {});
  },

  /**
   * Mark all notifications as read
   */
  async markAllNotificationsAsRead() {
    return apiClient.put('/notifications/read-all', {});
  },

  /**
   * Delete notification
   */
  async deleteNotification(id) {
    return apiClient.delete(`/notifications/${id}`);
  },
};

export default notificationAPI;
