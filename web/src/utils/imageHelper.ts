import { API_BASE_URL } from '../api/axiosClient';

/**
 * Normalizes destination & user image URLs to ensure they resolve properly
 * whether they are absolute URLs or backend static upload paths.
 */
export const getImageUrl = (url?: string): string => {
  if (!url) {
    return 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=800&q=80';
  }
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  if (url.startsWith('/uploads')) {
    return `${API_BASE_URL}${url}`;
  }
  if (url.startsWith('uploads/')) {
    return `${API_BASE_URL}/${url}`;
  }
  return url;
};
