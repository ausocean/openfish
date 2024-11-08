import { LitElement, css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { VideoStream } from '../api/videostream'
import resetcss from '../styles/reset.css?lit'
import btncss from '../styles/buttons.css?lit'

import './annotation-displayer'
import './annotation-card'
import './playback-controls'
import './observation-editor'
import './bounding-box-creator'

import vidstackcss from 'vidstack/player/styles/default/theme.css?lit'
import 'vidstack/player'
import 'vidstack/player/ui'

import type { MouseoverAnnotationEvent } from './annotation-displayer'
import type { ObservationEvent } from './observation-editor'
import type { UpdateBoundingBoxEvent } from './bounding-box-creator'

import { ref, type Ref, createRef } from 'lit/directives/ref.js'
import type { MediaPlayerElement } from 'vidstack/elements'
import { extractVideoID } from '../utils/youtube'
import { repeat } from 'lit/directives/repeat.js'
import { instanceToPlain, plainToInstance } from 'class-transformer'
import { Annotation, BoundingBox, Keypoint } from '../api/annotation'
import { formatVideoTime } from '../utils/datetime'

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
  private _keypoints: Keypoint[] = []

  @state()
  private _boundingBox: [number, number, number, number] | null = null

  private addKeyPoint() {
    if (this._boundingBox === null) {
      return
    }
    const box = new BoundingBox(
      this._boundingBox[0],
      this._boundingBox[1],
      this._boundingBox[2],
      this._boundingBox[3]
    )
    this._keypoints.push(new Keypoint(this._currentTime, box))

    this.requestUpdate()
  }

  private addAnnotation() {
    this._mode = 'editor'
    this.pause()
    this._keypoints = []
  }

  playerRef: Ref<MediaPlayerElement> = createRef()

  private async confirmAnnotation() {
    const payload = {
      videostreamId: this._videostream!.id,
      observer: 'user@placeholder.com',
      observation: this._observation,
      keypoints: instanceToPlain(this._keypoints),
    }

    // Make annotation.
    await fetch('/api/v1/annotations', {
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
      const res = await fetch(`/api/v1/videostreams/${id}`)
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
      const res = await fetch(`/api/v1/annotations?videostream=${id}&order=StartTime`)
      const json = await res.json()
      this._annotations = plainToInstance<Annotation, object[]>(Annotation, json.results)
      console.log(this._annotations)
    } catch (error) {
      console.error(error) // TODO: handle errors.
    }
  }

  render() {
    let filteredAnnotations: Annotation[] = []
    if (this._videostream != null) {
      // Filter annotations to only show those spanning the current playback time/position.
      filteredAnnotations = this._annotations.filter(
        (a) => a.start <= this._currentTime && this._currentTime <= a.end
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
            this._keypoints.length === 0 || Object.keys(this._observation).length === 0
          }>Done</button>
        </header>
        <observation-editor .observation=${this._observation} @observation=${(
          ev: ObservationEvent
        ) => {
          this._observation = ev.detail
        }}>
      </aside>`

    const overlay =
      this._mode === 'playback'
        ? html`
          <annotation-overlay
            .annotations=${filteredAnnotations}
            .activeAnnotation=${this._activeId}
            .currentTime=${this._currentTime}
            @mouseover-annotation=${(e: MouseoverAnnotationEvent) => (this._activeId = e.detail)}
          ></annotation-overlay>
          `
        : html`
          <bounding-box-creator @updateboundingbox=${(e: UpdateBoundingBoxEvent) =>
            (this._boundingBox = e.detail)}></bounding-box-creator>
          <div class="keypoint-container">
            <button class="btn btn-transparent" @click=${this.addKeyPoint} .disabled=${this._boundingBox === null || this._keypoints.map((k) => k.time).includes(this._currentTime)}>Add keypoint</button>
            ${repeat(
              this._keypoints,
              (k: Keypoint) => html`
                <span class="keypoint">
                  ${formatVideoTime(k.time, true)}
                  <button class="btn-icon btn-transparent" @click=${() => {
                    this._keypoints = this._keypoints.filter((v) => v.time !== k.time)
                  }}>✕</button>
                </span>
              `
            )}
          </div>
          `

    return html`
      <div class="root">
        <div class="row">  
          <div class="video">
            ${video}
            ${overlay}
          </div>
          ${aside}
        </div>
        ${playbackControls}
      </div>`
  }

  static styles = css`
    ${resetcss}
    ${btncss}
    ${vidstackcss}

    .root {
      --video-ratio: 4 / 3;
      
      height: calc(100vh - 12.75rem);
      border-radius: 0.5rem;
      overflow: clip;
      display: flex;
      flex-direction: column;
    }

    .row {
      display: flex;
      height: calc(100% - 3.5rem);
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
      background-color: var(--blue-700);
      display: flex;
      flex-direction: column;
      width: 100%;
    }
    .video {
      aspect-ratio: var(--video-ratio);
      height: 100%;
      position: relative;
    }
    media-player, annotation-overlay, bounding-box-creator {
      aspect-ratio: var(--video-ratio);
      height: 100%;
      position: absolute;
      inset: 0;
    }
    annotation-overlay, bounding-box-creator {
      z-index: 100;
    }
    .keypoint-container {
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      height: min-content;
      z-index: 200;
      display: flex;
      padding: 1rem;
      padding-top: 0;
      gap: 1rem;
      overflow-x: scroll;
      

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
      inset: 0;
    }

    annotation-list {
      width: 100%;
    }

    observation-editor {
      height: calc(100% - 4rem);
      overflow: hidden;
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
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'watch-stream': WatchStream
  }
}
