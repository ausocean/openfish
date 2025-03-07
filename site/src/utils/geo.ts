import L from 'leaflet'

export function parseLatLng(str: string): L.LatLng | null {
  try {
    const [lat, lng] = str
      .trim()
      .split(',')
      .map((s) => Number(s.trim()))
    return new L.LatLng(lat, lng)
  } catch (error) {
    console.error(error)
    return null
  }
}

export function fmtLatLng(latLng: L.LatLng | null): string {
  return latLng ? `${latLng.lat},${latLng.lng}` : ''
}
