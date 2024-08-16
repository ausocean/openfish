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

export type VideoTime = `${number}:${number}:${number}.${number}`

export type LatLong = `${number},${number}`

export type VideoStream = {
  id: number
  stream_url: string
  capturesource: string
  startTime: string
  endTime: string
}

export type CaptureSource = {
  id: number
  name: string
  location: LatLong
  camera_hardware: string
  site_id?: number
}
