import { expect, test } from 'vitest'
import { datetimeDifference } from './datetime'

test('datetimeDifference calculates time correctly', () => {
  const a = '2023-12-06T12:17:00.000Z'
  const b = '2023-12-06T12:19:00.000Z'
  const diff = datetimeDifference(b, a)
  expect(diff).toBe(120)
})
