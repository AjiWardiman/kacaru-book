/**
 * OpenAPI 3.0 / Swagger 2.0 Types generated from backend/docs/swagger.json
 * Single Source of Truth for Tourism Booking System API
 */

export interface APIResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

export interface UserResponse {
  id?: number | string;
  name: string;
  email: string;
  phone?: string;
  role: 'admin' | 'user' | string;
  created_at?: string;
  updated_at?: string;
}

export interface DestinationResponse {
  id: number;
  name: string;
  location: string;
  description?: string;
  category?: string;
  price: number;
  image_url?: string;
  is_active: boolean;
  rating?: number;
  created_at?: string;
  updated_at?: string;
}

export interface SlotResponse {
  id: number;
  destination_id: number;
  date: string;
  quota: number;
  booked_quota: number;
  status: string;
}

export interface VisitorInput {
  name: string;
  phone: string;
  email: string;
  age: number;
}

export interface CreateBookingRequest {
  slot_id: number;
  total_people: number;
  visitors: VisitorInput[];
}

export interface BookingResponse {
  id: number;
  booking_code: string;
  slot_id: number;
  total_people: number;
  total_price: number;
  status: 'PENDING' | 'PAID' | 'EXPIRED' | 'CANCELLED' | string;
  booked_at: string;
  expired_at?: string;
  destination_name?: string;
  destination_location?: string;
  destination_image?: string;
  user_name?: string;
  user_email?: string;
  payment_url?: string;
}

export interface PaymentResponse {
  id: number;
  booking_id: number;
  amount: number;
  status: string;
  payment_url: string;
  created_at: string;
}

export interface NotificationResponse {
  id: number;
  user_id: number;
  type: string;
  title: string;
  message: string;
  is_read: boolean;
  sent_at: string;
}
