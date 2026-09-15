import { create } from 'zustand';
import { axiosClient } from '../api/axiosClient';
import { 
  userCreateSchema, 
  userUpdateSchema, 
  type UserCreateInput, 
  type UserUpdateInput, 
  type UserItem 
} from '../schemas/schemas';

interface UserState {
  users: UserItem[];
  total: number;
  isLoading: boolean;
  error: string | null;
  successMessage: string | null;

  fetchUsers: (search?: string, role?: string) => Promise<void>;
  createUser: (data: UserCreateInput) => Promise<boolean>;
  updateUser: (id: number | string, data: UserUpdateInput) => Promise<boolean>;
  deleteUser: (id: number | string) => Promise<boolean>;
  clearMessages: () => void;
}

export const useUserStore = create<UserState>((set, get) => ({
  users: [],
  total: 0,
  isLoading: false,
  error: null,
  successMessage: null,

  fetchUsers: async (search?: string, role?: string) => {
    set({ isLoading: true, error: null });
    try {
      const params: Record<string, string> = { limit: '500' };
      if (search && search.trim() !== '') params.search = search.trim();
      if (role && role.toUpperCase() !== 'ALL') {
        params.role = role.toLowerCase() === 'customer' ? 'user' : role;
      }

      const response = await axiosClient.get('/admin/users', { params });
      if (response.data && response.data.success) {
        const list = response.data.data?.users || [];
        const total = response.data.data?.total || list.length;
        set({ users: list, total, isLoading: false });
      } else {
        set({ isLoading: false, error: response.data.error || 'Gagal mengambil data user' });
      }
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal memuat data pengguna';
      set({ isLoading: false, error: errMsg });
    }
  },

  createUser: async (data: UserCreateInput) => {
    set({ isLoading: true, error: null, successMessage: null });
    const validation = userCreateSchema.safeParse(data);
    if (!validation.success) {
      const firstErr = validation.error.issues[0]?.message || 'Input data tidak valid';
      set({ isLoading: false, error: firstErr });
      return false;
    }

    try {
      const response = await axiosClient.post('/admin/users', data);
      if (response.data && response.data.success) {
        set({ isLoading: false, successMessage: 'Pengguna berhasil ditambahkan' });
        await get().fetchUsers();
        return true;
      } else {
        set({ isLoading: false, error: response.data.error || 'Gagal menambah pengguna' });
        return false;
      }
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal menambah pengguna';
      set({ isLoading: false, error: errMsg });
      return false;
    }
  },

  updateUser: async (id: number | string, data: UserUpdateInput) => {
    set({ isLoading: true, error: null, successMessage: null });
    const validation = userUpdateSchema.safeParse(data);
    if (!validation.success) {
      const firstErr = validation.error.issues[0]?.message || 'Input data tidak valid';
      set({ isLoading: false, error: firstErr });
      return false;
    }

    try {
      const response = await axiosClient.put(`/admin/users/${id}`, data);
      if (response.data && response.data.success) {
        set({ isLoading: false, successMessage: 'Pengguna berhasil diperbarui' });
        await get().fetchUsers();
        return true;
      } else {
        set({ isLoading: false, error: response.data.error || 'Gagal memperbarui pengguna' });
        return false;
      }
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal memperbarui pengguna';
      set({ isLoading: false, error: errMsg });
      return false;
    }
  },

  deleteUser: async (id: number | string) => {
    set({ isLoading: true, error: null, successMessage: null });
    try {
      const response = await axiosClient.delete(`/admin/users/${id}`);
      if (response.data && response.data.success) {
        set({ isLoading: false, successMessage: 'Pengguna berhasil dihapus' });
        await get().fetchUsers();
        return true;
      } else {
        set({ isLoading: false, error: response.data.error || 'Gagal menghapus pengguna' });
        return false;
      }
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Gagal menghapus pengguna';
      set({ isLoading: false, error: errMsg });
      return false;
    }
  },

  clearMessages: () => set({ error: null, successMessage: null }),
}));
