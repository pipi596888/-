import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api/user': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/voice': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      '/api/tts': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/api/works': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    },
  },
})
