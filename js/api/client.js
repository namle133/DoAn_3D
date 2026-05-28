/**
 * API Client - Base HTTP client with JWT token management
 */

import { API_CONFIG, STORAGE_KEYS } from '../utils/constants.js';

class APIClient {
  constructor() {
    this.baseURL = API_CONFIG.BASE_URL;
    this.token = this.getToken();
    this.refreshToken = this.getRefreshToken();
  }

  /**
   * Get stored JWT token
   */
  getToken() {
    return localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN);
  }

  /**
   * Get stored refresh token
   */
  getRefreshToken() {
    return localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
  }

  /**
   * Set tokens
   */
  setTokens(token, refreshToken) {
    localStorage.setItem(STORAGE_KEYS.AUTH_TOKEN, token);
    if (refreshToken) {
      localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken);
    }
    this.token = token;
    this.refreshToken = refreshToken;
  }

  /**
   * Clear tokens
   */
  clearTokens() {
    localStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.CURRENT_USER);
    localStorage.removeItem(STORAGE_KEYS.USER_ROLE);
    this.token = null;
    this.refreshToken = null;
  }

  /**
   * Build headers
   */
  getHeaders(contentType = 'application/json') {
    const headers = {
      'Content-Type': contentType,
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    return headers;
  }

  /**
   * Fetch wrapper with error handling
   */
  async fetch(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      method: options.method || 'GET',
      headers: this.getHeaders(options.contentType),
      timeout: API_CONFIG.TIMEOUT,
      ...options,
    };

    if (options.body && typeof options.body === 'object') {
      config.body = JSON.stringify(options.body);
    }

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), config.timeout);

      const response = await fetch(url, {
        ...config,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      // Handle token expiration
      if (response.status === 401) {
        this.clearTokens();
        window.location.href = '/pages/login.html';
        throw new Error('Unauthorized: Please log in again');
      }

      // Handle other HTTP errors
      if (!response.ok) {
        let errorMessage = `HTTP ${response.status}`;
        try {
          const errorData = await response.json();
          errorMessage = errorData.message || errorData.error || errorMessage;
        } catch (e) {
          // Response body is not JSON
        }
        throw new Error(errorMessage);
      }

      // Parse response
      const contentType = response.headers.get('content-type');
      if (contentType && contentType.includes('application/json')) {
        return await response.json();
      }
      return response;
    } catch (error) {
      if (error.name === 'AbortError') {
        throw new Error('Request timeout');
      }
      throw error;
    }
  }

  /**
   * GET request
   */
  async get(endpoint) {
    return this.fetch(endpoint, { method: 'GET' });
  }

  /**
   * POST request
   */
  async post(endpoint, body) {
    return this.fetch(endpoint, {
      method: 'POST',
      body,
    });
  }

  /**
   * PUT request
   */
  async put(endpoint, body) {
    return this.fetch(endpoint, {
      method: 'PUT',
      body,
    });
  }

  /**
   * DELETE request
   */
  async delete(endpoint) {
    return this.fetch(endpoint, { method: 'DELETE' });
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated() {
    return !!this.token;
  }

  /**
   * Get current user from storage
   */
  getCurrentUser() {
    const userJson = localStorage.getItem(STORAGE_KEYS.CURRENT_USER);
    return userJson ? JSON.parse(userJson) : null;
  }

  /**
   * Store current user
   */
  setCurrentUser(user) {
    localStorage.setItem(STORAGE_KEYS.CURRENT_USER, JSON.stringify(user));
  }

  /**
   * Get user role
   */
  getUserRole() {
    return localStorage.getItem(STORAGE_KEYS.USER_ROLE);
  }

  /**
   * Set user role
   */
  setUserRole(role) {
    localStorage.setItem(STORAGE_KEYS.USER_ROLE, role);
  }

  /**
   * Check if user has role
   */
  hasRole(roles) {
    const userRole = this.getUserRole();
    if (typeof roles === 'string') {
      return userRole === roles;
    }
    return roles.includes(userRole);
  }
}

// Export singleton instance
const apiClient = new APIClient();
export default apiClient;
