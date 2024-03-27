export type Result<T> = {
  results: T[]
  offset: number
  limit: number
  total: number
}

export type Species = {
  id: number
  species: string
  common_name: string
  images?: Image[]
}

export type Image = { src: string; attribution: string }

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
  start: Date | string
  end: Date | string
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

export interface CaptureSource {
  id: number
  name: string
  location: `${number},${number}`
  camera_hardware: string
  site_id?: number
}
