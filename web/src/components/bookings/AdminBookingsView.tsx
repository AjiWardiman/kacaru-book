import React, { useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight, Ticket } from 'lucide-react';
import { axiosClient } from '../../api/axiosClient';

interface AdminBookingsViewProps {
  filteredBookings: any[];
  bookingStatusFilter: string;
  setBookingStatusFilter: (status: string) => void;
  formatIDR: (amount: number) => string;
}

export const AdminBookingsView: React.FC<AdminBookingsViewProps> = ({
  filteredBookings,
  bookingStatusFilter,
  setBookingStatusFilter,
  formatIDR,
}) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);

  // Reset to page 1 when filter or page size changes
  useEffect(() => {
    setCurrentPage(1);
  }, [bookingStatusFilter, itemsPerPage]);

  const totalItems = filteredBookings.length;
  const totalPages = Math.ceil(totalItems / itemsPerPage) || 1;
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = Math.min(startIndex + itemsPerPage, totalItems);
  const paginatedBookings = filteredBookings.slice(startIndex, endIndex);

  return (
    <div className="space-y-6">
      {/* Header & Filter Bar */}
      <div className="bg-white dark:bg-slate-900 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h3 className="text-lg font-extrabold text-slate-900 dark:text-white flex items-center gap-2">
            <Ticket className="w-5 h-5 text-emerald-500" />
            Manajemen Pemesanan Platform
          </h3>
          <p className="text-xs font-bold text-slate-600 dark:text-slate-300 mt-0.5">
            Daftar reservasi tiket wisatawan ({totalItems} Pemesanan Terdaftar).
          </p>
        </div>

        {/* Status Filter Pills */}
        <div className="flex flex-wrap items-center gap-2">
          {['ALL', 'Lunas', 'Menunggu', 'Dibatalkan'].map((st) => (
            <button
              key={st}
              onClick={() => setBookingStatusFilter(st)}
              className={`px-3 py-1.5 rounded-xl text-xs font-extrabold transition-all border ${
                bookingStatusFilter === st
                  ? 'bg-emerald-600 text-white border-emerald-600 shadow-md shadow-emerald-600/20'
                  : 'bg-slate-50 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800'
              }`}
            >
              {st}
            </button>
          ))}
        </div>
      </div>

      {/* Info & Items Per Page Controls */}
      <div className="bg-white dark:bg-slate-900 p-4 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm flex flex-col sm:flex-row items-center justify-between gap-3 text-xs font-bold">
        <div className="text-slate-600 dark:text-slate-300">
          Menampilkan <span className="text-emerald-600 dark:text-emerald-400 font-black">{totalItems > 0 ? startIndex + 1 : 0}-{endIndex}</span> dari <span className="font-black text-slate-900 dark:text-white">{totalItems}</span> Transaksi Pemesanan
        </div>

        <div className="flex items-center gap-2">
          <span className="text-slate-400 font-medium text-[11px]">Tampilkan per halaman:</span>
          <select
            value={itemsPerPage}
            onChange={(e) => setItemsPerPage(Number(e.target.value))}
            className="px-3 py-1.5 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-xs font-extrabold text-slate-700 dark:text-slate-200 focus:outline-none"
          >
            <option value={5}>5 Data</option>
            <option value={10}>10 Data</option>
            <option value={20}>20 Data</option>
            <option value={50}>50 Data</option>
          </select>
        </div>
      </div>

      {/* Bookings Table */}
      <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-slate-200 dark:border-slate-800 text-[11px] font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider bg-slate-50/50 dark:bg-slate-800/50">
                <th className="py-3.5 px-5">ID Tiket</th>
                <th className="py-3.5 px-5">Nama Pelanggan</th>
                <th className="py-3.5 px-5">Destinasi</th>
                <th className="py-3.5 px-5">Tanggal Kunjungan</th>
                <th className="py-3.5 px-5">Jumlah Orang</th>
                <th className="py-3.5 px-5">Total Bayar</th>
                <th className="py-3.5 px-5">Status Pembayaran</th>
                <th className="py-3.5 px-5 text-right">Aksi Admin</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 dark:divide-slate-800 text-xs font-semibold">
              {paginatedBookings.length === 0 ? (
                <tr>
                  <td colSpan={8} className="py-8 text-center text-slate-400 font-medium">
                    Tidak ada pemesanan yang cocok dengan filter ini.
                  </td>
                </tr>
              ) : (
                paginatedBookings.map((b) => (
                  <tr key={b.id} className="hover:bg-slate-50/50 dark:hover:bg-slate-800/50 transition-colors">
                    <td className="py-3.5 px-5 font-bold text-emerald-600 dark:text-emerald-400">{b.id || b.bookingCode}</td>
                    <td className="py-3.5 px-5">
                      <div className="font-extrabold text-slate-900 dark:text-slate-100">{b.customerName || 'Customer'}</div>
                      <div className="text-[10px] text-slate-400 font-medium">{b.customerEmail || 'user@example.com'}</div>
                    </td>
                    <td className="py-3.5 px-5 text-slate-800 dark:text-slate-200 font-bold">{b.destinationName}</td>
                    <td className="py-3.5 px-5 text-slate-600 dark:text-slate-300">{b.date || b.bookedAt}</td>
                    <td className="py-3.5 px-5 text-slate-700 dark:text-slate-300 font-bold">{b.guests || b.totalPeople} Orang</td>
                    <td className="py-3.5 px-5 font-extrabold text-slate-900 dark:text-white">{formatIDR(b.totalPrice)}</td>
                    <td className="py-3.5 px-5">
                      <span
                        className={`px-2.5 py-1 rounded-full text-[10px] font-extrabold ${
                          String(b.status).toUpperCase() === 'PAID' || b.status === 'Lunas'
                            ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-700'
                            : String(b.status).toUpperCase() === 'PENDING' || b.status === 'Menunggu'
                            ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/80 dark:text-amber-300 border border-amber-300 dark:border-amber-700'
                            : 'bg-rose-100 text-rose-800 dark:bg-rose-950/80 dark:text-rose-300 border border-rose-300 dark:border-rose-700'
                        }`}
                      >
                        {String(b.status).toUpperCase() === 'PAID' || b.status === 'Lunas'
                          ? '✅ Lunas'
                          : String(b.status).toUpperCase() === 'PENDING' || b.status === 'Menunggu'
                          ? '⏳ Menunggu'
                          : String(b.status).toUpperCase() === 'EXPIRED'
                          ? '⌛ Expired'
                          : '❌ Dibatalkan'}
                      </span>
                    </td>
                    <td className="py-3.5 px-5 text-right">
                      {String(b.status).toUpperCase() === 'PENDING' || b.status === 'Menunggu' ? (
                        <button
                          type="button"
                          onClick={async () => {
                            try {
                              const code = b.bookingCode || b.rawBookingId || b.id;
                              await axiosClient.post(
                                '/webhooks/xendit',
                                {
                                  external_id: code,
                                  status: 'PAID',
                                  amount: b.totalPrice || 15000,
                                  payment_method: 'ADMIN_MANUAL',
                                },
                                {
                                  headers: {
                                    'x-callback-token': 'xnd_callback_dummytoken',
                                  },
                                }
                              );
                              const { toast } = await import('sonner');
                              const { useBookingStore } = await import('../../store/useBookingStore');
                              await useBookingStore.getState().fetchAllAdminBookings();
                              toast.success(`Pemesanan ${code} berhasil dikonfirmasi LUNAS!`);
                            } catch {
                              const { toast } = await import('sonner');
                              toast.error('Gagal memperbarui status ke Lunas.');
                            }
                          }}
                          className="px-3 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-black text-[11px] transition shadow-md shadow-emerald-600/20 active:scale-95"
                        >
                          ⚡ Konfirmasi Lunas
                        </button>
                      ) : (
                        <span className="text-[10px] text-slate-400 font-medium">Sudah Selesai</span>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination Controls */}
      {totalPages > 1 && (
        <div className="bg-white dark:bg-slate-900 p-4 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="text-xs font-extrabold text-slate-500 dark:text-slate-400">
            Halaman <span className="text-emerald-600 dark:text-emerald-400 font-black">{currentPage}</span> dari {totalPages}
          </div>

          <div className="flex items-center gap-1.5">
            <button
              onClick={() => setCurrentPage((p) => Math.max(p - 1, 1))}
              disabled={currentPage === 1}
              className="p-2 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
              title="Halaman Sebelumnya"
            >
              <ChevronLeft className="w-4 h-4" />
            </button>

            {Array.from({ length: totalPages }, (_, i) => i + 1).map((num) => (
              <button
                key={num}
                onClick={() => setCurrentPage(num)}
                className={`w-8 h-8 rounded-xl text-xs font-extrabold transition-all border ${
                  currentPage === num
                    ? 'bg-emerald-600 text-white border-emerald-600 shadow-md shadow-emerald-600/20'
                    : 'bg-slate-50 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-700'
                }`}
              >
                {num}
              </button>
            ))}

            <button
              onClick={() => setCurrentPage((p) => Math.min(p + 1, totalPages))}
              disabled={currentPage === totalPages}
              className="p-2 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
              title="Halaman Selanjutnya"
            >
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

