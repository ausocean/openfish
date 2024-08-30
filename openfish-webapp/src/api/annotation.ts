import 'reflect-metadata'

import { Transform, Type } from 'class-transformer'
import { formatVideoTime, parseVideoTime } from '../utils/datetime'

export class Annotation {
  id: number
  videostreamId: number

  @Type(() => Keypoint)
  keypoints: Keypoint[]

  observer: string

  // TODO: use Map instead
  observation: Observation

  get start(): number {
    return this.keypoints.at(0)!.time
  }
  get end(): number {
    return this.keypoints.at(-1)!.time
  }
}

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

export type Observation = {
  [key: string]: string
}

export class BoundingBox {
  x1: number
  y1: number
  x2: number
  y2: number

  constructor(x1: number, y1: number, x2: number, y2: number) {
    this.x1 = Math.round(x1)
    this.y1 = Math.round(y1)
    this.x2 = Math.round(x2)
    this.y2 = Math.round(y2)
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
