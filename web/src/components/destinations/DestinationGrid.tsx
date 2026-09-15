import React, { useState, useEffect } from 'react';
import { Plus, Sparkles, ChevronLeft, ChevronRight, Search, X } from 'lucide-react';
import { DestinationCard } from './DestinationCard';
import { WeatherWidget } from '../weather/WeatherWidget';
import { AIRecommendationModal } from '../ai/AIRecommendationModal';
import type { UserResponse } from '../../api/schema';

interface DestinationGridProps {
  user: UserResponse | null;
  filteredDestinations: any[];
  destCategoryFilter: string;
  setDestCategoryFilter: (cat: string) => void;
  destSearch: string;
  setDestSearch: (q: string) => void;
  onOpenAddModal: () => void;
  onBookTicket: (dest: any) => void;
  onEditDestination?: (dest: any) => void;
  onDeleteDestination: (id: string | number) => void;
  formatIDR: (amount: number) => string;
}

export const DestinationGrid: React.FC<DestinationGridProps> = ({
  user,
  filteredDestinations,
  destCategoryFilter,
  setDestCategoryFilter,
  destSearch,
  setDestSearch,
  onOpenAddModal,
  onBookTicket,
  onEditDestination,
  onDeleteDestination,
  formatIDR,
}) => {
  const isAdmin = user?.role === 'admin';
  const [isAIModalOpen, setIsAIModalOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const ITEMS_PER_PAGE = 6;

  // Reset ke halaman 1 saat kategori filter atau kata kunci pencarian berubah
  useEffect(() => {
    setCurrentPage(1);
  }, [destCategoryFilter, destSearch]);

  const totalItems = filteredDestinations.length;
  const totalPages = Math.max(1, Math.ceil(totalItems / ITEMS_PER_PAGE));
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = Math.min(startIndex + ITEMS_PER_PAGE, totalItems);

  const paginatedDestinations = filteredDestinations.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  return (
    <div className="space-y-6">
      {/* Weather Forecast Banner */}
      <WeatherWidget location="Bima, NTB" />

      {/* View Header */}
      <div className="bg-white dark:bg-slate-900 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
        <div>
          <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">
            {isAdmin ? 'Kelola Destinasi Wisata' : 'Eksplorasi Destinasi Wisata'}
          </h3>
          <p className="text-xs font-bold text-slate-600 dark:text-slate-300 mt-0.5">
            {isAdmin
              ? 'Kelola informasi lokasi, harga tiket, dan ketersediaan slot destinasi.'
              : 'Jelajahi tempat wisata terbaik di Kabupaten Bima & NTB dan pesan tiket online.'}
          </p>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={() => setIsAIModalOpen(true)}
            className="py-3 px-5 rounded-2xl bg-gradient-to-r from-sky-600 via-indigo-600 to-purple-600 hover:from-sky-500 hover:to-purple-500 text-white font-extrabold text-xs shadow-lg shadow-sky-600/20 active:scale-95 transition-all flex items-center justify-center gap-2"
          >
            <Sparkles className="w-4 h-4 text-amber-300 animate-pulse" />
            <span>✨ Rekomendasi AI</span>
          </button>

          {isAdmin && (
            <button
              onClick={onOpenAddModal}
              className="py-3 px-5 rounded-2xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white font-extrabold text-xs shadow-lg shadow-emerald-600/20 active:scale-95 transition-all flex items-center justify-center gap-2"
            >
              <Plus className="w-4 h-4" />
              <span>Tambah Destinasi</span>
            </button>
          )}
        </div>
      </div>

      <AIRecommendationModal
        isOpen={isAIModalOpen}
        onClose={() => setIsAIModalOpen(false)}
        onSelectDestination={(destId) => {
          const matched = filteredDestinations.find((d) => d.id === destId);
          if (matched) onBookTicket(matched);
        }}
      />

      {/* Search Input Bar & Category Filter Pills */}
      <div className="space-y-4">
        {/* Search Bar */}
        <div className="relative">
          <Search className="w-4 h-4 text-slate-400 absolute left-4 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={destSearch}
            onChange={(e) => setDestSearch(e.target.value)}
            placeholder="Cari nama destinasi, lokasi, atau kata kunci (misal: Lariti, Tambora, Sape)..."
            className="w-full pl-11 pr-10 py-3 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white shadow-sm"
          />
          {destSearch && (
            <button
              onClick={() => setDestSearch('')}
              className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-slate-400 hover:text-slate-600 dark:hover:text-white"
              title="Bersihkan Pencarian"
            >
              <X className="w-4 h-4" />
            </button>
          )}
        </div>

        {/* Filter Pills & Counter */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-center gap-2 overflow-x-auto pb-1">
            {['ALL', 'Bahari', 'Gunung & Alam', 'Budaya & Geowisata'].map((cat) => (
              <button
                key={cat}
                onClick={() => setDestCategoryFilter(cat)}
                className={`px-4 py-2 rounded-2xl text-xs font-extrabold transition-all border whitespace-nowrap ${
                  destCategoryFilter === cat
                    ? 'bg-emerald-600 text-white border-emerald-600 shadow-md shadow-emerald-600/20'
                    : 'bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800'
                }`}
              >
                {cat === 'ALL' ? 'Semua Kategori' : cat}
              </button>
            ))}
          </div>

          <div className="text-xs font-bold text-slate-500 dark:text-slate-400 whitespace-nowrap">
            Menampilkan <span className="text-emerald-600 dark:text-emerald-400 font-black">{totalItems > 0 ? startIndex + 1 : 0}-{endIndex}</span> dari <span className="font-black text-slate-800 dark:text-slate-200">{totalItems}</span> Destinasi
          </div>
        </div>
      </div>

      {/* Destinations Grid */}
      {paginatedDestinations.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {paginatedDestinations.map((dest) => (
            <DestinationCard
              key={dest.id}
              dest={dest}
              isAdmin={isAdmin}
              onBookTicket={onBookTicket}
              onEditDestination={onEditDestination}
              onDeleteDestination={onDeleteDestination}
              formatIDR={formatIDR}
            />
          ))}
        </div>
      ) : (
        <div className="bg-white dark:bg-slate-900 p-12 rounded-3xl border border-slate-200/80 dark:border-slate-800 text-center space-y-3">
          <p className="text-sm font-bold text-slate-500 dark:text-slate-400">
            Tidak ada destinasi ditemukan untuk kategori ini.
          </p>
        </div>
      )}

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
              className="p-2.5 rounded-xl border border-slate-200 dark:border-slate-800 text-xs font-bold disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-100 dark:hover:bg-slate-800 transition flex items-center gap-1 text-slate-700 dark:text-slate-200"
            >
              <ChevronLeft className="w-4 h-4" />
              <span>Sebelumnya</span>
            </button>

            <div className="flex items-center gap-1 px-1">
              {Array.from({ length: totalPages }, (_, i) => i + 1).map((pageNum) => (
                <button
                  key={pageNum}
                  onClick={() => setCurrentPage(pageNum)}
                  className={`w-8 h-8 rounded-xl text-xs font-black transition-all ${
                    currentPage === pageNum
                      ? 'bg-emerald-600 text-white shadow-md shadow-emerald-600/20'
                      : 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'
                  }`}
                >
                  {pageNum}
                </button>
              ))}
            </div>

            <button
              onClick={() => setCurrentPage((p) => Math.min(p + 1, totalPages))}
              disabled={currentPage === totalPages}
              className="p-2.5 rounded-xl border border-slate-200 dark:border-slate-800 text-xs font-bold disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-100 dark:hover:bg-slate-800 transition flex items-center gap-1 text-slate-700 dark:text-slate-200"
            >
              <span>Berikutnya</span>
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
