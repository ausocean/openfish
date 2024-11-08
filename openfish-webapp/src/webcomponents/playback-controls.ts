import { LitElement, css, html, svg } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { VideoStream } from '../api/videostream.ts'
import { repeat } from 'lit/directives/repeat.js'
import { formatVideoTime } from '../utils/datetime'
import resetcss from '../styles/reset.css?lit'
import btncss from '../styles/buttons.css?lit'
import zoomIn from '../icons/magnifying-glass-plus.svg'
import zoomOut from '../icons/magnifying-glass-minus.svg'
import type { Annotation } from '../api/annotation'

export type SeekEvent = CustomEvent<number>

@customElement('playback-controls')
export class PlaybackControls extends LitElement {
  @property({ type: Number })
  duration = 0

  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Object })
  videostream: VideoStream | null = null

  @property({ type: Number })
  currentTime = 0

  @property({ type: Boolean })
  playing = false

  @property({ type: Boolean })
  editMode = false

  @state()
  zoom = 1

  // Play / pause the video when the user clicks the button.
  private togglePlayback() {
    if (this.playing) {
      this.dispatchEvent(new Event('pause'))
    } else {
      this.dispatchEvent(new Event('play'))
    }
  }

  dispatchSeekEvent(time: number) {
    this.dispatchEvent(new CustomEvent('seek', { detail: time }) as SeekEvent)
  }

  // Emit seek events when user drags the slider.
  private seek(e: InputEvent & { target: HTMLInputElement }) {
    this.currentTime = Number(e.target.value)
    this.dispatchSeekEvent(Number(e.target.value))
  }

  private fwd(seconds: number) {
    this.currentTime = Math.min(this.duration, this.currentTime + seconds)
    this.dispatchSeekEvent(this.currentTime)
  }

  private bwd(seconds: number) {
    this.currentTime = Math.max(0, this.currentTime - seconds)
    this.dispatchSeekEvent(this.currentTime)
  }

  // Jump to when the next annotation occurs in the video.
  private next() {
    const nextAnnotation = this.annotations.find((a) => a.start > this.currentTime)
    if (nextAnnotation !== undefined) {
      this.dispatchSeekEvent(nextAnnotation.start)
    }
  }

  // Jump to when the previous annotation occurs in the video.
  private prev() {
    const idx = this.annotations.findLastIndex((a) => a.start < this.currentTime)
    if (idx > 0) {
      this.dispatchSeekEvent(this.annotations[idx - 1].start)
    }
  }

  render() {
    const svgContents = repeat(this.annotations, (a) => {
      const duration = a.end - a.start

      const x = (a.start / this.duration) * 100
      const width = Math.max((duration / this.duration) * 100, 0.25) // Give them a min width of 0.25% so they are legible.

      return svg`<rect x="${x}%" y="0%" width="${width}%" height="100%" />`
    })

    const heatmap = this.editMode
      ? html``
      : html`
      <svg>
        ${svgContents}
      </svg>
    `

    return html`
      <div class="root">
        <button class="btn-sm btn-orange btn-wide" @click="${this.togglePlayback}">${
          this.playing ? 'Pause' : 'Play'
        }</button>
        
        <button class="btn-sm btn-blue btn-compact" @click="${() => this.bwd(5)}">-5s</button>
        <button class="btn-sm btn-blue btn-compact" @click="${() => this.bwd(1)}">-1s</button>
        <button class="btn-sm btn-blue btn-compact" @click="${() => this.fwd(1)}">+1s</button>
        <button class="btn-sm btn-blue btn-compact" @click="${() => this.fwd(5)}">+5s</button>

        <div class="rangecontrols" style="--zoom-percent: ${this.zoom * 100}%">
          ${heatmap}
          <input
            type="range" 
            min="0" .max="${this.duration}" step="1"
            .value="${this.currentTime}" 
            @input="${this.seek}" 
          />
        </div>

        <button class="btn-sm btn-blue btn-compact" @click="${() => {
          this.zoom *= 2
        }}"><img src=${zoomIn} alt="Zoom in" /></button>
        <button class="btn-sm btn-blue btn-compact" @click="${() => {
          this.zoom = Math.max(1, this.zoom * 0.5)
        }}"><img src=${zoomOut} alt="Zoom out" /></button>

        <button class="btn-sm btn-blue btn-wide btn-compact" @click="${this.prev}">&lt;&nbsp;Prev</button>
        <button class="btn-sm btn-blue btn-wide btn-compact" @click="${this.next}">Next&nbsp;&gt;</button>

        <span class="time">${formatVideoTime(this.currentTime)} / ${formatVideoTime(this.duration)}</span>
      </div>`
  }

  static styles = css`
    ${resetcss}
    ${btncss}
    .root {
      display: flex;
      width: 100%;
      align-items: center;
      padding: 0.5rem 1rem;
      gap: 0.5rem;
      background-color: var(--blue-600);
      color: var(--gray-50);
    }
    input {
      width: 100%;
      height: 100%;
    }
    .btn-wide {
      width: 12ch;
    }
    .btn-compact {
      padding-left: 0.5rem;
      padding-right: 0.5rem;
    }
    span {
      white-space: nowrap;
      padding: 0.25rem;
      width: 15rem;
      text-align: right;
    }
    .rangecontrols {
      display: flex;
      width: 100%;
      align-items: center;
      padding: 0.5rem 1rem;
      gap: 0.5rem;
      height: 2.5rem;
      background-color: var(--gray-50);
      position: relative;
      z-index: 0;
      overflow-x: scroll;
      overflow-y: hidden;
      scroll-behavior: smooth;
    }
    svg {
      top: 0;
      left: 0;
      width: var(--zoom-percent);
      height: 100%;
      border-radius: 0.25rem;
      position: absolute;
      z-index: 10;
      pointer-events: none;
    }
    svg rect {
      fill: var(--green-600);
      opacity: 0.5;
    }
    input[type="range"] {
      top: 0;
      left: 0;
      width: var(--zoom-percent);
      height: 100%;
      position: absolute;
      background-color: transparent;
      z-index: 20;
      -webkit-appearance: none;
    }
    input[type="range"]::-moz-range-thumb {
      background-color: var(--red-400);
      width: 2px;
      border: none;
      height: 3rem;
      cursor: ew-resize;
      transform: translate(-1px, -2px);
    }
    input[type="range"]::-webkit-slider-thumb {
      -webkit-appearance: none;
      box-sizing: content-box;
      background-color: var(--red-400);
      width: 2px;
      border: none;
      height: 2rem;
      cursor: ew-resize;
      transform: translate(-1px, -2px);
    }
    .time {
      width: 28ch;
    }  
    `
}

declare global {
  interface HTMLElementTagNameMap {
    'playback-controls': PlaybackControls
  }
}
