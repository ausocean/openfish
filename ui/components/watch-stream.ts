import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { ref, type Ref, createRef } from 'lit/directives/ref.js'
import type { MediaPlayerElement } from 'vidstack/elements'
import { extractVideoID } from '../utils/youtube'
import { repeat } from 'lit/directives/repeat.js'
import { instanceToPlain } from 'class-transformer'
import { BoundingBox, Keypoint } from '../utils/keypoints.ts'
import { formatVideoTime, parseVideoTime } from '../utils/datetime'
import { unsafeSVG } from 'lit/directives/unsafe-svg.js'

import type { AnnotationWithJoins, OpenfishClient, VideoStreamWithJoins } from '@openfish/client'
import type { MouseoverAnnotationEvent } from './annotation-displayer'
import type { SpeciesSelectionEvent } from './species-selection'
import type { UpdateBoundingBoxEvent } from './bounding-box-creator'

import './annotation-displayer'
import './annotation-card'
import './timeline-heatmap'
import './species-selection'
import './bounding-box-creator'

import vidstackcss from 'vidstack/player/styles/default/theme.css?lit'
import 'vidstack/player'
import 'vidstack/player/ui'

import caretLeft from '../icons/caret-left.svg?raw'
import caretRight from '../icons/caret-right.svg?raw'
import caretDoubleLeft from '../icons/caret-double-left.svg?raw'
import caretDoubleRight from '../icons/caret-double-right.svg?raw'
import play from '../icons/play.svg?raw'
import pause from '../icons/pause.svg?raw'
import replay from '../icons/replay.svg?raw'
import x from '../icons/x.svg?raw'
import { clientContext } from '../utils/context'
import { consume } from '@lit/context'

@customElement('watch-stream')
export class WatchStream extends TailwindElement {
  @consume({ context: clientContext, subscribe: true })
  client!: OpenfishClient

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
  private _duration = 0

  @state()
  private _seekTo: number | null = null

  @state()
  private _mode: 'playback' | 'editor' = 'playback'

  private play() {
    this.playerRef.value?.play()
  }

  private pause() {
    this.playerRef.value?.pause()
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
    await this.client.POST('/api/v1/annotations', {
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

  private fwd(seconds: number) {
    this._seekTo = Math.min(this._duration, this._currentTime + seconds)
  }

  private bwd(seconds: number) {
    this._seekTo = Math.max(0, this._currentTime - seconds)
  }

  async fetchVideoStream(id: number) {
    // Fetch video stream with ID.
    const { data, error } = await this.client.GET('/api/v1/videostreams/{id}', {
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
    const { data, error } = await this.client.GET('/api/v1/annotations', {
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

    const playbackControls = html`
		<div class="flex w-full px-4 py-1 gap-2 bg-blue-600 text-slate-50 items-center">
        <media-play-button class="btn size-sm variant-orange w-28 justify-center group with-icon gap-2">
          <span class="contents not-group-data-[paused]:hidden group-data-[ended]:hidden whitespace-nowrap">${unsafeSVG(play)} Play</span>
          <span class="contents not-group-data-[ended]:hidden whitespace-nowrap">${unsafeSVG(replay)} Replay</span>
          <span class="contents group-data-[paused]:hidden whitespace-nowrap">${unsafeSVG(pause)} Pause</span>
        </media-play-button>

        <div class="flex w-32 *:px-0 rounded-md overflow-clip">
          <button id="bwd-5" class="btn size-sm variant-blue with-icon rounded-none w-18" @click="${() => this.bwd(5)}">${unsafeSVG(caretDoubleLeft)}</button>
          <button id="bwd-1" class="btn size-sm variant-blue with-icon rounded-none w-18" @click="${() => this.bwd(1)}">${unsafeSVG(caretLeft)}</button>
          <button id="fwd-1" class="btn size-sm variant-blue with-icon rounded-none w-18" @click="${() => this.fwd(1)}">${unsafeSVG(caretRight)}</button>
          <button id="fwd-5" class="btn size-sm variant-blue with-icon rounded-none w-18" @click="${() => this.fwd(5)}">${unsafeSVG(caretDoubleRight)}</button>

          <tooltip-elem for="bwd-5" trigger="hover" placement="top" class="text-nowrap">Jump back 5 seconds</tooltip-elem>
          <tooltip-elem for="bwd-1" trigger="hover" placement="top" class="text-nowrap">Jump back 1 second</tooltip-elem>
          <tooltip-elem for="fwd-1" trigger="hover" placement="top" class="text-nowrap">Jump forward 1 second</tooltip-elem>
          <tooltip-elem for="fwd-5" trigger="hover" placement="top" class="text-nowrap">Jump forward 5 seconds</tooltip-elem>
        </div>


        <div class="w-full">
            <media-time-slider
                class="group relative inline-flex h-10 w-full cursor-pointer touch-none select-none items-center outline-none aria-hidden:hidden"
                >
                <!-- Track -->
                <div
                    class="relative h-6 w-full rounded-sm bg-blue-500 group-data-[focus]:ring-[3px] overflow-clip"
                >
                    <div
                        class="absolute z-0 h-full w-[var(--slider-progress)] bg-blue-400 will-change-[width]"
                    ></div>
                    <div
                      class="absolute z-10 h-full w-[var(--slider-fill)] bg-blue-300/50 will-change-[width]"
                    ></div>
                    <timeline-heatmap class="absolute z-20 h-full w-full" .annotations=${this._annotations} .duration=${this._duration}/>
                </div>

                <!-- Preview -->
                <media-slider-preview
                    class="pointer-events-none flex flex-col items-center opacity-0 transition-opacity duration-200 data-[visible]:opacity-100"
                    noClamp
                >
                    <media-slider-thumbnail
                    class="block h-[var(--thumbnail-height)] max-h-[160px] min-h-[80px] w-[var(--thumbnail-width)] min-w-[120px] max-w-[180px] overflow-hidden border border-white bg-black"
                    src="/your_thumbnails.vtt"
                    ></media-slider-thumbnail>
                    <media-slider-value
                    class="rounded-sm bg-black px-2 py-px text-[13px] font-medium text-white"
                    ></media-slider-value>
                </media-slider-preview>

            <!-- Thumb -->
            <div
                class="absolute left-[var(--slider-fill)] top-1/2 z-30 h-8 w-1 -translate-x-1/2 -translate-y-1/2 bg-red-400 ring-white/40 opacity-75 transition-opacity group-data-[active]:opacity-100  will-change-[left]  group-data-[dragging]:ring-4"
            ></div>
            </media-time-slider>
        </div>

        <span class="p-1 w-48 whitespace-nowrap text-right text-blue-50">
            <media-time class="inline" type="current"></media-time>
            <span class="mx-1 text-blue-200">/</span>
            <media-time class="inline" type="duration"></media-time>
        </span>
    </div>`

    const asideContents =
      this._mode === 'playback'
        ? html` <header class="bg-blue-600 flex p-4 align-center shadow-sm border-b border-b-blue-500">
              <h3 class="text-blue-50 text-lg flex-1">Annotations</h3>
              <button class="btn variant-orange" @click=${this.addAnnotation}>
                + Add annotation
              </button>
            </header>
            <annotation-list
              class="h-full w-full overflow-y-scroll"
              .annotations=${this._annotations}
              .currentTime=${this._currentTime}
              .activeAnnotation=${this._activeId}
              @mouseover-annotation=${(e: MouseoverAnnotationEvent) => (this._activeId = e.detail)}
              @seek=${this.onSeek}
            >
            </annotation-list>
            `
        : html`
          <header class="bg-blue-600 flex p-4 align-center gap-2">
            <h3 class="text-blue-50 text-lg flex-1">Add Annotation</h3>
            <button class="btn variant-slate" @click=${this.cancelAnnotation}>Cancel</button>
            <button class="btn variant-orange" @click=${this.confirmAnnotation} .disabled=${this._keypoints.length === 0}>Done</button>
          </header>
          <species-selection
            class="h-full w-full"
            @selection=${(e: SpeciesSelectionEvent) => {
              this._identification = e.detail
              console.log(e.detail)
            }}
          >
          </species-selection>
          `

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
            <div class="keypoint-contain absolute h-min-content flex gap-4 p-4 pt-0 left-0 right-0 bottom-0">
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
  		<media-player
        class="h-full flex flex-col"
        ${ref(this.playerRef)}
        title="Openfish Video"
        src="youtube/${videoID}"
        .currentTime="${this._seekTo}"
        @time-update=${(e: CustomEvent<{ currentTime: number }>) =>
          (this._currentTime = e.detail.currentTime)}
        @duration-change=${(e: CustomEvent<number>) => (this._duration = e.detail)}
        .muted=${true}
        >
        <div class="flex h-[calc(100%-3rem)] w-full">

            <div class="aspect-[4/3] relative">
                <media-provider>
                    <media-poster
                        class="blur-xl absolute inset-0 block h-full w-full bg-blue-950 opacity-0 transition-opacity data-[visible]:opacity-100 [&>img]:h-full [&>img]:w-full [&>img]:object-cover"
                        ?src=${videoID !== null ? `https://i.ytimg.com/vi/${videoID}/maxresdefault.jpg` : null}
                    ></media-poster>
                </media-provider>
                <media-video-layout ></media-video-layout>
                <div class="absolute inset-0 z-10">
                    ${overlay}
                </div>
            </div>

            <aside class="flex flex-col bg-blue-700 overflow-y-hidden w-full">
              ${asideContents}
            </aside>
        </div>
        ${playbackControls}
    </media-player>
		</main>
    `
  }

  static styles = [TailwindElement.styles!, vidstackcss]
}

declare global {
  interface HTMLElementTagNameMap {
    'watch-stream': WatchStream
  }
}
