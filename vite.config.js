import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,  // Modifies the Origin header to appear as if the request came from the backend's domain
        secure: false
      },
      '/api/': {
        target: 'http://localhost:8080',
        changeOrigin: true,  // Modifies the Origin header to appear as if the request came from the backend's domain
        secure: false
      }
    }
  },
})
