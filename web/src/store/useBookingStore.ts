import { create } from 'zustand';
import { axiosClient } from '../api/axiosClient';

export interface VisitorInput {
  name: string;
  phone: string;
  email: string;
  age: number;
}

export interface BookingModel {
  id: number;
  booking_code: string;
  slot_id: number;
  total_people: number;
  total_price: number;
  status: 'PENDING' | 'PAID' | 'EXPIRED' | 'CANCELLED' | string;
  booked_at: string;
  expired_at?: string;
  destination_name?: string;
  destination_location?: string;
  destination_image?: string;
  user_name?: string;
  payment_url?: string;
}

interface BookingState {
  bookings: BookingModel[];
  isLoading: boolean;
  error: string | null;

  fetchMyBookings: () => Promise<void>;
  fetchAllAdminBookings: () => Promise<void>;
  createBooking: (slotId: number, totalPeople: number, visitors: VisitorInput[]) => Promise<BookingModel | null>;
  payBooking: (bookingId: number) => Promise<string | null>;
  markBookingAsPaid: (identifier: string | number) => void;
}

export const useBookingStore = create<BookingState>((set, get) => ({
  bookings: [],
  isLoading: false,
  error: null,

  fetchMyBookings: async () => {
    set({ isLoading: true, error: null });
    try {
      const res = await axiosClient.get('/bookings/my-bookings');
      if (res.data.success) {
        set({
          bookings: res.data.data || [],
          isLoading: false,
        });
      } else {
        set({ isLoading: false, error: 'Gagal mengambil riwayat pemesanan' });
      }
    } catch (err: any) {
      console.error('Failed to fetch my bookings:', err);
      set({ isLoading: false, error: err.response?.data?.error || 'Gagal memuat pemesanan' });
    }
  },

  fetchAllAdminBookings: async () => {
    set({ isLoading: true, error: null });
    try {
      const res = await axiosClient.get('/admin/bookings');
      if (res.data.success) {
        set({
          bookings: res.data.data || [],
          isLoading: false,
        });
      } else {
        set({ isLoading: false, error: 'Gagal mengambil seluruh pemesanan platform' });
      }
    } catch (err: any) {
      console.error('Failed to fetch all admin bookings:', err);
      set({ isLoading: false, error: err.response?.data?.error || 'Gagal memuat pemesanan platform' });
    }
  },

  createBooking: async (slotId: number, totalPeople: number, visitors: VisitorInput[]) => {
    set({ isLoading: true, error: null });
    try {
      const payload = {
        slot_id: slotId,
        total_people: totalPeople,
        visitors: visitors,
      };
      const res = await axiosClient.post('/bookings', payload);
      if (res.data.success) {
        const newBooking: BookingModel = res.data.data;
        await get().fetchMyBookings();
        set({ isLoading: false });
        return newBooking;
      } else {
        set({ isLoading: false, error: res.data.error || 'Gagal membuat pemesanan' });
        return null;
      }
    } catch (err: any) {
      console.error('Failed to create booking:', err);
      const errMsg = err.response?.data?.error || 'Gagal membuat pemesanan';
      set({ isLoading: false, error: errMsg });
      return null;
    }
  },

  payBooking: async (bookingId: number) => {
    set({ isLoading: true, error: null });
    try {
      const res = await axiosClient.post(`/bookings/${bookingId}/pay`);
      if (res.data.success && res.data.data?.payment_url) {
        set({ isLoading: false });
        return res.data.data.payment_url;
      } else {
        set({ isLoading: false, error: res.data.error || 'Gagal membuat link pembayaran' });
        return null;
      }
    } catch (err: any) {
      console.error('Failed to pay booking:', err);
      const errMsg = err.response?.data?.error || 'Gagal membuat link pembayaran';
      set({ isLoading: false, error: errMsg });
      return null;
    }
  },

  markBookingAsPaid: (identifier: string | number) => {
    set((state) => ({
      bookings: state.bookings.map((b) =>
        String(b.id) === String(identifier) || b.booking_code === String(identifier)
          ? { ...b, status: 'PAID' }
          : b
      ),
    }));
  },
}));
