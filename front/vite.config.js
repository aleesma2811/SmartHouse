import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      // Lets `npm run dev` talk to a backend running locally (e.g. via
      // `docker compose up back`) without hitting CORS issues.
      '/rooms': 'http://localhost:4000',
      '/plugs': 'http://localhost:4000',
    },
  },
})
