import { test, expect } from '@playwright/test';

test.describe('Modul 4: Admin Management, Analytics & Reporting', () => {
  test('Login Admin, Manajemen Pengguna, dan Ekspor Laporan Excel/PDF', async ({ page }) => {
    // 1. Akses Web Application
    await page.goto('/');

    // 2. Quick Login Admin
    await page.click('button:has-text("Quick Login Admin")');

    // 3. Pastikan Masuk ke Dashboard Admin Overview
    await expect(page.locator('text=Total Booking')).toBeVisible();

    // 4. Navigasi ke Tab Laporan Analytics
    await page.click('button:has-text("Laporan Analytics")');
    await expect(page.locator('text=📊 Laporan & Analytics Sistem')).toBeVisible();

    // 5. Verifikasi Tombol Export Excel & PDF Ada
    await expect(page.locator('button:has-text("Export Excel")')).toBeVisible();
    await expect(page.locator('button:has-text("Export PDF")')).toBeVisible();

    // 6. Navigasi ke Tab Pengguna System
    await page.click('button:has-text("Pengguna System")');
    await expect(page.locator('text=Manajemen Pengguna Sistem')).toBeVisible();
  });
});
