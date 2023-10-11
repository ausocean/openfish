// vite.config.js
import { defineConfig } from 'vite'
import path from 'path'
import { glob } from 'glob'

const root = path.join(__dirname, 'src')
const publicDir = path.join(__dirname, 'public')
const outDir = path.join(__dirname, 'dist')

export default defineConfig({
  root,
  publicDir,
  build: {
    outDir,
    rollupOptions: {
      input: glob.sync(path.resolve(root, '*.html')),
    },
  },
})
