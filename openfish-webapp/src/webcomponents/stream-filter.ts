import { LitElement, css, html } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import resetcss from '../styles/reset.css?lit'
import '@fooloomanzoo/datetime-picker/datetime-picker.js'
import type { SelectCaptureSourceEvent } from './capture-source-dropdown.ts'

// TODO: add support for start and end times.
export type Filter = {
  captureSource?: number
  'timespan[start]'?: string
  'timespan[end]'?: string
}

export type FilterUpdateEvent = CustomEvent<Filter>

@customElement('stream-filter')
export class StreamFilter extends LitElement {
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

  onSelectCaptureSource(event: SelectCaptureSourceEvent) {
    this._captureSource = event.detail
    this.dispatchFilterUpdateEvent()
  }

  onCloseStartDatetime(event: CustomEvent) {
    this._startTime = event.detail.datetime
    this.dispatchFilterUpdateEvent()
  }

  onCloseEndDatetime(event: CustomEvent) {
    this._endTime = event.detail.datetime
    this.dispatchFilterUpdateEvent()
  }

  render() {
    return html`
    <aside>
        <h3>Filter Options</h3>
        <form>
            <fieldset>
            <legend>Filter by time of stream</legend>
            <label>From:</label>
            <datetime-picker locale="en-AU" @input-picker-closed=${this.onCloseStartDatetime}></datetime-picker>

            <label>Until:</label>
            <datetime-picker locale="en-AU" @input-picker-closed=${this.onCloseEndDatetime}></datetime-picker>
            </fieldset>

            <fieldset>
            <legend>Filter by capture source</legend>
            <label>Capture source:</label>
            <capture-source-dropdown @selectcapturesource=${this.onSelectCaptureSource}></capture-source-dropdown>
            </fieldset>
        </form>
    </aside>
    `
  }

  static styles = css`
    ${resetcss}
    aside {
      background-color: var(--gray-50);
      border-radius: 0.25rem;
      border: 1px solid var(--gray-200);
    }

    h3 {
        margin-top: 0;
        margin-bottom: 0;
        padding: .5rem 1.5rem; 
        background-color: var(--gray-50);
        border-bottom: 1px solid var(--gray-200);
    }

    form {
        padding: .5rem .5rem; 
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    datetime-picker {
        margin-top: .25rem;
    }

    capture-source-dropdown {
        width: 100%;
    }

    fieldset {
        border-radius: 0.25rem;
    }
    `
}

declare global {
  interface HTMLElementTagNameMap {
    'stream-filter': StreamFilter
  }
}
