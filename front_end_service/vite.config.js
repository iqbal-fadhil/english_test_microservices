import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { webcrypto } from 'crypto'

// Polyfill Web Crypto so Vite/plugins that call getRandomValues don't crash on older Node.
if (!globalThis.crypto) {
  globalThis.crypto = webcrypto
}

export default defineConfig({
  plugins: [react()],
})