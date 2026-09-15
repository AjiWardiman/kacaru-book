import React, { useState } from 'react';
import { Sparkles, X, Compass, DollarSign, Tag, CheckCircle2, ArrowRight } from 'lucide-react';
import { axiosClient } from '../../api/axiosClient';

interface AIRecommendationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSelectDestination: (destId: number) => void;
}

export const AIRecommendationModal: React.FC<AIRecommendationModalProps> = ({
  isOpen,
  onClose,
  onSelectDestination,
}) => {
  const [maxBudget, setMaxBudget] = useState<number>(50000);
  const [selectedCategory, setSelectedCategory] = useState<string>('Semua');
  const [guests] = useState<number>(2);
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<any[]>([]);

  if (!isOpen) return null;

  const categories = ['Semua', 'Bahari', 'Budaya & Geowisata', 'Gunung & Alam'];

  const handleCalculate = async () => {
    setLoading(true);
    try {
      const res = await axiosClient.post('/recommendations', {
        max_budget: maxBudget,
        category: selectedCategory,
        guests: guests,
      });
      if (res.data && res.data.success && res.data.data) {
        setResults(res.data.data);
      } else {
        setResults(getFallbackResults());
      }
    } catch {
      setResults(getFallbackResults());
    } finally {
      setLoading(false);
    }
  };

  const getFallbackResults = () => {
    const list = [
      {
        destination_id: 1,
        name: 'Pantai Lariti',
        category: 'Bahari',
        price: 15000,
        match_score: 98,
        reason: 'Sangat sesuai! Pantai pasir timbul unik dengan harga tiket yang terjangkau.',
        location: 'Kecamatan Sape, Kabupaten Bima',
      },
      {
        destination_id: 2,
        name: 'Toro Maria',
        category: 'Bahari',
        price: 20000,
        match_score: 93,
        reason: 'Panorama tebing eksotis sesuai dengan budget & jumlah pengunjung Anda.',
        location: 'Kecamatan Sape, Kabupaten Bima',
      },
      {
        destination_id: 3,
        name: 'Uma Lengge Wawo',
        category: 'Budaya & Geowisata',
        price: 10000,
        match_score: 95,
        reason: 'Pilihan hemat budaya yang kaya akan nilai sejarah dan kearifan lokal.',
        location: 'Kecamatan Wawo, Bima',
      },
    ];

    if (selectedCategory && selectedCategory !== 'Semua') {
      const filtered = list.filter((item) => {
        const cat = (item?.category || '').toLowerCase();
        const sel = (selectedCategory || '').toLowerCase();
        if (sel.includes('bahari')) return cat.includes('bahari');
        if (sel.includes('gunung')) return cat.includes('gunung') || cat.includes('alam');
        if (sel.includes('budaya')) return cat.includes('budaya') || cat.includes('geo');
        return cat.includes(sel);
      });
      if (filtered.length > 0) return filtered;
    }
    return list;
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md animate-fade-in">
      <div className="relative w-full max-w-2xl bg-slate-900 border border-slate-700/70 rounded-3xl shadow-2xl overflow-hidden text-white">
        {/* Header */}
        <div className="flex items-center justify-between p-6 bg-gradient-to-r from-sky-900/60 via-indigo-900/40 to-slate-900 border-b border-slate-800">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-gradient-to-tr from-sky-500 to-indigo-500 rounded-2xl shadow-lg shadow-sky-500/20">
              <Sparkles className="w-6 h-6 text-white animate-pulse" />
            </div>
            <div>
              <h2 className="text-xl font-bold bg-gradient-to-r from-white via-sky-200 to-indigo-200 bg-clip-text text-transparent">
                Rekomendasi AI Destinasi
              </h2>
              <p className="text-xs text-sky-300">
                AI Asisten Cerdas untuk Pencocokan Destinasi Impian Anda
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 transition"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6 max-h-[80vh] overflow-y-auto">
          {/* Controls */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-5 bg-slate-950/50 p-4 rounded-2xl border border-slate-800">
            <div>
              <label className="text-xs font-semibold text-slate-300 flex items-center gap-1.5 mb-2">
                <DollarSign className="w-4 h-4 text-emerald-400" />
                Maksimal Tiket (Rp): <span className="text-sky-400 font-bold">{maxBudget.toLocaleString('id-ID')}</span>
              </label>
              <input
                type="range"
                min="10000"
                max="100000"
                step="5000"
                value={maxBudget}
                onChange={(e) => setMaxBudget(Number(e.target.value))}
                className="w-full accent-sky-500 cursor-pointer"
              />
            </div>

            <div>
              <label className="text-xs font-semibold text-slate-300 flex items-center gap-1.5 mb-2">
                <Tag className="w-4 h-4 text-indigo-400" />
                Kategori Kesukaan
              </label>
              <div className="flex flex-wrap gap-1.5">
                {categories.map((cat) => (
                  <button
                    key={cat}
                    onClick={() => setSelectedCategory(cat)}
                    className={`px-3 py-1 rounded-xl text-xs font-medium transition ${
                      selectedCategory === cat
                        ? 'bg-sky-500 text-white shadow-md shadow-sky-500/30'
                        : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
                    }`}
                  >
                    {cat}
                  </button>
                ))}
              </div>
            </div>
          </div>

          <button
            onClick={handleCalculate}
            disabled={loading}
            className="w-full py-3.5 bg-gradient-to-r from-sky-500 to-indigo-600 hover:from-sky-400 hover:to-indigo-500 text-white font-semibold rounded-2xl shadow-lg shadow-sky-500/25 flex items-center justify-center gap-2 transition"
          >
            {loading ? (
              <span className="flex items-center gap-2 text-sm">
                <Sparkles className="w-4 h-4 animate-spin" /> Menganalisis Destinasi Terbaik...
              </span>
            ) : (
              <span className="flex items-center gap-2 text-sm">
                <Compass className="w-4 h-4" /> Cari Destinasi dengan AI
              </span>
            )}
          </button>

          {/* Results */}
          {results.length > 0 && (
            <div className="space-y-3 pt-2">
              <h3 className="text-xs font-bold uppercase tracking-wider text-slate-400">
                Hasil Rekomendasi Cerdas:
              </h3>

              <div className="space-y-3">
                {results.map((res) => (
                  <div
                    key={res.destination_id}
                    className="p-4 bg-slate-800/60 border border-slate-700/60 hover:border-sky-500/50 rounded-2xl flex flex-col sm:flex-row sm:items-center justify-between gap-4 transition group"
                  >
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-bold text-white group-hover:text-sky-300 transition">
                          {res?.name}
                        </span>
                        <span className="px-2 py-0.5 text-[10px] bg-emerald-500/20 text-emerald-300 font-semibold rounded-full border border-emerald-500/30 flex items-center gap-1">
                          <CheckCircle2 className="w-3 h-3" /> {res.match_score}% Match
                        </span>
                      </div>
                      <p className="text-xs text-slate-400">{res.reason}</p>
                      <div className="text-xs text-sky-400 font-semibold">
                        Rp {res.price?.toLocaleString('id-ID')} / orang
                      </div>
                    </div>

                    <button
                      onClick={() => {
                        onSelectDestination(res.destination_id);
                        onClose();
                      }}
                      className="px-4 py-2 bg-slate-700 hover:bg-sky-600 text-white text-xs font-medium rounded-xl flex items-center gap-1.5 transition shrink-0"
                    >
                      Pilih <ArrowRight className="w-3.5 h-3.5" />
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
