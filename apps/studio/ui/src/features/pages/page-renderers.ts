import type { Component } from 'vue'

import EntityIndexRenderer from '@/renderers/pages/EntityIndexRenderer.vue'
import type { StudioPageDescriptor } from './pages.api'
import { pageViewMatches } from './page-view-matches'

const renderers = new Map<string, Component>([
  ['entity-index', EntityIndexRenderer],
])

const pageViews = import.meta.glob<Component>([
  '../../../../../**/pages/*/*.vue',
  '../../../../../../../../apps/**/pages/*/*.vue',
], { eager: true, import: 'default' })

export function resolvePageRenderer(page: StudioPageDescriptor): Component | null {
  const builtin = pageRenderer(page.renderer)
  if (builtin) return builtin
  if (page.renderer !== 'vue') return null
  for (const [path, view] of Object.entries(pageViews)) {
    if (pageViewMatches(path, page.app.name, page.key)) return view
  }
  return null
}

export function pageRenderer(name: string): Component | null {
  return renderers.get(name.trim()) ?? null
}

export function registerPageRenderer(name: string, renderer: Component): () => void {
  const key = name.trim()
  if (key === '') {
    throw new Error('Page renderer name is required')
  }
  if (renderers.has(key)) {
    throw new Error(`Page renderer ${key} is already registered`)
  }

  renderers.set(key, renderer)
  return () => {
    if (renderers.get(key) === renderer) {
      renderers.delete(key)
    }
  }
}
