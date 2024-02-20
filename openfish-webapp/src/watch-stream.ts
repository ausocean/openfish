import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { Annotation, VideoStream } from './api.types'
import { videotimeToDatetime } from './datetime'
import { buttonStyles, resetcss } from './reset.css'

import './annotation-displayer'
import './youtube-player'
import './annotation-card'
import './playback-controls'
import './observation-editor'
import { MouseoverAnnotationEvent } from './annotation-displayer'
import { DurationChangeEvent, TimeUpdateEvent } from './youtube-player'
import { ObservationEvent } from './observation-editor'

@customElement('watch-stream')
export class WatchStream extends LitElement {
  @property({ type: Number })
  set streamID(val: number) {
    this.fetchVideoStream(val)
    this.fetchAnnotations(val)
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

  @state()
  private _mode: 'playback' | 'editor' = 'playback'

  @state()
  private _editorMode: 'simple' | 'advanced' = 'simple'

  private play() {
    this._playing = true
  }

  private pause() {
    this._playing = false
  }

  @state()
  private _observation: Record<string, string> = {}

  private addAnnotation() {
    this._mode = 'editor'
    this.pause()
  }

  private async confirmAnnotation() {
    // TODO: let user adjust timespan.
    const start = videotimeToDatetime(this._videostream!.startTime, this._currentTime)
    const end = videotimeToDatetime(this._videostream!.startTime, this._currentTime + 10)

    // TODO: set observer to user's name.
    // TODO: let user include bounding box (optional).
    const payload = {
      videostreamId: this._videostream!.id,
      observer: 'user@placeholder.com',
      observation: this._observation,
      timespan: { start, end },
    }

    // Make annotation.
    await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/annotations`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    // Refetch annotations.
    await this.fetchAnnotations(this._videostream!.id)

    // Start playing.
    this._mode = 'playback'
    this.play()
  }

  private cancelAnnotation() {
    this._mode = 'playback'
    this._observation = {}
    this.play()
  }

  async fetchVideoStream(id: number) {
    try {
      // Fetch video stream with ID.
      const res = await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/videostreams/${id}`)
      this._videostream = (await res.json()) as VideoStream
    } catch (error) {
      console.error(error) // TODO: handle errors.
    }
  }
  async fetchAnnotations(id: number) {
    try {
      // Fetch annotations for this video stream.
      // TODO: We should only fetch a small portion of the annotations near the current playback position.
      //       When the user plays the video we can fetch in more as needed.
      const res = await fetch(
        `${import.meta.env.VITE_API_HOST}/api/v1/annotations?videostream=${id}`
      )
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
          @loadeddata=${this.play}
          >
        </video-player>`

    const playbackControls = html`
    <playback-controls 
      .playing=${this._playing} 
      .duration=${this._duration} 
      .currentTime=${this._currentTime}
      .annotations=${this._annotations}
      .videostream=${this._videostream}
      @play=${this.play} 
      @pause=${this.pause}
      @seek=${(e: CustomEvent) => (this._seekTo = e.detail)}
    ></playback-controls>`

    const observationEditor = html`
    <div class="scrollable">
      <div class="observation-editor ">
        <menu>
            <span>Observation mode</span>
            <button class="btn-sm ${
              this._editorMode === 'simple' ? ' btn-secondary' : 'btn-outline'
            }"  @click=${() => (this._editorMode = 'simple')}>Simple</button>
            <button class="btn-sm ${
              this._editorMode === 'advanced' ? ' btn-secondary' : 'btn-outline'
            }" @click=${() => (this._editorMode = 'advanced')}>Advanced</button>
        </menu>
        <div>
        ${
          this._editorMode === 'simple'
            ? html`<species-selection .observation=${this._observation}           @observation=${(
                ev: ObservationEvent
              ) => (this._observation = ev.detail)}></species-selection>`
            : html`<advanced-editor .observation=${this._observation}           @observation=${(
                ev: ObservationEvent
              ) => (this._observation = ev.detail)}></advanced-editor>`
        }
        </div>
      </div>
    </div>
    `

    const aside =
      this._mode === 'playback'
        ? html`
      <aside>
        <header>
          <h3>Annotations</h3>
          <button class="btn btn-orange" @click=${this.addAnnotation}>+ Add annotation</button>
        </header>
        <div class="scrollable">
          <annotation-list
            .annotations=${filteredAnnotations}
            .activeAnnotation=${this._activeId}
            @mouseover-annotation=${(e: MouseoverAnnotationEvent) => (this._activeId = e.detail)}
            >
          </annotation-list>
        </div>
      </aside>`
        : html`
      <aside>
        <header>
          <h3>Add Annotation</h3>
          <button class="btn btn-secondary" @click=${this.cancelAnnotation}>Cancel</button>
          <button class="btn btn-orange" @click=${this.confirmAnnotation}>Done</button>
        </header>
        
          ${observationEditor}
      </aside>`

    const overlay =
      this._mode === 'playback'
        ? html`
    <annotation-overlay
      .annotations=${filteredAnnotations}
      .activeAnnotation=${this._activeId}
      @mouseover-annotation=${(e: MouseoverAnnotationEvent) => (this._activeId = e.detail)}
    ></annotation-overlay>
    `
        : html``

    return html`
      <div class="root">
        ${video}
        ${overlay}
        ${aside}
        ${playbackControls}
      </div>`
  }

  static styles = css`
    ${resetcss} ${buttonStyles}

    .root {
      --video-ratio: 4 / 3;
      --aside-width: 45ch;

      border-radius: 0.5rem;
      overflow: clip;
      display: grid;
      grid-template-rows: min-content 1fr;
      grid-template-columns: 1fr var(--aside-width);
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
      background-color: var(--blue-700);
      padding: 0 1rem;
      display: flex;
      flex-direction: column;
    }
    youtube-player {
      grid-area: video-player;
      aspect-ratio: var(--video-ratio);
      height: 100%;
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
      justify-content: end;
      gap: 0.5rem;

      & :first-child {
        margin-right: auto
      }
    }

    aside .scrollable {
      position: relative;
      height: 100%;
      overflow-y: scroll;
    }
    aside .scrollable > * {
      position: absolute;
      left: 0;
      top: 0;
    }

    h3 {
      margin-top: 0;
      margin-bottom: 0;
      margin-left: 0.5rem;
      color: var(--blue-50)
    }

    ${buttonStyles}
    .observation-editor {
        display: flex;
        flex-direction: column;
        padding: 0.5rem 0;
        gap: 0.5rem;
        width: 100%;
    }

    .observation-editor menu {
        display: flex;
        justify-content: end;
        margin: 0;
        padding: 0 0.5rem;
        gap: 0.5rem;

        & > span {
          color: var(--gray-50);
          margin-right: auto;
        }  
  
        & button[data-active="true"] {
            background-color: var(--gray-50);
            color: var(--gray-900);
        }
    
        & button[data-active="false"] {
            background-color: transparent;
            color: var(--gray-50);
        }
    }



  `
}

declare global {
  interface HTMLElementTagNameMap {
    'watch-stream': WatchStream
  }
}
