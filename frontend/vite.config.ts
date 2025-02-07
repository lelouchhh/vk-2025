import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3001, // Убедитесь, что порт совпадает с настройками
    strictPort: true,
    hmr: true,
  },
  base: '/', // Укажите базовый путь, если ваше приложение развернуто в поддиректории
});