/**
 * Users API — UC 2.7.7
 */
import apiClient from './client.js';

const usersAPI = {
  async getUsers(params = {}) {
    const qs = new URLSearchParams(params).toString();
    return apiClient.get('/users' + (qs ? '?' + qs : ''));
  },

  async createUser(data) {
    return apiClient.post('/users', data);
  },

  async updateUser(id, data) {
    return apiClient.put(`/users/${id}`, data);
  },

  async deleteUser(id) {
    return apiClient.delete(`/users/${id}`);
  },
};

export default usersAPI;
