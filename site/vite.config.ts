import config from '@openfish/ui/vite.config'
import { defineConfig } from 'vite'
import { globSync } from 'glob'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const input = Object.fromEntries(
  globSync('{admin/*,*}.html').map((file) => [
    file.slice(0, file.length - path.extname(file).length),
    fileURLToPath(new URL(file, import.meta.url)),
  ])
)

export default defineConfig({
  ...config,
  build: {
    rollupOptions: {
      input,
    },
  },
  resolve: {
    alias: {
      '@openfish/site': fileURLToPath(new URL('./src', import.meta.url)),
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
