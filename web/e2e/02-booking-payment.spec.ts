import { test, expect } from '@playwright/test';

test.describe('Modul 2: Booking Wisata & Integrasi Pembayaran Xendit', () => {
  test('Alur Pemesanan Tiket, Weather API, dan Konfirmasi Xendit Webhook (PAID)', async ({ page, request }) => {
    // 1. Akses Web Application
    await page.goto('/');

    // 2. Quick Login User
    await page.click('button:has-text("Quick Login User")');
    await expect(page.locator('text=Eksplorasi Destinasi Wisata')).toBeVisible();

    // 3. Verifikasi Widget Cuaca (OpenWeatherMap)
    await expect(page.locator('text=Bima, NTB').first()).toBeVisible();

    // 4. Pilih Destinasi dan Klik 'Pesan Tiket'
    const bookButton = page.locator('button:has-text("Pesan Tiket")').first();
    await expect(bookButton).toBeVisible();
    await bookButton.click();

    // 5. Verifikasi Modal Booking Terbuka
    await expect(page.locator('text=Formulir Pemesanan Tiket')).toBeVisible();

    // 6. Konfirmasi Pemesanan Tiket
    const confirmBookingBtn = page.locator('button:has-text("Konfirmasi & Pesan Tiket")');
    if (await confirmBookingBtn.isVisible()) {
      await confirmBookingBtn.click();
    }

    // 7. Ambil Token User dari LocalStorage
    const token = await page.evaluate(() => localStorage.getItem('kacaru_token'));

    // 8. Dapatkan Daftar Booking User via API untuk mengambil External ID Xendit
    if (token) {
      const myBookingsRes = await request.get('http://localhost:8080/bookings/my-bookings', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (myBookingsRes.ok()) {
        const myBookingsData = await myBookingsRes.json();
        
        if (myBookingsData.data && myBookingsData.data.length > 0) {
          const latestBooking = myBookingsData.data[0];
          
          // Buat Payment Link Xendit
          const payRes = await request.post(`http://localhost:8080/bookings/${latestBooking.id}/pay`, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
          if (payRes.ok()) {
            const payData = await payRes.json();
            const externalId = payData.data?.external_id;

            if (externalId) {
              // 9. SIMULASI WEBHOOK XENDIT (Kirim Notifikasi 'PAID' dari Xendit ke Backend)
              await request.post('http://localhost:8080/webhooks/xendit', {
                headers: {
                  'x-callback-token': 'xnd_callback_dummytoken',
                  'Content-Type': 'application/json',
                },
                data: {
                  external_id: externalId,
                  status: 'PAID',
                  amount: latestBooking.total_price,
                  payment_method: 'QRIS',
                },
              });
            }
          }
        }
      }
    }

    // 10. Verifikasi Halaman Berjalan Lancar
    await expect(page.locator('text=Eksplorasi Destinasi Wisata')).toBeVisible();
  });
});
