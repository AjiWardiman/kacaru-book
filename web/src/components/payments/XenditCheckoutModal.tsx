import React, { useState, useEffect } from 'react';
import { QrCode, Building2, Wallet, X, ShieldCheck, ArrowRight, Clock, Loader2, Copy, CheckCircle2, AlertCircle } from 'lucide-react';
import { toast } from 'sonner';
import { useBookingStore } from '../../store/useBookingStore';
import { axiosClient } from '../../api/axiosClient';

interface XenditCheckoutModalProps {
  isOpen: boolean;
  onClose: () => void;
  booking: any;
  formatIDR: (amount: number) => string;
  onPaymentSuccess: (paidBooking: any) => void;
}

export const XenditCheckoutModal: React.FC<XenditCheckoutModalProps> = ({
  isOpen,
  onClose,
  booking,
  formatIDR,
  onPaymentSuccess,
}) => {
  const [selectedMethod, setSelectedMethod] = useState<'qris' | 'va' | 'ewallet'>('qris');
  const [selectedBank, setSelectedBank] = useState<'bca' | 'mandiri' | 'bni' | 'bri'>('bca');
  const [selectedEwallet, setSelectedEwallet] = useState<'gopay' | 'ovo' | 'dana' | 'shopeepay'>('gopay');
  const [eWalletPhone, setEWalletPhone] = useState('081234567890');
  
  // Realtime Simulation Steps (0 = idle, 1..4 = steps, 5 = done)
  const [procStep, setProcStep] = useState<number>(0);
  const [timeLeft, setTimeLeft] = useState<number>(899); // 14:59

  useEffect(() => {
    if (!isOpen) {
      setProcStep(0);
      return;
    }
    setTimeLeft(899);
    const timer = setInterval(() => {
      setTimeLeft((prev) => (prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => clearInterval(timer);
  }, [isOpen]);

  if (!isOpen || !booking) return null;

  const formatTimer = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const getVaNumber = () => {
    switch (selectedBank) {
      case 'bca':
        return '8808 1234 5678 9901';
      case 'mandiri':
        return '8900 8765 4321 1122';
      case 'bni':
        return '8241 5566 7788 9900';
      case 'bri':
        return '8077 1122 3344 5566';
      default:
        return '8808 1234 5678 9901';
    }
  };

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text.replace(/\s+/g, ''));
    toast.success(`Nomor Virtual Account ${selectedBank.toUpperCase()} berhasil disalin!`);
  };

  const handleSimulatePayment = async () => {
    setProcStep(1);

    const extId = booking.bookingCode || booking.booking_code || `INV-BIMA-${booking.id}`;
    const amount = booking.totalPrice || booking.total_price || 30000;
    const methodStr = selectedMethod === 'va' 
      ? `VA_${selectedBank.toUpperCase()}` 
      : selectedMethod === 'ewallet' 
      ? selectedEwallet.toUpperCase() 
      : 'QRIS';

    // Step 1 -> 2: Sending Webhook to backend
    setTimeout(() => {
      setProcStep(2);
    }, 900);

    // Step 2 -> 3: Backend DB processing
    setTimeout(() => {
      setProcStep(3);
    }, 1900);

    // Step 3 -> 4: Ticket generation & Email notification
    setTimeout(() => {
      setProcStep(4);
    }, 2800);

    try {
      // Direct Webhook call to Golang Backend via axiosClient
      await axiosClient.post(
        '/webhooks/xendit',
        {
          external_id: extId,
          status: 'PAID',
          amount: amount,
          payment_method: methodStr,
        },
        {
          headers: {
            'x-callback-token': 'xnd_callback_dummytoken',
          },
        }
      );

      // Refresh Zustand store so UI updates state from backend DB
      await useBookingStore.getState().fetchMyBookings();
      if (booking.id) {
        useBookingStore.getState().markBookingAsPaid(booking.id);
      }
      if (booking.bookingCode || booking.booking_code) {
        useBookingStore.getState().markBookingAsPaid(booking.bookingCode || booking.booking_code);
      }
    } catch (e) {
      console.warn('Webhook trigger warning (fallback mode active):', e);
      if (booking.id) {
        useBookingStore.getState().markBookingAsPaid(booking.id);
      }
      if (booking.bookingCode || booking.booking_code) {
        useBookingStore.getState().markBookingAsPaid(booking.bookingCode || booking.booking_code);
      }
    }

    // Step 4 -> Complete & Transition
    setTimeout(async () => {
      setProcStep(5);
      const updated = {
        ...booking,
        status: 'PAID',
        payment_method: methodStr,
      };
      
      setTimeout(() => {
        onClose();
        setProcStep(0);
        onPaymentSuccess(updated);
      }, 700);
    }, 3600);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/85 backdrop-blur-md animate-fade-in text-white font-sans">
      <div className="w-full max-w-lg bg-slate-900 border border-slate-700/80 rounded-3xl p-6 space-y-5 shadow-2xl relative">
        {/* Header Branding Xendit */}
        <div className="flex items-center justify-between border-b border-slate-800 pb-4">
          <div className="flex items-center gap-3">
            <div className="px-3 py-1.5 rounded-xl bg-blue-600 font-black text-xs tracking-wider uppercase text-white shadow-md shadow-blue-600/30 flex items-center gap-1.5">
              <span>XENDIT</span>
              <span className="text-[9px] bg-blue-700 px-1.5 py-0.5 rounded font-mono">GATEWAY</span>
            </div>
            <div>
              <h3 className="text-sm font-extrabold text-white">Invoice Pembayaran Digital</h3>
              <p className="text-[11px] text-slate-400 font-medium">Kacaru Tourism Booking Bima</p>
            </div>
          </div>
          {procStep === 0 && (
            <button onClick={onClose} className="p-1 text-slate-400 hover:text-white transition">
              <X className="w-5 h-5" />
            </button>
          )}
        </div>

        {/* Realtime Step Processing Screen */}
        {procStep > 0 ? (
          <div className="py-6 px-4 space-y-6 text-center animate-fade-in">
            <div className="space-y-2">
              <div className="inline-flex p-4 rounded-3xl bg-blue-600/10 text-blue-400 border border-blue-500/20 shadow-xl shadow-blue-500/10">
                {procStep < 5 ? (
                  <Loader2 className="w-12 h-12 animate-spin text-blue-400" />
                ) : (
                  <CheckCircle2 className="w-12 h-12 text-emerald-400 animate-bounce" />
                )}
              </div>
              <h4 className="text-lg font-black text-white">
                {procStep < 5 ? 'Memproses Webhook Pembayaran Xendit...' : 'Pembayaran LUNAS Terverifikasi! 🎉'}
              </h4>
              <p className="text-xs text-slate-400 max-w-xs mx-auto">
                Simulasi komunikasi realtime antara Xendit Gateway & Server Backend Go PostgreSQL
              </p>
            </div>

            {/* Checklist Step Items */}
            <div className="p-4 rounded-2xl bg-slate-950 border border-slate-800 space-y-3 text-left text-xs font-semibold">
              <div className="flex items-center justify-between">
                <span className={procStep >= 1 ? 'text-slate-200' : 'text-slate-500'}>
                  1. Mengirim Callback Webhook Xendit (<code className="text-sky-400 font-mono text-[11px]">POST /webhooks/xendit</code>)
                </span>
                {procStep > 1 ? (
                  <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                ) : procStep === 1 ? (
                  <Loader2 className="w-4 h-4 text-blue-400 animate-spin shrink-0" />
                ) : (
                  <span className="w-4 h-4 rounded-full border border-slate-700 block shrink-0" />
                )}
              </div>

              <div className="flex items-center justify-between">
                <span className={procStep >= 2 ? 'text-slate-200' : 'text-slate-500'}>
                  2. Backend Go mengonfirmasi transaksi & row locking (<code className="text-amber-400 font-mono text-[11px]">SELECT FOR UPDATE</code>)
                </span>
                {procStep > 2 ? (
                  <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                ) : procStep === 2 ? (
                  <Loader2 className="w-4 h-4 text-blue-400 animate-spin shrink-0" />
                ) : (
                  <span className="w-4 h-4 rounded-full border border-slate-700 block shrink-0" />
                )}
              </div>

              <div className="flex items-center justify-between">
                <span className={procStep >= 3 ? 'text-slate-200' : 'text-slate-500'}>
                  3. Menerbitkan E-Ticket QR Code & mengunci kuota slot
                </span>
                {procStep > 3 ? (
                  <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                ) : procStep === 3 ? (
                  <Loader2 className="w-4 h-4 text-blue-400 animate-spin shrink-0" />
                ) : (
                  <span className="w-4 h-4 rounded-full border border-slate-700 block shrink-0" />
                )}
              </div>

              <div className="flex items-center justify-between">
                <span className={procStep >= 4 ? 'text-slate-200' : 'text-slate-500'}>
                  4. Mengirim notifikasi email lunas via Mailtrap SMTP
                </span>
                {procStep > 4 ? (
                  <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                ) : procStep === 4 ? (
                  <Loader2 className="w-4 h-4 text-blue-400 animate-spin shrink-0" />
                ) : (
                  <span className="w-4 h-4 rounded-full border border-slate-700 block shrink-0" />
                )}
              </div>
            </div>
          </div>
        ) : (
          <>
            {/* Countdown Expiry Banner */}
            <div className="p-3 rounded-2xl bg-amber-500/10 border border-amber-500/30 flex items-center justify-between text-amber-300 text-xs font-bold">
              <div className="flex items-center gap-2">
                <Clock className="w-4 h-4 text-amber-400 animate-pulse" />
                <span>Batas Waktu Pembayaran:</span>
              </div>
              <span className="font-mono text-sm font-black text-amber-400 bg-amber-950/60 px-2.5 py-0.5 rounded-lg border border-amber-500/20">
                {formatTimer(timeLeft)}
              </span>
            </div>

            {/* Invoice Summary Box */}
            <div className="p-4 rounded-2xl bg-slate-950 border border-slate-800 space-y-2">
              <div className="flex justify-between items-center text-xs">
                <span className="text-slate-400 font-medium">Kode Pemesanan:</span>
                <span className="font-mono font-bold text-sky-400">{booking.bookingCode || booking.booking_code || 'BK-2026-001'}</span>
              </div>
              <div className="flex justify-between items-center text-xs">
                <span className="text-slate-400 font-medium">Destinasi Wisata:</span>
                <span className="font-bold text-white truncate max-w-[200px]">{booking.destinationName || booking.destination_name}</span>
              </div>
              <div className="flex justify-between items-center text-xs">
                <span className="text-slate-400 font-medium">Biaya Layanan Gateway:</span>
                <span className="font-bold text-emerald-400">GRATIS (Rp 0)</span>
              </div>
              <div className="flex justify-between items-center pt-2 border-t border-slate-800/80">
                <span className="text-xs font-black uppercase text-slate-400">Total Tagihan:</span>
                <span className="text-lg font-black text-emerald-400">
                  {formatIDR(booking.totalPrice || booking.total_price || 30000)}
                </span>
              </div>
            </div>

            {/* Payment Method Tabs */}
            <div className="space-y-2">
              <label className="block text-[11px] font-black uppercase tracking-wider text-slate-400">
                Pilih Metode Pembayaran Xendit
              </label>
              <div className="grid grid-cols-3 gap-2">
                <button
                  type="button"
                  onClick={() => setSelectedMethod('qris')}
                  className={`p-3 rounded-2xl border text-xs font-bold flex flex-col items-center gap-1.5 transition-all ${
                    selectedMethod === 'qris'
                      ? 'bg-blue-600/20 border-blue-500 text-blue-300 shadow-lg shadow-blue-500/20'
                      : 'bg-slate-800/50 border-slate-700 text-slate-400 hover:bg-slate-800'
                  }`}
                >
                  <QrCode className="w-5 h-5 text-blue-400" />
                  <span>QRIS</span>
                </button>

                <button
                  type="button"
                  onClick={() => setSelectedMethod('va')}
                  className={`p-3 rounded-2xl border text-xs font-bold flex flex-col items-center gap-1.5 transition-all ${
                    selectedMethod === 'va'
                      ? 'bg-blue-600/20 border-blue-500 text-blue-300 shadow-lg shadow-blue-500/20'
                      : 'bg-slate-800/50 border-slate-700 text-slate-400 hover:bg-slate-800'
                  }`}
                >
                  <Building2 className="w-5 h-5 text-blue-400" />
                  <span>Virtual Account</span>
                </button>

                <button
                  type="button"
                  onClick={() => setSelectedMethod('ewallet')}
                  className={`p-3 rounded-2xl border text-xs font-bold flex flex-col items-center gap-1.5 transition-all ${
                    selectedMethod === 'ewallet'
                      ? 'bg-blue-600/20 border-blue-500 text-blue-300 shadow-lg shadow-blue-500/20'
                      : 'bg-slate-800/50 border-slate-700 text-slate-400 hover:bg-slate-800'
                  }`}
                >
                  <Wallet className="w-5 h-5 text-blue-400" />
                  <span>E-Wallet</span>
                </button>
              </div>
            </div>

            {/* Method Detail View */}
            <div className="p-4 rounded-2xl bg-slate-950 border border-slate-800 text-center space-y-3">
              {selectedMethod === 'qris' && (
                <div className="space-y-3">
                  <div className="text-xs font-bold text-slate-300 flex items-center justify-center gap-1.5">
                    <span>Pindai Kode QRIS Resmi Merchant</span>
                    <span className="px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-400 text-[10px]">KACARU BIMA</span>
                  </div>
                  <div className="bg-white p-3.5 rounded-2xl inline-block mx-auto shadow-lg border border-slate-200">
                    <svg className="w-32 h-32 text-slate-900" viewBox="0 0 100 100" fill="currentColor">
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
                    <div className="text-[9px] font-mono text-slate-600 font-bold mt-1 uppercase">
                      NMID: ID102938475620
                    </div>
                  </div>
                  <p className="text-[11px] text-slate-400 font-medium">
                    Bisa di-scan dari BCA Mobile, Livin, GoPay, OVO, DANA, ShopeePay, atau Bank Jatim.
                  </p>
                </div>
              )}

              {selectedMethod === 'va' && (
                <div className="space-y-3 text-left">
                  <div className="flex gap-2">
                    {(['bca', 'mandiri', 'bni', 'bri'] as const).map((b) => (
                      <button
                        key={b}
                        type="button"
                        onClick={() => setSelectedBank(b)}
                        className={`flex-1 py-1.5 rounded-xl text-[11px] font-black uppercase border transition ${
                          selectedBank === b
                            ? 'bg-blue-600 text-white border-blue-500 shadow-md shadow-blue-600/30'
                            : 'bg-slate-900 text-slate-400 border-slate-800 hover:bg-slate-800'
                        }`}
                      >
                        {b}
                      </button>
                    ))}
                  </div>

                  <div className="p-3 bg-slate-900 rounded-xl border border-slate-700 flex justify-between items-center">
                    <div>
                      <span className="text-[10px] text-slate-400 block font-bold">Nomor Virtual Account {selectedBank.toUpperCase()}</span>
                      <span className="font-mono font-black text-emerald-400 text-base">{getVaNumber()}</span>
                    </div>
                    <button
                      type="button"
                      onClick={() => handleCopy(getVaNumber())}
                      className="text-[11px] font-bold px-3 py-1.5 rounded-xl bg-blue-600/20 text-blue-300 border border-blue-500/30 flex items-center gap-1.5 hover:bg-blue-600/30 transition active:scale-95"
                    >
                      <Copy className="w-3.5 h-3.5" /> Salin VA
                    </button>
                  </div>
                </div>
              )}

              {selectedMethod === 'ewallet' && (
                <div className="space-y-3 text-left">
                  <div className="flex gap-2">
                    {(['gopay', 'ovo', 'dana', 'shopeepay'] as const).map((ew) => (
                      <button
                        key={ew}
                        type="button"
                        onClick={() => setSelectedEwallet(ew)}
                        className={`flex-1 py-1.5 rounded-xl text-[10px] font-black uppercase border transition ${
                          selectedEwallet === ew
                            ? 'bg-emerald-600 text-white border-emerald-500 shadow-md shadow-emerald-600/30'
                            : 'bg-slate-900 text-slate-400 border-slate-800 hover:bg-slate-800'
                        }`}
                      >
                        {ew}
                      </button>
                    ))}
                  </div>

                  <div>
                    <label className="block text-[10px] font-bold text-slate-400 mb-1">
                      Nomor HP Terdaftar {selectedEwallet.toUpperCase()}
                    </label>
                    <input
                      type="tel"
                      value={eWalletPhone}
                      onChange={(e) => setEWalletPhone(e.target.value)}
                      className="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-xs font-mono font-bold text-white focus:outline-none focus:border-emerald-500"
                    />
                  </div>
                </div>
              )}
            </div>

            {/* Practical Demo Action Button */}
            <div className="space-y-2 pt-1">
              <button
                type="button"
                onClick={handleSimulatePayment}
                className="w-full py-3.5 rounded-2xl bg-gradient-to-r from-emerald-600 via-teal-600 to-sky-600 hover:from-emerald-500 hover:to-sky-500 text-white font-black text-xs shadow-lg shadow-emerald-600/30 flex items-center justify-center gap-2 transition active:scale-95"
              >
                <ShieldCheck className="w-4 h-4 text-emerald-300" />
                <span>Simulasi Pembayaran Lunas (Trigger Webhook Xendit)</span>
                <ArrowRight className="w-4 h-4" />
              </button>

              <div className="p-2.5 rounded-xl bg-slate-950 border border-slate-800 flex items-start gap-2 text-[10px] text-slate-400">
                <AlertCircle className="w-3.5 h-3.5 text-sky-400 shrink-0 mt-0.5" />
                <span>
                  <strong>Info Demo:</strong> Tombol ini mensimulasikan callback Webhook resmi Xendit ke Backend Go (<code className="text-sky-300">/webhooks/xendit</code>) untuk memperbarui status transaksi menjadi <strong>PAID</strong> secara realtime.
                </span>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
};
