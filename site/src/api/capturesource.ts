export type CaptureSource = {
  id: number
  name: string
  location: LatLong
  camera_hardware: string
  site_id?: number
}

export type LatLong = `${number},${number}`
