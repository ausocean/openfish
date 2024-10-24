import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'

import L from 'leaflet'
import leafletcss from 'leaflet/dist/leaflet.css'
import iconUrl from 'leaflet/dist/images/marker-icon.png'
import { fmtLatLng, parseLatLng } from '../utils/geo'

const tilesURL = 'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
const attribution = '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'

export type PickLocationEvent = CustomEvent<[number, number]>

type TextInputEvent = InputEvent & { target: HTMLInputElement }

const defaultIcon = new L.Icon({
  iconUrl,
  iconAnchor: [25 / 2, 41],
})

@customElement('location-picker')
export class LocationPicker extends LitElement {
  static formAssociated = true

  private _internals: ElementInternals
  private _element: HTMLElement | null = null
  private _pin: L.Marker | null = null
  private _map: L.Map | null = null

  @property({ type: String })
  value = ''

  setLocation(location: L.LatLng | null) {
    console.log(location)

    // If location is null, remove pin from the map.
    if (location == null) {
      this._pin?.removeFrom(this._map!)
      this._pin = null
      return
    }

    // Add marker.
    if (this._pin) {
      this._pin.setLatLng(location)
    } else {
      this._pin = L.marker(location, { icon: defaultIcon })
      this._pin.addTo(this._map!)
    }

    // Center the map around the point.
    this._map?.setView(location)

    // Set lat/long in the text input.
    this.value = fmtLatLng(location)

    // Update form value.
    this._internals.setFormValue(this.value)

    // Dispatch event.
    this.dispatchEvent(new Event('change', { bubbles: true }))
    // this.dispatchEvent(new CustomEvent('pick-location', { detail: [location.lat, location.lng] }))
  }

  constructor() {
    super()

    // Enable use in forms.
    this._internals = this.attachInternals()

    // Create element to use as root for map.
    this._element = document.createElement('div')
    this._element.id = 'map'

    // Create leaflet map.
    this._map = L.map(this._element, {
      center: L.latLng(-34.9991715, 137.9873683),
      zoom: 10,
    })

    L.tileLayer(tilesURL, {
      attribution: attribution,
      maxZoom: 19,
    }).addTo(this._map)

    // If the element is resized, invalidate the map's size.
    const resizeObserver = new ResizeObserver((_) => {
      this._map?.invalidateSize()
    })

    resizeObserver.observe(this._element)

    // On click, set location.
    this._map.on('click', (e) => this.setLocation(e.latlng))
  }

  render() {
    return html`
    <div><small>Click to pick location or paste lat-long coordinates</small></div>
    <input 
        type="text" 
        placeholder="Lat-long coordinates"
        name="location"
        .value=${this.value} 
        @input=${(e: TextInputEvent) => this.setLocation(parseLatLng(e.target.value))} />
    ${this._element}
    `
  }

  static styles = css`
    ${leafletcss} 
    #map {
      aspect-ratio: 4/3;
    }
    input {
      width: 100%;
      margin-bottom: 0.5rem;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'location-picker': LocationPicker
  }
}
