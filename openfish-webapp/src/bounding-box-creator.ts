import { LitElement, css, html, svg } from 'lit'
import { customElement, state } from 'lit/decorators.js'

export type UpdateBoundingBoxEvent = CustomEvent<[number, number, number, number]>

@customElement('bounding-box-creator')
export class BoundingBoxCreator extends LitElement {
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
    const rect = this._box
      ? svg`
        <g> 
          <line x1="${this._box!.at(0)}%" x2="${this._box!.at(0)}%" y1="${this._box!.at(
            1
          )}%" y2="${this._box!.at(3)}%" stroke-width="4" class="ew" @mousedown=${(e: MouseEvent) =>
            this.startEditEdge(0, e)}/>
          <line x1="${this._box!.at(0)}%" x2="${this._box!.at(2)}%" y1="${this._box!.at(
            1
          )}%" y2="${this._box!.at(1)}%" stroke-width="4" class="ns" @mousedown=${(e: MouseEvent) =>
            this.startEditEdge(1, e)}/>
          <line x1="${this._box!.at(2)}%" x2="${this._box!.at(2)}%" y1="${this._box!.at(
            1
          )}%" y2="${this._box!.at(3)}%" stroke-width="4" class="ew" @mousedown=${(e: MouseEvent) =>
            this.startEditEdge(2, e)}/>
          <line x1="${this._box!.at(0)}%" x2="${this._box!.at(2)}%" y1="${this._box!.at(
            3
          )}%" y2="${this._box!.at(3)}%" stroke-width="4" class="ns" @mousedown=${(e: MouseEvent) =>
            this.startEditEdge(3, e)}/>
          
          <rect 
            x="${this._box!.at(0)}%" y="${this._box!.at(1)}%" width="${
              this._box!.at(2)! - this._box!.at(0)!
            }%" height="${this._box!.at(3)! - this._box!.at(1)!}%" 
            fill="#00000000" 
            @mousedown=${this.startDrag}
          />

          <circle cx="${this._box!.at(0)}%" cy="${this._box!.at(
            1
          )}%" r="5" class="nwse" @mousedown=${(e: MouseEvent) => this.startEditCorner(0, e)}/>
          <circle cx="${this._box!.at(0)}%" cy="${this._box!.at(
            3
          )}%" r="5" class="nesw" @mousedown=${(e: MouseEvent) => this.startEditCorner(1, e)}/>
          <circle cx="${this._box!.at(2)}%" cy="${this._box!.at(
            1
          )}%" r="5" class="nesw" @mousedown=${(e: MouseEvent) => this.startEditCorner(2, e)}/>
          <circle cx="${this._box!.at(2)}%" cy="${this._box!.at(
            3
          )}%" r="5" class="nwse" @mousedown=${(e: MouseEvent) => this.startEditCorner(3, e)}/>
        </g>`
      : ''

    return html`
          <svg width="100%" height="100%" @mousedown=${this.startDrawBox} @mouseup=${this.endBox} @mousemove=${this.updateBox}>
            ${rect}
          </svg>`
  }

  static styles = css`


  svg {
    cursor: crosshair;
  }

  .nesw {
    cursor: nesw-resize;
  }

  .nwse {
    cursor: nwse-resize;
  }
  
  .ew {
    cursor: ew-resize;
  }

  .ns {
    cursor: ns-resize;
  }

  line {
    stroke: var(--bright-blue-400);

    &:hover {
      stroke: var(--bright-blue-600);
    }

  }
  
  rect {
    cursor: move;
  }

  circle {
    fill: var(--blue-800);
  }

  `
}

declare global {
  interface HTMLElementTagNameMap {
    'bounding-box-creator': BoundingBoxCreator
  }
}
