import { test, expect } from '@playwright/test';

test.describe('Modul 3: AI Recommendation Engine', () => {
  test('Pencarian rekomendasi wisata berbasis AI (Budget & Kategori)', async ({ page }) => {
    // 1. Akses Web Application
    await page.goto('/');

    // 2. Quick Login User
    await page.click('button:has-text("Quick Login User")');
    await expect(page.locator('text=Eksplorasi Destinasi Wisata')).toBeVisible();

    // 3. Klik Tombol ✨ Rekomendasi AI
    const aiBtn = page.getByRole('button', { name: /Rekomendasi AI/i }).first();
    await expect(aiBtn).toBeVisible();
    await aiBtn.click();

    // 4. Verifikasi Modal AI Rekomendasi Terbuka
    await expect(page.locator('text=Rekomendasi AI Destinasi')).toBeVisible();

    // 5. Klik Tombol Cari Destinasi dengan AI
    const getRecBtn = page.locator('button:has-text("Cari Destinasi dengan AI")').first();
    await expect(getRecBtn).toBeVisible();
    await getRecBtn.click();

    // 6. Memastikan Hasil Rekomendasi AI Tampil
    await expect(page.locator('text=Pantai Lariti').first()).toBeVisible({ timeout: 10000 });
  });
});
