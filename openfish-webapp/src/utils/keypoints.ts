import { BoundingBox, type Keypoint } from '../api/annotation'

// Keypoints must be sorted.
export function findClosestKeypointPair(
  keypoints: Keypoint[],
  currentTime: number
): [Keypoint, Keypoint] {
  for (let i = 0; i < keypoints.length - 1; i++) {
    const kp = keypoints[i]
    const kpNext = keypoints[i + 1]

    if (kp.time === currentTime) {
      return [kp, kp]
    }
    if (kpNext.time > currentTime) {
      return [kp, kpNext]
    }
  }

  throw new Error('current time outside range of keypoints')
}

export function lerp(t: number, t0: number, t1: number, v0: number, v1: number) {
  return (v0 * (t1 - t) + v1 * (t - t0)) / (t1 - t0)
}

export function interpolateKeypoints(
  keypoints: [Keypoint, Keypoint],
  currentTime: number
): BoundingBox {
  const t0 = keypoints[0].time
  const t1 = keypoints[1].time

  return new BoundingBox(
    lerp(currentTime, t0, t1, keypoints[0].box.x1, keypoints[1].box.x1),
    lerp(currentTime, t0, t1, keypoints[0].box.y1, keypoints[1].box.y1),
    lerp(currentTime, t0, t1, keypoints[0].box.x2, keypoints[1].box.x2),
    lerp(currentTime, t0, t1, keypoints[0].box.y2, keypoints[1].box.y2)
  )
}
