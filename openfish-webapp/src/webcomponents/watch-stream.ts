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
import { repeat } from 'lit/directives/repeat.js'

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
  private _keypoints: VideoTime[] = []

  @state()
  private _boundingBox: [number, number, number, number] | null = null

  private addKeyPoint() {
    if (this._keypoints.length < 2) {
      this._keypoints.push(formatVideoTime(this._currentTime))
    }
    this.requestUpdate()
  }

  private addAnnotation() {
    this._mode = 'editor'
    this.pause()
    this._keypoints = []
  }

  playerRef: Ref<MediaPlayerElement> = createRef()

  private async confirmAnnotation() {
    const payload: Omit<Annotation, 'id'> = {
      videostreamId: this._videostream!.id,
      observer: 'user@placeholder.com',
      observation: this._observation,
      timespan: { start: this._keypoints[0], end: this._keypoints[1] },
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
        `${import.meta.env.VITE_API_HOST}/api/v1/annotations?videostream=${id}&order=Timespan.Start`
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
      .editMode=${this._mode === 'editor'}
      @play=${this.play} 
      @pause=${this.pause}
      @seek=${(e: CustomEvent) => {
        console.log(e)
        this._seekTo = e.detail
      }}
    ></playback-controls>`

    const observationEditor = html`


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
            .annotations=${this._annotations}
            .currentTime=${this._currentTime}
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
            this._keypoints.length !== 2 || Object.keys(this._observation).length === 0
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
          <div class="add-keypoint">
            <button class="btn btn-transparent" @click=${this.addKeyPoint}>Add keypoint</button>
            ${repeat(
              this._keypoints,
              (k: VideoTime) => html`
                <span class="keypoint">
                  ${k}
                  <button class="btn-icon btn-transparent" @click=${() => {
                    this._keypoints = this._keypoints.filter((v) => v !== k)
                  }}>✕</button>
                </span>
              `
            )}
          </div>
          `

    return html`
      <div class="root">
        <div class="row">  
          ${video}
          ${overlay}
          ${aside}
        </div>
        ${playbackControls}
      </div>`
  }

  static styles = css`
    ${unsafeCSS(resetcss)}
    ${unsafeCSS(btncss)}
    ${unsafeCSS(vidstackcss)}

    .root {
      --video-ratio: 4 / 3;
      
      height: calc(100vh - 12.75rem);
      border-radius: 0.5rem;
      overflow: clip;
      display: flex;
      flex-direction: column;
    }

    .row {
      display: grid;
      height: calc(100% - 3.5rem);
      grid-template-rows: 1fr;
      grid-template-columns: min-content 1fr;
      grid-template-areas: "video-player annotations";
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
      display: flex;
      flex-direction: column;
    }
    media-player {
      grid-area: video-player;
      aspect-ratio: var(--video-ratio);
      height: 100%;
    }
    annotation-overlay, bounding-box-creator{
      grid-area: video-player;
      z-index: 100;
    }
    .add-keypoint {
      grid-area: video-player;
      z-index: 200;
      width: 100%;
      height: 100%;
      display: flex;
      align-items: end;
      padding: 1rem;
      pointer-events: none;
      gap: 1rem;

      & button {
        pointer-events: auto;
      }
    }
    .keypoint {
      width: min-content;
      height: 2.5rem;
      border-radius: 999px;
      font-size: 1rem;
      padding-left: 1rem;
      white-space: nowrap;
      border: 1px solid;
      display: flex;
      gap: 1rem;
      align-items: center;
      background-color: rgba(0, 0, 0, 0.75);
      color: var(--gray-100);
      border-color: rgba(0, 0, 0, 0.75);
    }
    aside header {
      padding: 0.75rem 1rem;
      background: var(--blue-500);
      display: flex;
      align-items: center;
      justify-content: end;
      gap: 0.5rem;
      box-shadow: var(--shadow-md);

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

    annotation-list {
      width: 100%;
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
