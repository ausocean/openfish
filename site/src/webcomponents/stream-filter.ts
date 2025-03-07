import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import '@fooloomanzoo/datetime-picker/datetime-picker.js'
import './data-select.ts'

// TODO: add support for start and end times.
export type Filter = {
  captureSource?: number
  'timespan[start]'?: string
  'timespan[end]'?: string
}

export type FilterUpdateEvent = CustomEvent<Filter>

@customElement('stream-filter')
export class StreamFilter extends TailwindElement {
  @state()
  protected _captureSource: number | null = null

  @state()
  protected _startTime: string | undefined

  @state()
  protected _endTime: string | undefined

  dispatchFilterUpdateEvent() {
    const detail: Filter = {}
    if (this._captureSource !== null) {
      detail.captureSource = this._captureSource
    }

    if (this._startTime != null) {
      detail['timespan[start]'] = this._startTime
    }
    if (this._endTime != null) {
      detail['timespan[end]'] = this._endTime
    }

    const options = {
      detail: detail,
      bubbles: true,
      composed: true,
    }
    this.dispatchEvent(new CustomEvent('filterupdate', options))
  }

  onCloseStartDatetime(event: CustomEvent) {
    this._startTime = event.detail.datetime
    this.dispatchFilterUpdateEvent()
  }

  onCloseEndDatetime(event: CustomEvent) {
    this._endTime = event.detail.datetime
    this.dispatchFilterUpdateEvent()
  }

  onSelect(event: InputEvent & { target: HTMLSelectElement }) {
    console.log(event.target.value)
    if (event.target.value !== '') {
      this._captureSource = Number(event.target.value)
    } else {
      this._captureSource = null
    }

    this.dispatchFilterUpdateEvent()
  }

  render() {
    return html`
    <aside class="">
        <h3 class="text-blue-800 font-bold">Filter Options</h3>
        <form class="flex flex-col gap-2">
            <fieldset>
              <legend class="text-blue-800">Filter by time of stream</legend>
              <label class="block text-slate-800">From:</label>
              <datetime-picker locale="en-AU" @input-picker-closed=${this.onCloseStartDatetime}></datetime-picker>

              <label class="block text-slate-800">Until:</label>
              <datetime-picker locale="en-AU" @input-picker-closed=${this.onCloseEndDatetime}></datetime-picker>
            </fieldset>

            <fieldset>
              <legend class="text-blue-800">Filter by capture source</legend>
              <label class="block text-slate-800">Capture source:</label>
              <data-select name="capturesource" src="/api/v1/capturesources" defaultText="Any" @input=${this.onSelect}></data-select>
            </fieldset>
        </form>
    </aside>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'stream-filter': StreamFilter
  }
}
