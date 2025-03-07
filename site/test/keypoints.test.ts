import { expect, test } from 'vitest'
import { findClosestKeypointPair, interpolateKeypoints, lerp } from '../src/utils/keypoints'
import { BoundingBox, Keypoint } from '../src/api/annotation'
import { instanceToPlain } from 'class-transformer'

test(' instanceToPlain formats time as string', () => {
  const plain = instanceToPlain(new Keypoint(123.456, new BoundingBox(0, 0, 0, 0)))

  expect(plain.time).toBe('00:02:03.456')
})

test.each([
  [3, [2, 5]],
  [3.1, [2, 5]],
  [200, [162, 600]],
])('findClosestKeypointPair', (t, expected) => {
  const keypoints = [1, 2, 5, 7, 162, 600, 620, 621].map(
    (time) => new Keypoint(time, new BoundingBox(0, 0, 0, 0))
  )

  const pair = findClosestKeypointPair(keypoints, t)

  expect(pair[0].time).toBe(expected[0])
  expect(pair[1].time).toBe(expected[1])
})

test.each([
  [0, 0, 100, 80, 40, 80],
  [100, 0, 100, 80, 40, 40],
  [50, 0, 100, 80, 40, 60],
  [120, 110, 130, 1, 2, 1.5],
])('lerp', (t, t0, t1, v0, v1, expected) => {
  const v = lerp(t, t0, t1, v0, v1)

  expect(v).toBe(expected)
})

// TODO: table tests
test('interpolate', () => {
  const keypoints: [Keypoint, Keypoint] = [
    {
      time: 1001,
      box: new BoundingBox(10, 10, 50, 80),
    },
    {
      time: 1003,
      box: new BoundingBox(20, 0, 30, 90),
    },
  ]

  const t = 1002

  const box = interpolateKeypoints(keypoints, t)

  expect(box).toStrictEqual(new BoundingBox(15, 5, 40, 85))
})
