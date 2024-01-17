import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { Annotation, VideoStream } from './api.types'
import { videotimeToDatetime } from './datetime'
import { resetcss } from './reset.css'
import './annotation-list'
import { MouseoverAnnotationEvent } from './annotation-list'

import './annotation-overlay'
import './youtube-player'
import './annotation-card'
import './playback-controls'
import { MouseoverAnnotationEvent } from './annotation-overlay'
import { DurationChangeEvent, TimeUpdateEvent } from './youtube-player'

@customElement('watch-stream')
export class WatchStream extends LitElement {
  @property({ type: Number })
  set streamID(val: number) {
    this.fetchData(val)
  }

  @state()
  private _videostream: VideoStream | null = null

  @state()
  private _annotations: Annotation[] = []

  @state()
  private _activeId: number | null = null

  @state()
  private _currentTime = 0

  @state()
  private _playing = false

  @state()
  private _duration = 0

  @state()
  private _seekTo: number | null = null

  private play() {
    this._playing = true
  }

  private pause() {
    this._playing = false
  }

  async fetchData(id: number) {
    try {
      // Fetch video stream with ID.
      const res = await fetch(`http://localhost:3000/api/v1/videostreams/${id}`)
      this._videostream = (await res.json()) as VideoStream
    } catch (error) {
      console.error(error) // TODO: handle errors.
    }
    try {
      // Fetch annotations for this video stream.
      // TODO: We should only fetch a small portion of the annotations near the current playback position.
      //       When the user plays the video we can fetch in more as needed.
      const res = await fetch(`http://localhost:3000/api/v1/annotations?videostream=${id}`)
      const json = await res.json()
      this._annotations = json.results as Annotation[]
    } catch (error) {
      console.error(error) // TODO: handle errors.
    }
  }

  render() {
    let filteredAnnotations: Annotation[] = []
    if (this._videostream != null) {
      // Convert playback time in seconds to a datetime.
      const playbackDatetime = videotimeToDatetime(this._videostream?.startTime, this._currentTime)

      // Filter annotations to only show those spanning the current playback time/position.
      filteredAnnotations = this._annotations.filter(
        (an: Annotation) =>
          new Date(an.timespan.start).getTime() <= playbackDatetime.getTime() &&
          playbackDatetime.getTime() <= new Date(an.timespan.end).getTime()
      )
    }

    const video =
      this._videostream == null
        ? html``
        : html`
        <youtube-player 
          .url=${this._videostream?.stream_url}
          .seekTo=${this._seekTo}
          .playing=${this._playing}
          @timeupdate=${(e: TimeUpdateEvent) => (this._currentTime = e.detail)} 
          @durationchange=${(e: DurationChangeEvent) => (this._duration = e.detail)}
          @loadeddata=${() => (this._playing = true)}
          >
        </video-player>`

    return html`
      <div class="root">
        ${video}

        <annotation-overlay
          .annotations=${filteredAnnotations}
          .activeAnnotation=${this._activeId}
          @mouseover-annotation=${(e: MouseoverAnnotationEvent) => (this._activeId = e.detail)}
        ></annotation-overlay>

        <aside>
          <header>
            <h3>Annotations</h3>
            <button class="add-btn" @click=${() => console.error('Not implemented')}>+ Add annotation</button>
          </header>
          <annotation-list
            .annotations=${filteredAnnotations}
            .activeAnnotation=${this._activeId}
            @mouseover-annotation=${(e: MouseoverAnnotationEvent) => (this._activeId = e.detail)}
            >
          </annotation-list>
        </aside>

        <playback-controls 
          .playing=${this._playing} 
          .duration=${this._duration} 
          .currentTime=${this._currentTime}
          .annotations=${this._annotations}
          .videostream=${this._videostream}
          @play=${this.play} 
          @pause=${this.pause}
          @seek=${(e: CustomEvent) => (this._seekTo = e.detail)}
        ></playback-controls>

      </div>`
  }

  static styles = css`
    ${resetcss}

    .root {
      border-radius: 0.5rem;
      overflow: hidden;
      display: grid;
      grid-template-rows: min-content 1fr min-content;
      grid-template-columns: 1fr 28rem;
      grid-template-areas:
        "video-player annotations"
        "controls controls";
    }

    h2 {
      margin: 0;
      padding: .25rem .5rem; 
      border-bottom: 1px solid var(--gray-100);
    }
    span.subtitle {
      color: var(--gray-300);
      font-weight: normal;
    }
    aside {
      grid-area: annotations;
      overflow-y: hidden;
      background-color: var(--blue-700);
      padding: 0 1rem;
    }
   youtube-player {
      grid-area: video-player;
      aspect-ratio: 4 / 3;  
      height: 100%
    }
    annotation-overlay {
      grid-area: video-player;
      z-index: 100;
    }
    playback-controls  {
      grid-area: controls;
    }
    aside header {
      padding: 0.75rem 0;
      border-bottom: 1px solid var(--blue-200);
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
    h3 {
      margin-top: 0;
      margin-bottom: 0;
      margin-left: 0.5rem;
      color: var(--blue-50)
    }

    .add-btn {
      width: min-content;
      height: 2.5rem;
      border-radius: 999px;
      font-size: 1rem;
      padding: 0 1rem;
      white-space: nowrap;
      border: none;
      cursor: pointer;
      
      background-color: var(--orange-400);
      color: var(--orange-800);

      &:hover {
        background-color: var(--orange-500);
      }
    }

  `
}

declare global {
  interface HTMLElementTagNameMap {
    'watch-stream': WatchStream
  }
}
