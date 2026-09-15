import { create } from 'zustand';
import { toast } from 'sonner';
import { axiosClient, API_BASE_URL } from '../api/axiosClient';
import { type NotificationModel } from '../schemas/schemas';

// Helper extractor untuk menangani casing field id, title, message, is_read, sent_at
export const getNotifId = (n: any): number => Number(n?.id ?? n?.ID ?? 0);
export const getNotifTitle = (n: any): string => String(n?.title ?? n?.Title ?? '');
export const getNotifMessage = (n: any): string => String(n?.message ?? n?.Message ?? '');
export const getNotifIsRead = (n: any): boolean => Boolean(n?.is_read ?? n?.IsRead ?? false);
export const getNotifSentAt = (n: any): string => String(n?.sent_at ?? n?.SentAt ?? new Date().toISOString());

interface NotificationState {
  notifications: NotificationModel[];
  unreadCount: number;
  isConnected: boolean;
  activeToast: NotificationModel | null;
  eventSource: EventSource | null;

  fetchNotifications: () => Promise<void>;
  connectSSE: (token: string) => void;
  disconnectSSE: () => void;
  markAsRead: (id: number) => Promise<void>;
  markAllAsRead: () => Promise<void>;
  dismissToast: () => void;
}

export const useNotificationStore = create<NotificationState>((set, get) => ({
  notifications: [],
  unreadCount: 0,
  isConnected: false,
  activeToast: null,
  eventSource: null,

  fetchNotifications: async () => {
    try {
      const res = await axiosClient.get('/notifications');
      if (res.data.success) {
        set({
          notifications: res.data.notifications || [],
          unreadCount: res.data.unread_count || 0,
        });
      }
    } catch (err) {
      console.error('Failed to fetch notifications:', err);
    }
  },

  connectSSE: (token: string) => {
    get().disconnectSSE();

    const url = `${API_BASE_URL}/notifications/stream?token=${encodeURIComponent(token)}`;
    const es = new EventSource(url);

    es.onopen = () => {
      set({ isConnected: true, eventSource: es });
      console.log('⚡ SSE Connected to Kacaru Book Notification Stream');
    };

    es.addEventListener('notification', (event: MessageEvent) => {
      try {
        const notif: NotificationModel = JSON.parse(event.data);
        const incomingId = getNotifId(notif);
        const isRead = getNotifIsRead(notif);
        const title = getNotifTitle(notif);
        const message = getNotifMessage(notif);

        // Show elegant Sonner toast notification
        toast.info(title || 'Notifikasi Baru', {
          description: message,
          duration: 4500,
        });

        set((state) => {
          const exists = state.notifications.some((n) => getNotifId(n) === incomingId);
          const updated = exists
            ? state.notifications.map((n) => (getNotifId(n) === incomingId ? notif : n))
            : [notif, ...state.notifications];

          return {
            notifications: updated,
            unreadCount: isRead ? state.unreadCount : state.unreadCount + 1,
            activeToast: notif,
          };
        });
      } catch (err) {
        console.error('Error parsing SSE notification payload:', err);
      }
    });

    es.onerror = (err) => {
      console.warn('SSE EventSource error/reconnecting:', err);
      set({ isConnected: false });
    };
  },

  disconnectSSE: () => {
    const { eventSource } = get();
    if (eventSource) {
      eventSource.close();
    }
    set({
      notifications: [],
      unreadCount: 0,
      isConnected: false,
      activeToast: null,
      eventSource: null,
    });
  },

  markAsRead: async (id: number) => {
    try {
      await axiosClient.patch(`/notifications/${id}/read`);
      set((state) => ({
        notifications: state.notifications.map((n) =>
          getNotifId(n) === id ? { ...n, is_read: true, IsRead: true } : n
        ),
        unreadCount: Math.max(0, state.unreadCount - 1),
      }));
    } catch (err) {
      console.error('Failed to mark notification as read:', err);
    }
  },

  markAllAsRead: async () => {
    try {
      await axiosClient.patch('/notifications/read-all');
      set((state) => ({
        notifications: state.notifications.map((n) => ({ ...n, is_read: true, IsRead: true })),
        unreadCount: 0,
      }));
    } catch (err) {
      console.error('Failed to mark all as read:', err);
    }
  },

  dismissToast: () => set({ activeToast: null }),
}));
