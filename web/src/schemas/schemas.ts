import { z } from 'zod';

// Skema validasi Form Login
export const loginSchema = z.object({
  identifier: z.string().min(1, 'Email atau Nomor HP wajib diisi'),
  password: z.string().min(6, 'Password minimal 6 karakter'),
});

export type LoginInput = z.infer<typeof loginSchema>;

// Skema validasi objek Notifikasi (fleksibel untuk snake_case dan PascalCase)
export const notificationSchema = z.object({
  id: z.number().or(z.string()).optional(),
  ID: z.number().or(z.string()).optional(),
  user_id: z.number().optional(),
  UserID: z.number().optional(),
  type: z.string().optional(),
  Type: z.string().optional(),
  title: z.string().optional(),
  Title: z.string().optional(),
  message: z.string().optional(),
  Message: z.string().optional(),
  is_read: z.boolean().optional(),
  IsRead: z.boolean().optional(),
  sent_at: z.string().optional(),
  SentAt: z.string().optional(),
});

export type NotificationModel = z.infer<typeof notificationSchema>;

// Skema validasi respons Auth Login
export const authUserSchema = z.object({
  id: z.number().or(z.string()),
  name: z.string(),
  email: z.string(),
  role: z.string(),
  phone: z.string().optional().nullable(),
});

export const authDataSchema = z.object({
  token: z.string(),
  user: authUserSchema,
});

export const authResponseSchema = z.object({
  success: z.boolean(),
  message: z.string().optional(),
  data: authDataSchema.optional(),
  error: z.string().optional(),
});

export type AuthUser = z.infer<typeof authUserSchema>;

// Skema Validasi Pengguna (User Management Admin)
export const userCreateSchema = z.object({
  name: z.string().min(1, 'Nama lengkap wajib diisi'),
  email: z.string().min(1, 'Email wajib diisi').email('Format email tidak valid'),
  phone: z.string().min(1, 'Nomor HP/WA wajib diisi'),
  role: z.enum(['admin', 'customer']),
  password: z.string().min(6, 'Password minimal 6 karakter'),
});

export const userUpdateSchema = z.object({
  name: z.string().min(1, 'Nama lengkap wajib diisi'),
  email: z.string().min(1, 'Email wajib diisi').email('Format email tidak valid'),
  phone: z.string().min(1, 'Nomor HP/WA wajib diisi'),
  role: z.enum(['admin', 'customer']),
  password: z.string().optional().refine(val => !val || val.length >= 6, {
    message: 'Password minimal 6 karakter jika ingin diubah',
  }),
});

export const userItemSchema = z.object({
  id: z.number().or(z.string()),
  name: z.string(),
  email: z.string(),
  phone: z.string().optional().nullable(),
  role: z.string(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export type UserCreateInput = z.infer<typeof userCreateSchema>;
export type UserUpdateInput = z.infer<typeof userUpdateSchema>;
export type UserItem = z.infer<typeof userItemSchema>;

