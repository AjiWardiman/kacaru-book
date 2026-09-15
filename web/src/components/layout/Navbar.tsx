import React, { useState } from 'react';
import {
  Sun,
  Moon,
  Bell,
  Menu
} from 'lucide-react';
import type { UserResponse } from '../../api/schema';
import { getNotifId, getNotifTitle, getNotifMessage, getNotifIsRead, getNotifSentAt } from '../../store/useNotificationStore';

interface NavbarProps {
  user: UserResponse | null;
  theme: 'light' | 'dark';
  setTheme: (theme: 'light' | 'dark') => void;
  notifications: any[];
  unreadCount: number;
  isConnected: boolean;
  onMarkAsRead: (id: number) => void;
  onMarkAllAsRead: () => void;
  setIsMobileOpen: (open: boolean) => void;
}

export const Navbar: React.FC<NavbarProps> = ({
  user,
  theme,
  setTheme,
  notifications,
  unreadCount,
  isConnected,
  onMarkAsRead,
  onMarkAllAsRead,
  setIsMobileOpen,
}) => {
  const [isNotifOpen, setIsNotifOpen] = useState(false);

  return (
    <header className="sticky top-0 z-30 bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl border-b border-slate-200/80 dark:border-slate-800/80 px-6 py-4 flex items-center justify-between transition-colors">
      <div className="flex items-center gap-4">
        <button
          onClick={() => setIsMobileOpen(true)}
          className="p-2.5 rounded-2xl bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-200 lg:hidden"
        >
          <Menu className="w-5 h-5" />
        </button>

        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-extrabold text-slate-900 dark:text-white">
              Selamat Datang, {user?.name || 'Pengguna'}! 👋
            </h2>
            <span
              className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider shadow-sm transition-all ${
                user?.role === 'admin'
                  ? 'bg-gradient-to-r from-amber-500/20 to-orange-500/20 text-amber-600 dark:text-amber-300 border border-amber-500/40 shadow-amber-500/10'
                  : 'bg-gradient-to-r from-emerald-500/20 to-teal-500/20 text-emerald-600 dark:text-emerald-300 border border-emerald-500/40 shadow-emerald-500/10'
              }`}
            >
              {user?.role === 'admin' ? '👑 Admin Utama' : '👤 Customer'}
            </span>
          </div>
          <p className="text-xs text-slate-500 dark:text-slate-400 font-bold">
            Sistem Reservasi & Tiket Wisata Kabupaten Bima, NTB
          </p>
        </div>
      </div>

      <div className="flex items-center gap-3">
        {/* SSE Stream Live Indicator */}
        <div
          className={`hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-full text-[11px] font-extrabold border transition-all duration-300 ${
            isConnected
              ? 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/30 shadow-md shadow-emerald-500/10'
              : 'bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30 shadow-md shadow-amber-500/10'
          }`}
          title={isConnected ? 'Koneksi SSE Realtime Aktif' : 'Menghubungkan SSE Stream...'}
        >
          <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-emerald-500 animate-ping' : 'bg-amber-500 animate-pulse'}`} />
          <span>{isConnected ? 'SSE Live' : 'Connecting'}</span>
        </div>

        {/* Notification Bell Dropdown */}
        <div className="relative">
          <button
            onClick={() => setIsNotifOpen(!isNotifOpen)}
            className="p-2.5 rounded-2xl bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 transition-all duration-300 hover:scale-105 active:scale-95 border border-slate-200/60 dark:border-slate-700/60 relative group"
            title="Notifikasi Sistem"
          >
            <Bell className="w-5 h-5 group-hover:rotate-12 transition-transform" />
            {unreadCount > 0 && (
              <span className="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-gradient-to-r from-rose-500 to-pink-500 text-white text-[10px] font-black flex items-center justify-center shadow-lg shadow-rose-500/40 animate-bounce">
                {unreadCount > 9 ? '9+' : unreadCount}
              </span>
            )}
          </button>

          {isNotifOpen && (
            <div className="absolute right-0 mt-3 w-80 sm:w-96 bg-white/95 dark:bg-slate-900/95 backdrop-blur-2xl rounded-3xl shadow-2xl border border-slate-200/80 dark:border-slate-800 p-4 space-y-3 z-50 animate-in fade-in slide-in-from-top-2 duration-200">
              <div className="flex items-center justify-between pb-3 border-b border-slate-100 dark:border-slate-800">
                <div className="flex items-center gap-2">
                  <h4 className="text-sm font-black text-slate-900 dark:text-white">Pemberitahuan</h4>
                  <span className="px-2.5 py-0.5 rounded-full bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 text-[10px] font-extrabold border border-emerald-500/30">
                    {unreadCount} Baru
                  </span>
                </div>
                {unreadCount > 0 && (
                  <button
                    onClick={onMarkAllAsRead}
                    className="text-[11px] font-black text-emerald-600 dark:text-emerald-400 hover:underline"
                  >
                    Tandai Semua Dibaca
                  </button>
                )}
              </div>

              <div className="max-h-80 overflow-y-auto space-y-2 pr-1">
                {notifications.length === 0 ? (
                  <div className="p-8 text-center text-xs text-slate-400 font-medium">
                    Belum ada pemberitahuan baru.
                  </div>
                ) : (
                  notifications.map((n) => {
                    const nid = getNotifId(n);
                    const ntitle = getNotifTitle(n);
                    const nmsg = getNotifMessage(n);
                    const nisRead = getNotifIsRead(n);
                    const nsent = getNotifSentAt(n);
                    return (
                      <div
                        key={nid}
                        onClick={() => !nisRead && onMarkAsRead(nid)}
                        className={`p-3 rounded-2xl border transition-all cursor-pointer hover:scale-[1.01] ${
                          !nisRead
                            ? 'bg-gradient-to-r from-emerald-50/80 to-teal-50/80 dark:from-emerald-950/40 dark:to-teal-950/30 border-emerald-500/40 shadow-sm'
                            : 'bg-slate-50/50 dark:bg-slate-800/40 border-slate-100 dark:border-slate-800'
                        }`}
                      >
                        <div className="flex items-start justify-between gap-2">
                          <h5 className="text-xs font-black text-slate-900 dark:text-white">{ntitle}</h5>
                          {!nisRead && <span className="w-2 h-2 rounded-full bg-emerald-500 shrink-0 mt-1 shadow-sm shadow-emerald-500" />}
                        </div>
                        <p className="text-[11px] text-slate-600 dark:text-slate-300 font-medium mt-1">
                          {nmsg}
                        </p>
                        <span className="text-[9px] text-slate-400 font-semibold block mt-2">
                          {nsent ? new Date(nsent).toLocaleString('id-ID') : 'Baru saja'}
                        </span>
                      </div>
                    );
                  })
                )}
              </div>
            </div>
          )}
        </div>

        {/* Theme Toggle Button */}
        <button
          onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
          className="p-2.5 rounded-2xl bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 transition-all duration-300 hover:scale-105 active:scale-95 border border-slate-200/60 dark:border-slate-700/60 shadow-sm"
          title="Ganti Tema Dark/Light"
        >
          {theme === 'light' ? (
            <Moon className="w-5 h-5 text-indigo-600 hover:rotate-12 transition-transform" />
          ) : (
            <Sun className="w-5 h-5 text-amber-400 hover:rotate-45 transition-transform" />
          )}
        </button>
      </div>
    </header>
  );
};
