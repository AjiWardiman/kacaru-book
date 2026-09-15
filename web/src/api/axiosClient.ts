import axios from 'axios';
import { useAuthStore } from '../store/useAuthStore';

export const API_BASE_URL = 'http://localhost:8080';

export const axiosClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor untuk menyisipkan token JWT ke setiap request
axiosClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('kacaru_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Interceptor untuk menangani error global (misal 401 Unauthorized)
axiosClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Jika token expired / invalid, logout & hapus token lokal
      const isLoginRequest = error.config?.url?.includes('/auth/login');
      if (!isLoginRequest) {
        useAuthStore.getState().logout();
      }
    }
    return Promise.reject(error);
  }
);
