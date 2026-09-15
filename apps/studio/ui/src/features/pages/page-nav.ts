import type { StudioPageClaim } from '../boot/boot.api.ts'
import { normalizePageClaimPath } from '../../router/routes.ts'
import { humanizeEntity } from '../../stores/metadata.identity.ts'

export type StudioPageNavItem = {
  label: string
  to: string
  icon?: string
  current: boolean
}

export function studioPageNavItems(pages: readonly StudioPageClaim[] | null | undefined, currentPath: string): StudioPageNavItem[] {
  if (!pages) return []
  return pages.flatMap((page) => {
    if (!page) return []
    const path = normalizePageClaimPath(page.path)
    if (!path || path === '/') return []
    const label = page.label?.trim() || humanizeEntity(page.key)
    if (!label) return []
    return [{
      label,
      to: path,
      ...(page.icon ? { icon: page.icon } : {}),
      current: currentPath === path,
    }]
  })
}
