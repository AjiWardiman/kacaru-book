import React from 'react';
import { Toaster } from 'sonner';

export const ToastProvider: React.FC = () => {
  return (
    <Toaster
      position="bottom-right"
      theme="system"
      richColors
      closeButton
      toastOptions={{
        style: {
          borderRadius: '1rem',
          padding: '0.85rem 1rem',
          fontSize: '0.85rem',
          fontWeight: 500,
          boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
        },
      }}
    />
  );
};
