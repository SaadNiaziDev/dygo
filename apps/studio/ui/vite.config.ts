import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig, searchForWorkspaceRoot } from 'vite'

const studioRoot = fileURLToPath(new URL('..', import.meta.url))
const appsRoot = fileURLToPath(new URL('../..', import.meta.url))
const studioPages = fileURLToPath(new URL('../pages', import.meta.url))
const studioUIImporter = fileURLToPath(new URL('./src/app/main.ts', import.meta.url))

function isAppPageModule(importer: string | undefined): boolean {
  if (!importer) {
    return false
  }

  return /\/pages\/[^/]+\/[^/]+\.vue(?:\?|$)/.test(importer.replaceAll('\\', '/'))
}

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    {
      name: 'dygo-watch-app-pages',
      configureServer(server) {
        server.watcher.add(studioPages)
        server.watcher.add(appsRoot)
      },
    },
    {
      name: 'dygo-resolve-app-page-modules',
      enforce: 'pre',
      resolveId(id, importer, options) {
        if (!isAppPageModule(importer)) {
          return null
        }
        if (id.startsWith('\0') || id.startsWith('.') || id.startsWith('/') || id.startsWith('@/') || id.startsWith('@dygo/')) {
          return null
        }
        return this.resolve(id, studioUIImporter, { ...options, skipSelf: true })
      },
    },
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@dygo/ui': fileURLToPath(new URL('./src/design/index.ts', import.meta.url)),
      '@dygo/ui/': fileURLToPath(new URL('./src/design/', import.meta.url)),
      'vue-router': fileURLToPath(new URL('./node_modules/vue-router', import.meta.url)),
      vue: fileURLToPath(new URL('./node_modules/vue', import.meta.url)),
    },
  },
  server: {
    port: 6791,
    strictPort: true,
    fs: {
      allow: [searchForWorkspaceRoot(fileURLToPath(new URL('.', import.meta.url))), studioRoot, appsRoot],
    },
    proxy: {
      '/api': 'http://127.0.0.1:6790',
    },
  },
})
