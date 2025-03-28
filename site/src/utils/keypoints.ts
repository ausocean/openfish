import 'reflect-metadata'

import { Transform, Type } from 'class-transformer'
import { formatVideoTime, parseVideoTime } from '../utils/datetime'

export class Keypoint {
  @Transform(({ value }) => parseVideoTime(value), { toClassOnly: true })
  @Transform(({ value }) => formatVideoTime(value, true), { toPlainOnly: true })
  time: number

  @Type(() => BoundingBox)
  box: BoundingBox

  constructor(time: number, box: BoundingBox) {
    this.time = time
    this.box = box
  }
}

export class BoundingBox {
  x1: number
  y1: number
  x2: number
  y2: number

  constructor(x1: number, y1: number, x2: number, y2: number) {
    this.x1 = x1
    this.y1 = y1
    this.x2 = x2
    this.y2 = y2
  }

  get w(): number {
    return Math.abs(this.x1 - this.x2)
  }
  get h(): number {
    return Math.abs(this.y1 - this.y2)
  }

  get xmin(): number {
    return Math.min(this.x1, this.x2)
  }
  get ymin(): number {
    return Math.min(this.y1, this.y2)
  }
}

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
