import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  preview: {
    host: true,
    port: parseInt(process.env.PORT || '8080'),
    strictPort: true,
    allowedHosts: [
      'abundant-caring-production-f136.up.railway.app',
      '.railway.app',
      'localhost',
      '127.0.0.1',
      'all'
    ]
  },
  server: {
    host: '0.0.0.0',
    port: parseInt(process.env.PORT || '5173')
  }
})
