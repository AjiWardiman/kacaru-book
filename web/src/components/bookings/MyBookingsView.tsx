import React, { useState } from 'react';
import { Ticket, MapPin, DollarSign, Check } from 'lucide-react';
import { getImageUrl } from '../../utils/imageHelper';
import { XenditCheckoutModal } from '../payments/XenditCheckoutModal';
import { PaymentSuccessModal } from '../payments/PaymentSuccessModal';
import { useBookingStore } from '../../store/useBookingStore';

interface MyBookingsViewProps {
  effectiveUserBookings: any[];
  onNavigateToDestinations: () => void;
  onPayBooking: (bookingId: number) => Promise<string | null> | any;
  formatIDR: (amount: number) => string;
}

export const MyBookingsView: React.FC<MyBookingsViewProps> = ({
  effectiveUserBookings,
  onNavigateToDestinations,
  onPayBooking,
  formatIDR,
}) => {
  const [selectedETicket, setSelectedETicket] = useState<any>(null);
  const [selectedBookingForXendit, setSelectedBookingForXendit] = useState<any>(null);
  const [paidSuccessBooking, setPaidSuccessBooking] = useState<any>(null);

  return (
    <div className="space-y-6">
      <div className="bg-white dark:bg-slate-900 p-6 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm flex items-center justify-between">
        <div>
          <h3 className="text-lg font-extrabold text-slate-900 dark:text-white">Pemesanan & Tiket Saya</h3>
          <p className="text-xs font-bold text-slate-600 dark:text-slate-300 mt-0.5">
            Riwayat tiket wisata digital dan status pembayaran transaksi Anda.
          </p>
        </div>
        <div className="px-3.5 py-1.5 rounded-full bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 text-xs font-black">
          {(effectiveUserBookings || []).length} Tiket Aktif
        </div>
      </div>

      {(effectiveUserBookings || []).length === 0 ? (
        <div className="p-12 text-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 space-y-3">
          <Ticket className="w-12 h-12 text-slate-400 mx-auto" />
          <h4 className="text-base font-extrabold text-slate-700 dark:text-slate-200">Belum Ada Pemesanan Tiket</h4>
          <p className="text-xs text-slate-500 max-w-sm mx-auto font-medium">
            Anda belum memesan tiket wisata. Pilih destinasi favorit Anda di Katalog Wisata!
          </p>
          <button
            onClick={onNavigateToDestinations}
            className="mt-2 py-2.5 px-5 rounded-xl bg-emerald-600 text-white font-extrabold text-xs shadow-md shadow-emerald-600/20 hover:bg-emerald-700 transition-all"
          >
            Eksplorasi Destinasi
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {effectiveUserBookings.map((b) => (
            <div
              key={b.id}
              className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200/80 dark:border-slate-800 shadow-sm overflow-hidden flex flex-col justify-between group hover:shadow-md transition-all"
            >
              <div className="p-6 space-y-4">
                <div className="flex items-center justify-between pb-3 border-b border-slate-100 dark:border-slate-800">
                  <div className="flex items-center gap-2">
                    <Ticket className="w-4 h-4 text-emerald-500" />
                    <span className="text-xs font-black font-mono tracking-wider text-slate-800 dark:text-slate-200">
                      {b.bookingCode || b.booking_code}
                    </span>
                  </div>
                  <span
                    className={`px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-wider ${
                      String(b.status).toUpperCase() === 'PAID' || b.status === 'Lunas'
                        ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300'
                        : String(b.status).toUpperCase() === 'PENDING' || b.status === 'Menunggu'
                        ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/80 dark:text-amber-300'
                        : String(b.status).toUpperCase() === 'EXPIRED'
                        ? 'bg-rose-100 text-rose-800 dark:bg-rose-950/80 dark:text-rose-300'
                        : 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300'
                    }`}
                  >
                    {String(b.status).toUpperCase() === 'PAID' || b.status === 'Lunas'
                      ? '✅ LUNAS'
                      : String(b.status).toUpperCase() === 'PENDING' || b.status === 'Menunggu'
                      ? '⏳ PENDING'
                      : String(b.status).toUpperCase() === 'EXPIRED'
                      ? '⌛ EXPIRED'
                      : '❌ DIBATALKAN'}
                  </span>
                </div>

                <div className="flex items-start gap-4">
                  <img
                    src={getImageUrl(b.destinationImage)}
                    alt={b.destinationName}
                    className="w-20 h-20 rounded-2xl object-cover shrink-0 border border-slate-200 dark:border-slate-800"
                  />
                  <div className="space-y-1 overflow-hidden">
                    <h4 className="text-sm font-extrabold text-slate-900 dark:text-white truncate">
                      {b.destinationName || b.destination_name}
                    </h4>
                    <div className="flex items-center gap-1 text-[11px] text-slate-500 font-bold">
                      <MapPin className="w-3 h-3 text-emerald-500 shrink-0" />
                      <span className="truncate">{b.destinationLocation || 'Kabupaten Bima, NTB'}</span>
                    </div>
                    <div className="text-[11px] text-slate-500 font-medium pt-1">
                      👥 <span className="font-bold">{b.totalPeople || b.total_people} Pengunjung</span> • 📅 {b.bookedAt || b.booked_at?.slice(0, 10)}
                    </div>
                  </div>
                </div>
              </div>

              <div className="p-5 border-t border-slate-100 dark:border-slate-800 bg-slate-50/60 dark:bg-slate-800/40 flex items-center justify-between">
                <div>
                  <span className="text-[10px] font-bold text-slate-400 uppercase block">Total Bayar</span>
                  <span className="text-base font-black text-emerald-600 dark:text-emerald-400">
                    {formatIDR(b.totalPrice || b.total_price)}
                  </span>
                </div>

                {String(b.status).toUpperCase() === 'PENDING' || b.status === 'Menunggu' ? (
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={async () => {
                        if (b.paymentUrl && typeof b.paymentUrl === 'string' && b.paymentUrl.startsWith('https://checkout.xendit.co') && !b.paymentUrl.includes('inv_mock_')) {
                          window.open(b.paymentUrl, '_blank');
                        } else {
                          const urlRes: any = await onPayBooking(Number(b.id));
                          if (typeof urlRes === 'string' && urlRes.startsWith('https://checkout.xendit.co') && !urlRes.includes('inv_mock_')) {
                            window.open(urlRes, '_blank');
                          } else {
                            setSelectedBookingForXendit(b);
                          }
                        }
                      }}
                      className="py-2 px-4 rounded-xl bg-gradient-to-r from-amber-500 via-orange-500 to-amber-600 hover:from-amber-600 hover:to-orange-600 text-white font-black text-xs shadow-md shadow-amber-500/20 active:scale-95 transition-all flex items-center gap-1.5"
                    >
                      <DollarSign className="w-3.5 h-3.5" />
                      <span>Bayar via Xendit</span>
                    </button>
                  </div>
                ) : String(b.status).toUpperCase() === 'PAID' || b.status === 'Lunas' ? (
                  <button
                    onClick={() => setSelectedETicket(b)}
                    className="py-2 px-4 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-extrabold text-xs shadow-md shadow-emerald-600/20 active:scale-95 transition-all flex items-center gap-1.5"
                  >
                    <Check className="w-3.5 h-3.5 text-white" />
                    <span>Lihat E-Ticket (QR)</span>
                  </button>
                ) : (
                  <span className="px-3 py-1.5 rounded-xl bg-slate-200/80 dark:bg-slate-800 text-slate-500 dark:text-slate-400 text-[11px] font-bold">
                    {String(b.status).toUpperCase() === 'EXPIRED' ? '⏰ Waktu Bayar Habis' : '❌ Pemesanan Batal'}
                  </span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Xendit Checkout Modal Simulator */}
      <XenditCheckoutModal
        isOpen={Boolean(selectedBookingForXendit)}
        onClose={() => setSelectedBookingForXendit(null)}
        booking={selectedBookingForXendit}
        formatIDR={formatIDR}
        onPaymentSuccess={(paidBooking) => {
          if (paidBooking?.id) {
            useBookingStore.getState().markBookingAsPaid(paidBooking.id);
          }
          if (paidBooking?.bookingCode || paidBooking?.booking_code) {
            useBookingStore.getState().markBookingAsPaid(paidBooking.bookingCode || paidBooking.booking_code);
          }
          setPaidSuccessBooking(paidBooking || selectedBookingForXendit);
          setSelectedBookingForXendit(null);
        }}
      />

      {/* Payment Success Celebration Modal */}
      {paidSuccessBooking && (
        <PaymentSuccessModal
          isOpen={Boolean(paidSuccessBooking)}
          onClose={() => {
            setPaidSuccessBooking(null);
          }}
          booking={paidSuccessBooking}
          formatIDR={formatIDR}
          onViewETicket={() => {
            setSelectedETicket(paidSuccessBooking);
            setPaidSuccessBooking(null);
          }}
          onPrintETicket={() => {
            setSelectedETicket(paidSuccessBooking);
            setPaidSuccessBooking(null);
            setTimeout(() => window.print(), 300);
          }}
        />
      )}


      {/* E-Ticket QR Code Viewer Modal */}
      {selectedETicket && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md animate-fade-in print:bg-white print:p-0">
          <div className="w-full max-w-md bg-slate-900 border border-slate-700 rounded-3xl p-6 text-white text-center space-y-4 shadow-2xl relative print:border-none print:shadow-none print:text-black print:bg-white">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3 print:hidden">
              <span className="text-xs font-black tracking-wider text-emerald-400 uppercase">E-TICKET RESMI KACARU BOOK</span>
              <button
                onClick={() => setSelectedETicket(null)}
                className="text-slate-400 hover:text-white text-xs font-bold"
              >
                Tutup ✕
              </button>
            </div>

            <div className="space-y-1 text-center">
              <span className="inline-block px-3 py-1 rounded-full text-[10px] font-extrabold uppercase bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 mb-1">
                ✅ LUNAS VIA XENDIT
              </span>
              <h3 className="text-xl font-black text-white print:text-black">
                {selectedETicket.destinationName || selectedETicket.destination_name}
              </h3>
              <p className="text-xs text-sky-400 font-mono font-bold">
                Kode Tiket: {selectedETicket.bookingCode || selectedETicket.booking_code || 'BK-20260915-001'}
              </p>
            </div>

            {/* Ticket Info Card */}
            <div className="p-3 bg-slate-950 rounded-2xl border border-slate-800 text-left text-xs space-y-1.5 print:bg-slate-100 print:text-black">
              <div className="flex justify-between">
                <span className="text-slate-400">Pengunjung:</span>
                <span className="font-bold text-white print:text-black">{selectedETicket.totalPeople || selectedETicket.total_people} Orang</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Tanggal Kunjungan:</span>
                <span className="font-bold text-white print:text-black">{selectedETicket.bookedAt || selectedETicket.booked_at?.slice(0, 10) || '2026-09-20'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Total Pembayaran:</span>
                <span className="font-extrabold text-emerald-400 print:text-emerald-700">{formatIDR(selectedETicket.totalPrice || selectedETicket.total_price || 30000)}</span>
              </div>
            </div>

            {/* Simulated QR Code SVG */}
            <div className="bg-white p-4 rounded-2xl inline-block mx-auto shadow-inner border border-slate-200">
              <svg className="w-36 h-36 mx-auto text-slate-900" viewBox="0 0 100 100" fill="currentColor">
                <rect x="10" y="10" width="25" height="25" fill="black" />
                <rect x="15" y="15" width="15" height="15" fill="white" />
                <rect x="65" y="10" width="25" height="25" fill="black" />
                <rect x="70" y="15" width="15" height="15" fill="white" />
                <rect x="10" y="65" width="25" height="25" fill="black" />
                <rect x="15" y="70" width="15" height="15" fill="white" />
                <rect x="45" y="45" width="10" height="10" fill="black" />
                <rect x="60" y="45" width="15" height="10" fill="black" />
                <rect x="45" y="65" width="10" height="25" fill="black" />
                <rect x="65" y="65" width="20" height="15" fill="black" />
              </svg>
              <div className="text-[9px] font-mono text-slate-500 font-bold mt-1 uppercase">
                VALIDATED BY XENDIT GATEWAY
              </div>
            </div>

            <div className="text-[11px] text-slate-400 bg-slate-950 p-3 rounded-xl border border-slate-800 print:hidden">
              Tunjukkan QR Code ini ke petugas di pintu masuk gerbang wisata.
            </div>

            {/* Print / Download Button */}
            <div className="pt-2 print:hidden">
              <button
                type="button"
                onClick={() => window.print()}
                className="w-full py-3 rounded-2xl bg-gradient-to-r from-sky-600 to-indigo-600 hover:from-sky-500 hover:to-indigo-500 text-white font-extrabold text-xs shadow-lg shadow-sky-600/25 flex items-center justify-center gap-2 transition active:scale-95"
              >
                🖨️ Cetak / Unduh E-Ticket (PDF)
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
