// vite.config.js
import { defineConfig } from 'vite'
import { vite as vidstack } from 'vidstack/plugins'
import litcss from 'vite-plugin-lit-css'
import { globSync } from 'glob'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const input = Object.fromEntries(
  globSync('*.html').map((file) => [
    file.slice(0, file.length - path.extname(file).length),
    fileURLToPath(new URL(file, import.meta.url)),
  ])
)

export default defineConfig({
  plugins: [vidstack(), litcss()],
  build: {
    rollupOptions: {
      input,
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
