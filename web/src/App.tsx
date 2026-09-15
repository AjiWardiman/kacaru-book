import React, { useState, useEffect, useMemo } from 'react';
import {
  Sun,
  Moon,
  Eye,
  EyeOff
} from 'lucide-react';

import { useAuthStore } from './store/useAuthStore';
import { useNotificationStore } from './store/useNotificationStore';
import { useUserStore } from './store/useUserStore';
import { useDestinationStore } from './store/useDestinationStore';
import { useBookingStore } from './store/useBookingStore';
import { toast } from 'sonner';
import { ToastProvider } from './components/ToastProvider';

// Modular Components
import { Sidebar } from './components/layout/Sidebar';
import { Navbar } from './components/layout/Navbar';
import { DashboardOverview } from './components/dashboard/DashboardOverview';
import { DestinationGrid } from './components/destinations/DestinationGrid';
import { DestinationModal } from './components/destinations/DestinationModal';
import { TicketModal } from './components/bookings/TicketModal';
import { MyBookingsView } from './components/bookings/MyBookingsView';
import { AdminBookingsView } from './components/bookings/AdminBookingsView';
import { UserManagementView } from './components/users/UserManagementView';
import { UserModal } from './components/users/UserModal';
import { SettingsView } from './components/settings/SettingsView';
import { ReportAnalyticsView } from './components/reports/ReportAnalyticsView';
import { OTPVerificationModal } from './components/auth/OTPVerificationModal';

// Constants & Fallbacks


const INITIAL_DESTINATIONS = [
  {
    id: 1,
    name: 'Pantai Lariti',
    location: 'Kecamatan Sape, Kabupaten Bima, NTB',
    price: 15000,
    category: 'Bahari',
    rating: 4.9,
    image: '/uploads/destinations/lariti.jpg',
    status: 'Tersedia',
    description: 'Pantai unik eksotis dengan fenomena laut membelah dan pasir timbul yang memukau saat laut surut.',
  },
  {
    id: 2,
    name: 'Toro Maria',
    location: 'Kecamatan Sape, Kabupaten Bima, NTB',
    price: 20000,
    category: 'Bahari',
    rating: 4.8,
    image: '/uploads/destinations/toro_maria.jpg',
    status: 'Tersedia',
    description: 'Tebing karang memukau dengan lanskap bukit dan pemandangan lautan lepas perbatasan Flores.',
  },
  {
    id: 3,
    name: 'Uma Lengge Wawo',
    location: 'Kecamatan Wawo, Kabupaten Bima, NTB',
    price: 10000,
    category: 'Budaya & Geowisata',
    rating: 4.9,
    image: '/uploads/destinations/uma_lengge.jpg',
    status: 'Tersedia',
    description: 'Situs cagar budaya rumah adat kuno suku Bima (Mbojo) dengan struktur kerucut yang kaya sejarah.',
  },
  {
    id: 4,
    name: 'Pulau Sangiang',
    location: 'Kecamatan Wera, Kabupaten Bima, NTB',
    price: 50000,
    category: 'Gunung & Alam',
    rating: 4.9,
    image: '/uploads/destinations/sangiang.jpg',
    status: 'Tersedia',
    description: 'Pulau gunung berapi aktif di tengah lautan Bima dengan keindahan alam bahari dan savana pesisir.',
  },
  {
    id: 5,
    name: 'Air Terjun Oi Fanda',
    location: 'Kecamatan Ambalawi, Kabupaten Bima, NTB',
    price: 15000,
    category: 'Gunung & Alam',
    rating: 4.7,
    image: '/uploads/destinations/oi_fanda.jpg',
    status: 'Tersedia',
    description: 'Air terjun alami berair jernih segar yang tersembunyi di rimbunnya hutan lereng pegunungan Ambalawi.',
  },
  {
    id: 6,
    name: 'Pantai Pink Lambu',
    location: 'Kecamatan Lambu, Kabupaten Bima, NTB',
    price: 25000,
    category: 'Bahari',
    rating: 4.9,
    image: '/uploads/destinations/pantai_pink.jpg',
    status: 'Tersedia',
    description: 'Hamparan pasir pantai berwarna merah muda alami dengan perairan jernih tenang nan fotogenik.',
  },
  {
    id: 7,
    name: 'Pulau Kelapa',
    location: 'Kecamatan Lambu, Kabupaten Bima, NTB',
    price: 75000,
    category: 'Bahari',
    rating: 4.9,
    image: '/uploads/destinations/pulau_kelapa.jpg',
    status: 'Tersedia',
    description: 'Gugusan pulau eksotis nan sepi dengan lanskap perbukitan hijau dan bukit pandang bahari mempesona.',
  },
  {
    id: 8,
    name: 'Tanjung Meriam',
    location: 'Kecamatan Lambu, Kabupaten Bima, NTB',
    price: 30000,
    category: 'Budaya & Geowisata',
    rating: 4.8,
    image: '/uploads/destinations/tanjung_meriam.jpg',
    status: 'Tersedia',
    description: 'Keajaiban geologi batuan basal tersusun rapi segiempat/segienam alami menyerupai monumen meriam.',
  },
  {
    id: 9,
    name: 'Nisa Ndoko',
    location: 'Kecamatan Langgudu, Kabupaten Bima, NTB',
    price: 40000,
    category: 'Bahari',
    rating: 4.8,
    image: '/uploads/destinations/nisa_ndoko.jpg',
    status: 'Tersedia',
    description: 'Pulau kecil (Nisa) berpasir putih bersih dengan keindahan bawah laut terumbu karang yang terjaga.',
  },
  {
    id: 10,
    name: 'Pantai Kampa',
    location: 'Kecamatan Sape, Kabupaten Bima, NTB',
    price: 20000,
    category: 'Bahari',
    rating: 4.9,
    image: '/uploads/destinations/pantai_kampa.jpg',
    status: 'Tersedia',
    description: 'Pantai teluk tersembunyi berpasir putih jernih dengan air laut tenang dan pemandangan bukit karang eksotis.',
  },
  {
    id: 11,
    name: 'Air Terjun Pawa',
    location: 'Kecamatan Langgudu, Kabupaten Bima, NTB',
    price: 15000,
    category: 'Gunung & Alam',
    rating: 4.8,
    image: '/uploads/destinations/air_terjun_pawa.jpg',
    status: 'Tersedia',
    description: 'Air terjun bertingkat alami di tengah rimbunnya hutan tropis pegunungan Langgudu berair jernih segar.',
  },
  {
    id: 12,
    name: 'Gunung Tambora',
    location: 'Kabupaten Bima - Dompu, NTB',
    price: 50000,
    category: 'Gunung & Alam',
    rating: 5.0,
    image: '/uploads/destinations/gunung_tambora.jpg',
    status: 'Tersedia',
    description: 'Kaldera gunung berapi raksasa legendaris dunia dengan lanskap padang savana dan panorama samudra luas.',
  },
  {
    id: 13,
    name: 'Pantai Kalaki',
    location: 'Kecamatan Palibelo, Kabupaten Bima, NTB',
    price: 10000,
    category: 'Bahari',
    rating: 4.7,
    image: '/uploads/destinations/pantai_kalaki.jpg',
    status: 'Tersedia',
    description: 'Pesisir pantai favorit sunset wisatawan lokal dengan lanskap Teluk Bima dan kuliner pesisir khas Mbojo.',
  },
  {
    id: 14,
    name: 'Bukit Matompo',
    location: 'Kecamatan Wawo, Kabupaten Bima, NTB',
    price: 15000,
    category: 'Gunung & Alam',
    rating: 4.8,
    image: '/uploads/destinations/bukit_matompo.jpg',
    status: 'Tersedia',
    description: 'Puncak perbukitan hijau tempat menyaksikan hamparan awan pagi dan lanskap lembah pegunungan Wawo.',
  },
  {
    id: 15,
    name: 'Mata Air Madapangga',
    location: 'Kecamatan Madapangga, Kabupaten Bima, NTB',
    price: 10000,
    category: 'Gunung & Alam',
    rating: 4.7,
    image: '/uploads/destinations/madapangga.jpg',
    status: 'Tersedia',
    description: 'Pemandian alam sumber air pegunungan yang sangat jernih, sejuk, dan dikelilingi pepohonan rimbun alami.',
  },
  {
    id: 16,
    name: 'Pantai Wane',
    location: 'Kecamatan Monta, Kabupaten Bima, NTB',
    price: 20000,
    category: 'Bahari',
    rating: 4.9,
    image: '/uploads/destinations/pantai_wane.jpg',
    status: 'Tersedia',
    description: 'Pantai pesisir selatan berombak eksotis dengan hamparan pasir putih dan atraksi lumba-lumba liar.',
  },
  {
    id: 17,
    name: 'Pulau Utop',
    location: 'Kecamatan Sape, Kabupaten Bima, NTB',
    price: 60000,
    category: 'Bahari',
    rating: 4.9,
    image: '/uploads/destinations/pulau_utop.jpg',
    status: 'Tersedia',
    description: 'Pulau nirwana mungil berpasir halus dengan keindahan taman laut bawah air yang ideal untuk snorkeling.',
  },
  {
    id: 18,
    name: 'Istana Asi Mbojo',
    location: 'Kota / Kabupaten Bima, NTB',
    price: 15000,
    category: 'Budaya & Geowisata',
    rating: 4.9,
    image: '/uploads/destinations/asi_mbojo.jpg',
    status: 'Tersedia',
    description: 'Museum cagar budaya bekas Istana Kesultanan Bima dengan koleksi mahkota emas dan sejarah Mbojo.',
  },
  {
    id: 19,
    name: 'Pantai Kiri Wera',
    location: 'Kecamatan Wera, Kabupaten Bima, NTB',
    price: 15000,
    category: 'Bahari',
    rating: 4.8,
    image: '/uploads/destinations/pantai_kiri.jpg',
    status: 'Tersedia',
    description: 'Pantai eksotis berpasir hitam berkilau khas lereng gunung api Sangiang dengan panorama perahu tradisional.',
  },
  {
    id: 20,
    name: 'Nisa Dora',
    location: 'Kecamatan Sape, Kabupaten Bima, NTB',
    price: 35000,
    category: 'Bahari',
    rating: 4.9,
    image: '/uploads/destinations/nisa_dora.jpg',
    status: 'Tersedia',
    description: 'Pulau nirwana kecil berpasir putih jernih untuk tempat bersantai, berenang, dan aktivitas diving.',
  },
  {
    id: 21,
    name: 'Toro Jambu',
    location: 'Kecamatan Lambu, Kabupaten Bima, NTB',
    price: 20000,
    category: 'Bahari',
    rating: 4.8,
    image: '/uploads/destinations/toro_jambu.jpg',
    status: 'Tersedia',
    description: 'Tanjung bukit perbukitan hijau nan megah berhadapan dengan samudra lepas, tempat favorit sunrise & camping.',
  },
];

const INITIAL_BOOKINGS = [
  { id: 'BK-9901', customerName: 'Budi Santoso', customerEmail: 'budi.santoso@gmail.com', destinationName: 'Pantai Lariti', date: '2026-09-15', guests: 2, totalPrice: 30000, status: 'Lunas' },
  { id: 'BK-9902', customerName: 'Dewi Lestari', customerEmail: 'dewi.lestari@yahoo.com', destinationName: 'Pulau Sangiang', date: '2026-09-18', guests: 1, totalPrice: 50000, status: 'Menunggu' },
  { id: 'BK-9903', customerName: 'Siti Rahmawati', customerEmail: 'siti.rahma@gmail.com', destinationName: 'Uma Lengge Wawo', date: '2026-09-20', guests: 4, totalPrice: 40000, status: 'Lunas' },
];

const normalizeRoute = (rawHash: string): string => {
  let clean = rawHash.replace(/^#\/?/, '').trim();
  if (!clean) return 'destinations';
  if (clean === 'dashboard') return 'admin/dashboard';
  if (clean === 'reports') return 'admin/reports';
  if (clean === 'bookings') return 'admin/bookings';
  if (clean === 'my-bookings') return 'bookings/my-bookings';
  if (clean === 'users') return 'admin/users';
  if (clean === 'settings') return 'settings';
  return clean;
};

export default function App() {
  const [theme, setTheme] = useState<'light' | 'dark'>('dark');
  const [currentRouteState, setCurrentRouteState] = useState<string>(() => {
    return normalizeRoute(window.location.hash);
  });
  const [isMobileOpen, setIsMobileOpen] = useState(false);

  const setCurrentRoute = (route: string) => {
    const norm = normalizeRoute(route);
    setCurrentRouteState(norm);
    const targetHash = '#/' + norm;
    if (window.location.hash !== targetHash) {
      window.location.hash = targetHash;
    }
  };

  const currentRoute = currentRouteState;

  // Fallback state for local user bookings
  const [userBookings, setUserBookings] = useState<any[]>([]);

  // Zustand Stores
  const { user, token, isAuthenticated, isLoading: isLoggingIn, error: loginError, login, register, logout, initAuth } = useAuthStore();
  const { notifications, unreadCount, isConnected, fetchNotifications, connectSSE, disconnectSSE, markAsRead, markAllAsRead } = useNotificationStore();
  const { users: apiUsers, total: userTotal, isLoading: isUserLoading, fetchUsers, createUser, updateUser, deleteUser } = useUserStore();
  const { destinations: apiDestinations, fetchDestinations, createDestination, updateDestination, deleteDestination, fetchSlots } = useDestinationStore();
  const { bookings: apiBookings, fetchMyBookings, fetchAllAdminBookings, createBooking, payBooking } = useBookingStore();

  // Auth View State (Login vs Register)
  const [authMode, setAuthMode] = useState<'login' | 'register'>('login');
  const [loginEmail, setLoginEmail] = useState('aji@example.com');
  const [loginPassword, setLoginPassword] = useState('password123');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(true);

  // Register inputs & OTP state
  const [regName, setRegName] = useState('');
  const [regEmail, setRegEmail] = useState('');
  const [regPhone, setRegPhone] = useState('');
  const [regPassword, setRegPassword] = useState('');
  const [isOtpModalOpen, setIsOtpModalOpen] = useState(false);
  const [pendingRegData, setPendingRegData] = useState<{ name: string; email: string; phone: string; password: string } | null>(null);

  // Filters & Modals
  const [destinations, setDestinations] = useState<any[]>(INITIAL_DESTINATIONS);
  const [userSearch, setUserSearch] = useState('');
  const [userRoleFilter, setUserRoleFilter] = useState<string>('ALL');
  const [isUserModalOpen, setIsUserModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<any | null>(null);

  const [formName, setFormName] = useState('');
  const [formEmail, setFormEmail] = useState('');
  const [formPhone, setFormPhone] = useState('');
  const [formRole, setFormRole] = useState<'admin' | 'customer'>('customer');
  const [formPassword, setFormPassword] = useState('');

  const [destCategoryFilter, setDestCategoryFilter] = useState<string>('ALL');
  const [destSearch, setDestSearch] = useState<string>('');
  const [isDestModalOpen, setIsDestModalOpen] = useState(false);
  const [editingDest, setEditingDest] = useState<any | null>(null);
  const [destName, setDestName] = useState('');
  const [destLocation, setDestLocation] = useState('');
  const [destPrice, setDestPrice] = useState(15000);
  const [destCategory, setDestCategory] = useState<string>('Bahari');
  const [destDescription, setDestDescription] = useState('');
  const [destImageUrl, setDestImageUrl] = useState('');

  const [bookingStatusFilter, setBookingStatusFilter] = useState<string>('ALL');
  const [selectedDestForBooking, setSelectedDestForBooking] = useState<any | null>(null);
  const [bookingDate, setBookingDate] = useState<string>('2026-09-20');
  const [bookingTotalPeople, setBookingTotalPeople] = useState<number>(2);

  // Inisialisasi Auth & Listeners
  useEffect(() => {
    initAuth();
    fetchDestinations();

    const handleHashChange = () => {
      const norm = normalizeRoute(window.location.hash);
      setCurrentRouteState(norm);
    };

    if (!window.location.hash || window.location.hash === '#' || window.location.hash === '#/') {
      const defaultRoute = useAuthStore.getState().user?.role === 'admin' ? 'admin/dashboard' : 'destinations';
      window.location.hash = '#/' + defaultRoute;
      setCurrentRouteState(defaultRoute);
    } else {
      handleHashChange();
    }

    window.addEventListener('hashchange', handleHashChange);
    return () => window.removeEventListener('hashchange', handleHashChange);
  }, []);

  useEffect(() => {
    if (isAuthenticated && token) {
      fetchNotifications();
      if (user?.role === 'admin') {
        fetchAllAdminBookings();
      } else {
        fetchMyBookings();
      }
      connectSSE(token);
      if (user?.role !== 'admin' && (currentRoute === 'admin/dashboard' || currentRoute === 'dashboard')) {
        setCurrentRoute('destinations');
      }
    } else {
      disconnectSSE();
    }
  }, [isAuthenticated, token, user, currentRoute, fetchAllAdminBookings, fetchMyBookings]);

  useEffect(() => {
    if (isAuthenticated) {
      fetchUsers(userSearch, userRoleFilter);
    }
  }, [isAuthenticated, userSearch, userRoleFilter]);

  useEffect(() => {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [theme]);

  // Merged List Calculations
  const allDestinationsList = useMemo(() => {
    if (apiDestinations && apiDestinations.length > 0) {
      const formattedApi = (apiDestinations || []).filter(Boolean).map((d) => ({
        id: d?.id,
        name: d?.name || 'Destinasi Wisata',
        location: d?.location || 'Kabupaten Bima',
        price: d?.price || 15000,
        category: d?.category || 'Bahari',
        rating: d?.rating || 4.8,
        image: d?.image_url || '/uploads/destinations/lariti.jpg',
        status: d?.is_active ? 'Tersedia' : 'Penuh',
        description: d?.description || 'Destinasi wisata unggulan Kabupaten Bima.',
      }));
      return formattedApi;
    }
    return destinations;
  }, [apiDestinations, destinations]);

  const filteredDestinations = useMemo(() => {
    return allDestinationsList.filter((d) => {
      // 1. Kategori Filter
      if (destCategoryFilter !== 'ALL') {
        const fCat = destCategoryFilter.toLowerCase();
        const dCat = (d?.category || '').toLowerCase();
        const matchCat = fCat.includes('gunung')
          ? dCat.includes('gunung') || dCat.includes('alam')
          : fCat.includes('budaya')
          ? dCat.includes('budaya') || dCat.includes('geo')
          : fCat.includes('bahari')
          ? dCat.includes('bahari')
          : dCat.includes(fCat);
        if (!matchCat) return false;
      }

      // 2. Search Keyword Filter
      if (destSearch.trim() !== '') {
        const q = destSearch.toLowerCase().trim();
        const matchName = (d?.name || '').toLowerCase().includes(q);
        const matchLoc = (d?.location || '').toLowerCase().includes(q);
        const matchDesc = (d?.description || '').toLowerCase().includes(q);
        const matchCategory = (d?.category || '').toLowerCase().includes(q);
        if (!matchName && !matchLoc && !matchDesc && !matchCategory) return false;
      }

      return true;
    });
  }, [allDestinationsList, destCategoryFilter, destSearch]);

  const allAdminBookings = useMemo(() => {
    if (apiBookings && apiBookings.length > 0) {
      const mapped = apiBookings.map((b) => {
        const rawStatus = (b.status || '').toUpperCase();
        let statusLabel = 'Lunas';
        if (rawStatus === 'EXPIRED') {
          statusLabel = 'Expired';
        } else if (rawStatus === 'CANCELLED') {
          statusLabel = 'Dibatalkan';
        }

        return {
          id: b.booking_code || `BK-${b.id}`,
          customerName: b.user_name || user?.name || 'Wisatawan Bima',
          customerEmail: (b as any).user_email || user?.email || 'customer@example.com',
          destinationName: b.destination_name || 'Wisata Kabupaten Bima',
          date: b.booked_at ? b.booked_at.slice(0, 10) : '2026-09-20',
          guests: b.total_people,
          totalPrice: b.total_price,
          status: statusLabel,
          rawBookingId: b.id,
        };
      });

      // Sort newest booking ID first so latest bookings show at the top
      return mapped.sort((a, b) => (b.rawBookingId || 0) - (a.rawBookingId || 0));
    }
    return INITIAL_BOOKINGS;
  }, [apiBookings, user]);

  const filteredBookings = useMemo(() => {
    return allAdminBookings.filter((b) => {
      if (bookingStatusFilter === 'ALL') return true;
      return (b?.status || '').toLowerCase() === (bookingStatusFilter || '').toLowerCase();
    });
  }, [allAdminBookings, bookingStatusFilter]);

  const computedCategoryChartData = useMemo(() => {
    const counts: Record<string, number> = { 'Bahari': 0, 'Gunung & Alam': 0, 'Budaya & Geowisata': 0 };
    allDestinationsList.forEach((d) => {
      const cat = (d?.category || '').toLowerCase();
      if (cat.includes('bahari')) counts['Bahari']++;
      else if (cat.includes('gunung') || cat.includes('alam')) counts['Gunung & Alam']++;
      else counts['Budaya & Geowisata']++;
    });

    const total = allDestinationsList.length || 1;
    return [
      { name: 'Bahari', value: Math.round((counts['Bahari'] / total) * 100) || 50, color: '#10b981' },
      { name: 'Gunung & Alam', value: Math.round((counts['Gunung & Alam'] / total) * 100) || 30, color: '#0284c7' },
      { name: 'Budaya & Geowisata', value: Math.round((counts['Budaya & Geowisata'] / total) * 100) || 20, color: '#f59e0b' },
    ];
  }, [allDestinationsList]);

  const computedRevenueChartData = useMemo(() => {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];
    const monthStats: Record<number, { revenue: number; bookings: number }> = {};
    for (let i = 0; i < 12; i++) {
      monthStats[i] = { revenue: 0, bookings: 0 };
    }

    allAdminBookings.forEach((b) => {
      if (b.date) {
        const parts = b.date.split('-');
        if (parts.length >= 2) {
          const monthIdx = parseInt(parts[1], 10) - 1;
          if (monthIdx >= 0 && monthIdx < 12) {
            monthStats[monthIdx].bookings += 1;
            if (b.status === 'Lunas') {
              monthStats[monthIdx].revenue += b.totalPrice || 0;
            }
          }
        }
      }
    });

    return months.slice(0, 9).map((m, idx) => ({
      month: m,
      revenue: monthStats[idx].revenue,
      bookings: monthStats[idx].bookings,
    }));
  }, [allAdminBookings]);

  const effectiveUserBookings = useMemo(() => {
    if (apiBookings) {
      return apiBookings.map((b) => {
        const rawStatus = (b.status || '').toUpperCase();
        const normStatus = (rawStatus === 'PAID' || rawStatus === 'LUNAS') ? 'PAID' : rawStatus === 'PENDING' ? 'PENDING' : rawStatus === 'EXPIRED' ? 'EXPIRED' : 'CANCELLED';

        const destNameClean = (b.destination_name || '').toLowerCase().trim();
        const matchedDest = allDestinationsList.find((d) => {
          const dName = (d?.name || '').toLowerCase().trim();
          return Boolean(dName && destNameClean && (dName === destNameClean || destNameClean.includes(dName) || dName.includes(destNameClean)));
        });

        return {
          id: b.id,
          bookingCode: b.booking_code,
          customerName: user?.name || 'Customer',
          customerEmail: (b as any).user_email || user?.email || 'customer@example.com',
          destinationName: b.destination_name || matchedDest?.name || 'Wisata Kabupaten Bima',
          destinationLocation: b.destination_location || matchedDest?.location || 'Kabupaten Bima, NTB',
          destinationImage: b.destination_image || matchedDest?.image || '/uploads/destinations/lariti.jpg',
          totalPeople: b.total_people,
          totalPrice: b.total_price,
          status: normStatus,
          bookedAt: b.booked_at ? b.booked_at.slice(0, 10) : '2026-09-20',
          paymentUrl: b.payment_url,
        };
      });
    }
    return userBookings.length > 0 ? userBookings : [];
  }, [apiBookings, userBookings, user, allDestinationsList]);

  const formatIDR = (amount: number) => {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(amount);
  };

  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!loginEmail || !loginPassword) return;
    const success = await login(loginEmail, loginPassword, rememberMe);
    if (success) {
      toast.success('Login berhasil! Selamat datang di Kacaru Book.');
      const loggedUser = useAuthStore.getState().user;
      setCurrentRoute(loggedUser?.role === 'admin' ? 'dashboard' : 'destinations');
    } else {
      toast.error(useAuthStore.getState().error || 'Login gagal! Periksa email/no. HP dan password Anda.');
    }
  };

  const handleRegisterSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!regName || !regEmail || !regPassword) return;

    // Validasi format email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(regEmail)) {
      toast.error('Format email tidak valid! Gunakan format email resmi (contoh: bima@example.com)');
      return;
    }

    // Validasi format nomor HP jika diisi
    if (regPhone && (regPhone.length < 10 || !/^[0-9]+$/.test(regPhone))) {
      toast.error('Nomor HP tidak valid! Masukkan minimal 10 digit angka nomor HP (contoh: 081234567890)');
      return;
    }

    // 1. Process Register ke Backend Go untuk membuat akun & token JWT
    const success = await register({
      name: regName,
      email: regEmail,
      phone: regPhone || `0812${Math.floor(10000000 + Math.random() * 90000000)}`,
      password: regPassword,
    });

    if (success) {
      // 2. Minta Kode OTP 6-Digit asli dari Backend Go (tersimpan di otps DB & terkirim ke Mailtrap)
      await useAuthStore.getState().requestOTP();

      setPendingRegData({
        name: regName,
        email: regEmail,
        phone: regPhone,
        password: regPassword,
      });
      setIsOtpModalOpen(true);
      toast.info(`Kode OTP verifikasi 6-digit telah dikirimkan ke Mailtrap & email ${regEmail}`);
    }
  };

  const handleOtpVerifySuccess = () => {
    useAuthStore.getState().completeAuth();
    setIsOtpModalOpen(false);
    setPendingRegData(null);
    setRegName('');
    setRegEmail('');
    setRegPhone('');
    setRegPassword('');
    toast.success('Akun & No. HP Berhasil Diverifikasi OTP! Selamat datang di Kacaru Book.');
    setCurrentRoute('destinations');
  };

  const handleOpenUserModal = (userToEdit?: any) => {
    if (userToEdit) {
      setEditingUser(userToEdit);
      setFormName(userToEdit?.name || '');
      setFormEmail(userToEdit?.email || '');
      setFormPhone(userToEdit?.phone || '');
      setFormRole(userToEdit?.role === 'admin' ? 'admin' : 'customer');
      setFormPassword('');
    } else {
      setEditingUser(null);
      setFormName('');
      setFormEmail('');
      setFormPhone('');
      setFormRole('customer');
      setFormPassword('');
    }
    setIsUserModalOpen(true);
  };

  const handleSaveUser = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formName || !formEmail) return;

    if (editingUser) {
      const ok = await updateUser(editingUser.id, {
        name: formName,
        email: formEmail,
        phone: formPhone,
        role: formRole,
        password: formPassword || undefined,
      });
      if (ok) {
        toast.success('Data pengguna berhasil diperbarui!');
        setIsUserModalOpen(false);
      }
    } else {
      const ok = await createUser({
        name: formName,
        email: formEmail,
        phone: formPhone,
        role: formRole,
        password: formPassword || 'password123',
      });
      if (ok) {
        toast.success('Pengguna baru berhasil ditambahkan!');
        setIsUserModalOpen(false);
      }
    }
  };

  const handleDeleteUser = async (id: number | string) => {
    if (window.confirm('Apakah Anda yakin ingin menghapus pengguna ini dari sistem?')) {
      const ok = await deleteUser(id);
      if (ok) {
        toast.success('Pengguna telah dihapus dari sistem.');
      }
    }
  };

  const handleDeleteDestination = async (id: string | number) => {
    if (window.confirm('Apakah Anda yakin ingin menghapus destinasi wisata ini dari database?')) {
      const numId = Number(id);
      if (!isNaN(numId) && numId > 0) {
        const ok = await deleteDestination(numId);
        if (ok) {
          toast.success('Destinasi wisata berhasil dihapus dari database!');
          return;
        }
      }
      setDestinations((prev) => prev.filter((d) => String(d.id) !== String(id)));
      toast.success('Destinasi wisata telah dihapus.');
    }
  };

  const handleOpenAddModal = () => {
    setEditingDest(null);
    setDestName('');
    setDestLocation('');
    setDestPrice(15000);
    setDestCategory('Bahari');
    setDestDescription('');
    setDestImageUrl('');
    setIsDestModalOpen(true);
  };

  const handleEditDestination = (dest: any) => {
    setEditingDest(dest);
    setDestName(dest?.name || '');
    setDestLocation(dest?.location || '');
    setDestPrice(dest?.price || 15000);
    setDestCategory(dest?.category || 'Bahari');
    setDestDescription(dest?.description || '');
    setDestImageUrl(dest?.image || dest?.image_url || '');
    setIsDestModalOpen(true);
  };

  const handleSaveDestination = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!destName || !destLocation) return;
    const finalImage = destImageUrl || 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=800&q=80';

    if (editingDest) {
      const numId = Number(editingDest.id);
      if (!isNaN(numId) && numId > 0) {
        await updateDestination(numId, {
          name: destName,
          location: destLocation,
          price: destPrice,
          category: destCategory,
          image_url: finalImage,
          description: destDescription,
        });
      }
      setDestinations((prev) =>
        prev.map((d) =>
          d.id === editingDest.id
            ? {
                ...d,
                name: destName,
                location: destLocation,
                price: destPrice,
                category: destCategory,
                image: finalImage,
                description: destDescription,
              }
            : d
        )
      );
      toast.success('Destinasi wisata berhasil diperbarui!');
    } else {
      await createDestination({
        name: destName,
        location: destLocation,
        price: destPrice,
        category: destCategory,
        image_url: finalImage,
        description: destDescription || 'Destinasi wisata unggulan Kabupaten Bima.',
      });

      const newDestObj = {
        id: Date.now(),
        name: destName,
        location: destLocation,
        price: destPrice,
        category: destCategory,
        rating: 4.9,
        image: finalImage,
        status: 'Tersedia',
        description: destDescription || 'Destinasi wisata unggulan Kabupaten Bima.',
      };

      setDestinations((prev) => [newDestObj, ...prev]);
      toast.success('Destinasi wisata baru berhasil dibuat & tampil di katalog!');
    }

    setIsDestModalOpen(false);
    setEditingDest(null);
    setDestName('');
    setDestLocation('');
    setDestDescription('');
    setDestImageUrl('');
    fetchDestinations();
  };

  const handleConfirmBookingSubmit = async () => {
    if (!selectedDestForBooking) return;
    const numDestId = Number(selectedDestForBooking.id);

    let slotId = 1;
    if (!isNaN(numDestId) && numDestId > 0) {
      const slots = await fetchSlots(numDestId);
      if (slots && slots.length > 0) {
        const foundSlot = slots.find((s) => s.date.startsWith(bookingDate)) || slots[0];
        slotId = foundSlot.id;
      }
    }

    const visitors = Array.from({ length: bookingTotalPeople }, (_, i) => ({
      name: i === 0 ? (user?.name || 'Visitor 1') : `${user?.name || 'Visitor'} (Rombongan ${i + 1})`,
      phone: user?.phone || '08123456789',
      email: user?.email || 'customer@example.com',
      age: 25,
    }));

    const createdBooking = await createBooking(slotId, bookingTotalPeople, visitors);

    if (createdBooking) {
      await fetchDestinations();
      if (!isNaN(numDestId) && numDestId > 0) {
        await fetchSlots(numDestId);
      }
      const payUrl = await payBooking(createdBooking.id);
      if (payUrl && typeof payUrl === 'string' && payUrl.startsWith('https://checkout.xendit.co') && !payUrl.includes('inv_mock_')) {
        toast.info('Membuka Invoice Xendit Payment Gateway...');
        window.open(payUrl, '_blank');
      } else {
        toast.success('Pemesanan tiket berhasil! Silakan klik Bayar via Xendit.');
      }
      setSelectedDestForBooking(null);
      setCurrentRoute('my-bookings');
    } else {
      const newB = {
        id: Date.now(),
        bookingCode: `BK-2026-${Math.floor(1000 + Math.random() * 9000)}`,
        destinationName: selectedDestForBooking?.name || 'Wisata Kabupaten Bima',
        destinationLocation: selectedDestForBooking?.location || 'Kabupaten Bima',
        destinationImage: selectedDestForBooking?.image || '/uploads/destinations/lariti.jpg',
        totalPeople: bookingTotalPeople,
        totalPrice: (selectedDestForBooking?.price || 15000) * bookingTotalPeople,
        status: 'PENDING',
        bookedAt: new Date().toISOString().slice(0, 16).replace('T', ' '),
        paymentUrl: 'https://checkout.xendit.co/web/inv_xnd_demo',
      };
      setUserBookings((prev) => [newB, ...prev]);
      setSelectedDestForBooking(null);
      toast.success('Pemesanan tiket berhasil dibuat! Silakan bayar via Xendit.');
      setCurrentRoute('my-bookings');
    }
  };

  const handlePayBookingAction = async (bookingId: number) => {
    const numId = Number(bookingId);
    if (!isNaN(numId) && numId > 0) {
      const url = await payBooking(numId);
      if (url && typeof url === 'string' && url.startsWith('https://checkout.xendit.co') && !url.includes('inv_mock_')) {
        return url;
      }
    }
    return null;
  };

  // Login View Render
  if (!isAuthenticated) {
    return (
      <div className={`min-h-screen transition-colors duration-300 font-sans flex items-center justify-center p-4 ${theme === 'dark' ? 'bg-slate-950 text-white' : 'bg-slate-100 text-slate-900'}`}>
        <ToastProvider />
        <div className="absolute top-6 right-6 z-20">
          <button
            onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
            className="p-3 rounded-2xl bg-white/80 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-800 shadow-lg backdrop-blur-md hover:scale-105 transition-all text-slate-700 dark:text-slate-200"
            title="Ganti Tema"
          >
            {theme === 'light' ? <Moon className="w-5 h-5 text-indigo-600" /> : <Sun className="w-5 h-5 text-amber-400" />}
          </button>
        </div>

        <div className="w-full max-w-md bg-white/90 dark:bg-slate-900/90 rounded-3xl p-8 border border-slate-200/80 dark:border-slate-800 shadow-2xl backdrop-blur-xl space-y-6">
          <div className="text-center space-y-2">
            <div className="inline-flex mb-2">
              <img
                src="/kacaru-logo.png"
                alt="Kacaru Book Logo"
                className="w-28 h-28 rounded-full object-cover shadow-2xl hover:scale-105 transition-all duration-300 mx-auto"
              />
            </div>
            <h1 className="text-2xl font-black bg-gradient-to-r from-emerald-600 to-teal-500 bg-clip-text text-transparent">
              Kacaru Book Bima
            </h1>
            <p className="text-xs text-slate-500 font-bold">Sistem Reservasi Wisata Kabupaten Bima, NTB</p>
          </div>

          {/* Mode Switcher Tabs */}
          <div className="grid grid-cols-2 p-1 bg-slate-100 dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700/80">
            <button
              type="button"
              onClick={() => setAuthMode('login')}
              className={`py-2 text-xs font-black rounded-xl transition-all ${
                authMode === 'login'
                  ? 'bg-white dark:bg-slate-900 text-emerald-600 dark:text-emerald-400 shadow-md'
                  : 'text-slate-500 hover:text-slate-700 dark:hover:text-white'
              }`}
            >
              🔑 Masuk (Login)
            </button>
            <button
              type="button"
              onClick={() => setAuthMode('register')}
              className={`py-2 text-xs font-black rounded-xl transition-all ${
                authMode === 'register'
                  ? 'bg-white dark:bg-slate-900 text-emerald-600 dark:text-emerald-400 shadow-md'
                  : 'text-slate-500 hover:text-slate-700 dark:hover:text-white'
              }`}
            >
              📝 Daftar Akun (Register)
            </button>
          </div>

          {loginError && (
            <div className="p-3 rounded-2xl bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs font-bold text-center">
              {loginError}
            </div>
          )}

          {authMode === 'login' ? (
            <form onSubmit={handleLoginSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-black uppercase text-slate-500 mb-1.5">Email / Nomor HP</label>
                <input
                  type="text"
                  required
                  value={loginEmail}
                  onChange={(e) => setLoginEmail(e.target.value)}
                  placeholder="Nomor HP (contoh: 081234567890) atau Email"
                  className="w-full px-4 py-3 rounded-2xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-xs font-black uppercase text-slate-500 mb-1.5">Password</label>
                <div className="relative">
                  <input
                    type={showPassword ? 'text' : 'password'}
                    required
                    value={loginPassword}
                    onChange={(e) => setLoginPassword(e.target.value)}
                    placeholder="••••••••"
                    className="w-full px-4 py-3 rounded-2xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                  >
                    {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
              </div>

              <div className="flex items-center justify-between text-xs font-bold text-slate-500">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={rememberMe}
                    onChange={(e) => setRememberMe(e.target.checked)}
                    className="rounded text-emerald-600 focus:ring-emerald-500"
                  />
                  <span>Ingat Saya</span>
                </label>
              </div>

              <button
                type="submit"
                disabled={isLoggingIn}
                className="w-full py-3.5 px-4 rounded-2xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white font-black text-xs shadow-lg shadow-emerald-600/30 active:scale-95 transition-all disabled:opacity-50"
              >
                {isLoggingIn ? 'Memverifikasi Access...' : 'Masuk Aplikasi (Login)'}
              </button>
            </form>
          ) : (
            <form onSubmit={handleRegisterSubmit} autoComplete="off" className="space-y-3.5">
              <div>
                <label className="block text-xs font-black uppercase text-slate-500 mb-1">Nama Lengkap</label>
                <input
                  type="text"
                  required
                  autoComplete="off"
                  value={regName}
                  onChange={(e) => setRegName(e.target.value)}
                  placeholder="Nama Lengkap Anda"
                  className="w-full px-4 py-2.5 rounded-2xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-xs font-black uppercase text-slate-500 mb-1">Alamat Email</label>
                <input
                  type="email"
                  required
                  autoComplete="off"
                  value={regEmail}
                  onChange={(e) => setRegEmail(e.target.value)}
                  placeholder="nama@email.com"
                  className="w-full px-4 py-2.5 rounded-2xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-xs font-black uppercase text-slate-500 mb-1">Nomor WhatsApp / HP</label>
                <input
                  type="tel"
                  autoComplete="off"
                  value={regPhone}
                  onChange={(e) => setRegPhone(e.target.value)}
                  placeholder="081234567890"
                  className="w-full px-4 py-2.5 rounded-2xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-xs font-black uppercase text-slate-500 mb-1">Password</label>
                <input
                  type="password"
                  required
                  autoComplete="new-password"
                  value={regPassword}
                  onChange={(e) => setRegPassword(e.target.value)}
                  placeholder="Kata sandi (min. 8 karakter)"
                  className="w-full px-4 py-2.5 rounded-2xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white"
                />
              </div>

              <button
                type="submit"
                disabled={isLoggingIn}
                className="w-full py-3.5 px-4 rounded-2xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white font-black text-xs shadow-lg shadow-emerald-600/30 active:scale-95 transition-all disabled:opacity-50 mt-2"
              >
                {isLoggingIn ? 'Memproses Pendaftaran...' : 'Daftar Akun Sekarang'}
              </button>
            </form>
          )}

          <div className="p-4 rounded-2xl bg-gradient-to-br from-slate-100 to-emerald-50/50 dark:from-slate-800/80 dark:to-slate-800/40 text-[11px] font-bold text-slate-600 dark:text-slate-300 space-y-3 border border-emerald-500/20 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="font-extrabold text-slate-800 dark:text-slate-100 uppercase tracking-wider text-[11px] flex items-center gap-1.5">
                ⚡ Quick Login (1-Click Demo)
              </span>
              <span className="text-[10px] text-emerald-600 dark:text-emerald-400 font-black px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20">
                Tanpa Isian
              </span>
            </div>
            <div className="grid grid-cols-2 gap-3 pt-1">
              <button
                type="button"
                onClick={async () => {
                  setLoginEmail('aji@example.com');
                  setLoginPassword('password123');
                  const ok = await login('aji@example.com', 'password123');
                  if (ok) {
                    toast.success('Berhasil Masuk sebagai Administrator!');
                    setCurrentRoute('dashboard');
                  } else {
                    toast.error(useAuthStore.getState().error || 'Gagal masuk sebagai Admin.');
                  }
                }}
                className="py-3 px-3 rounded-2xl bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600 text-white font-black text-xs transition-all active:scale-95 text-center flex items-center justify-center gap-1.5 shadow-md shadow-amber-500/20"
              >
                <span>👑 Masuk Sebagai Admin</span>
              </button>
              <button
                type="button"
                onClick={async () => {
                  setLoginEmail('budi@example.com');
                  setLoginPassword('password123');
                  const ok = await login('budi@example.com', 'password123');
                  if (ok) {
                    toast.success('Berhasil Masuk sebagai Wisatawan!');
                    setCurrentRoute('destinations');
                  } else {
                    toast.error(useAuthStore.getState().error || 'Gagal masuk sebagai Wisatawan.');
                  }
                }}
                className="py-3 px-3 rounded-2xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white font-black text-xs transition-all active:scale-95 text-center flex items-center justify-center gap-1.5 shadow-md shadow-emerald-600/20"
              >
                <span>👤 Masuk Sebagai Wisatawan</span>
              </button>
            </div>
          </div>
        </div>

        <OTPVerificationModal
          isOpen={isOtpModalOpen}
          onClose={() => setIsOtpModalOpen(false)}
          email={pendingRegData?.email || ''}
          phone={pendingRegData?.phone}
          onVerifySuccess={handleOtpVerifySuccess}
        />
      </div>
    );
  }

  // Main Dashboard Application Render
  return (
    <div className={`min-h-screen transition-colors duration-300 font-sans ${theme === 'dark' ? 'bg-slate-950 text-white' : 'bg-slate-50 text-slate-900'}`}>
      <ToastProvider />

      <Sidebar
        user={user}
        currentRoute={currentRoute}
        setCurrentRoute={setCurrentRoute}
        onLogout={logout}
        isMobileOpen={isMobileOpen}
        setIsMobileOpen={setIsMobileOpen}
        effectiveUserBookingsCount={effectiveUserBookings.length}
      />

      <div className="lg:pl-72 flex flex-col min-h-screen">
        <Navbar
          user={user}
          theme={theme}
          setTheme={setTheme}
          notifications={notifications}
          unreadCount={unreadCount}
          isConnected={isConnected}
          onMarkAsRead={markAsRead}
          onMarkAllAsRead={markAllAsRead}
          setIsMobileOpen={setIsMobileOpen}
        />

        <main className="flex-1 p-6 space-y-6 max-w-7xl w-full mx-auto">
          {(currentRoute === 'admin/dashboard' || currentRoute === 'dashboard') && user?.role === 'admin' && (
            <DashboardOverview
              filteredBookingsCount={allAdminBookings.length}
              paidBookingsCount={allAdminBookings.filter((b) => b.status === 'Lunas').length}
              totalRevenue={allAdminBookings.reduce((acc, curr) => curr.status === 'Lunas' ? acc + curr.totalPrice : acc, 0)}
              destinationsCount={allDestinationsList.length}
              formatIDR={formatIDR}
              revenueChartData={computedRevenueChartData}
              categoryChartData={computedCategoryChartData}
            />
          )}

          {(currentRoute === 'destinations' || (user?.role !== 'admin' && (currentRoute === 'admin/dashboard' || currentRoute === 'dashboard'))) && (
            <DestinationGrid
              user={user}
              filteredDestinations={filteredDestinations}
              destCategoryFilter={destCategoryFilter}
              setDestCategoryFilter={setDestCategoryFilter}
              destSearch={destSearch}
              setDestSearch={setDestSearch}
              onOpenAddModal={handleOpenAddModal}
              onBookTicket={(dest) => setSelectedDestForBooking(dest)}
              onEditDestination={handleEditDestination}
              onDeleteDestination={handleDeleteDestination}
              formatIDR={formatIDR}
            />
          )}

          {(currentRoute === 'bookings/my-bookings' || currentRoute === 'my-bookings') && user?.role !== 'admin' && (
            <MyBookingsView
              effectiveUserBookings={effectiveUserBookings}
              onNavigateToDestinations={() => setCurrentRoute('destinations')}
              onPayBooking={handlePayBookingAction}
              formatIDR={formatIDR}
            />
          )}

          {(currentRoute === 'admin/bookings' || currentRoute === 'bookings') && user?.role === 'admin' && (
            <AdminBookingsView
              filteredBookings={filteredBookings}
              bookingStatusFilter={bookingStatusFilter}
              setBookingStatusFilter={setBookingStatusFilter}
              formatIDR={formatIDR}
            />
          )}

          {(currentRoute === 'admin/users' || currentRoute === 'users') && user?.role === 'admin' && (
            <UserManagementView
              apiUsers={apiUsers}
              userTotal={userTotal}
              isUserLoading={isUserLoading}
              userSearch={userSearch}
              setUserSearch={setUserSearch}
              userRoleFilter={userRoleFilter}
              setUserRoleFilter={setUserRoleFilter}
              onOpenUserModal={handleOpenUserModal}
              onDeleteUser={handleDeleteUser}
            />
          )}

          {(currentRoute === 'admin/reports' || currentRoute === 'reports') && <ReportAnalyticsView />}

          {(currentRoute === 'admin/settings' || currentRoute === 'settings') && <SettingsView user={user} />}
        </main>
      </div>

      {/* Modals */}
      <UserModal
        isOpen={isUserModalOpen}
        onClose={() => setIsUserModalOpen(false)}
        editingUser={editingUser}
        formName={formName}
        setFormName={setFormName}
        formEmail={formEmail}
        setFormEmail={setFormEmail}
        formPhone={formPhone}
        setFormPhone={setFormPhone}
        formRole={formRole}
        setFormRole={setFormRole}
        formPassword={formPassword}
        setFormPassword={setFormPassword}
        onSave={handleSaveUser}
      />

      <DestinationModal
        isOpen={isDestModalOpen}
        onClose={() => setIsDestModalOpen(false)}
        editingDestination={editingDest}
        destName={destName}
        setDestName={setDestName}
        destLocation={destLocation}
        setDestLocation={setDestLocation}
        destPrice={destPrice}
        setDestPrice={setDestPrice}
        destCategory={destCategory}
        setDestCategory={setDestCategory}
        destDescription={destDescription}
        setDestDescription={setDestDescription}
        destImageUrl={destImageUrl}
        setDestImageUrl={setDestImageUrl}
        onSave={handleSaveDestination}
      />

      {Boolean(selectedDestForBooking) && (
        <TicketModal
          selectedDest={selectedDestForBooking}
          onClose={() => setSelectedDestForBooking(null)}
          bookingDate={bookingDate}
          setBookingDate={setBookingDate}
          bookingTotalPeople={bookingTotalPeople}
          setBookingTotalPeople={setBookingTotalPeople}
          onConfirm={handleConfirmBookingSubmit}
          formatIDR={formatIDR}
        />
      )}
    </div>
  );
}
