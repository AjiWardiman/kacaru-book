import { create } from 'zustand';
import { axiosClient } from '../api/axiosClient';
import { loginSchema, type LoginInput, type AuthUser } from '../schemas/schemas';

interface AuthState {
  user: AuthUser | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  
  login: (credentials: LoginInput | string, passwordArg?: string, rememberMe?: boolean) => Promise<boolean>;
  register: (data: { name: string; email: string; phone?: string; password: string }) => Promise<boolean>;
  completeAuth: () => void;
  requestOTP: () => Promise<boolean>;
  updateProfile: (data: { name: string; email: string; phone?: string }) => Promise<boolean>;
  updatePassword: (data: { old_password?: string; new_password: string }) => Promise<boolean>;
  logout: () => void;
  initAuth: () => void;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,

  initAuth: () => {
    const token = localStorage.getItem('kacaru_token');
    const userStr = localStorage.getItem('kacaru_user');
    if (token && userStr && userStr !== 'undefined' && userStr !== 'null') {
      try {
        const user = JSON.parse(userStr);
        if (user && user.id && user.role) {
          set({ token, user, isAuthenticated: true });
          return;
        }
      } catch {
        // reset invalid session
      }
    }
    localStorage.removeItem('kacaru_token');
    localStorage.removeItem('kacaru_user');
    set({ token: null, user: null, isAuthenticated: false });
  },

  login: async (credentials: LoginInput | string, passwordArg?: string, _rememberMe?: boolean) => {
    set({ isLoading: true, error: null });

    const identifier = typeof credentials === 'string' ? credentials : (credentials.identifier || (credentials as any).email || '');
    const password = typeof credentials === 'string' ? passwordArg || '' : credentials.password;

    // Validasi input dengan Zod (dukung Email atau Nomor HP)
    const validation = loginSchema.safeParse({ identifier, password });
    if (!validation.success) {
      const firstErr = validation.error.issues[0]?.message || 'Input tidak valid';
      set({ isLoading: false, error: firstErr });
      return false;
    }

    try {
      const response = await axiosClient.post('/auth/login', {
        identifier: identifier.trim(),
        password: password,
      });
      const resData = response.data;

      if (resData && resData.success && resData.data && resData.data.token && resData.data.user) {
        const { token, user } = resData.data;
        localStorage.setItem('kacaru_token', token);
        localStorage.setItem('kacaru_user', JSON.stringify(user));

        set({
          token,
          user,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        });
        return true;
      }

      const errMsg = resData?.error || 'Kredensial tidak valid. Silakan periksa kembali email/no HP dan kata sandi Anda.';
      set({ isLoading: false, error: errMsg });
      return false;
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal terhubung ke server backend';
      set({ isLoading: false, error: errMsg });
      return false;
    }
  },

  register: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosClient.post('/auth/register', {
        name: data.name,
        email: data.email,
        phone: data.phone || '',
        password: data.password,
      });

      const resData = response.data;
      if (resData.success && resData.data) {
        const { token, user } = resData.data;
        localStorage.setItem('kacaru_token', token);
        localStorage.setItem('kacaru_user', JSON.stringify(user));
        set({
          token: token,
          user: user,
          isAuthenticated: false, // Menahan login otomatis agar Modal OTP wajib selesai lebih dulu!
          isLoading: false,
          error: null,
        });
        return true;
      }
      set({ isLoading: false, error: resData.error || 'Registrasi gagal' });
      return false;
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal mendaftar';
      set({ isLoading: false, error: errMsg });
      return false;
    }
  },

  completeAuth: () => {
    set({ isAuthenticated: true });
  },

  requestOTP: async () => {
    try {
      const token = get().token || localStorage.getItem('kacaru_token');
      if (!token) return false;
      const res = await axiosClient.post(
        '/otp/request',
        {},
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
      return res.data?.success || false;
    } catch {
      return false;
    }
  },

  updateProfile: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const res = await axiosClient.put('/users/profile', data);
      if (res.data && res.data.success) {
        const updatedData = res.data.data;
        const currentUser = get().user;
        const newUser: AuthUser = {
          ...currentUser,
          id: updatedData?.id || currentUser?.id || 1,
          name: updatedData?.name || data.name,
          email: updatedData?.email || data.email,
          phone: updatedData?.phone || data.phone,
          role: updatedData?.role || currentUser?.role || 'user',
        };
        localStorage.setItem('kacaru_user', JSON.stringify(newUser));
        set({ user: newUser, isLoading: false });
        return true;
      }
      set({ isLoading: false, error: res.data?.error || 'Gagal memperbarui profil' });
      return false;
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal memperbarui profil di database';
      set({ isLoading: false, error: errMsg });
      return false;
    }
  },

  updatePassword: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const res = await axiosClient.put('/users/password', data);
      if (res.data && res.data.success) {
        set({ isLoading: false });
        return true;
      }
      set({ isLoading: false, error: res.data?.error || 'Gagal memperbarui password' });
      return false;
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal memperbarui kata sandi di database';
      set({ isLoading: false, error: errMsg });
      return false;
    }
  },

  logout: () => {
    localStorage.removeItem('kacaru_token');
    localStorage.removeItem('kacaru_user');
    set({ user: null, token: null, isAuthenticated: false, error: null });
  },

  clearError: () => set({ error: null }),
}));
