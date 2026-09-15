import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig, searchForWorkspaceRoot } from 'vite'

const studioRoot = fileURLToPath(new URL('..', import.meta.url))
const studioPages = fileURLToPath(new URL('../pages', import.meta.url))

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    {
      name: 'dygo-watch-app-pages',
      configureServer(server) {
        server.watcher.add(studioPages)
      },
    },
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@dygo/ui': fileURLToPath(new URL('./src/design/index.ts', import.meta.url)),
      '@dygo/ui/': fileURLToPath(new URL('./src/design/', import.meta.url)),
    },
  },
  server: {
    port: 6791,
    strictPort: true,
    fs: {
      allow: [searchForWorkspaceRoot(fileURLToPath(new URL('.', import.meta.url))), studioRoot],
    },
    proxy: {
      '/api': 'http://127.0.0.1:6790',
    },
  },
})
