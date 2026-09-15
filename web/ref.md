import React, { useState, useEffect, useMemo } from 'react';
import {
  LayoutDashboard,
  Users,
  Compass,
  Ticket,
  Settings,
  LogOut,
  Sun,
  Moon,
  Eye,
  EyeOff,
  ChevronRight,
  Search,
  Plus,
  Edit3,
  Trash2,
  CheckCircle,
  Clock,
  XCircle,
  Shield,
  ArrowUpRight,
  TrendingUp,
  Sparkles,
  MapPin,
  Calendar,
  DollarSign,
  Filter,
  Check,
  Bell,
  Menu,
  X,
  ChevronDown,
  Briefcase,
  Layers,
  Star,
  Activity,
  UserCheck,
  Lock
} from 'lucide-react';
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  PieChart,
  Pie,
  Cell,
  Legend
} from 'recharts';

type Role = 'Admin' | 'Staff' | 'Customer';
type ThemeMode = 'light' | 'dark';

interface UserItem {
  id: string;
  name: string;
  email: string;
  role: Role;
  status: 'Aktif' | 'Nonaktif' | 'Pending';
  avatar: string;
  joinedDate: string;
}

interface DestinationItem {
  id: string;
  name: string;
  location: string;
  price: number;
  category: 'Bahari' | 'Gunung' | 'Budaya' | 'Kuliner';
  rating: number;
  image: string;
  status: 'Tersedia' | 'Penuh';
}

interface BookingItem {
  id: string;
  customerName: string;
  destinationName: string;
  date: string;
  totalPrice: number;
  status: 'Lunas' | 'Menunggu' | 'Dibatalkan';
}

const INITIAL_USERS: UserItem[] = [
  {
    id: 'USR-001',
    name: 'Rizky Pratama',
    email: 'rizky.admin@kacarubook.com',
    role: 'Admin',
    status: 'Aktif',
    avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=150',
    joinedDate: '12 Jan 2025'
  },
  {
    id: 'USR-002',
    name: 'Siti Rahmawati',
    email: 'siti.staff@kacarubook.com',
    role: 'Staff',
    status: 'Aktif',
    avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&q=80&w=150',
    joinedDate: '03 Feb 2025'
  },
  {
    id: 'USR-003',
    name: 'Budi Santoso',
    email: 'budi.traveler@gmail.com',
    role: 'Customer',
    status: 'Aktif',
    avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=150',
    joinedDate: '18 Mar 2025'
  },
  {
    id: 'USR-004',
    name: 'Dewi Lestari',
    email: 'dewi.lestari@yahoo.com',
    role: 'Customer',
    status: 'Pending',
    avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&q=80&w=150',
    joinedDate: '01 Apr 2025'
  },
  {
    id: 'USR-005',
    name: 'Andi Wijaya',
    email: 'andi.wijaya@kacarubook.com',
    role: 'Staff',
    status: 'Nonaktif',
    avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&q=80&w=150',
    joinedDate: '10 Mei 2025'
  }
];

const INITIAL_DESTINATIONS: DestinationItem[] = [
  {
    id: 'DEST-01',
    name: 'Raja Ampat Islands',
    location: 'Papua Barat',
    price: 4500000,
    category: 'Bahari',
    rating: 4.9,
    image: 'https://images.unsplash.com/photo-1516690561799-46d8f74f9abf?auto=format&fit=crop&q=80&w=600',
    status: 'Tersedia'
  },
  {
    id: 'DEST-02',
    name: 'Gunung Bromo Sunrise',
    location: 'Jawa Timur',
    price: 1250000,
    category: 'Gunung',
    rating: 4.8,
    image: 'https://images.unsplash.com/photo-1588668214407-6ea9a6d8c272?auto=format&fit=crop&q=80&w=600',
    status: 'Tersedia'
  },
  {
    id: 'DEST-03',
    name: 'Candi Borobudur Heritage',
    location: 'Magelang, Jawa Tengah',
    price: 750000,
    category: 'Budaya',
    rating: 4.7,
    image: 'https://images.unsplash.com/photo-1609949279531-cf48d64bed89?auto=format&fit=crop&q=80&w=600',
    status: 'Tersedia'
  },
  {
    id: 'DEST-04',
    name: 'Nusa Penida Snorkeling',
    location: 'Bali',
    price: 2100000,
    category: 'Bahari',
    rating: 4.9,
    image: 'https://images.unsplash.com/photo-1537996194471-e657df975ab4?auto=format&fit=crop&q=80&w=600',
    status: 'Penuh'
  }
];

const INITIAL_BOOKINGS: BookingItem[] = [
  { id: 'BK-9901', customerName: 'Budi Santoso', destinationName: 'Raja Ampat Islands', date: '2026-09-15', totalPrice: 4500000, status: 'Lunas' },
  { id: 'BK-9902', customerName: 'Dewi Lestari', destinationName: 'Gunung Bromo Sunrise', date: '2026-09-18', totalPrice: 1250000, status: 'Menunggu' },
  { id: 'BK-9903', customerName: 'Rudi Hartono', destinationName: 'Candi Borobudur Heritage', date: '2026-09-20', totalPrice: 750000, status: 'Lunas' },
  { id: 'BK-9904', customerName: 'Maya Indah', destinationName: 'Nusa Penida Snorkeling', date: '2026-09-22', totalPrice: 2100000, status: 'Dibatalkan' }
];

const REVENUE_CHART_DATA = [
  { month: 'Jan', Pendapatan: 42000000, Pemesanan: 120 },
  { month: 'Feb', Pendapatan: 58000000, Pemesanan: 180 },
  { month: 'Mar', Pendapatan: 65000000, Pemesanan: 210 },
  { month: 'Apr', Pendapatan: 51000000, Pemesanan: 160 },
  { month: 'Mei', Pendapatan: 89000000, Pemesanan: 290 },
  { month: 'Jun', Pendapatan: 105000000, Pemesanan: 340 },
  { month: 'Jul', Pendapatan: 128000000, Pemesanan: 410 },
  { month: 'Agu', Pendapatan: 115000000, Pemesanan: 380 }
];

const CATEGORY_CHART_DATA = [
  { name: 'Bahari (Wisata Bahari)', value: 45, color: '#10b981' },
  { name: 'Gunung & Alam', value: 30, color: '#0284c7' },
  { name: 'Budaya & Sejarah', value: 15, color: '#f59e0b' },
  { name: 'Kuliner & Kota', value: 10, color: '#ec4899' }
];

export default function App() {
  // Theme & Layout state
  const [theme, setTheme] = useState<ThemeMode>('light');
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(true);
  const [currentRoute, setCurrentRoute] = useState<'dashboard' | 'users' | 'destinations' | 'bookings' | 'settings'>('dashboard');
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(true);

  // Authentication State
  const [loginEmail, setLoginEmail] = useState('admin@kacarubook.com');
  const [loginPassword, setLoginPassword] = useState('12345678');
  const [showPassword, setShowPassword] = useState(false);
  const [selectedRole, setSelectedRole] = useState<Role>('Admin');
  const [rememberMe, setRememberMe] = useState(true);
  const [isLoggingIn, setIsLoggingIn] = useState(false);

  // Data Store
  const [users, setUsers] = useState<UserItem[]>(INITIAL_USERS);
  const [destinations, setDestinations] = useState<DestinationItem[]>(INITIAL_DESTINATIONS);
  const [bookings, setBookings] = useState<BookingItem[]>(INITIAL_BOOKINGS);

  // User Management State
  const [userSearch, setUserSearch] = useState('');
  const [userRoleFilter, setUserRoleFilter] = useState<string>('ALL');
  const [userStatusFilter, setUserStatusFilter] = useState<string>('ALL');
  const [isUserModalOpen, setIsUserModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<UserItem | null>(null);

  // User Form State
  const [formName, setFormName] = useState('');
  const [formEmail, setFormEmail] = useState('');
  const [formRole, setFormRole] = useState<Role>('Staff');
  const [formStatus, setFormStatus] = useState<'Aktif' | 'Nonaktif' | 'Pending'>('Aktif');

  // Dark mode effect toggler
  useEffect(() => {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme]);

  // Login handler with smooth spinner
  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoggingIn(true);
    setTimeout(() => {
      setIsLoggingIn(false);
      setIsAuthenticated(true);
    }, 700);
  };

  // Filtered Users computation
  const filteredUsers = useMemo(() => {
    return users.filter((u) => {
      const matchSearch = u.name.toLowerCase().includes(userSearch.toLowerCase()) || u.email.toLowerCase().includes(userSearch.toLowerCase());
      const matchRole = userRoleFilter === 'ALL' || u.role === userRoleFilter;
      const matchStatus = userStatusFilter === 'ALL' || u.status === userStatusFilter;
      return matchSearch && matchRole && matchStatus;
    });
  }, [users, userSearch, userRoleFilter, userStatusFilter]);

  // User Modal CRUD Handlers
  const handleOpenUserModal = (user?: UserItem) => {
    if (user) {
      setEditingUser(user);
      setFormName(user.name);
      setFormEmail(user.email);
      setFormRole(user.role);
      setFormStatus(user.status);
    } else {
      setEditingUser(null);
      setFormName('');
      setFormEmail('');
      setFormRole('Staff');
      setFormStatus('Aktif');
    }
    setIsUserModalOpen(true);
  };

  const handleSaveUser = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formName || !formEmail) return;

    if (editingUser) {
      setUsers((prev) =>
        prev.map((u) =>
          u.id === editingUser.id
            ? { ...u, name: formName, email: formEmail, role: formRole, status: formStatus }
            : u
        )
      );
    } else {
      const newUser: UserItem = {
        id: `USR-00${users.length + 1}`,
        name: formName,
        email: formEmail,
        role: formRole,
        status: formStatus,
        avatar: `https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&q=80&w=150`,
        joinedDate: new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
      };
      setUsers((prev) => [newUser, ...prev]);
    }
    setIsUserModalOpen(false);
  };

  const handleDeleteUser = (id: string) => {
    if (window.confirm('Apakah Anda yakin ingin menghapus pengguna ini dari sistem?')) {
      setUsers((prev) => prev.filter((u) => u.id !== id));
    }
  };

  const formatIDR = (amount: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(amount);
  };

  if (!isAuthenticated) {
    return (
      <div className={`min-h-screen transition-colors duration-300 font-sans flex items-center justify-center p-4 ${theme === 'dark' ? 'bg-slate-950 text-white' : 'bg-slate-100 text-slate-900'}`}>
        {/* Floating Theme Toggle in Login */}
        <div className="absolute top-6 right-6 z-20">
          <button
            onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
            className="p-3 rounded-2xl bg-white/80 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-800 shadow-lg backdrop-blur-md hover:scale-105 transition-all text-slate-700 dark:text-slate-200"
            title="Ganti Tema"
          >
            {theme === 'light' ? <Moon className="w-5 h-5 text-indigo-600" /> : <Sun className="w-5 h-5 text-amber-400" />}
          </button>
        </div>

        {/* Split Login Container */}
        <div className="w-full max-w-5xl bg-white dark:bg-slate-900 rounded-3xl shadow-2xl overflow-hidden border border-slate-200/80 dark:border-slate-800 grid grid-cols-1 lg:grid-cols-12 min-h-[580px] transition-all">
          
          {/* Left Decorative Branding Section */}
          <div className="lg:col-span-6 relative p-8 lg:p-12 flex flex-col justify-between overflow-hidden bg-gradient-to-br from-emerald-700 via-teal-800 to-emerald-900 text-white">
            {/* Background Image Layer */}
            <div 
              className="absolute inset-0 opacity-20 mix-blend-overlay bg-cover bg-center pointer-events-none transform scale-105 transition-transform duration-1000"
              style={{ backgroundImage: `url('https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&q=80&w=1200')` }}
            />
            
            {/* Subtle Gradient Overlay */}
            <div className="absolute inset-0 bg-gradient-to-t from-emerald-950/80 via-teal-900/40 to-transparent" />

            <div className="relative z-10 space-y-6">
              {/* Brand Logo Header */}
              <div className="flex items-center gap-3">
                <div className="p-3 rounded-2xl bg-white/10 backdrop-blur-md border border-white/20 shadow-inner group">
                  <Compass className="w-8 h-8 text-emerald-300 animate-spin-slow group-hover:rotate-45 transition-transform duration-500" />
                </div>
                <div>
                  <h2 className="text-2xl font-black tracking-tight text-white">Kacaru Book</h2>
                  <p className="text-[10px] uppercase font-extrabold tracking-widest text-emerald-300">TRAVEL & TOURISM ADMIN</p>
                </div>
              </div>

              {/* Tagline Headline */}
              <div className="pt-6 space-y-3">
                <h1 className="text-3xl lg:text-4xl font-extrabold leading-tight text-white drop-shadow-md">
                  Kelola Destinasi & Pemesanan Wisata Indonesia.
                </h1>
                <p className="text-sm font-medium text-emerald-100/90 leading-relaxed max-w-md">
                  Platform modern manajemen pemesanan paket wisata, pengguna, serta analitik bisnis pariwisata terpadu.
                </p>
              </div>
            </div>

            {/* Social Proof Footer */}
            <div className="relative z-10 pt-8 border-t border-emerald-500/30 flex items-center gap-4">
              <div className="flex -space-x-3 overflow-hidden">
                <img className="inline-block h-9 w-9 rounded-full ring-2 ring-emerald-400 object-cover" src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=150" alt="User 1" />
                <img className="inline-block h-9 w-9 rounded-full ring-2 ring-emerald-400 object-cover" src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=150" alt="User 2" />
                <img className="inline-block h-9 w-9 rounded-full ring-2 ring-emerald-400 object-cover" src="https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&q=80&w=150" alt="User 3" />
              </div>
              <p className="text-xs font-semibold text-emerald-100">
                Dipercaya oleh <span className="font-extrabold text-white">500+</span> agen perjalanan lokal
              </p>
            </div>
          </div>

          {/* Right Form Login Section */}
          <div className="lg:col-span-6 p-8 lg:p-12 flex flex-col justify-between bg-white dark:bg-slate-900 transition-colors">
            <div>
              {/* Heading */}
              <div className="mb-8">
                <h2 className="text-2xl lg:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">
                  Selamat Datang Kembali!
                </h2>
                <p className="text-xs lg:text-sm font-bold text-slate-600 dark:text-slate-300 mt-1">
                  Silakan masuk ke akun manajemen Kacaru Book Anda.
                </p>
              </div>

              {/* Login Form */}
              <form onSubmit={handleLogin} className="space-y-5">
                {/* Email Input */}
                <div>
                  <label className="block text-xs font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
                    EMAIL
                  </label>
                  <div className="relative">
                    <input
                      type="email"
                      required
                      value={loginEmail}
                      onChange={(e) => setLoginEmail(e.target.value)}
                      placeholder="admin@kacarubook.com"
                      className="w-full px-4 py-3.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50/50 dark:bg-slate-800/80 text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-slate-500 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-all"
                    />
                  </div>
                </div>

                {/* Password Input */}
                <div>
                  <label className="block text-xs font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
                    KATA SANDI
                  </label>
                  <div className="relative">
                    <input
                      type={showPassword ? 'text' : 'password'}
                      required
                      value={loginPassword}
                      onChange={(e) => setLoginPassword(e.target.value)}
                      placeholder="••••••••••••"
                      className="w-full px-4 py-3.5 pr-12 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50/50 dark:bg-slate-800/80 text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-slate-500 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-all"
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 dark:text-slate-300 hover:text-slate-700 dark:hover:text-white transition-colors"
                    >
                      {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                    </button>
                  </div>
                </div>

                {/* Role Switcher Pills */}
                <div>
                  <label className="block text-xs font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
                    AKSES PERAN (DEMO)
                  </label>
                  <div className="grid grid-cols-3 gap-2">
                    {(['Admin', 'Staff', 'Customer'] as Role[]).map((r) => (
                      <button
                        key={r}
                        type="button"
                        onClick={() => setSelectedRole(r)}
                        className={`py-2 px-3 rounded-xl text-xs font-extrabold transition-all duration-200 border ${
                          selectedRole === r
                            ? 'bg-emerald-600 text-white border-emerald-600 shadow-md shadow-emerald-600/30'
                            : 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-700 hover:bg-slate-200 dark:hover:bg-slate-700'
                        }`}
                      >
                        {r}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Remember Me & Forgot Password Bar */}
                <div className="flex items-center justify-between pt-1 text-xs">
                  <label className="flex items-center gap-2 cursor-pointer select-none">
                    <div
                      onClick={() => setRememberMe(!rememberMe)}
                      className={`w-4 h-4 rounded border flex items-center justify-center transition-all ${
                        rememberMe ? 'bg-emerald-600 border-emerald-600 text-white' : 'border-slate-400 dark:border-slate-600 bg-transparent'
                      }`}
                    >
                      {rememberMe && <Check className="w-3 h-3 stroke-[3]" />}
                    </div>
                    <span className="font-bold text-slate-700 dark:text-slate-300">Ingat Saya</span>
                  </label>
                  <a href="#" onClick={(e) => e.preventDefault()} className="font-extrabold text-emerald-600 dark:text-emerald-400 hover:underline">
                    Lupa Kata Sandi?
                  </a>
                </div>

                {/* Submit Button */}
                <button
                  type="submit"
                  disabled={isLoggingIn}
                  className="w-full mt-2 py-3.5 px-6 rounded-xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white font-extrabold text-sm shadow-lg shadow-emerald-600/30 hover:shadow-emerald-600/50 active:scale-[0.99] transition-all flex items-center justify-center gap-2"
                >
                  {isLoggingIn ? (
                    <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  ) : (
                    <>
                      <span>Masuk ke Portal</span>
                      <ChevronRight className="w-4 h-4" />
                    </>
                  )}
                </button>
              </form>
            </div>

            {/* Login Footer Note */}
            <div className="mt-8 text-center border-t border-slate-100 dark:border-slate-800 pt-4">
              <p className="text-[11px] font-bold text-slate-500 dark:text-slate-400">
                Pengembangan internal © 2026 Kacaru Book Platform. All rights reserved.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen font-sans transition-colors duration-300 flex bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-white`}>
      
      {/* MOBILE BACKDROP */}
      {isSidebarOpen && (
        <div
          onClick={() => setIsSidebarOpen(false)}
          className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-30 lg:hidden"
        />
      )}

      {/* SIDEBAR NAVIGATION */}
      <aside
        className={`fixed lg:static inset-y-0 left-0 z-40 w-72 bg-white dark:bg-slate-900 border-r border-slate-200/80 dark:border-slate-800 flex flex-col justify-between transition-transform duration-300 ${
          isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        }`}
      >
        <div>
          {/* Sidebar Header Logo */}
          <div className="p-6 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="p-2.5 rounded-2xl bg-gradient-to-tr from-emerald-600 to-teal-500 text-white shadow-md shadow-emerald-600/30">
                <Compass className="w-6 h-6 animate-pulse" />
              </div>
              <div>
                <h1 className="text-lg font-black tracking-tight bg-gradient-to-r from-emerald-600 via-teal-600 to-cyan-600 bg-clip-text text-transparent">
                  Kacaru Book
                </h1>
                <p className="text-[10px] font-extrabold text-slate-600 dark:text-slate-200 uppercase tracking-widest">
                  ADMIN PORTAL
                </p>
              </div>
            </div>
            <button
              onClick={() => setIsSidebarOpen(false)}
              className="lg:hidden p-1.5 rounded-lg text-slate-500 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          {/* Navigation Links with Varied Colors & Interactive Micro-Animations */}
          <nav className="p-4 space-y-2">
            {[
              {
                route: 'dashboard',
                label: 'Dashboard Overview',
                icon: LayoutDashboard,
                badge: 'Live',
                color: 'from-emerald-500 to-teal-500',
                iconColor: 'text-emerald-500 dark:text-emerald-400',
                activeBg: 'bg-gradient-to-r from-emerald-600 to-teal-600 text-white shadow-lg shadow-emerald-600/25'
              },
              {
                route: 'users',
                label: 'Pengguna System',
                icon: Users,
                badge: `${users.length} User`,
                color: 'from-indigo-500 to-purple-500',
                iconColor: 'text-indigo-500 dark:text-indigo-400',
                activeBg: 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-lg shadow-indigo-600/25'
              },
              {
                route: 'destinations',
                label: 'Katalog Wisata',
                icon: Compass,
                badge: `${destinations.length} Lokasi`,
                color: 'from-cyan-500 to-blue-500',
                iconColor: 'text-cyan-500 dark:text-cyan-400',
                activeBg: 'bg-gradient-to-r from-cyan-600 to-blue-600 text-white shadow-lg shadow-cyan-600/25'
              },
              {
                route: 'bookings',
                label: 'Pemesanan Tiket',
                icon: Ticket,
                badge: 'Baru',
                color: 'from-amber-500 to-orange-500',
                iconColor: 'text-amber-500 dark:text-amber-400',
                activeBg: 'bg-gradient-to-r from-amber-500 to-orange-600 text-white shadow-lg shadow-amber-600/25'
              },
              {
                route: 'settings',
                label: 'Pengaturan Akun',
                icon: Settings,
                color: 'from-rose-500 to-pink-500',
                iconColor: 'text-rose-500 dark:text-rose-400',
                activeBg: 'bg-gradient-to-r from-rose-600 to-pink-600 text-white shadow-lg shadow-rose-600/25'
              }
            ].map((item) => {
              const Icon = item.icon;
              const isActive = currentRoute === item.route;
              return (
                <button
                  key={item.route}
                  onClick={() => {
                    setCurrentRoute(item.route as any);
                    if (window.innerWidth < 1024) setIsSidebarOpen(false);
                  }}
                  className={`group relative w-full flex items-center justify-between px-3.5 py-3 rounded-2xl font-bold text-xs transition-all duration-300 transform active:scale-95 ${
                    isActive
                      ? item.activeBg
                      : theme === 'dark'
                      ? 'text-slate-100 hover:bg-slate-800/90 hover:text-white'
                      : 'text-slate-700 hover:bg-slate-100/90 hover:text-slate-900'
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <div className={`p-2 rounded-xl transition-all duration-300 group-hover:scale-110 group-hover:rotate-6 ${
                      isActive
                        ? 'bg-white/20 text-white shadow-inner'
                        : theme === 'dark'
                        ? 'bg-slate-800 text-slate-100 group-hover:bg-slate-700 group-hover:text-white'
                        : 'bg-slate-100 text-slate-600 group-hover:bg-white shadow-sm'
                    }`}>
                      <Icon className={`w-4 h-4 transition-all duration-300 ${isActive ? 'text-white' : item.iconColor}`} />
                    </div>
                    <span className="tracking-wide text-slate-800 dark:text-slate-100 group-hover:dark:text-white font-extrabold">{item.label}</span>
                  </div>

                  <div className="flex items-center gap-2">
                    {item.badge && !isActive && (
                      <span className={`text-[10px] font-black px-2 py-0.5 rounded-full bg-gradient-to-r ${item.color} text-white shadow-sm transition-transform duration-300 group-hover:scale-105`}>
                        {item.badge}
                      </span>
                    )}
                    {isActive ? (
                      <span className="w-2 h-2 rounded-full bg-white animate-ping" />
                    ) : (
                      <ChevronRight className="w-3.5 h-3.5 text-slate-400 dark:text-slate-200 opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300" />
                    )}
                  </div>
                </button>
              );
            })}
          </nav>
        </div>

        {/* Sidebar Footer User Info */}
        <div className="p-4 border-t border-slate-100 dark:border-slate-800">
          <div className="p-3 rounded-xl bg-slate-50 dark:bg-slate-800/80 flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              <img
                src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=150"
                alt="Profile"
                className="w-8 h-8 rounded-full object-cover ring-2 ring-emerald-500"
              />
              <div className="flex flex-col overflow-hidden">
                <span className="text-xs font-bold text-slate-800 dark:text-slate-100 truncate">Rizky Admin</span>
                <span className="text-[10px] text-emerald-600 dark:text-emerald-400 font-extrabold uppercase">{selectedRole}</span>
              </div>
            </div>
            <button
              onClick={() => setIsAuthenticated(false)}
              title="Keluar"
              className="p-1.5 text-slate-400 dark:text-slate-300 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/50 rounded-lg transition-colors"
            >
              <LogOut className="w-4 h-4" />
            </button>
          </div>
        </div>
      </aside>

      {/* MAIN CONTENT AREA */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* TOP HEADER */}
        <header className="sticky top-0 z-20 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200/80 dark:border-slate-800 px-4 lg:px-8 py-3.5 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setIsSidebarOpen(true)}
              className="lg:hidden p-2 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-300"
            >
              <Menu className="w-5 h-5" />
            </button>
            <div>
              <h2 className="text-base font-extrabold text-slate-900 dark:text-white capitalize">
                {currentRoute === 'dashboard' && 'Dashboard Analitik'}
                {currentRoute === 'users' && 'Manajemen Pengguna'}
                {currentRoute === 'destinations' && 'Katalog Destinasi Wisata'}
                {currentRoute === 'bookings' && 'Pemesanan & Tiket'}
                {currentRoute === 'settings' && 'Pengaturan Akun'}
              </h2>
              <p className="text-[11px] font-medium text-slate-500 dark:text-slate-300 hidden sm:block">
                Sistem Terpadu Wisata Indonesia • Kacaru Book Admin
              </p>
            </div>
          </div>

          {/* Top Actions */}
          <div className="flex items-center gap-2 sm:gap-3">
            {/* Theme Toggle Button */}
            <button
              onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
              className="p-2.5 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all"
              title="Ganti Mode Tampilan"
            >
              {theme === 'light' ? <Moon className="w-4 h-4 text-indigo-600" /> : <Sun className="w-4 h-4 text-amber-400" />}
            </button>

            {/* Notification Icon */}
            <button className="relative p-2.5 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all">
              <Bell className="w-4 h-4" />
              <span className="absolute top-2 right-2 w-2 h-2 bg-emerald-500 rounded-full ring-2 ring-white dark:ring-slate-900" />
            </button>

            <div className="h-6 w-px bg-slate-200 dark:bg-slate-800 my-auto hidden sm:block" />

            {/* Quick Status Pill */}
            <div className="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-50 dark:bg-emerald-950/60 border border-emerald-200 dark:border-emerald-800 text-xs font-bold text-emerald-700 dark:text-emerald-300">
              <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              System Active
            </div>
          </div>
        </header>

        {/* ROUTE PAGES BODY */}
        <main className="flex-1 p-4 lg:p-8 space-y-6 max-w-7xl w-full mx-auto">
          
          {/* 1. DASHBOARD VIEW */}
          {currentRoute === 'dashboard' && (
            <div className="space-y-6">
              
              {/* Vibrant KPI Cards Grid */}
              {}
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                
                {/* KPI Card 1: Revenue (Primary Emerald Gradient Accent) */}
                <div className="p-5 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-shadow relative overflow-hidden group">
                  <div className="absolute top-0 right-0 w-24 h-24 bg-emerald-500/10 rounded-bl-full pointer-events-none group-hover:scale-110 transition-transform" />
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs font-extrabold text-slate-700 dark:text-slate-200 uppercase tracking-wider">Total Pendapatan</span>
                    <div className="p-2.5 rounded-2xl bg-emerald-500/15 text-emerald-700 dark:text-emerald-300">
                      <DollarSign className="w-5 h-5" />
                    </div>
                  </div>
                  <div className="text-2xl font-black text-slate-900 dark:text-white">IDR 693M</div>
                  <div className="mt-2 flex items-center gap-1.5 text-xs font-extrabold text-emerald-700 dark:text-emerald-300">
                    <TrendingUp className="w-3.5 h-3.5" />
                    <span>+18.4% bulan ini</span>
                  </div>
                </div>

                {/* KPI Card 2: Bookings (Indigo/Violet Accent) */}
                <div className="p-5 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-shadow relative overflow-hidden group">
                  <div className="absolute top-0 right-0 w-24 h-24 bg-indigo-500/10 rounded-bl-full pointer-events-none group-hover:scale-110 transition-transform" />
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs font-extrabold text-slate-700 dark:text-slate-200 uppercase tracking-wider">Total Tiket Dipesan</span>
                    <div className="p-2.5 rounded-2xl bg-indigo-500/15 text-indigo-700 dark:text-indigo-300">
                      <Ticket className="w-5 h-5" />
                    </div>
                  </div>
                  <div className="text-2xl font-black text-slate-900 dark:text-white">2,090 Tiket</div>
                  <div className="mt-2 flex items-center gap-1.5 text-xs font-extrabold text-indigo-700 dark:text-indigo-300">
                    <TrendingUp className="w-3.5 h-3.5" />
                    <span>+12.1% dari target</span>
                  </div>
                </div>

                {/* KPI Card 3: Active Users (Sky Blue Accent) */}
                <div className="p-5 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-shadow relative overflow-hidden group">
                  <div className="absolute top-0 right-0 w-24 h-24 bg-sky-500/10 rounded-bl-full pointer-events-none group-hover:scale-110 transition-transform" />
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs font-extrabold text-slate-700 dark:text-slate-200 uppercase tracking-wider">Pengguna Terdaftar</span>
                    <div className="p-2.5 rounded-2xl bg-sky-500/15 text-sky-700 dark:text-sky-300">
                      <Users className="w-5 h-5" />
                    </div>
                  </div>
                  <div className="text-2xl font-black text-slate-900 dark:text-white">{users.length} Akun</div>
                  <div className="mt-2 flex items-center gap-1.5 text-xs font-extrabold text-sky-700 dark:text-sky-300">
                    <UserCheck className="w-3.5 h-3.5" />
                    <span>89% Akun Aktif</span>
                  </div>
                </div>

                {/* KPI Card 4: Destinations & Rating (Amber/Gold Accent) */}
                <div className="p-5 rounded-3xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-shadow relative overflow-hidden group">
                  <div className="absolute top-0 right-0 w-24 h-24 bg-amber-500/10 rounded-bl-full pointer-events-none group-hover:scale-110 transition-transform" />
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-xs font-extrabold text-slate-700 dark:text-slate-200 uppercase tracking-wider">Destinasi Favorit</span>
                    <div className="p-2.5 rounded-2xl bg-amber-500/15 text-amber-700 dark:text-amber-300">
                      <Star className="w-5 h-5 fill-amber-500" />
                    </div>
                  </div>
                  <div className="text-2xl font-black text-slate-900 dark:text-white">4.88 / 5.0</div>
                  <div className="mt-2 flex items-center gap-1.5 text-xs font-extrabold text-amber-700 dark:text-amber-300">
                    <Compass className="w-3.5 h-3.5" />
                    <span>{destinations.length} Lokasi Terdaftar</span>
                  </div>
                </div>
              </div>

              {/* Charts Grid */}
              {}
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                
                {/* Main Area Chart */}
                <div className="lg:col-span-2 bg-white dark:bg-slate-900 rounded-3xl p-6 border border-slate-200/80 dark:border-slate-800 shadow-sm">
                  <div className="flex items-center justify-between mb-6">
                    <div>
                      <h3 className="text-base font-bold text-slate-900 dark:text-white">Tren Pendapatan Bulanan</h3>
                      <p className="text-xs text-slate-600 dark:text-slate-200 font-medium">Statistik performa penjualan paket wisata 2026</p>
                    </div>
                    <span className="text-xs font-bold text-emerald-700 dark:text-emerald-300 bg-emerald-100 dark:bg-emerald-950/80 px-3 py-1 rounded-full border border-emerald-300/40 dark:border-emerald-700/50">
                      IDR (Rupiah)
                    </span>
                  </div>

                  <div className="h-72 w-full">
                    <ResponsiveContainer width="100%" height="100%">
                      <AreaChart data={REVENUE_CHART_DATA} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
                        <defs>
                          <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#10b981" stopOpacity={0.4} />
                            <stop offset="95%" stopColor="#10b981" stopOpacity={0.0} />
                          </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" stroke={theme === 'dark' ? '#334155' : '#e2e8f0'} vertical={false} />
                        <XAxis dataKey="month" stroke={theme === 'dark' ? '#94a3b8' : '#64748b'} fontSize={12} tickLine={false} />
                        <YAxis stroke={theme === 'dark' ? '#94a3b8' : '#64748b'} fontSize={12} tickLine={false} tickFormatter={(v) => `${v / 1000000}M`} />
                        <Tooltip
                          contentStyle={{
                            backgroundColor: theme === 'dark' ? '#0f172a' : '#ffffff',
                            borderColor: theme === 'dark' ? '#334155' : '#e2e8f0',
                            borderRadius: '16px',
                            color: theme === 'dark' ? '#ffffff' : '#0f172a',
                            fontWeight: 'bold'
                          }}
                          formatter={(value: any) => [formatIDR(value), 'Pendapatan']}
                        />
                        <Area type="monotone" dataKey="Pendapatan" stroke="#10b981" strokeWidth={3} fillOpacity={1} fill="url(#colorRevenue)" />
                      </AreaChart>
                    </ResponsiveContainer>
                  </div>
                </div>

                {/* Pie Chart Category Breakdown */}
                <div className="bg-white dark:bg-slate-900 rounded-3xl p-6 border border-slate-200/80 dark:border-slate-800 shadow-sm flex flex-col justify-between">
                  <div>
                    <h3 className="text-base font-bold text-slate-900 dark:text-white">Kategori Terpopuler</h3>
                    <p className="text-xs text-slate-600 dark:text-slate-200 font-medium">Distribusi minat wisatawan</p>
                  </div>

                  <div className="h-56 w-full my-auto">
                    <ResponsiveContainer width="100%" height="100%">
                      <PieChart>
                        <Pie
                          data={CATEGORY_CHART_DATA}
                          cx="50%"
                          cy="50%"
                          innerRadius={55}
                          outerRadius={80}
                          paddingAngle={4}
                          dataKey="value"
                        >
                          {CATEGORY_CHART_DATA.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={entry.color} />
                          ))}
                        </Pie>
                        <Tooltip />
                      </PieChart>
                    </ResponsiveContainer>
                  </div>

                  <div className="grid grid-cols-2 gap-2 pt-2 border-t border-slate-100 dark:border-slate-800 text-[11px]">
                    {CATEGORY_CHART_DATA.map((c) => (
                      <div key={c.name} className="flex items-center gap-2">
                        <span className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: c.color }} />
                        <span className="font-bold text-slate-700 dark:text-slate-200 truncate">{c.name}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Recent Bookings Table Preview */}
              <div className="bg-white dark:bg-slate-900 rounded-3xl p-6 border border-slate-200/80 dark:border-slate-800 shadow-sm">
                <div className="flex items-center justify-between mb-4">
                  <div>
                    <h3 className="text-base font-bold text-slate-900 dark:text-white">Pemesanan Terbaru</h3>
                    <p className="text-xs text-slate-600 dark:text-slate-200 font-medium">Tiket yang baru saja dipesan oleh pengguna</p>
                  </div>
                  <button
                    onClick={() => setCurrentRoute('bookings')}
                    className="text-xs font-bold text-emerald-700 dark:text-emerald-400 hover:underline flex items-center gap-1"
                  >
                    Lihat Semua <ArrowUpRight className="w-3.5 h-3.5" />
                  </button>
                </div>

                <div className="overflow-x-auto">
                  <table className="w-full text-left border-collapse">
                    <thead>
                      <tr className="border-b border-slate-200 dark:border-slate-800 text-[11px] font-extrabold text-slate-700 dark:text-slate-200 uppercase tracking-wider">
                        <th className="py-3 px-4">ID Booking</th>
                        <th className="py-3 px-4">Pemesan</th>
                        <th className="py-3 px-4">Destinasi</th>
                        <th className="py-3 px-4">Tanggal</th>
                        <th className="py-3 px-4">Total</th>
                        <th className="py-3 px-4">Status</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100 dark:divide-slate-800 text-xs font-semibold">
                      {bookings.slice(0, 4).map((b) => (
                        <tr key={b.id} className="hover:bg-slate-50/50 dark:hover:bg-slate-800/50 transition-colors">
                          <td className="py-3 px-4 font-bold text-emerald-700 dark:text-emerald-400">{b.id}</td>
                          <td className="py-3 px-4 font-extrabold text-slate-900 dark:text-slate-100">{b.customerName}</td>
                          <td className="py-3 px-4 text-slate-700 dark:text-slate-200">{b.destinationName}</td>
                          <td className="py-3 px-4 text-slate-600 dark:text-slate-300">{b.date}</td>
                          <td className="py-3 px-4 font-extrabold text-slate-900 dark:text-slate-100">{formatIDR(b.totalPrice)}</td>
                          <td className="py-3 px-4">
                            <span className={`px-2.5 py-1 rounded-full text-[10px] font-extrabold ${
                              b.status === 'Lunas'
                                ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-700'
                                : b.status === 'Menunggu'
                                ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/80 dark:text-amber-300 border border-amber-300 dark:border-amber-700'
                                : 'bg-rose-100 text-rose-800 dark:bg-rose-950/80 dark:text-rose-300 border border-rose-300 dark:border-rose-700'
                            }`}>
                              {b.status}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

          {/* 2. USER MANAGEMENT VIEW */}
          {}
          {currentRoute === 'users' && (
            <div className="space-y-6">
              
              {/* Header & Controls */}
              <div className="bg-white dark:bg-slate-900 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
                <div>
                  <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">Kelola Pengguna Sistem</h3>
                  <p className="text-xs font-bold text-slate-600 dark:text-slate-200 mt-0.5">Atur hak akses, tambah staf baru, dan status aktivasi pengguna.</p>
                </div>

                <button
                  onClick={() => handleOpenUserModal()}
                  className="py-3 px-5 rounded-2xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white font-extrabold text-xs shadow-lg shadow-emerald-600/20 active:scale-95 transition-all flex items-center justify-center gap-2"
                >
                  <Plus className="w-4 h-4" />
                  <span>Tambah User Baru</span>
                </button>
              </div>

              {/* Filters Bar */}
              <div className="grid grid-cols-1 md:grid-cols-12 gap-3">
                {/* Search Bar */}
                <div className="md:col-span-6 relative">
                  <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-500 dark:text-slate-300" />
                  <input
                    type="text"
                    value={userSearch}
                    onChange={(e) => setUserSearch(e.target.value)}
                    placeholder="Cari nama atau email pengguna..."
                    className="w-full pl-10 pr-4 py-2.5 rounded-2xl border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white placeholder-slate-400 dark:placeholder-slate-500 text-xs font-semibold focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  />
                </div>

                {/* Role Filter */}
                <div className="md:col-span-3">
                  <select
                    value={userRoleFilter}
                    onChange={(e) => setUserRoleFilter(e.target.value)}
                    className="w-full px-3 py-2.5 rounded-2xl border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white text-xs font-extrabold focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="ALL">Semua Peran (Role)</option>
                    <option value="Admin">Admin</option>
                    <option value="Staff">Staff</option>
                    <option value="Customer">Customer</option>
                  </select>
                </div>

                {/* Status Filter */}
                <div className="md:col-span-3">
                  <select
                    value={userStatusFilter}
                    onChange={(e) => setUserStatusFilter(e.target.value)}
                    className="w-full px-3 py-2.5 rounded-2xl border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white text-xs font-extrabold focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="ALL">Semua Status</option>
                    <option value="Aktif">Aktif</option>
                    <option value="Pending">Pending</option>
                    <option value="Nonaktif">Nonaktif</option>
                  </select>
                </div>
              </div>

              {/* Users Data Table */}
              <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full text-left border-collapse">
                    <thead>
                      <tr className="bg-slate-50/80 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-800 text-[11px] font-extrabold text-slate-700 dark:text-slate-200 uppercase tracking-wider">
                        <th className="py-4 px-6">Pengguna</th>
                        <th className="py-4 px-4">Peran (Role)</th>
                        <th className="py-4 px-4">Status</th>
                        <th className="py-4 px-4">Tanggal Bergabung</th>
                        <th className="py-4 px-6 text-right">Aksi</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100 dark:divide-slate-800 text-xs font-medium">
                      {filteredUsers.length === 0 ? (
                        <tr>
                          <td colSpan={5} className="py-8 text-center text-slate-600 dark:text-slate-200 font-bold">
                            Tidak ada pengguna yang cocok dengan kriteria pencarian.
                          </td>
                        </tr>
                      ) : (
                        filteredUsers.map((u) => (
                          <tr key={u.id} className="hover:bg-slate-50/50 dark:hover:bg-slate-800/50 transition-colors">
                            {/* User Avatar & Info */}
                            <td className="py-4 px-6">
                              <div className="flex items-center gap-3">
                                <img
                                  src={u.avatar}
                                  alt={u.name}
                                  className="w-10 h-10 rounded-full object-cover ring-2 ring-emerald-500/30"
                                />
                                <div>
                                  <div className="font-extrabold text-slate-900 dark:text-white text-sm">{u.name}</div>
                                  <div className="text-slate-700 dark:text-slate-300 font-semibold text-[11px]">{u.email}</div>
                                </div>
                              </div>
                            </td>

                            {/* Role Badge */}
                            <td className="py-4 px-4">
                              <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-extrabold ${
                                u.role === 'Admin'
                                  ? 'bg-purple-100 text-purple-800 dark:bg-purple-950/80 dark:text-purple-300 border border-purple-300 dark:border-purple-700'
                                  : u.role === 'Staff'
                                  ? 'bg-blue-100 text-blue-800 dark:bg-blue-950/80 dark:text-blue-300 border border-blue-300 dark:border-blue-700'
                                  : 'bg-slate-100 text-slate-800 dark:bg-slate-800 dark:text-slate-300 border border-slate-300 dark:border-slate-700'
                              }`}>
                                <Shield className="w-3 h-3" />
                                {u.role}
                              </span>
                            </td>

                            {/* Status Badge */}
                            <td className="py-4 px-4">
                              <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-extrabold ${
                                u.status === 'Aktif'
                                  ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-700'
                                  : u.status === 'Pending'
                                  ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/80 dark:text-amber-300 border border-amber-300 dark:border-amber-700'
                                  : 'bg-rose-100 text-rose-800 dark:bg-rose-950/80 dark:text-rose-300 border border-rose-300 dark:border-rose-700'
                              }`}>
                                {u.status === 'Aktif' && <CheckCircle className="w-3 h-3" />}
                                {u.status === 'Pending' && <Clock className="w-3 h-3" />}
                                {u.status === 'Nonaktif' && <XCircle className="w-3 h-3" />}
                                {u.status}
                              </span>
                            </td>

                            {/* Date */}
                            <td className="py-4 px-4 text-slate-700 dark:text-slate-200 font-bold">{u.joinedDate}</td>

                            {/* Actions */}
                            <td className="py-4 px-6 text-right space-x-2">
                              <button
                                onClick={() => handleOpenUserModal(u)}
                                className="p-2 rounded-xl text-slate-600 dark:text-slate-200 hover:text-emerald-600 hover:bg-emerald-50 dark:hover:bg-slate-800 transition-colors"
                                title="Edit Pengguna"
                              >
                                <Edit3 className="w-4 h-4" />
                              </button>
                              <button
                                onClick={() => handleDeleteUser(u.id)}
                                className="p-2 rounded-xl text-slate-600 dark:text-slate-200 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-slate-800 transition-colors"
                                title="Hapus Pengguna"
                              >
                                <Trash2 className="w-4 h-4" />
                              </button>
                            </td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

          {/* 3. DESTINATIONS VIEW */}
          {}
          {currentRoute === 'destinations' && (
            <div className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">Katalog Wisata Populer</h3>
                  <p className="text-xs font-bold text-slate-600 dark:text-slate-200">Daftar paket perjalanan destinasi yang aktif dijual di platform Kacaru Book.</p>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {destinations.map((item) => (
                  <div key={item.id} className="bg-white dark:bg-slate-900 rounded-3xl overflow-hidden border border-slate-200/80 dark:border-slate-800 shadow-sm hover:shadow-md transition-shadow group flex flex-col justify-between">
                    <div>
                      <div className="relative h-44 overflow-hidden">
                        <img
                          src={item.image}
                          alt={item.name}
                          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                        />
                        <div className="absolute top-3 right-3 px-2.5 py-1 rounded-full bg-slate-900/70 backdrop-blur-md text-amber-300 text-xs font-black flex items-center gap-1">
                          <Star className="w-3.5 h-3.5 fill-amber-400 text-amber-400" />
                          <span>{item.rating}</span>
                        </div>
                        <span className="absolute bottom-3 left-3 px-2.5 py-1 rounded-full bg-emerald-600/90 backdrop-blur-md text-white text-[10px] font-extrabold">
                          {item.category}
                        </span>
                      </div>

                      <div className="p-4 space-y-2">
                        <h4 className="font-extrabold text-slate-900 dark:text-white text-base leading-snug">{item.name}</h4>
                        <div className="flex items-center gap-1 text-xs font-bold text-slate-600 dark:text-slate-300">
                          <MapPin className="w-3.5 h-3.5 text-rose-500" />
                          <span>{item.location}</span>
                        </div>
                      </div>
                    </div>

                    <div className="p-4 pt-0 border-t border-slate-100 dark:border-slate-800/80 mt-2 flex items-center justify-between">
                      <div>
                        <span className="text-[10px] text-slate-600 dark:text-slate-300 font-extrabold block">Mulai dari</span>
                        <span className="text-sm font-black text-emerald-700 dark:text-emerald-400">{formatIDR(item.price)}</span>
                      </div>
                      <span className={`px-2.5 py-1 rounded-full text-[10px] font-extrabold ${
                        item.status === 'Tersedia' ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-300' : 'bg-rose-100 text-rose-800 dark:bg-rose-950 dark:text-rose-300'
                      }`}>
                        {item.status}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* 4. BOOKINGS VIEW */}
          {}
          {currentRoute === 'bookings' && (
            <div className="space-y-6">
              <div className="bg-white dark:bg-slate-900 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm">
                <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">Daftar Transaksi Tiket Wisata</h3>
                <p className="text-xs font-bold text-slate-600 dark:text-slate-200 mt-1">Status verifikasi pembayaran dan keberangkatan turis.</p>
              </div>

              <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full text-left border-collapse">
                    <thead>
                      <tr className="bg-slate-50/80 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-800 text-[11px] font-extrabold text-slate-700 dark:text-slate-200 uppercase tracking-wider">
                        <th className="py-4 px-6">ID Tiket</th>
                        <th className="py-4 px-4">Nama Pelanggan</th>
                        <th className="py-4 px-4">Destinasi</th>
                        <th className="py-4 px-4">Tgl Keberangkatan</th>
                        <th className="py-4 px-4">Total Biaya</th>
                        <th className="py-4 px-6">Status Pembayaran</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100 dark:divide-slate-800 text-xs font-semibold">
                      {bookings.map((b) => (
                        <tr key={b.id} className="hover:bg-slate-50/50 dark:hover:bg-slate-800/50 transition-colors">
                          <td className="py-4 px-6 font-bold text-emerald-700 dark:text-emerald-400">{b.id}</td>
                          <td className="py-4 px-4 font-extrabold text-slate-900 dark:text-white">{b.customerName}</td>
                          <td className="py-4 px-4 text-slate-700 dark:text-slate-200">{b.destinationName}</td>
                          <td className="py-4 px-4 text-slate-600 dark:text-slate-300">{b.date}</td>
                          <td className="py-4 px-4 font-extrabold text-slate-900 dark:text-white">{formatIDR(b.totalPrice)}</td>
                          <td className="py-4 px-6">
                            <span className={`px-3 py-1 rounded-full text-[10px] font-extrabold ${
                              b.status === 'Lunas'
                                ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-700'
                                : b.status === 'Menunggu'
                                ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/80 dark:text-amber-300 border border-amber-300 dark:border-amber-700'
                                : 'bg-rose-100 text-rose-800 dark:bg-rose-950/80 dark:text-rose-300 border border-rose-300 dark:border-rose-700'
                            }`}>
                              {b.status}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

          {/* 5. SETTINGS VIEW */}
          {currentRoute === 'settings' && (
            <div className="max-w-2xl bg-white dark:bg-slate-900 rounded-3xl p-6 border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-6">
              <div>
                <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">Pengaturan Akun Admin</h3>
                <p className="text-xs font-bold text-slate-600 dark:text-slate-200">Perbarui informasi profil operasional Anda.</p>
              </div>

              <div className="space-y-4 text-xs font-semibold">
                <div>
                  <label className="block font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">Nama Administrator</label>
                  <input
                    type="text"
                    defaultValue="Rizky Admin"
                    className="w-full px-4 py-3 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-bold"
                  />
                </div>
                <div>
                  <label className="block font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">Email Operasional</label>
                  <input
                    type="email"
                    defaultValue="admin@kacarubook.com"
                    className="w-full px-4 py-3 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-bold"
                  />
                </div>
                <button
                  onClick={() => alert('Pengaturan profil berhasil disimpan!')}
                  className="py-3 px-6 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-extrabold text-xs shadow-md shadow-emerald-600/20"
                >
                  Simpan Perubahan
                </button>
              </div>
            </div>
          )}

        </main>
      </div>

      {/* USER FORM MODAL (ADD / EDIT) */}
      {}
      {isUserModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm animate-fade-in">
          <div className="bg-white dark:bg-slate-900 rounded-3xl max-w-md w-full border border-slate-200 dark:border-slate-800 shadow-2xl p-6 space-y-6">
            <div className="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4">
              <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
                {editingUser ? 'Edit Data Pengguna' : 'Tambah Pengguna Baru'}
              </h3>
              <button
                onClick={() => setIsUserModalOpen(false)}
                className="p-1 rounded-lg text-slate-400 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleSaveUser} className="space-y-4 text-xs">
              <div>
                <label className="block font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1.5">Nama Lengkap</label>
                <input
                  type="text"
                  required
                  value={formName}
                  onChange={(e) => setFormName(e.target.value)}
                  placeholder="Contoh: Budi Santoso"
                  className="w-full px-3.5 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-semibold focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div>
                <label className="block font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1.5">Alamat Email</label>
                <input
                  type="email"
                  required
                  value={formEmail}
                  onChange={(e) => setFormEmail(e.target.value)}
                  placeholder="budi@gmail.com"
                  className="w-full px-3.5 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-semibold focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1.5">Peran (Role)</label>
                  <select
                    value={formRole}
                    onChange={(e) => setFormRole(e.target.value as Role)}
                    className="w-full px-3 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-extrabold focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="Admin">Admin</option>
                    <option value="Staff">Staff</option>
                    <option value="Customer">Customer</option>
                  </select>
                </div>

                <div>
                  <label className="block font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1.5">Status Akun</label>
                  <select
                    value={formStatus}
                    onChange={(e) => setFormStatus(e.target.value as any)}
                    className="w-full px-3 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-extrabold focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="Aktif">Aktif</option>
                    <option value="Pending">Pending</option>
                    <option value="Nonaktif">Nonaktif</option>
                  </select>
                </div>
              </div>

              <div className="pt-4 flex items-center justify-end gap-2 border-t border-slate-100 dark:border-slate-800">
                <button
                  type="button"
                  onClick={() => setIsUserModalOpen(false)}
                  className="px-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 font-extrabold hover:bg-slate-100 dark:hover:bg-slate-800"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="px-5 py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-extrabold shadow-md shadow-emerald-600/20"
                >
                  Simpan Data
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}