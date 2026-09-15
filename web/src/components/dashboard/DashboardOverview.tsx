import React from 'react';
import {
  Ticket,
  DollarSign,
  Compass,
  TrendingUp,
  ArrowUpRight,
  Calendar
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
  Cell
} from 'recharts';

interface DashboardOverviewProps {
  filteredBookingsCount: number;
  paidBookingsCount?: number;
  totalRevenue: number;
  destinationsCount: number;
  formatIDR: (amount: number) => string;
  revenueChartData: any[];
  categoryChartData: any[];
}

// Custom Tooltip for Revenue & Booking Chart
const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0]?.payload || {};
    return (
      <div className="bg-slate-900/95 backdrop-blur-xl p-4 rounded-2xl border border-emerald-500/30 shadow-2xl space-y-2 text-xs text-white min-w-[200px]">
        <div className="flex items-center justify-between border-b border-slate-800 pb-2">
          <span className="font-extrabold text-slate-300 flex items-center gap-1.5">
            <Calendar className="w-3.5 h-3.5 text-emerald-400" />
            {label} 2026
          </span>
          <span className="text-[10px] font-black bg-emerald-500/20 text-emerald-300 px-2.5 py-0.5 rounded-full border border-emerald-500/30">
            {data.bookings || 0} Booking
          </span>
        </div>
        <div className="space-y-1 pt-1">
          <div className="text-[11px] text-slate-400 font-medium">Total Pendapatan</div>
          <div className="text-base font-black bg-gradient-to-r from-emerald-400 to-teal-200 bg-clip-text text-transparent">
            {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(data.revenue || 0)}
          </div>
        </div>
      </div>
    );
  }
  return null;
};

// Formatter for Y-Axis ticks in readable IDR format (Rp 19,5 jt / Rp 500 rb)
const formatYAxisIDR = (value: number) => {
  if (value === 0) return 'Rp 0';
  if (value >= 1_000_000_000) {
    return `Rp ${(value / 1_000_000_000).toLocaleString('id-ID', { maximumFractionDigits: 1 })} M`;
  }
  if (value >= 1_000_000) {
    return `Rp ${(value / 1_000_000).toLocaleString('id-ID', { maximumFractionDigits: 1 })} jt`;
  }
  if (value >= 1_000) {
    return `Rp ${(value / 1_000).toLocaleString('id-ID')} rb`;
  }
  return `Rp ${value}`;
};

export const DashboardOverview: React.FC<DashboardOverviewProps> = ({
  filteredBookingsCount,
  paidBookingsCount = 0,
  totalRevenue,
  destinationsCount,
  formatIDR,
  revenueChartData,
  categoryChartData,
}) => {
  const paidRate = filteredBookingsCount > 0 ? Math.round((paidBookingsCount / filteredBookingsCount) * 100) : 0;

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Top 4 Stat Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        {/* Card 1: Total Booking */}
        <div className="bg-gradient-to-br from-emerald-950/40 via-slate-900 to-slate-900 p-5 rounded-3xl border border-emerald-500/30 hover:border-emerald-400 shadow-lg shadow-emerald-500/10 hover:shadow-emerald-500/20 hover:-translate-y-1.5 transition-all duration-300 space-y-3 cursor-pointer group">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-black uppercase tracking-widest text-emerald-400">
              Total Booking DB
            </span>
            <div className="p-2.5 rounded-2xl bg-emerald-500/15 text-emerald-400 border border-emerald-500/30 shadow-md group-hover:scale-110 transition-transform">
              <Ticket className="w-5 h-5 animate-pulse" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-black text-white group-hover:text-emerald-300 transition-colors">
              {filteredBookingsCount} Pemesanan
            </div>
            <div className="flex items-center gap-1 text-[11px] font-extrabold text-emerald-400 mt-1">
              <ArrowUpRight className="w-3.5 h-3.5" />
              <span>{paidBookingsCount} Lunas di Database</span>
            </div>
          </div>
        </div>

        {/* Card 2: Total Pendapatan */}
        <div className="bg-gradient-to-br from-teal-950/40 via-slate-900 to-slate-900 p-5 rounded-3xl border border-teal-500/30 hover:border-teal-400 shadow-lg shadow-teal-500/10 hover:shadow-teal-500/20 hover:-translate-y-1.5 transition-all duration-300 space-y-3 cursor-pointer group">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-black uppercase tracking-widest text-teal-400">
              Total Pendapatan Riil
            </span>
            <div className="p-2.5 rounded-2xl bg-teal-500/15 text-teal-400 border border-teal-500/30 shadow-md group-hover:scale-110 transition-transform">
              <DollarSign className="w-5 h-5" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-black bg-gradient-to-r from-white via-teal-200 to-emerald-300 bg-clip-text text-transparent">
              {formatIDR(totalRevenue)}
            </div>
            <div className="flex items-center gap-1 text-[11px] font-extrabold text-teal-300 mt-1">
              <ArrowUpRight className="w-3.5 h-3.5" />
              <span>Omzet dari {paidBookingsCount} Transaksi Lunas</span>
            </div>
          </div>
        </div>

        {/* Card 3: Destinasi Aktif */}
        <div className="bg-gradient-to-br from-sky-950/40 via-slate-900 to-slate-900 p-5 rounded-3xl border border-sky-500/30 hover:border-sky-400 shadow-lg shadow-sky-500/10 hover:shadow-sky-500/20 hover:-translate-y-1.5 transition-all duration-300 space-y-3 cursor-pointer group">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-black uppercase tracking-widest text-sky-400">
              Destinasi Aktif DB
            </span>
            <div className="p-2.5 rounded-2xl bg-sky-500/15 text-sky-400 border border-sky-500/30 shadow-md group-hover:scale-110 transition-transform">
              <Compass className="w-5 h-5 group-hover:rotate-45 transition-transform duration-500" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-black text-white group-hover:text-sky-300 transition-colors">
              {destinationsCount} Lokasi Wisata
            </div>
            <div className="text-[11px] font-extrabold text-sky-400 mt-1">
              Kabupaten Bima & NTB
            </div>
          </div>
        </div>

        {/* Card 4: Direct Sales */}
        <div className="bg-gradient-to-br from-amber-950/40 via-slate-900 to-slate-900 p-5 rounded-3xl border border-amber-500/30 hover:border-amber-400 shadow-lg shadow-amber-500/10 hover:shadow-amber-500/20 hover:-translate-y-1.5 transition-all duration-300 space-y-3 cursor-pointer group">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-black uppercase tracking-widest text-amber-400">
              Rate Transaksi Lunas
            </span>
            <div className="p-2.5 rounded-2xl bg-amber-500/15 text-amber-400 border border-amber-500/30 shadow-md group-hover:scale-110 transition-transform">
              <TrendingUp className="w-5 h-5" />
            </div>
          </div>
          <div>
            <div className="text-2xl font-black text-white group-hover:text-amber-300 transition-colors">
              {paidRate}% Lunas
            </div>
            <div className="text-[11px] font-extrabold text-amber-400 mt-1">
              Terverifikasi Payment Gateway
            </div>
          </div>
        </div>
      </div>

      {/* Analytics Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-white dark:bg-slate-900 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h3 className="text-base font-black text-slate-900 dark:text-white flex items-center gap-2">
                Grafik Pendapatan & Booking (2026)
              </h3>
              <p className="text-xs text-slate-500 font-medium">
                Tren omzet bulanan dalam Mata Uang Rupiah (IDR).
              </p>
            </div>
            <div className="flex items-center gap-2">
              <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 text-xs font-extrabold border border-emerald-500/20">
                <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                Format Rupiah (IDR)
              </span>
            </div>
          </div>

          <div className="h-72 w-full pt-2">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={revenueChartData} margin={{ top: 10, right: 10, left: 15, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorRev" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.45} />
                    <stop offset="95%" stopColor="#10b981" stopOpacity={0.02} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#334155" opacity={0.15} />
                <XAxis 
                  dataKey="month" 
                  stroke="#94a3b8" 
                  fontSize={11} 
                  tickLine={false} 
                  dy={5}
                />
                <YAxis 
                  stroke="#94a3b8" 
                  fontSize={11} 
                  tickLine={false} 
                  axisLine={false} 
                  tickFormatter={formatYAxisIDR}
                  width={75}
                />
                <Tooltip content={<CustomTooltip />} />
                <Area 
                  type="monotone" 
                  dataKey="revenue" 
                  stroke="#10b981" 
                  strokeWidth={3} 
                  fillOpacity={1} 
                  fill="url(#colorRev)" 
                  dot={{ r: 4, fill: '#10b981', stroke: '#ffffff', strokeWidth: 1.5 }}
                  activeDot={{ r: 7, fill: '#10b981', stroke: '#0f172a', strokeWidth: 2 }}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-900 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm space-y-4">
          <div>
            <h3 className="text-base font-black text-slate-900 dark:text-white">
              Sebaran Kategori Wisata
            </h3>
            <p className="text-xs text-slate-500 font-medium">
              Persentase minat destinasi pengunjung.
            </p>
          </div>

          <div className="h-48 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie data={categoryChartData} cx="50%" cy="50%" innerRadius={50} outerRadius={75} paddingAngle={5} dataKey="value">
                  {categoryChartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="space-y-2 pt-2 border-t border-slate-100 dark:border-slate-800">
            {(categoryChartData || []).map((c, idx) => (
              <div key={c?.name || idx} className="flex items-center justify-between text-xs font-bold">
                <div className="flex items-center gap-2">
                  <span className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: c?.color || '#10b981' }} />
                  <span className="text-slate-700 dark:text-slate-300">{c?.name || '-'}</span>
                </div>
                <span className="text-slate-900 dark:text-white font-black">{c?.value || 0}%</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

