import { describe, it, expect } from 'vitest';
import { getImageUrl } from './imageHelper';

describe('getImageUrl helper', () => {
  it('returns default fallback image when url is undefined or empty', () => {
    expect(getImageUrl(undefined)).toContain('unsplash.com');
    expect(getImageUrl('')).toContain('unsplash.com');
  });

  it('returns full absolute http/https url unchanged', () => {
    const absoluteUrl = 'https://example.com/photo.jpg';
    expect(getImageUrl(absoluteUrl)).toBe(absoluteUrl);
  });

  it('prefixes backend API_BASE_URL for relative /uploads paths', () => {
    const path = '/uploads/destinations/bali.jpg';
    const result = getImageUrl(path);
    expect(result).toContain(path);
  });

  it('prefixes backend API_BASE_URL for uploads/ without leading slash', () => {
    const path = 'uploads/destinations/bromo.jpg';
    const result = getImageUrl(path);
    expect(result).toContain('/uploads/destinations/bromo.jpg');
  });
});
