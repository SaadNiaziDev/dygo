import assert from 'node:assert/strict'
import test from 'node:test'

import { studioPageNavItems } from './page-nav.ts'

test('studioPageNavItems lists App Pages and skips Home', () => {
  assert.deepEqual(studioPageNavItems([
    { app: 'studio', key: 'home', path: '/', label: 'Home' },
    { app: 'sales', key: 'board', path: '/board', label: 'Board', icon: 'layout-dashboard' },
    { app: 'sales', key: 'dispatch-board', path: '/dispatch-board' },
  ], '/board'), [
    { label: 'Board', to: '/board', icon: 'layout-dashboard', current: true },
    { label: 'Dispatch Board', to: '/dispatch-board', current: false },
  ])
})

test('studioPageNavItems skips invalid Page paths', () => {
  assert.deepEqual(studioPageNavItems([
    { app: 'studio', key: 'outside', path: '//example.com' },
    { app: 'studio', key: 'search', path: '/search?q=1' },
  ], '/board'), [])
})
