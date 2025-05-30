import { defineConfig } from 'vite'
import { vite as vidstack } from 'vidstack/plugins'
import litcss from 'vite-plugin-lit-css'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  esbuild: {
    target: 'es2022',
  },
  plugins: [
    tailwindcss(),
    vidstack(),
    litcss({
      include: /[?&]lit\b/,
    }),
  ],
})
