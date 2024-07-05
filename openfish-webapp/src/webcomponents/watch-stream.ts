import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { Annotation, VideoStream, VideoTime } from '../utils/api.types'
import { formatVideoTime, parseVideoTime } from '../utils/datetime'
import resetcss from '../styles/reset.css?raw'
import btncss from '../styles/buttons.css?raw'

import './annotation-displayer'
import './annotation-card'
import './playback-controls'
import './observation-editor'
import './bounding-box-creator'

import vidstackcss from 'vidstack/player/styles/default/theme.css?raw'
import 'vidstack/player'
import 'vidstack/player/ui'

import type { MouseoverAnnotationEvent } from './annotation-displayer'
import type { ObservationEvent } from './observation-editor'
import type { UpdateBoundingBoxEvent } from './bounding-box-creator'

import { ref, type Ref, createRef } from 'lit/directives/ref.js'
import type { MediaPlayerElement } from 'vidstack/elements'
import { extractVideoID } from '../utils/youtube'

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
    this.playerRef.value?.play()
    this._playing = true
  }

  private pause() {
    this.playerRef.value?.pause()
    this._playing = false
  }

  @state()
  private _observation: Record<string, string> = {}

  @state()
  private _start: VideoTime | null = null

  @state()
  private _end: VideoTime | null = null

  @state()
  private _boundingBox: [number, number, number, number] | null = null

  private setStart() {
    this._start = formatVideoTime(this._currentTime)
  }

  private setEnd() {
    this._end = formatVideoTime(this._currentTime)
  }

  private addAnnotation() {
    this._mode = 'editor'
    this.pause()
    this._start = null
    this._end = null
  }

  playerRef: Ref<MediaPlayerElement> = createRef()

  private async confirmAnnotation() {
    const payload: Omit<Annotation, 'id'> = {
      videostreamId: this._videostream!.id,
      observer: 'user@placeholder.com',
      observation: this._observation,
      timespan: { start: this._start!, end: this._end! },
    }

    if (this._boundingBox) {
      payload.boundingBox = {
        x1: Math.round(this._boundingBox[0]),
        y1: Math.round(this._boundingBox[1]),
        x2: Math.round(this._boundingBox[2]),
        y2: Math.round(this._boundingBox[3]),
      }
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
      // Filter annotations to only show those spanning the current playback time/position.
      filteredAnnotations = this._annotations.filter(
        (an: Annotation) =>
          parseVideoTime(an.timespan.start) <= this._currentTime &&
          this._currentTime <= parseVideoTime(an.timespan.end)
      )
    }

    // Render video.
    const videoID = extractVideoID(this._videostream?.stream_url)

    const video = html`
    <media-player
      ${ref(this.playerRef)} 
      title="Openfish Video" 
      src="youtube/${videoID}"
      .currentTime="${this._seekTo}"
      @time-update=${(e: CustomEvent<{ currentTime: number }>) =>
        (this._currentTime = e.detail.currentTime)} 
      @duration-change=${(e: CustomEvent<number>) => (this._duration = e.detail)}
      @can-play=${this.play}
      .muted=${true}
    >
      <media-provider></media-provider>
      <media-video-layout></media-video-layout>
    </media-player>
    `

    const playbackControls = html`
    <playback-controls 
      .playing=${this._playing} 
      .duration=${this._duration} 
      .currentTime=${this._currentTime}
      .annotations=${this._annotations}
      .videostream=${this._videostream}
      @play=${this.play} 
      @pause=${this.pause}
      @seek=${(e: CustomEvent) => {
        console.log(e)
        this._seekTo = e.detail
      }}
    ></playback-controls>`

    const observationEditor = html`

    <section>
      <h4>Annotation times</h4>
      <table>
        <tbody>
        <tr>
          <td>Start:</td>
          <td>${this._start == null ? '' : this._start}</td>
          <td><button class="btn-sm btn-blue w-full" @click=${this.setStart}>Set start time</button></td>
        </tr>
        <tr>
          <td>End:</td>
          <td>${this._end == null ? '' : this._end}</td>
          <td><button class="btn-sm btn-blue w-full" @click=${this.setEnd}>Set end time</button></td>
        </tr>
        </tbody>
      </table>
    </section>

    <section>
      <h4>Bounding Box</h4>
      <span>${this._boundingBox ? 'Yes' : 'No bounding box given'}</span>
    </section>

    <menu>
    <h4>Observation</h4>
    <button class="btn-sm ${
      this._editorMode === 'simple' ? ' btn-secondary' : 'btn-outline'
    }"  @click=${() => (this._editorMode = 'simple')}>Simple</button>
    <button class="btn-sm ${
      this._editorMode === 'advanced' ? ' btn-secondary' : 'btn-outline'
    }" @click=${() => (this._editorMode = 'advanced')}>Advanced</button>
    </menu>
    <div class="scrollable">
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
          <button class="btn btn-orange" @click=${this.confirmAnnotation} .disabled=${
            !this._start || !this._end || Object.keys(this._observation).length === 0
          }>Done</button>
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
        : html`
          <bounding-box-creator @updateboundingbox=${(e: UpdateBoundingBoxEvent) =>
            (this._boundingBox = e.detail)}></bounding-box-creator>    
          `

    return html`
      <div class="root">
        ${video}
        ${overlay}
        ${aside}
        ${playbackControls}
      </div>`
  }

  static styles = css`
    ${unsafeCSS(resetcss)}
    ${unsafeCSS(btncss)}
    ${unsafeCSS(vidstackcss)}

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
    media-player {
      grid-area: video-player;
      aspect-ratio: var(--video-ratio);
      height: 100%;
    }
    annotation-overlay, bounding-box-creator {
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
      padding: 0.5rem;
    }

    h3 {
      margin-top: 0;
      margin-bottom: 0;
      margin-left: 0.5rem;
      color: var(--blue-50)
    }

    h4 {
      margin: 0;
      color: var(--gray-50);
    }

    section {
      padding: 0.5rem;
      color: var(--gray-50);
    }

    table {
      width: 100%;
      
      & tbody tr :nth-child(3) {
        width: 0;
      }
    }

    .w-full { 
      width: 100%;
    }

    menu {
        display: flex;
        justify-content: end;
        margin: 0;
        padding: 0.5rem;
        gap: 0.5rem;

        & > h4 {
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
