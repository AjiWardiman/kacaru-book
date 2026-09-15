import React from 'react';
import {
  LayoutDashboard,
  Compass,
  Ticket,
  Users,
  Settings,
  LogOut,
  X,
  BarChart3
} from 'lucide-react';
import type { UserResponse } from '../../api/schema';

interface SidebarProps {
  user: UserResponse | null;
  currentRoute: string;
  setCurrentRoute: (route: string) => void;
  onLogout: () => void;
  isMobileOpen: boolean;
  setIsMobileOpen: (open: boolean) => void;
  effectiveUserBookingsCount: number;
}

export const Sidebar: React.FC<SidebarProps> = ({
  user,
  currentRoute,
  setCurrentRoute,
  onLogout,
  isMobileOpen,
  setIsMobileOpen,
  effectiveUserBookingsCount,
}) => {
  const isAdmin = user?.role === 'admin';

  const navItems = isAdmin
    ? [
        { id: 'admin/dashboard', label: 'Dashboard Overview', icon: LayoutDashboard, color: 'text-sky-400', activeGlow: 'from-sky-600 via-indigo-600 to-sky-600 shadow-sky-500/30' },
        { id: 'admin/reports', label: 'Laporan Analytics', icon: BarChart3, color: 'text-indigo-400', activeGlow: 'from-indigo-600 via-purple-600 to-indigo-600 shadow-indigo-500/30' },
        { id: 'destinations', label: 'Kelola Destinasi', icon: Compass, color: 'text-emerald-400', activeGlow: 'from-emerald-600 via-teal-600 to-emerald-600 shadow-emerald-500/30' },
        { id: 'admin/bookings', label: 'Pemesanan Platform', icon: Ticket, color: 'text-amber-400', activeGlow: 'from-amber-600 via-orange-600 to-amber-600 shadow-amber-500/30' },
        { id: 'admin/users', label: 'Pengguna System', icon: Users, color: 'text-purple-400', activeGlow: 'from-purple-600 via-pink-600 to-purple-600 shadow-purple-500/30' },
        { id: 'admin/settings', label: 'Pengaturan Akun', icon: Settings, color: 'text-teal-400', activeGlow: 'from-teal-600 via-emerald-600 to-teal-600 shadow-teal-500/30' },
      ]
    : [
        { id: 'destinations', label: 'Eksplorasi Wisata', icon: Compass, color: 'text-emerald-400', activeGlow: 'from-emerald-600 via-teal-600 to-emerald-600 shadow-emerald-500/30' },
        {
          id: 'bookings/my-bookings',
          label: 'Tiket & Pemesanan Saya',
          icon: Ticket,
          color: 'text-amber-400',
          activeGlow: 'from-amber-600 via-orange-600 to-amber-600 shadow-amber-500/30',
          badge: `${effectiveUserBookingsCount} Tiket`,
        },
        { id: 'settings', label: 'Pengaturan Profil', icon: Settings, color: 'text-teal-400', activeGlow: 'from-teal-600 via-emerald-600 to-teal-600 shadow-teal-500/30' },
      ];

  return (
    <>
      {/* Backdrop overlay for mobile */}
      {isMobileOpen && (
        <div
          onClick={() => setIsMobileOpen(false)}
          className="fixed inset-0 z-40 bg-slate-950/70 backdrop-blur-md lg:hidden animate-fade-in"
        />
      )}

      <aside
        className={`fixed top-0 left-0 bottom-0 z-50 w-72 bg-white/95 dark:bg-slate-900/95 border-r border-slate-200/80 dark:border-slate-800/90 shadow-2xl backdrop-blur-2xl transition-all duration-300 ease-in-out flex flex-col justify-between p-6 ${
          isMobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        }`}
      >
        <div className="space-y-8">
          {/* Header Logo */}
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3.5 group cursor-pointer">
              <img
                src="/kacaru-logo.png"
                alt="Kacaru Book Logo"
                className="w-14 h-14 rounded-full object-cover shadow-xl shadow-sky-500/20 group-hover:scale-105 transition-all duration-300 shrink-0"
              />
              <div>
                <h1 className="text-xl font-black tracking-tight bg-gradient-to-r from-emerald-400 via-teal-300 to-cyan-300 bg-clip-text text-transparent drop-shadow-sm">
                  Kacaru Book
                </h1>
                <p className="text-[10px] font-extrabold uppercase tracking-widest text-emerald-500/90 dark:text-emerald-400">
                  Tourism Bima NTB
                </p>
              </div>
            </div>
            <button
              onClick={() => setIsMobileOpen(false)}
              className="lg:hidden p-2 rounded-2xl text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          {/* Navigation Links */}
          <nav className="space-y-2">
            <div className="px-3 pb-2 text-[10px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 flex items-center justify-between">
              <span>Menu Navigasi ({isAdmin ? 'Admin' : 'Customer'})</span>
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-ping" />
            </div>
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = currentRoute === item.id;
              return (
                <button
                  key={item.id}
                  onClick={() => {
                    setCurrentRoute(item.id);
                    setIsMobileOpen(false);
                  }}
                  className={`w-full relative flex items-center justify-between px-4 py-3 rounded-2xl font-extrabold text-xs transition-all duration-300 group ${
                    isActive
                      ? `bg-gradient-to-r ${item.activeGlow} text-white shadow-xl scale-[1.02]`
                      : 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800/80 hover:text-slate-900 dark:hover:text-white hover:translate-x-1.5'
                  }`}
                >
                  {/* Left Active Glow Indicator */}
                  {isActive && (
                    <span className="w-1.5 h-6 bg-white rounded-r-full absolute left-0 shadow-md shadow-white/50" />
                  )}

                  <div className="flex items-center gap-3">
                    <div
                      className={`p-1.5 rounded-xl transition-all duration-300 ${
                        isActive
                          ? 'bg-white/20 text-white'
                          : `bg-slate-100 dark:bg-slate-800 ${item.color} group-hover:scale-110 group-hover:rotate-6`
                      }`}
                    >
                      <Icon className="w-4 h-4" />
                    </div>
                    <span>{item.label}</span>
                  </div>

                  {item.badge && (
                    <span
                      className={`px-2.5 py-0.5 rounded-full text-[10px] font-black tracking-wide transition-all ${
                        isActive
                          ? 'bg-white/25 text-white backdrop-blur-md shadow-sm'
                          : 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30'
                      }`}
                    >
                      {item.badge}
                    </span>
                  )}
                </button>
              );
            })}
          </nav>
        </div>

        {/* User Card & Logout */}
        <div className="pt-6 border-t border-slate-200/80 dark:border-slate-800/80 space-y-4">
          <div className="flex items-center gap-3 p-3 rounded-2xl bg-gradient-to-r from-slate-100 via-slate-50 to-slate-100 dark:from-slate-800/80 dark:via-slate-900/60 dark:to-slate-800/80 border border-slate-200/80 dark:border-slate-700/60 shadow-md group hover:border-emerald-500/40 transition-all duration-300">
            <div className="w-10 h-10 rounded-2xl bg-gradient-to-tr from-emerald-500 via-teal-500 to-cyan-500 text-white font-black flex items-center justify-center text-sm shadow-md ring-2 ring-emerald-500/30 group-hover:scale-105 transition-all">
              {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
            </div>
            <div className="overflow-hidden flex-1">
              <h4 className="text-xs font-black text-slate-900 dark:text-white truncate group-hover:text-emerald-400 transition-colors">
                {user?.name || 'User'}
              </h4>
              <p className="text-[10px] text-slate-400 font-bold truncate">
                {user?.email || 'user@example.com'}
              </p>
            </div>
          </div>

          <button
            onClick={onLogout}
            className="w-full flex items-center justify-center gap-2 py-3 px-4 rounded-2xl bg-rose-500/10 hover:bg-rose-500/20 text-rose-600 dark:text-rose-400 font-black text-xs transition-all duration-300 active:scale-95 border border-rose-500/30 shadow-sm hover:shadow-rose-500/10 group"
          >
            <LogOut className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
            <span>Keluar Aplikasi</span>
          </button>
        </div>
      </aside>
    </>
  );
};
