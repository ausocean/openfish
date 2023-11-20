export type Result<T> = {
  results: T[]
  offset: number
  limit: number
  total: number
}

export interface Annotation {
  id: number
  videostreamId: number
  timespan: Timespan
  boundingBox?: BoundingBox
  observer: string
  observation: Observation
}

export interface Observation {
  [key: string]: string
}

export interface Timespan {
  start: string
  end: string
}

export interface BoundingBox {
  x1: number
  y1: number
  x2: number
  y2: number
}

export interface VideoStream {
  id: number
  stream_url: string
  capturesource: string
  startTime: string
  endTime: string
}
