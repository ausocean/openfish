import { expect, test } from 'vitest'
import { zip } from './array-utils'

test('zip', () => {
  const a = [1, 2, 3, 4]
  const b = ['a', 'b', 'c', 'd']
  const zipped = zip(a, b)
  expect(zipped).toStrictEqual([
    [1, 'a'],
    [2, 'b'],
    [3, 'c'],
    [4, 'd'],
  ])
})
