import React from 'react';
import { CheckCircle2, Ticket, Printer, ArrowRight, X } from 'lucide-react';

interface PaymentSuccessModalProps {
  isOpen: boolean;
  onClose: () => void;
  booking: any;
  formatIDR: (amount: number) => string;
  onViewETicket: () => void;
  onPrintETicket: () => void;
}

export const PaymentSuccessModal: React.FC<PaymentSuccessModalProps> = ({
  isOpen,
  onClose,
  booking,
  formatIDR,
  onViewETicket,
  onPrintETicket,
}) => {
  if (!isOpen || !booking) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/85 backdrop-blur-md animate-fade-in text-white font-sans">
      <div className="w-full max-w-md bg-slate-900 border border-slate-700/80 rounded-3xl p-6 text-center space-y-5 shadow-2xl relative">
        <button
          onClick={onClose}
          className="absolute top-5 right-5 text-slate-400 hover:text-white p-1 rounded-lg"
        >
          <X className="w-5 h-5" />
        </button>

        {/* Success Icon Badge */}
        <div className="pt-2">
          <div className="inline-flex p-4 bg-emerald-500/10 text-emerald-400 rounded-3xl border border-emerald-500/20 shadow-xl shadow-emerald-500/10 mb-1">
            <CheckCircle2 className="w-14 h-14 animate-bounce text-emerald-400" />
          </div>
          <h2 className="text-2xl font-black bg-gradient-to-r from-white via-emerald-200 to-teal-200 bg-clip-text text-transparent mt-2">
            Pembayaran Berhasil!
          </h2>
          <p className="text-xs text-emerald-400 font-bold mt-0.5">
            Dikonfirmasi Resmi oleh Xendit Payment Gateway
          </p>
        </div>

        {/* Transaction Summary Card */}
        <div className="p-4 rounded-2xl bg-slate-950 border border-slate-800 text-left text-xs space-y-2">
          <div className="flex justify-between items-center pb-2 border-b border-slate-800">
            <span className="text-slate-400 font-medium">Kode Transaksi:</span>
            <span className="font-mono font-black text-sky-400">{booking.bookingCode || booking.booking_code || 'BK-2026-001'}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-slate-400 font-medium">Destinasi Wisata:</span>
            <span className="font-bold text-white truncate max-w-[180px]">{booking.destinationName || booking.destination_name}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-slate-400 font-medium">Jumlah Pengunjung:</span>
            <span className="font-bold text-white">{booking.totalPeople || booking.total_people} Orang</span>
          </div>
          <div className="flex justify-between items-center pt-2 border-t border-slate-800">
            <span className="text-xs font-black uppercase text-slate-400">Total Dibayar:</span>
            <span className="text-base font-black text-emerald-400">
              {formatIDR(booking.totalPrice || booking.total_price || 30000)}
            </span>
          </div>
        </div>

        {/* Actions */}
        <div className="space-y-2.5 pt-1">
          <button
            type="button"
            onClick={() => {
              onClose();
              onViewETicket();
            }}
            className="w-full py-3.5 rounded-2xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 text-white font-extrabold text-xs shadow-lg shadow-emerald-600/30 flex items-center justify-center gap-2 transition active:scale-95"
          >
            <Ticket className="w-4 h-4 text-white" />
            <span>Lihat E-Ticket Digital (QR Code)</span>
            <ArrowRight className="w-4 h-4" />
          </button>

          <button
            type="button"
            onClick={() => {
              onClose();
              onPrintETicket();
            }}
            className="w-full py-3 rounded-2xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold text-xs border border-slate-700 flex items-center justify-center gap-2 transition active:scale-95"
          >
            <Printer className="w-4 h-4 text-sky-400" />
            <span>Cetak / Unduh PDF Tiket</span>
          </button>
        </div>
      </div>
    </div>
  );
};
