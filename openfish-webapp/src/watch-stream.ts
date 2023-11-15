import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { Annotation, VideoStream } from './api.types'
import { repeat } from 'lit/directives/repeat.js'
import { videotimeToDatetime } from './datetime'
import { resetcss } from './reset.css'

@customElement('watch-stream')
export class WatchStream extends LitElement {
  @property({ type: Object })
  videostream: VideoStream | null = null

  @property({ type: Array })
  annotations: Annotation[] = []

  @property({ type: Number })
  activeId: number | null = null

  @property({ type: Number })
  currentTime = 0

  @property({ type: Boolean })
  playing = false

  @property({ type: Number })
  duration = 0

  @property({ type: Number })
  seekTo: number | null = null

  private play() {
    this.playing = true
  }

  private pause() {
    this.playing = false
  }

  connectedCallback() {
    super.connectedCallback()
    const url = new URL(document.location.href)
    const id = url.searchParams.get('id')

    if (id) {
      this.fetchData(id)
    }
  }

  async fetchData(id: string) {
    try {
      // Fetch video stream with ID.
      const res = await fetch(`http://localhost:3000/api/v1/videostreams/${id}`)
      this.videostream = (await res.json()) as VideoStream
    } catch (error) {
      console.error(error) // TODO: handle errors.
    }
    try {
      // Fetch annotations for this video stream.
      // TODO: We should only fetch a small portion of the annotations near the current playback position.
      //       When the user plays the video we can fetch in more as needed.
      const res = await fetch(`http://localhost:3000/api/v1/annotations?videostream=${id}`)
      const json = await res.json()
      this.annotations = json.results as Annotation[]
    } catch (error) {
      console.error(error) // TODO: handle errors.
    }
  }

  render() {
    let filteredAnnotations: Annotation[] = []
    if (this.videostream != null) {
      // Convert playback time in seconds to a datetime.
      const playbackDatetime = videotimeToDatetime(this.videostream?.startTime, this.currentTime)

      // Filter annotations to only show those spanning the current playback time/position.
      filteredAnnotations = this.annotations.filter(
        (an: Annotation) =>
          new Date(an.timespan.start).getTime() <= playbackDatetime.getTime() &&
          playbackDatetime.getTime() <= new Date(an.timespan.end).getTime()
      )
    }

    const annotationList = repeat(filteredAnnotations, (annotation: Annotation) => {
      return html`
      <div>
        <annotation-card 
          @mouseover-annotation=${(e: CustomEvent) => (this.activeId = e.detail)} 
          .annotation=${annotation} 
          .outline=${this.activeId === annotation.id}
        />
      </div>`
    })

    return html`
      <div class="grid">
        <header>
          <h1>Video Playback: ${this.videostream?.id ?? ''}</h1>
        </header>

        <main>
          <video-player 
            .videostream=${this.videostream}
            .annotations=${filteredAnnotations}
            .activeAnnotation=${this.activeId}
            .seekTo=${this.seekTo}
            .playing=${this.playing}
            @mouseover-annotation=${(e: CustomEvent) => (this.activeId = e.detail)}
            @timeupdate=${(e: CustomEvent) => (this.currentTime = e.detail)} 
            @durationchange=${(e: CustomEvent) => (this.duration = e.detail)}
            @loadeddata=${() => (this.playing = true)}
            >
          </video-player>
        </main>

        <aside>
          <h2>Annotations</h2>
          <div class="annotation-list">
            ${annotationList}
          </div>
        </aside>

        <footer>
          <playback-controls 
            .playing=${this.playing} 
            .duration=${this.duration} 
            .currentTime=${this.currentTime}
            .annotations=${this.annotations}
            .videostream=${this.videostream}
            @play=${this.play} 
            @pause=${this.pause}
            @seek=${(e: CustomEvent) => (this.seekTo = e.detail)}
          />
        </footer>
      </div>`
  }

  static styles = css`
    ${resetcss}

    :host {
      width: min(100vw, 95rem);
      height: 100vh;
    }

    .grid {
      padding: 2rem;
      height: 100%; 
      width: 100%; 
      display: grid;
      grid-template-rows: min-content 1fr min-content;
      grid-template-columns: 1fr 32rem;
      grid-template-areas:
        "header header"
        "video-player annotations"
        "controls controls";
      gap: 2rem;
    }

    h1 {
      margin: 0;
      padding: .25rem .5rem; 
      border-bottom: 1px solid var(--gray1);
    }
    span.subtitle {
      color: var(--gray3);
      font-weight: normal;
    }
    header {
      grid-area: header;
    }
    aside {
      grid-area: annotations;
      overflow-y: hidden;
      background-color: var(--gray1);
      border-radius: 0.25rem;
      border: 1px solid var(--gray2);
    }
    main {
      grid-area: video-player;
    }
    footer {
      grid-area: controls;
    }
    h2 {
      margin-top: 0;
      margin-bottom: 0;
      padding: .5rem 1.5rem; 
      background-color: var(--gray0);
      border-bottom: 1px solid var(--gray2);
    }
    .annotation-list {
      height: 100%;
      overflow-y: scroll;
      display: flex;
      flex-direction: column;
      gap: 1rem;
      padding: 1rem;
    }
    annotation-card.active {
        border: 1px solid var(--secondary)
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'watch-stream': WatchStream
  }
}
