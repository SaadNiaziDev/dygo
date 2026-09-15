export function pageViewMatches(path: string, app: string, key: string): boolean {
  const normalized = path.replaceAll('\\', '/')
  const file = `${key}/${key}.vue`
  if (normalized.endsWith(`/${app}/pages/${file}`)) return true
  return app === 'studio'
    && normalized.endsWith(`/pages/${file}`)
    && !/\/[a-z0-9-]+\/pages\//.test(normalized)
}
