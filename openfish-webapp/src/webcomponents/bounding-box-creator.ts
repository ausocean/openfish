import { TailwindElement } from './tailwind-element'
import { html, svg, type TemplateResult } from 'lit'
import { customElement, state } from 'lit/decorators.js'

export type UpdateBoundingBoxEvent = CustomEvent<[number, number, number, number]>

@customElement('bounding-box-creator')
export class BoundingBoxCreator extends TailwindElement {
  @state()
  private _box: [number, number, number, number] | null = null

  @state()
  private _isDrawing = false

  @state()
  private _editCorner: number | null = null

  @state()
  private _editEdge: number | null = null

  @state()
  private _drag = false

  startDrawBox(ev: MouseEvent) {
    const x = (ev.offsetX / this.offsetWidth) * 100
    const y = (ev.offsetY / this.offsetHeight) * 100

    this._box = [x, y, x, y]
    this._isDrawing = true
  }

  startEditCorner(corner: number, ev: MouseEvent) {
    ev.stopPropagation()
    this._editCorner = corner
  }

  startEditEdge(edge: number, ev: MouseEvent) {
    ev.stopPropagation()
    this._editEdge = edge
  }

  startDrag(ev: MouseEvent) {
    ev.stopPropagation()
    this._drag = true
  }

  endBox(ev: MouseEvent) {
    this.updateBox(ev)
    this._isDrawing = false
    this._editCorner = null
    this._editEdge = null
    this._drag = false

    this.dispatchEvent(
      new CustomEvent('updateboundingbox', { detail: this._box }) as UpdateBoundingBoxEvent
    )
  }

  updateBox(ev: MouseEvent) {
    ev.stopPropagation()
    const x = (ev.offsetX / this.offsetWidth) * 100
    const y = (ev.offsetY / this.offsetHeight) * 100

    if (!this._box) {
      return
    }

    if (this._isDrawing) {
      this._box[2] = x
      this._box[3] = y
    }

    if (this._drag) {
      const dx = (ev.movementX / this.offsetWidth) * 100
      const dy = (ev.movementY / this.offsetHeight) * 100

      this._box[0] += dx
      this._box[1] += dy
      this._box[2] += dx
      this._box[3] += dy
    }

    switch (this._editCorner) {
      case 0:
        this._box[0] = x
        this._box[1] = y
        break

      case 1:
        this._box[0] = x
        this._box[3] = y
        break

      case 2:
        this._box[2] = x
        this._box[1] = y
        break

      case 3:
        this._box[2] = x
        this._box[3] = y
        break

      default:
        break
    }

    switch (this._editEdge) {
      case 0:
        this._box[0] = x
        break

      case 1:
        this._box[1] = y
        break

      case 2:
        this._box[2] = x
        break

      case 3:
        this._box[3] = y
        break

      default:
        break
    }

    // Mutating an array doesn't trigger an update.
    // https://lit.dev/docs/components/properties/#mutating-properties
    this.requestUpdate()
  }

  render() {
    let rect: TemplateResult = html``

    if (this._box !== null) {
      rect = svg`
        <g> 
          <line 
            x1="${this._box[0]}%" x2="${this._box[0]}%" 
            y1="${this._box[1]}%" y2="${this._box[3]}%" 
            stroke-width="2" class="stroke-sky-400 hover:stroke-sky-600 cursor-ew-resize" 
            @mousedown=${(e: MouseEvent) => this.startEditEdge(0, e)}
          />
          <line 
            x1="${this._box[0]}%" x2="${this._box[2]}%" 
            y1="${this._box[1]}%" y2="${this._box[1]}%" 
            stroke-width="2" class="stroke-sky-400 hover:stroke-sky-600 cursor-ns-resize"
            @mousedown=${(e: MouseEvent) => this.startEditEdge(1, e)}
          />
          <line
            x1="${this._box[2]}%" x2="${this._box[2]}%" 
            y1="${this._box[1]}%" y2="${this._box[3]}%"
            stroke-width="2" class="stroke-sky-400 hover:stroke-sky-600 cursor-ew-resize"
            @mousedown=${(e: MouseEvent) => this.startEditEdge(2, e)}
          />
          <line x1="${this._box[0]}%" x2="${this._box[2]}%" 
            y1="${this._box[3]}%" y2="${this._box[3]}%"
            stroke-width="2" class="stroke-sky-400 hover:stroke-sky-600 cursor-ns-resize"
            @mousedown=${(e: MouseEvent) => this.startEditEdge(3, e)}
          />
          <rect
            class="cursor-move"
            x="${this._box[0]}%" y="${this._box[1]}%" 
            width="${this._box[2] - this._box[0]}%"
            height="${this._box[3] - this._box[1]}%" 
            fill="#00000000" 
            @mousedown=${this.startDrag}
          />

          <circle 
            cx="${this._box[0]}%" cy="${this._box[1]}%" r="5" 
            class="fill-blue-800 cursor-nwse-resize" 
            @mousedown=${(e: MouseEvent) => this.startEditCorner(0, e)}
          />
          <circle 
            cx="${this._box[0]}%" cy="${this._box[3]}%" r="5" 
            class="fill-blue-800 cursor-nesw-resize" 
            @mousedown=${(e: MouseEvent) => this.startEditCorner(1, e)}
          />
          <circle 
            cx="${this._box[2]}%" cy="${this._box[1]}%" r="5"
            class="fill-blue-800 cursor-nesw-resize" 
            @mousedown=${(e: MouseEvent) => this.startEditCorner(2, e)}
          />
          <circle 
            cx="${this._box[2]}%" cy="${this._box[3]}%" r="5"
            class="fill-blue-800 cursor-nwse-resize"
            @mousedown=${(e: MouseEvent) => this.startEditCorner(3, e)}
          />
        </g>`
    }

    return html`
      <svg 
        class="cursor-crosshair"
        width="100%" height="100%" 
        @mousedown=${this.startDrawBox} @mouseup=${this.endBox} @mousemove=${this.updateBox}
      >
        ${rect}
      </svg>`
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'bounding-box-creator': BoundingBoxCreator
  }
}
