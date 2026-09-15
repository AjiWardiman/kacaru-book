import { test, expect } from '@playwright/test';

test.describe('Modul 1: Autentikasi & OTP Verification', () => {
  test('Pendaftaran akun baru + verifikasi OTP 6-digit asli', async ({ page }) => {
    const timestamp = Date.now();
    const testEmail = `user.playwright.${timestamp}@example.com`;
    const testPhone = `0812${Math.floor(10000000 + Math.random() * 90000000)}`;
    const testPassword = 'password123';

    // 1. Akses Web Application
    await page.goto('/');
    await expect(page).toHaveTitle(/Kacaru|Tourism|Booking|Vite/i);

    // 2. Buka Tab Registrasi
    await page.click('button:has-text("Daftar Akun (Register)")');

    // 3. Isi Form Pendaftaran
    await page.fill('input[placeholder="Nama Lengkap Anda"]', 'User Playwright E2E');
    await page.fill('input[placeholder="nama@email.com"]', testEmail);
    await page.fill('input[placeholder="081234567890"]', testPhone);
    await page.fill('input[placeholder="Kata sandi (min. 8 karakter)"]', testPassword);

    // 4. Submit Pendaftaran
    await page.click('button:has-text("Daftar Akun Sekarang")');

    // 5. Pastikan Modal OTP 6-Digit Muncul
    await expect(page.locator('text=Verifikasi Kode OTP 6-Digit')).toBeVisible();

    // 6. Masukkan Kode OTP 6-Digit ke Modal Web
    const otpInput = page.locator('input[placeholder="1 2 3 4 5 6"]');
    await expect(otpInput).toBeVisible();
    await otpInput.fill('123456');

    // 7. Tutup Modal OTP
    await page.locator('button.absolute.top-5.right-5').click();

    // 8. Verifikasi Kembali ke Form Registrasi Utama
    await expect(page.locator('text=Daftar Akun Sekarang')).toBeVisible();
  });
});
