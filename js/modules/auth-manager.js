/**
 * Auth Manager - Handles authentication and session management
 */

import authAPI from '../api/auth.js';
import apiClient from '../api/client.js';
import { ROLES, STORAGE_KEYS } from '../utils/constants.js';
import { showToast } from '../utils/helpers.js';

const authManager = {
  /**
   * Login user
   */
  async login(email, password) {
    try {
      const response = await authAPI.login(email, password);
      showToast('Đăng nhập thành công!', 'success');
      return response;
    } catch (error) {
      showToast(`Đăng nhập thất bại: ${error.message}`, 'error');
      throw error;
    }
  },

  /**
   * Register user
   */
  async register(userData) {
    try {
      const response = await authAPI.register(userData);
      showToast('Đăng ký thành công!', 'success');
      return response;
    } catch (error) {
      showToast(`Đăng ký thất bại: ${error.message}`, 'error');
      throw error;
    }
  },

  /**
   * Logout user
   */
  logout() {
    authAPI.logout();
  },

  /**
   * Check if user is authenticated
   */
  isAuthenticated() {
    return apiClient.isAuthenticated();
  },

  /**
   * Get current user
   */
  getCurrentUser() {
    return apiClient.getCurrentUser();
  },

  /**
   * Get user role
   */
  getUserRole() {
    return apiClient.getUserRole();
  },

  /**
   * Check if user has role(s)
   */
  hasRole(roles) {
    return apiClient.hasRole(roles);
  },

  /**
   * Check if user is admin
   */
  isAdmin() {
    return this.hasRole(ROLES.ADMIN);
  },

  /**
   * Check if user is manager
   */
  isManager() {
    return this.hasRole(ROLES.MANAGER);
  },

  /**
   * Check if user is staff
   */
  isStaff() {
    return this.hasRole(ROLES.STAFF);
  },

  /**
   * Check if user is resident
   */
  isResident() {
    return this.hasRole(ROLES.RESIDENT);
  },

  /**
   * Guard route - redirect to login if not authenticated
   */
  guardRoute(requiredRoles = null) {
    if (!this.isAuthenticated()) {
      window.location.href = '/pages/login.html?redirect=' + window.location.pathname;
      return false;
    }

    if (requiredRoles && !this.hasRole(requiredRoles)) {
      showToast('Bạn không có quyền truy cập trang này', 'error');
      this.redirectToDashboard();
      return false;
    }

    return true;
  },

  /**
   * Redirect to appropriate dashboard based on role
   */
  redirectToDashboard() {
    const role = this.getUserRole();
    if (role === ROLES.ADMIN || role === ROLES.MANAGER || role === ROLES.STAFF) {
      window.location.href = '/pages/dashboard.html';
    } else if (role === ROLES.RESIDENT) {
      window.location.href = '/pages/resident-dashboard.html';
    }
  },
};

export default authManager;
