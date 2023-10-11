import { LitElement, css, html, svg } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Annotation, VideoStream } from '../utils/api.types'
import { repeat } from 'lit/directives/repeat.js'
import { datetimeDifference, datetimeToVideoTime } from '../utils/datetime'
import { resetcss } from '../css/reset.css'

@customElement('playback-controls')
export class PlaybackControls extends LitElement {
  @property({ type: Number })
  duration: number = 0

  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Object })
  videostream: VideoStream | null = null

  @property({ type: Number })
  currentTime: number = 0

  @property({ type: Boolean })
  playing: boolean = false

  // Play / pause the video when the user clicks the button.
  private togglePlayback() {
    if (this.playing) {
      this.dispatchEvent(new Event('pause'))
    } else {
      this.dispatchEvent(new Event('play'))
    }
  }

  // Emit seek events when user drags the slider.
  private seek(e: InputEvent & { target: HTMLInputElement }) {
    this.currentTime = Number(e.target.value)
    this.dispatchEvent(new CustomEvent('seek', { detail: Number(e.target.value) }))
  }

  // Jump to when the next annotation occurs in the video.
  private next() {
    const nextAnnotation = this.annotations.find(
      (a) => datetimeToVideoTime(this.videostream!.startTime, a.timespan.start) > this.currentTime
    )
    if (nextAnnotation != undefined) {
      const seekTo = datetimeToVideoTime(this.videostream!.startTime, nextAnnotation.timespan.start)
      this.dispatchEvent(new CustomEvent('seek', { detail: seekTo }))
    }
  }

  // Jump to when the previous annotation occurs in the video.
  private prev() {
    const idx = this.annotations.findLastIndex(
      (a) => datetimeToVideoTime(this.videostream!.startTime, a.timespan.start) < this.currentTime
    )
    if (idx > 0) {
      const prevAnnotation = this.annotations[idx - 1]
      const seekTo = datetimeToVideoTime(this.videostream!.startTime, prevAnnotation.timespan.start)
      this.dispatchEvent(new CustomEvent('seek', { detail: seekTo }))
    }
  }

  render() {
    const heatmap = repeat(this.annotations, (annotation) => {
      const start = datetimeToVideoTime(this.videostream!.startTime, annotation.timespan.start)
      const duration = datetimeDifference(annotation.timespan.end, annotation.timespan.start)

      const x = (start / this.duration) * 100
      const width = Math.max((duration / this.duration) * 100, 0.25) // Give them a min width of 0.25% so they are legible.

      return svg`<rect x="${x}%" y="0%" width="${width}%" height="100%" />`
    })

    return html`
      <div>
        <button @click="${this.togglePlayback}">${this.playing ? 'Pause' : 'Play'}</button>
        <button @click="${this.prev}">&lt;&nbsp;Prev</button>
        <button @click="${this.next}">Next&nbsp;&gt;</button>
        <div id="rangecontrols">
          <svg>
          ${heatmap}
          </svg>
          <input 
            type="range" 
            min="0" .max="${this.duration}" step="1"
            .value="${this.currentTime}" 
            @input="${this.seek}" 
          />
        </div>
        <span>${Math.floor(this.currentTime)} / ${Math.floor(this.duration)}</span>
      </div>`
  }

  static styles = css`
    ${resetcss}
    div {
      display: flex;
      width: 100%;
      align-items: center;
      padding: 0.5rem 1rem;
      gap: 0.5rem;
      background-color: var(--gray1);
      border-radius: 0.25rem;
      border: 1px solid var(--gray2);
    }
    input {
      width: 100%;
      height: 100%;
    }
    button {
      width: 4.5rem;
      text-align: center;
    }
    span {
      white-space: nowrap;
      padding: 0.25rem;
      width: 8rem;
      text-align: center;
    }
    #rangecontrols {
      height: 2rem;
      background-color: var(--gray0);
      position: relative;
      z-index: 0;
    }
    svg {
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      border-radius: 0.25rem;
      position: absolute;
      z-index: 10;
      pointer-events: none;
    }
    svg rect {
      fill: var(--success);
      opacity: 0.5;
    }
    input[type="range"] {
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      position: absolute;
      background-color: transparent;
      z-index: 20;
      -webkit-appearance: none;
    }
    input[type="range"]::-moz-range-thumb {
      background-color: var(--error);
      width: 2px;
      border: none;
      height: 2rem;
      cursor: ew-resize;
      transform: translate(-1px, -2px);
    }
    input[type="range"]::-webkit-slider-thumb {
      -webkit-appearance: none;
      box-sizing: content-box;
      background-color: var(--error);
      width: 2px;
      border: none;
      height: 2rem;
      cursor: ew-resize;
      transform: translate(-1px, -2px);
    }`
}

declare global {
  interface HTMLElementTagNameMap {
    'playback-controls': PlaybackControls
  }
}
