import React, { useState } from 'react';
import { ShieldCheck, X, ArrowRight, RotateCw } from 'lucide-react';
import { toast } from 'sonner';

interface OTPVerificationModalProps {
  isOpen: boolean;
  onClose: () => void;
  email: string;
  phone?: string;
  onVerifySuccess: () => void;
}

export const OTPVerificationModal: React.FC<OTPVerificationModalProps> = ({
  isOpen,
  onClose,
  email,
  phone,
  onVerifySuccess,
}) => {
  const [otpCode, setOtpCode] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [resending, setResending] = useState(false);

  if (!isOpen) return null;

  const handleVerify = async (e: React.FormEvent) => {
    e.preventDefault();
    if (otpCode.length < 6) {
      toast.error('Masukkan 6 digit kode OTP verifikasi!');
      return;
    }

    setLoading(true);
    try {
      // Panggil backend API /otp/verify
      const token = localStorage.getItem('kacaru_token');
      const res = await fetch('http://localhost:8080/otp/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          code: otpCode,
          purpose: 'phone_verification',
        }),
      });

      const data = await res.json();
      if (res.ok && data.success) {
        toast.success('Verifikasi OTP Berhasil! Akun & No. HP Anda telah terverifikasi.');
        onVerifySuccess();
      } else {
        toast.error(data.error || 'Kode OTP salah atau tidak valid! Silakan cek 6-digit kode di Mailtrap Inbox Anda.');
      }
    } catch {
      toast.error('Gagal terhubung ke server untuk verifikasi OTP. Pastikan koneksi backend aktif.');
    } finally {
      setLoading(false);
    }
  };

  const handleResendOTP = async () => {
    setResending(true);
    try {
      const token = localStorage.getItem('kacaru_token');
      const res = await fetch('http://localhost:8080/otp/request', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });
      const data = await res.json();
      if (res.ok && data.success) {
        toast.success(`Kode OTP 6-Digit baru telah dikirim ke Mailtrap (${email})!`);
      } else {
        toast.error(data.error || 'Gagal mengirim ulang kode OTP.');
      }
    } catch {
      toast.error('Gagal terhubung ke server untuk mengirim ulang OTP.');
    } finally {
      setResending(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md animate-fade-in">
      <div className="w-full max-w-md bg-slate-900 border border-slate-700/80 rounded-3xl p-6 text-white space-y-5 shadow-2xl relative">
        <button
          onClick={onClose}
          className="absolute top-5 right-5 text-slate-400 hover:text-white p-1 rounded-lg"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="text-center space-y-2 pt-2">
          <div className="inline-flex p-3.5 bg-gradient-to-tr from-emerald-500 to-teal-500 rounded-2xl shadow-lg shadow-emerald-500/20 text-white mb-1">
            <ShieldCheck className="w-8 h-8 animate-pulse" />
          </div>
          <h3 className="text-xl font-extrabold bg-gradient-to-r from-white via-emerald-200 to-teal-200 bg-clip-text text-transparent">
            Verifikasi Kode OTP 6-Digit
          </h3>
          <p className="text-xs text-slate-400 px-2 leading-relaxed">
            Kode verifikasi telah dikirimkan ke email <strong className="text-emerald-400">{email}</strong>
            {phone && <span> & No. HP <strong className="text-emerald-400">{phone}</strong></span>}.
          </p>
        </div>

        <form onSubmit={handleVerify} className="space-y-4 pt-2">
          <div>
            <label className="block text-[11px] font-black uppercase tracking-wider text-slate-400 mb-2 text-center">
              Masukkan 6 Digit Kode OTP
            </label>
            <input
              type="text"
              maxLength={6}
              value={otpCode}
              onChange={(e) => setOtpCode(e.target.value.replace(/[^0-9]/g, ''))}
              placeholder="1 2 3 4 5 6"
              className="w-full py-3.5 text-center text-2xl font-mono font-black tracking-[0.5em] bg-slate-950 border border-slate-700 rounded-2xl focus:outline-none focus:ring-2 focus:ring-emerald-500 text-emerald-400 placeholder:text-slate-700 shadow-inner"
            />
            <p className="text-[10px] text-slate-400 text-center mt-1.5 font-medium">
              📧 <span className="text-emerald-400 font-bold">Mailtrap Integration:</span> Buka <strong className="text-white">Mailtrap Sandbox Inbox</strong> Anda untuk melihat & memasukkan 6-digit kode OTP asli.
            </p>
          </div>

          <button
            type="submit"
            disabled={loading || otpCode.length < 6}
            className="w-full py-3.5 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 disabled:opacity-50 text-white font-extrabold text-xs rounded-2xl shadow-lg shadow-emerald-600/25 flex items-center justify-center gap-2 transition active:scale-95"
          >
            {loading ? (
              <span>Memverifikasi Kode...</span>
            ) : (
              <span className="flex items-center gap-1.5">
                Verifikasi OTP & Masuk <ArrowRight className="w-4 h-4" />
              </span>
            )}
          </button>
        </form>

        <div className="border-t border-slate-800 pt-3 flex items-center justify-between text-xs">
          <span className="text-slate-400 text-[11px]">Tidak menerima kode?</span>
          <button
            type="button"
            onClick={handleResendOTP}
            disabled={resending}
            className="text-emerald-400 hover:text-emerald-300 font-bold flex items-center gap-1 text-[11px] transition"
          >
            <RotateCw className={`w-3.5 h-3.5 ${resending ? 'animate-spin' : ''}`} />
            <span>Kirim Ulang OTP</span>
          </button>
        </div>
      </div>
    </div>
  );
};
