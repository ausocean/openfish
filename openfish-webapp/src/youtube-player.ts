import { LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import useYoutubePlayer from 'youtube-player'
import type { YouTubePlayer } from 'youtube-player/dist/types'

export type TimeUpdateEvent = CustomEvent<number>
export type DurationChangeEvent = CustomEvent<number>
export type LoadedEvent = Event

@customElement('youtube-player')
export class YouTubePlayerElement extends LitElement {
  @property({ type: String })
  url = ''

  private _playing = false
  @property()
  get playing() {
    return this._playing
  }
  set playing(val: boolean) {
    this._playing = val
    if (val) {
      this._player?.playVideo()
    } else {
      this._player?.pauseVideo()
    }
  }

  private _player: YouTubePlayer | null = null
  private _element: HTMLElement | null = null

  seek(time: number) {
    this._player?.seekTo(time, true)
  }

  @property()
  get seekTo() {
    return null
  }
  set seekTo(time: number | null) {
    if (time != null) {
      this._player?.seekTo(time, true)
    }
  }

  constructor() {
    super()

    // Use youtube player api to create a Youtube player.
    this._element = document.createElement('div')
    this._player = useYoutubePlayer(this._element, {
      height: '100%',
      width: '100%',
      playerVars: {
        controls: 0,
        disablekb: 1,
      },
    })

    // Emit timeupdate events by polling the current time.
    setInterval(async () => {
      if (this._player != null) {
        const currentTime = await this._player?.getCurrentTime()
        this.dispatchEvent(new CustomEvent('timeupdate', { detail: currentTime }))
      }
    }, 50)

    // Emit play, pause and durationchange events on stateChange events.
    this._player.on('stateChange', async (e) => {
      if (e.data === 1) {
        const duration = (await this._player?.getDuration()) ?? 0
        this.dispatchEvent(new CustomEvent('durationchange', { detail: duration }))
        this.dispatchEvent(new Event('loadeddata'))
      }
    })
  }

  render() {
    // Parse youtube url for video ID.
    // https://www.youtube.com/watch?v=faolURG_uXQ
    const url = new URL(this.url)
    const videoID = url.searchParams.get('v')

    // Play video.
    if (videoID !== null) {
      this._player?.loadVideoById(videoID)
      this._player?.playVideo()
    }

    return this._element
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'youtube-player': YouTubePlayerElement
  }
}
