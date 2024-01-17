import { LitElement, css, html, svg } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Annotation, VideoStream } from './api.types.ts'
import { repeat } from 'lit/directives/repeat.js'
import { resetcss } from './reset.css.ts'

@customElement('video-player')
export class VideoPlayer extends LitElement {
  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Number })
  activeAnnotation: number | null = null

  @property({ type: Object })
  videostream: VideoStream | null = null

  @property({ type: Number })
  currentTime = 0

  @property({ type: Boolean })
  playing = false

  @property({ type: Number })
  seekTo: number | null = null

  hoverAnnotation(id: number | null) {
    this.dispatchEvent(new CustomEvent('mouseover-annotation', { detail: id }))
  }

  render() {
    const rects = repeat(this.annotations, (annotation: Annotation) => {
      const x1 = annotation.boundingBox?.x1 ?? 0
      const y1 = annotation.boundingBox?.y1 ?? 0

      const x2 = annotation.boundingBox?.x2 ?? 100
      const y2 = annotation.boundingBox?.y2 ?? 100

      return svg`
        <g 
          @mouseover="${() => this.hoverAnnotation(annotation.id)}" 
          @mouseout=${() => this.hoverAnnotation(null)}
        > 
          <rect 
            class="annotation-rect ${annotation.id === this.activeAnnotation ? 'active' : ''}" 
            x="${x1}%" y="${y1}%" width="${x2 - x1}%" height="${y2 - y1}%" 
            stroke-width="4" fill="#00000000" 
          />
          <foreignobject class="node" x="${x1}%" y="${y1}%" width="${x2 - x1}%" height="${y2 - y1}%" >
            <span class="annotation-label">Ann. #${annotation.id}</span>              
          </foreignobject>
        </g>`
    })

    let video = html``
    if (this.videostream?.stream_url != null) {
      video = html`
        <youtube-player 
          id="yt" 
          .url=${this.videostream.stream_url} 
          .playing=${this.playing}
          .seekTo=${this.seekTo}
          @timeupdate=${(e: CustomEvent) => {
            this.currentTime = e.detail
            this.seekTo = null
            this.dispatchEvent(new CustomEvent('timeupdate', { detail: e.detail }))
          }}
          @durationchange=${(e: CustomEvent) => {
            this.dispatchEvent(new CustomEvent('durationchange', { detail: e.detail }))
          }}
          @loadeddata=${() => this.dispatchEvent(new Event('loadeddata'))}
        />`
    }

    return html`
      <div class="video-container">
        ${video}
        <div class="annotation-overlay">
          <svg width="100%" height="100%">
            ${rects}
          </svg>
        </div>
      </div>`
  }

  static styles = css`
    ${resetcss}

    .no-video {
      font-weight: bold;
      color: var(--bg);
    }
    .video-container {
      width: 100%;
      aspect-ratio: 4 / 3;  
      background-color: var(--gray-300);
      position: relative
    }
    .annotation-overlay {
      pointer-events: none;
      position: absolute;
      width: 100%;
      height: 100%;
      z-index: 3;
    }
    youtube-player {
      position: absolute;
      width: 100%;
      height: 100%;
      z-index: 2;
    }
    .annotation-rect {
      transition: fill 0.25s;
      stroke: var(--bright-blue-400);
    }
    .annotation-rect.active {
      fill: #CCEEEE20;
    }
    .annotation-label {
      font-size: 0.75rem;
      padding: 0.4rem 0.5rem;
      background-color: var(--bright-blue-400);
      color: var(--content);
    }`
}

declare global {
  interface HTMLElementTagNameMap {
    'video-player': VideoPlayer
  }
}
