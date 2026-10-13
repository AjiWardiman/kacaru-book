import { describe, it, expect, beforeEach } from 'vitest';
import { useBookingStore } from './useBookingStore';

describe('useBookingStore', () => {
  beforeEach(() => {
    useBookingStore.setState({
      bookings: [
        {
          id: 1,
          booking_code: 'BK-001',
          slot_id: 10,
          total_people: 2,
          total_price: 100000,
          status: 'PENDING',
          booked_at: new Date().toISOString(),
        },
        {
          id: 2,
          booking_code: 'BK-002',
          slot_id: 11,
          total_people: 1,
          total_price: 50000,
          status: 'PENDING',
          booked_at: new Date().toISOString(),
        },
      ],
      isLoading: false,
      error: null,
    });
  });

  it('initializes with test bookings', () => {
    const state = useBookingStore.getState();
    expect(state.bookings.length).toBe(2);
  });

  it('markBookingAsPaid updates booking status by numeric id', () => {
    useBookingStore.getState().markBookingAsPaid(1);
    const updatedBookings = useBookingStore.getState().bookings;

    expect(updatedBookings[0].status).toBe('PAID');
    expect(updatedBookings[1].status).toBe('PENDING');
  });

  it('markBookingAsPaid updates booking status by booking_code string', () => {
    useBookingStore.getState().markBookingAsPaid('BK-002');
    const updatedBookings = useBookingStore.getState().bookings;

    expect(updatedBookings[0].status).toBe('PENDING');
    expect(updatedBookings[1].status).toBe('PAID');
  });
});
