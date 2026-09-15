import React, { useState, useEffect } from 'react';
import { Ticket, X, Check, ShieldCheck } from 'lucide-react';
import { useDestinationStore } from '../../store/useDestinationStore';

interface TicketModalProps {
  selectedDest: any;
  onClose: () => void;
  bookingDate: string;
  setBookingDate: (date: string) => void;
  bookingTotalPeople: number;
  setBookingTotalPeople: (count: number) => void;
  onConfirm: () => void;
  formatIDR: (amount: number) => string;
}

export const TicketModal: React.FC<TicketModalProps> = ({
  selectedDest,
  onClose,
  bookingDate,
  setBookingDate,
  bookingTotalPeople,
  setBookingTotalPeople,
  onConfirm,
  formatIDR,
}) => {
  if (!selectedDest) return null;

  const { fetchSlots } = useDestinationStore();
  const [slotInfo, setSlotInfo] = useState<{ total: number; used: number; remaining: number } | null>(null);
  const [isLoadingSlot, setIsLoadingSlot] = useState(false);

  useEffect(() => {
    if (selectedDest && selectedDest.id) {
      const numId = Number(selectedDest.id);
      if (!isNaN(numId) && numId > 0) {
        setIsLoadingSlot(true);
        fetchSlots(numId)
          .then((slots) => {
            if (slots && slots.length > 0) {
              const found = slots.find((s) => s.date.startsWith(bookingDate)) || slots[0];
              const tot = found.total_quota ?? found.quota ?? 50;
              const usd = found.used_quota ?? found.booked_quota ?? 0;
              const rem = found.remaining_quota ?? (tot - usd);
              setSlotInfo({
                total: tot,
                used: usd,
                remaining: rem,
              });
            } else {
              setSlotInfo({ total: 50, used: 0, remaining: 50 });
            }
          })
          .finally(() => setIsLoadingSlot(false));
      } else {
        setSlotInfo({ total: 50, used: 0, remaining: 50 });
      }
    }
  }, [selectedDest, bookingDate, fetchSlots]);

  const remainingQuota = slotInfo?.remaining ?? 50;
  const isSlotFull = remainingQuota <= 0 || (slotInfo as any)?.isAvailable === false;
  const isExceeding = remainingQuota < bookingTotalPeople;

  return (
    <div className="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="w-full max-w-lg bg-white dark:bg-slate-900 rounded-3xl shadow-2xl border border-slate-200/80 dark:border-slate-800 p-6 space-y-5">
        <div className="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-4">
          <div className="flex items-center gap-3">
            <div className="p-2.5 rounded-2xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
              <Ticket className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base font-extrabold text-slate-900 dark:text-white">Formulir Pemesanan Tiket</h3>
              <p className="text-xs text-slate-500 font-medium">{selectedDest?.name || 'Destinasi Wisata'}</p>
            </div>
          </div>
          <button onClick={onClose} className="p-1 text-slate-400 hover:text-slate-600 dark:hover:text-white">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-xs font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1.5">
              Tanggal Kunjungan
            </label>
            <input
              type="date"
              value={bookingDate}
              onChange={(e) => setBookingDate(e.target.value)}
              className="w-full px-4 py-3 rounded-2xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 text-xs font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 text-slate-900 dark:text-white"
            />
          </div>

          <div>
            <label className="block text-xs font-extrabold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1.5">
              Jumlah Pengunjung (Rombongan)
            </label>
            <div className="flex items-center gap-3">
              {[1, 2, 3, 4, 5].map((num) => (
                <button
                  key={num}
                  type="button"
                  onClick={() => setBookingTotalPeople(num)}
                  className={`flex-1 py-2.5 rounded-xl text-xs font-extrabold border transition-all ${
                    bookingTotalPeople === num
                      ? 'bg-emerald-600 text-white border-emerald-600 shadow-md shadow-emerald-600/20'
                      : 'bg-slate-50 dark:bg-slate-800/80 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-700 hover:bg-slate-100'
                  }`}
                >
                  {num} Orang
                </button>
              ))}
            </div>
          </div>

          <div className="p-4 rounded-2xl bg-emerald-50/50 dark:bg-emerald-950/30 border border-emerald-500/20 space-y-2">
            <div className="flex items-center justify-between text-xs font-bold text-slate-600 dark:text-slate-300">
              <span>Status Kuota Slot Tanggal Ini:</span>
              <span className={`px-2.5 py-1 rounded-full text-[10px] font-black border flex items-center gap-1 ${
                isSlotFull || isExceeding
                  ? 'bg-rose-100 text-rose-800 dark:bg-rose-950/80 dark:text-rose-300 border-rose-300 dark:border-rose-700'
                  : 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300 border-emerald-300 dark:border-emerald-700'
              }`}>
                <ShieldCheck className="w-3.5 h-3.5" />
                <span>
                  {isLoadingSlot
                    ? 'Memuat slot...'
                    : isSlotFull
                    ? `❌ Kuota Tiket Penuh (0 / ${slotInfo?.total ?? 50} Tiket)`
                    : isExceeding
                    ? `⚠️ Sisa Kuota (${remainingQuota} Tiket) Kurang!`
                    : `⚡ Sisa Kuota: ${remainingQuota} / ${slotInfo?.total ?? 50} Tiket`}
                </span>
              </span>
            </div>

            {(isSlotFull || isExceeding) && (
              <div className="p-2.5 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-[11px] font-extrabold text-center">
                {isSlotFull
                  ? '⚠️ Maaf, kuota pemesanan tiket untuk tanggal ini sudah habis. Silakan pilih tanggal lain.'
                  : `⚠️ Sisa kuota hanya ${remainingQuota} tiket. Kurangi jumlah rombongan atau pilih tanggal lain.`}
              </div>
            )}

            <div className="flex items-center justify-between text-xs font-bold text-slate-600 dark:text-slate-300">
              <span>Harga per Tiket:</span>
              <span>{formatIDR(selectedDest?.price || 0)}</span>
            </div>
            <div className="flex items-center justify-between text-sm font-black text-slate-900 dark:text-white pt-1.5 border-t border-emerald-500/20">
              <span>Total Pembayaran:</span>
              <span className="text-emerald-600 dark:text-emerald-400">
                {formatIDR((selectedDest?.price || 0) * bookingTotalPeople)}
              </span>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-end gap-3 pt-2">
          <button
            type="button"
            onClick={onClose}
            className="py-2.5 px-4 rounded-xl text-xs font-extrabold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all"
          >
            Batal
          </button>
          <button
            type="button"
            disabled={isSlotFull || isExceeding || isLoadingSlot}
            onClick={onConfirm}
            className={`py-2.5 px-5 rounded-xl text-white text-xs font-extrabold shadow-lg transition-all flex items-center gap-1.5 ${
              isSlotFull || isExceeding || isLoadingSlot
                ? 'bg-slate-400 dark:bg-slate-700 cursor-not-allowed opacity-60 shadow-none'
                : 'bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 shadow-emerald-600/20 active:scale-95'
            }`}
          >
            <Check className="w-4 h-4" />
            <span>
              {isSlotFull ? 'Kuota Tiket Habis' : isExceeding ? 'Kuota Tidak Cukup' : 'Konfirmasi & Pesan Tiket'}
            </span>
          </button>
        </div>
      </div>
    </div>
  );
};
