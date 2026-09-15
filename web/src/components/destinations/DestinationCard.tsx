import React from 'react';
import { MapPin, Ticket, Star, Trash2, Pencil } from 'lucide-react';
import { getImageUrl } from '../../utils/imageHelper';

interface DestinationCardProps {
  dest: any; // Flexible item with id, name, location, price, category, rating, image, status, description
  isAdmin: boolean;
  onBookTicket: (dest: any) => void;
  onEditDestination?: (dest: any) => void;
  onDeleteDestination: (id: string | number) => void;
  formatIDR: (amount: number) => string;
}

export const DestinationCard: React.FC<DestinationCardProps> = ({
  dest,
  isAdmin,
  onBookTicket,
  onEditDestination,
  onDeleteDestination,
  formatIDR,
}) => {
  return (
    <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm overflow-hidden group hover:shadow-lg transition-all duration-300 flex flex-col justify-between">
      <div>
        <div className="relative h-48 overflow-hidden">
          <img
            src={getImageUrl(dest?.image || dest?.image_url)}
            alt={dest?.name || 'Destinasi'}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
          />
          <div className="absolute top-3 left-3 px-3 py-1 rounded-full bg-slate-900/70 backdrop-blur-md text-white text-[10px] font-extrabold">
            {dest?.category || 'Bahari'}
          </div>
          <div className="absolute top-3 right-3 px-2.5 py-1 rounded-full bg-amber-500 text-white text-[10px] font-black flex items-center gap-1 shadow-md">
            <Star className="w-3 h-3 fill-white" />
            <span>{dest?.rating || 4.8}</span>
          </div>
        </div>

        <div className="p-5 space-y-2">
          <h4 className="text-base font-extrabold text-slate-900 dark:text-white group-hover:text-emerald-600 transition-colors">
            {dest?.name}
          </h4>
          <div className="flex items-center justify-between pt-1">
            <div className="flex items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400 font-bold">
              <MapPin className="w-3.5 h-3.5 text-emerald-500" />
              <span>{dest?.location || 'Kabupaten Bima'}</span>
            </div>
            <div className="px-2 py-0.5 rounded-lg text-[10px] font-black bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20">
              {dest?.remaining_quota !== undefined ? `⚡ Sisa Kuota: ${dest.remaining_quota} Tiket` : '⚡ Kuota Tersedia Harian'}
            </div>
          </div>
          <p className="text-xs text-slate-600 dark:text-slate-300 font-medium line-clamp-2 pt-1">
            {dest?.description || 'Destinasi unggulan Kabupaten Bima NTB.'}
          </p>
        </div>
      </div>

      <div className="p-5 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-800/30">
        <div>
          <div className="text-[10px] font-bold text-slate-400 uppercase">Harga Tiket</div>
          <div className="text-base font-black text-emerald-600 dark:text-emerald-400">
            {formatIDR(dest?.price || 0)}
          </div>
        </div>

        {!isAdmin ? (
          <button
            onClick={() => onBookTicket(dest)}
            className="py-2.5 px-4 rounded-xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 text-white font-extrabold text-xs shadow-md shadow-emerald-600/20 active:scale-95 transition-all flex items-center gap-1.5"
          >
            <Ticket className="w-3.5 h-3.5" />
            <span>Pesan Tiket</span>
          </button>
        ) : (
          <div className="flex items-center gap-2">
            <span
              className={`px-3 py-1 rounded-full text-[10px] font-extrabold ${
                dest?.is_active !== false && dest?.status !== 'Penuh'
                  ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300'
                  : 'bg-rose-100 text-rose-800 dark:bg-rose-950/80 dark:text-rose-300'
              }`}
            >
              {dest?.status || 'Tersedia'}
            </span>

            {onEditDestination && (
              <button
                onClick={() => onEditDestination(dest)}
                className="p-1.5 rounded-lg text-sky-500 hover:bg-sky-50 dark:hover:bg-sky-950/50 transition-colors"
                title="Edit Destinasi"
              >
                <Pencil className="w-4 h-4" />
              </button>
            )}

            <button
              onClick={() => onDeleteDestination(dest?.id)}
              className="p-1.5 rounded-lg text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/50 transition-colors"
              title="Hapus Destinasi"
            >
              <Trash2 className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>
    </div>
  );
};
