// vite.config.js
import { resolve } from 'node:path'
import { defineConfig } from 'vite'

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'),
        watch: resolve(__dirname, 'watch.html'),
        streams: resolve(__dirname, 'streams.html'),
      },
    },
  },
})
