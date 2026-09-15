import React from 'react';
import { X } from 'lucide-react';

interface DestinationModalProps {
  isOpen: boolean;
  onClose: () => void;
  editingDestination?: any | null;
  destName: string;
  setDestName: (val: string) => void;
  destLocation: string;
  setDestLocation: (val: string) => void;
  destPrice: number;
  setDestPrice: (val: number) => void;
  destCategory: string;
  setDestCategory: (val: string) => void;
  destDescription: string;
  setDestDescription: (val: string) => void;
  destImageUrl?: string;
  setDestImageUrl?: (val: string) => void;
  onSave: (e: React.FormEvent) => void;
}

export const DestinationModal: React.FC<DestinationModalProps> = ({
  isOpen,
  onClose,
  editingDestination,
  destName,
  setDestName,
  destLocation,
  setDestLocation,
  destPrice,
  setDestPrice,
  destCategory,
  setDestCategory,
  destDescription,
  setDestDescription,
  destImageUrl = '',
  setDestImageUrl,
  onSave,
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
      <div className="w-full max-w-md bg-white dark:bg-slate-900 rounded-3xl p-6 border border-slate-200 dark:border-slate-800 shadow-2xl space-y-4">
        <div className="flex items-center justify-between border-b border-slate-100 dark:border-slate-800 pb-3">
          <h3 className="text-base font-extrabold text-slate-900 dark:text-white">
            {editingDestination ? '✏️ Edit Destinasi Wisata' : '➕ Tambah Destinasi Wisata'}
          </h3>
          <button onClick={onClose} className="p-1 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-white">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={onSave} className="space-y-3.5 text-xs">
          <div>
            <label className="block font-bold text-slate-700 dark:text-slate-300 mb-1">Nama Destinasi</label>
            <input
              type="text"
              required
              value={destName}
              onChange={(e) => setDestName(e.target.value)}
              placeholder="Contoh: Pantai Lariti Sape"
              className="w-full px-3 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-semibold"
            />
          </div>

          <div>
            <label className="block font-bold text-slate-700 dark:text-slate-300 mb-1">Lokasi</label>
            <input
              type="text"
              required
              value={destLocation}
              onChange={(e) => setDestLocation(e.target.value)}
              placeholder="Kecamatan Sape, Kabupaten Bima, NTB"
              className="w-full px-3 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-semibold"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block font-bold text-slate-700 dark:text-slate-300 mb-1">Kategori</label>
              <select
                value={destCategory}
                onChange={(e) => setDestCategory(e.target.value)}
                className="w-full px-3 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-bold"
              >
                <option value="Bahari">Bahari</option>
                <option value="Gunung & Alam">Gunung & Alam</option>
                <option value="Budaya & Geowisata">Budaya & Geowisata</option>
              </select>
            </div>

            <div>
              <label className="block font-bold text-slate-700 dark:text-slate-300 mb-1">Harga Tiket (IDR)</label>
              <input
                type="number"
                value={destPrice}
                onChange={(e) => setDestPrice(Number(e.target.value))}
                className="w-full px-3 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-bold"
              />
            </div>
          </div>

          <div>
            <label className="block font-bold text-slate-700 dark:text-slate-300 mb-1">Foto / Gambar Destinasi</label>
            <div className="space-y-2">
              <input
                type="text"
                value={destImageUrl}
                onChange={(e) => setDestImageUrl && setDestImageUrl(e.target.value)}
                placeholder="URL foto (https://... atau /uploads/destinations/nama.jpg)"
                className="w-full px-3 py-2 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-medium"
              />
              <div className="flex items-center gap-2">
                <input
                  type="file"
                  accept="image/*"
                  onChange={(e) => {
                    const file = e.target.files?.[0];
                    if (file && setDestImageUrl) {
                      const objectUrl = URL.createObjectURL(file);
                      setDestImageUrl(objectUrl);
                    }
                  }}
                  className="block w-full text-xs text-slate-500 file:mr-3 file:py-2 file:px-3 file:rounded-xl file:border-0 file:text-xs file:font-extrabold file:bg-emerald-500/10 file:text-emerald-600 dark:file:text-emerald-400 hover:file:bg-emerald-500/20 cursor-pointer"
                />
              </div>
              {destImageUrl && (
                <div className="relative h-28 rounded-2xl overflow-hidden border border-slate-200 dark:border-slate-800 mt-2">
                  <img src={destImageUrl} alt="Preview" className="w-full h-full object-cover" />
                  <div className="absolute top-2 left-2 px-2 py-0.5 rounded-md bg-slate-900/80 text-white text-[10px] font-bold">
                    Preview Gambar
                  </div>
                </div>
              )}
            </div>
          </div>

          <div>
            <label className="block font-bold text-slate-700 dark:text-slate-300 mb-1">Deskripsi Singkat</label>
            <textarea
              rows={3}
              value={destDescription}
              onChange={(e) => setDestDescription(e.target.value)}
              placeholder="Jelaskan keindahan daya tarik wisata ini..."
              className="w-full px-3 py-2.5 rounded-xl border border-slate-300 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white font-medium"
            />
          </div>

          <div className="pt-2 flex justify-end gap-2">
            <button
              type="button"
              onClick={onClose}
              className="py-2.5 px-4 rounded-xl border border-slate-300 dark:border-slate-700 font-bold text-slate-600 dark:text-slate-300"
            >
              Batal
            </button>
            <button
              type="submit"
              className="py-2.5 px-4 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-extrabold shadow-md shadow-emerald-600/20"
            >
              Simpan Destinasi
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
