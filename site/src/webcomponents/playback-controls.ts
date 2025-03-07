import { TailwindElement } from './tailwind-element'
import { css, html, svg } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { VideoStream } from '../api/videostream.ts'
import { repeat } from 'lit/directives/repeat.js'
import { formatVideoTime } from '../utils/datetime'
import type { Annotation } from '../api/annotation'

export type SeekEvent = CustomEvent<number>

@customElement('playback-controls')
export class PlaybackControls extends TailwindElement {
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

  render() {
    const svgContents = repeat(this.annotations, (a) => {
      const duration = a.end - a.start

      const x = (a.start / this.duration) * 100
      const width = Math.max((duration / this.duration) * 100, 0.25) // Give them a min width of 0.25% so they are legible.

      return svg`<rect class="fill-green-500 opacity-50" x="${x}%" y="0%" width="${width}%" height="100%" />`
    })

    const heatmap = this.editMode
      ? html``
      : html`
      <svg class="absolute inset z-10 w-full h-6">
        ${svgContents}
      </svg>
    `

    return html`
      <div class="flex w-full px-4 py-2 gap-2 bg-blue-600 text-slate-50 items-center">
        <button class="btn size-sm variant-orange w-28 justify-center" @click="${this.togglePlayback}">${
          this.playing ? 'Pause' : 'Play'
        }</button>
        
        <div class="flex w-32 *:rounded-none *:px-0 *:w-18 rounded-md overflow-clip">
          <button class="btn size-sm variant-blue" @click="${() => this.bwd(5)}">≪</button>
          <button class="btn size-sm variant-blue" @click="${() => this.bwd(1)}">&lt;</button>
          <button class="btn size-sm variant-blue" @click="${() => this.fwd(1)}">&gt;</button>
          <button class="btn size-sm variant-blue" @click="${() => this.fwd(5)}">≫</button>
        </div>

        <div class="w-full h-6 px-1 bg-blue-500 rounded-md"> 
          <div class="relative">
            ${heatmap}
            <input
              class="absolute inset z-20 w-full h-6"
              type="range" 
              min="0" .max="${this.duration}" step="1"
              .value="${this.currentTime}" 
              @input="${this.seek}" 
            />
          </div>
        </div>

        <span class="p-1 w-60 text-right whitespace-nowrap">${formatVideoTime(this.currentTime)} / ${formatVideoTime(this.duration)}</span>
      </div>`
  }

  static styles = [
    TailwindElement.styles!,
    css`
      input[type="range"]::-moz-range-thumb,
      input[type="range"]::-webkit-slider-thumb {
        box-sizing: content-box;
        background-color: var(--color-red-400);
        width: 2px;
        height: calc(1.5rem + 4px);
        border: none;
        cursor: ew-resize;
        transform: translate(-1px, 0);
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'playback-controls': PlaybackControls
  }
}
