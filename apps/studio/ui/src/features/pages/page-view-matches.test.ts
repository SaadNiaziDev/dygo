import assert from 'node:assert/strict'
import test from 'node:test'

import { pageViewMatches } from './page-view-matches.ts'

test('pageViewMatches uses the app folder when the glob keeps it', () => {
  assert.equal(pageViewMatches('../../../../../sales/pages/board/board.vue', 'sales', 'board'), true)
  assert.equal(pageViewMatches('../../../../../sales/pages/board/board.vue', 'studio', 'board'), false)
})

test('pageViewMatches accepts Studio-relative Page globs for studio only', () => {
  assert.equal(pageViewMatches('../../../../pages/home/home.vue', 'studio', 'home'), true)
  assert.equal(pageViewMatches('../../../../pages/home/home.vue', 'sales', 'home'), false)
})
