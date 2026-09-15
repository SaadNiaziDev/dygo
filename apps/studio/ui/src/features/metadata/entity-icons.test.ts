import test from 'node:test'
import assert from 'node:assert/strict'
import { Box, House, Settings2 } from '@lucide/vue'

import { iconForEntity } from './entity-icons.ts'

test('iconForEntity resolves metadata names and falls back for unknown icons', () => {
  assert.equal(iconForEntity('settings-2'), Settings2)
  assert.equal(iconForEntity('Settings2'), Settings2)
  assert.equal(iconForEntity('house'), House)
  assert.equal(iconForEntity('missing-icon'), Box)
})
