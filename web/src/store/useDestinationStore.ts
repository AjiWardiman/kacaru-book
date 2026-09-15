import { create } from 'zustand';
import { axiosClient } from '../api/axiosClient';

export interface DestinationModel {
  id: number;
  name: string;
  location: string;
  description?: string;
  category?: string;
  price: number;
  image_url?: string;
  is_active?: boolean;
  rating?: number;
}

export interface SlotModel {
  id: number;
  destination_id: number;
  date: string;
  quota?: number;
  total_quota?: number;
  booked_quota?: number;
  used_quota?: number;
  remaining_quota?: number;
  status?: string;
  is_available?: boolean;
}

interface DestinationState {
  destinations: DestinationModel[];
  isLoading: boolean;
  error: string | null;

  fetchDestinations: () => Promise<void>;
  createDestination: (data: Partial<DestinationModel>, imageFile?: File) => Promise<boolean>;
  updateDestination: (id: number, data: Partial<DestinationModel>, imageFile?: File) => Promise<boolean>;
  deleteDestination: (id: number) => Promise<boolean>;
  fetchSlots: (destId: number) => Promise<SlotModel[]>;
}

export const useDestinationStore = create<DestinationState>((set, get) => ({
  destinations: [],
  isLoading: false,
  error: null,

  fetchDestinations: async () => {
    set({ isLoading: true, error: null });
    try {
      const res = await axiosClient.get('/destinations');
      if (res.data.success) {
        set({
          destinations: res.data.data || [],
          isLoading: false,
        });
      } else {
        set({ isLoading: false, error: 'Gagal mengambil data destinasi' });
      }
    } catch (err: any) {
      console.error('Failed to fetch destinations:', err);
      set({ isLoading: false, error: err.response?.data?.error || 'Gagal memuat destinasi' });
    }
  },

  createDestination: async (destData, imageFile) => {
    set({ isLoading: true, error: null });
    try {
      let finalImageUrl = destData.image_url || '';

      if (imageFile) {
        const formData = new FormData();
        formData.append('file', imageFile);
        const uploadRes = await axiosClient.post('/admin/destinations/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' },
        });
        if (uploadRes.data.success && uploadRes.data.data?.url) {
          finalImageUrl = uploadRes.data.data.url;
        }
      }

      const payload = {
        name: destData.name,
        location: destData.location,
        description: destData.description,
        category: destData.category,
        price: Number(destData.price),
        image_url: finalImageUrl,
        is_active: true,
      };

      const res = await axiosClient.post('/admin/destinations', payload);
      if (res.data.success) {
        await get().fetchDestinations();
        set({ isLoading: false });
        return true;
      } else {
        set({ isLoading: false, error: res.data.error || 'Gagal membuat destinasi' });
        return false;
      }
    } catch (err: any) {
      console.error('Failed to create destination:', err);
      set({ isLoading: false, error: err.response?.data?.error || 'Gagal membuat destinasi' });
      return false;
    }
  },

  updateDestination: async (id, destData, imageFile) => {
    set({ isLoading: true, error: null });
    try {
      let finalImageUrl = destData.image_url || '';

      if (imageFile) {
        const formData = new FormData();
        formData.append('file', imageFile);
        const uploadRes = await axiosClient.post('/admin/destinations/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' },
        });
        if (uploadRes.data.success && uploadRes.data.data?.url) {
          finalImageUrl = uploadRes.data.data.url;
        }
      }

      const payload = {
        name: destData.name,
        location: destData.location,
        description: destData.description,
        category: destData.category,
        price: Number(destData.price),
        image_url: finalImageUrl,
        is_active: true,
      };

      const res = await axiosClient.put(`/admin/destinations/${id}`, payload);
      if (res.data.success) {
        await get().fetchDestinations();
        set({ isLoading: false });
        return true;
      } else {
        set({ isLoading: false, error: res.data.error || 'Gagal mengubah destinasi' });
        return false;
      }
    } catch (err: any) {
      console.error('Failed to update destination:', err);
      set({ isLoading: false, error: err.response?.data?.error || 'Gagal mengubah destinasi' });
      return false;
    }
  },

  deleteDestination: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      const res = await axiosClient.delete(`/admin/destinations/${id}`);
      if (res.data.success) {
        await get().fetchDestinations();
        set({ isLoading: false });
        return true;
      } else {
        set({ isLoading: false, error: res.data.error || 'Gagal menghapus destinasi' });
        return false;
      }
    } catch (err: any) {
      console.error('Failed to delete destination:', err);
      set({ isLoading: false, error: err.response?.data?.error || 'Gagal menghapus destinasi' });
      return false;
    }
  },

  fetchSlots: async (destId: number) => {
    try {
      const res = await axiosClient.get(`/destinations/${destId}/slots`);
      if (res.data.success) {
        return res.data.data || [];
      }
      return [];
    } catch (err) {
      console.error('Failed to fetch slots:', err);
      return [];
    }
  },
}));
