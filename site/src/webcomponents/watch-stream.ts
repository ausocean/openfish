import { TailwindElement } from './tailwind-element'
import { css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { ref, type Ref, createRef } from 'lit/directives/ref.js'
import type { MediaPlayerElement } from 'vidstack/elements'
import { extractVideoID } from '../utils/youtube'
import { repeat } from 'lit/directives/repeat.js'
import { instanceToPlain } from 'class-transformer'
import { BoundingBox, Keypoint } from '../api/annotation'
import { formatVideoTime, parseVideoTime } from '../utils/datetime'
import { client } from '../api'
import { unsafeSVG } from 'lit/directives/unsafe-svg.js'

import type { AnnotationWithJoins, VideoStreamWithJoins } from '@openfish/client'
import type { MouseoverAnnotationEvent } from './annotation-displayer'
import type { SpeciesSelectionEvent } from './species-selection'
import type { UpdateBoundingBoxEvent } from './bounding-box-creator'

import './annotation-displayer'
import './annotation-card'
import './playback-controls'
import './species-selection'
import './bounding-box-creator'

import vidstackcss from 'vidstack/player/styles/default/theme.css?lit'
import 'vidstack/player'
import 'vidstack/player/ui'

import x from '../icons/x.svg?raw'

@customElement('watch-stream')
export class WatchStream extends TailwindElement {
  @property({ type: Number })
  set streamID(val: number) {
    this.fetchVideoStream(val)
    this.fetchAnnotations(val)
  }

  @state()
  private _videostream: VideoStreamWithJoins | null = null

  @state()
  private _annotations: AnnotationWithJoins[] = []

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
  private _identification: number | null = null

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
    if (this._identification === null) {
      console.error('attempted to create annotation without identification')
      return
    }

    // Make annotation.
    await client.POST('/api/v1/annotations', {
      body: {
        videostream_id: this._videostream!.id,
        keypoints: instanceToPlain(this._keypoints) as any[],
        identification: this._identification,
      },
    })

    // Refetch annotations.
    await this.fetchAnnotations(this._videostream!.id)

    // Start playing.
    this._mode = 'playback'
    this.play()
  }

  private cancelAnnotation() {
    this._mode = 'playback'
    this._identification = null
  }

  private onSeek(e: CustomEvent) {
    this._seekTo = e.detail
  }

  async fetchVideoStream(id: number) {
    // Fetch video stream with ID.
    const { data, error } = await client.GET('/api/v1/videostreams/{id}', {
      params: {
        path: { id },
      },
    })

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this._videostream = data
    }
  }
  async fetchAnnotations(id: number) {
    // Fetch annotations for this video stream.
    // TODO: We should only fetch a small portion of the annotations near the current playback position.
    //       When the user plays the video we can fetch in more as needed.
    const { data, error } = await client.GET('/api/v1/annotations', {
      params: {
        query: {
          videostream: id,
          order: 'StartTime',
        },
      },
    })

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this._annotations = data.results
    }
  }

  render() {
    let filteredAnnotations: AnnotationWithJoins[] = []
    if (this._videostream != null) {
      // Filter annotations to only show those spanning the current playback time/position.
      filteredAnnotations = this._annotations.filter(
        (a) =>
          parseVideoTime(a.start) <= this._currentTime && this._currentTime <= parseVideoTime(a.end)
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
        .muted=${true}
      >
        <media-provider></media-provider>
        <media-video-layout></media-video-layout>
      </media-player>
    `

    const playbackControls = html` <playback-controls
      .playing=${this._playing}
      .duration=${this._duration}
      .currentTime=${this._currentTime}
      .annotations=${this._annotations}
      .videostream=${this._videostream}
      .editMode=${this._mode === 'editor'}
      @play=${this.play}
      @pause=${this.pause}
      @seek=${this.onSeek}
    ></playback-controls>`

    const asideContents =
      this._mode === 'playback'
        ? html` <header class="bg-blue-600">
              <h3 class="text-blue-50">Annotations</h3>
              <button class="btn variant-orange" @click=${this.addAnnotation}>
                + Add annotation
              </button>
            </header>
            <div class="scrollable">
              <annotation-list
                .annotations=${this._annotations}
                .currentTime=${this._currentTime}
                .activeAnnotation=${this._activeId}
                @mouseover-annotation=${(e: MouseoverAnnotationEvent) =>
                  (this._activeId = e.detail)}
                @seek=${this.onSeek}
              >
              </annotation-list>
            </div>`
        : html`
          <header class="bg-blue-600">
            <h3 class="text-blue-50">Add Annotation</h3>
            <button class="btn variant-slate" @click=${this.cancelAnnotation}>Cancel</button>
            <button class="btn variant-orange" @click=${this.confirmAnnotation} .disabled=${this._keypoints.length === 0}>Done</button>
          </header>
          <species-selection
            @selection=${(e: SpeciesSelectionEvent) => {
              this._identification = e.detail
              console.log(e.detail)
            }}
          >
          </species-selection>`

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
            <bounding-box-creator
              @updateboundingbox=${(e: UpdateBoundingBoxEvent) => (this._boundingBox = e.detail)}
            ></bounding-box-creator>
            <div class="keypoint-contain">
              <button
                class="btn variant-slate"
                @click=${this.addKeyPoint}
                .disabled=${
                  this._boundingBox === null ||
                  this._keypoints.map((k) => k.time).includes(this._currentTime)
                }
              >
                Add keypoint
              </button>
              ${repeat(
                this._keypoints,
                (k: Keypoint) => html`
                  <span
                    class="flex overflow-clip rounded-lg"
                  >
                    <button
                      class="btn variant-slate rounded-none"
                      @click=${() => (this._seekTo = k.time)}
                    >
                        ${formatVideoTime(k.time, true)}
                    </button>
                    <button
                      class="btn variant-slate rounded-none with-icon p-2 aspect-square"
                      @click=${() => {
                        this._keypoints = this._keypoints.filter((v) => v.time !== k.time)
                      }}
                    >
                      ${unsafeSVG(x)}
                    </button>
                  </span>
                `
              )}
            </div>
          `

    return html`
      <main class="bg-blue-700 overflow-clip rounded-lg h-full">
        <div class="flex h-[calc(100%-3rem)]">
          <div class="video relative bg-blue-950">${video} ${overlay}</div>
          <aside class="w-full flex flex-col bg-blue-700 overflow-y-hidden">
            ${asideContents}
          </aside>
        </div>
        ${playbackControls}
      </main>
    `
  }

  static styles = [
    TailwindElement.styles!,
    vidstackcss,
    css`
      :host {
        --video-ratio: 4 / 3;
      }

      .video {
        aspect-ratio: var(--video-ratio);
      }
      media-player,
      annotation-overlay,
      bounding-box-creator {
        aspect-ratio: var(--video-ratio);
        height: 100%;
        position: absolute;
        inset: 0;
      }
      annotation-overlay,
      bounding-box-creator {
        z-index: 100;
      }
      .keypoint-contain {
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
        display: flex;
        align-items: center;
        justify-content: end;
        gap: 0.5rem;
        box-shadow: var(--shadow-md);

        & :first-child {
          margin-right: auto;
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

      species-selection {
        height: calc(100% - 4rem);
        overflow: hidden;
      }

      table {
        width: 100%;

        & tbody tr :nth-child(3) {
          width: 0;
        }
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'watch-stream': WatchStream
  }
}
