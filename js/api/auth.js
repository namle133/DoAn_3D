/**
 * Auth API endpoints
 */

import apiClient from './client.js';

const authAPI = {
  /**
   * Login
   */
  async login(email, password) {
    const response = await apiClient.post('/auth/login', { email, password });
    const data = response.data || response;

    if (data.token) {
      apiClient.setTokens(data.token, data.refresh_token);
      apiClient.setCurrentUser({
        id: data.user_id,
        username: data.username,
        email: data.email || email,
        name: data.name || data.username,
        role: data.role,
      });
      apiClient.setUserRole(data.role);
    }

    return response;
  },

  /**
   * Register
   */
  async register(userData) {
    const response = await apiClient.post('/auth/register', {
      username: userData.username,
      email: userData.email,
      password: userData.password,
      role: userData.role || 'resident',
      full_name: userData.full_name || '',
      id_card: userData.id_card || '',
      date_of_birth: userData.date_of_birth || '',
      gender: userData.gender || '',
      phone_number: userData.phone_number || '',
      permanent_address: userData.permanent_address || '',
      verify_code: userData.verify_code || '',
    });
    return response;
  },

  /**
   * Quên mật khẩu — UC 2.7.4
   */
  async forgotPassword(email) {
    return apiClient.post('/auth/forgot-password', { email });
  },

  /**
   * Đặt lại mật khẩu — UC 2.7.4
   */
  async resetPassword(email, token, newPassword) {
    return apiClient.post('/auth/reset-password', {
      email,
      token,
      new_password: newPassword,
    });
  },

  /**
   * Cập nhật hồ sơ
   */
  async updateProfile(data) {
    return apiClient.put('/profile', data);
  },

  /**
   * Logout
   */
  logout() {
    apiClient.clearTokens();
    window.location.href = '/pages/login.html';
  },

  /**
   * Get current user profile
   */
  async getProfile() {
    return apiClient.get('/profile');
  },

  /**
   * Check authentication status
   */
  isAuthenticated() {
    return apiClient.isAuthenticated();
  },
};

export default authAPI;
