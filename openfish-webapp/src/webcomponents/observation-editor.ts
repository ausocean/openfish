import { TailwindElement } from './tailwind-element'
import { css, html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { zip } from '../utils/array-utils'
import resetcss from '../styles/reset.css?lit'
import btncss from '../styles/buttons.css?lit'

type Species = { species: string; common_name: string; images?: Image[] }
type Image = { src: string; attribution: string }

export type ObservationEvent = CustomEvent<Record<string, string>>

abstract class AbstractObservationEditor extends TailwindElement {
  @state()
  protected _keys: string[] = []

  @state()
  protected _vals: string[] = []

  @property({ type: Object })
  set observation(val: Record<string, string>) {
    this._keys = Object.keys(val)
    this._vals = Object.values(val)
  }
  get observation() {
    return Object.fromEntries(zip(this._keys, this._vals))
  }

  dispatchObservationEvent() {
    this.dispatchEvent(
      new CustomEvent('observation', { detail: this.observation }) as ObservationEvent
    )
  }
}

@customElement('observation-editor')
export class ObservationEditor extends AbstractObservationEditor {
  @state()
  private _editorMode: 'simple' | 'advanced' = 'simple'

  _onObservation(ev: ObservationEvent) {
    this.observation = ev.detail
    this.dispatchObservationEvent()
  }

  render() {
    return html`
    <menu class="text-slate-50 flex py-2 px-4 gap-2 bg-blue-600">
      <h4 class="flex-1">Observation</h4>
      <button class="btn size-sm ${
        this._editorMode === 'simple' ? 'variant-slate' : 'variant-blue'
      }"  @click=${() => {
        this._editorMode = 'simple'
      }}>Select species</button>
      <button class="btn size-sm ${
        this._editorMode === 'advanced' ? 'variant-slate' : 'variant-blue'
      }" @click=${() => {
        this._editorMode = 'advanced'
      }}>Additional observations</button>
    </menu>

      ${
        this._editorMode === 'simple'
          ? html`<species-selection .observation=${this.observation} @observation=${this._onObservation} ></species-selection>`
          : html`<advanced-editor .observation=${this.observation} @observation=${this._onObservation}></advanced-editor>`
      }
    `
  }
}

@customElement('species-selection')
export class SpeciesSelection extends AbstractObservationEditor {
  @state()
  private _speciesList: Species[] = []

  @state()
  offset = 0

  @state()
  _search = ''

  private selectSpecies(val: string) {
    const { species, common_name } = this._speciesList.find((s) => s.species === val)!
    this.observation = { species, common_name }
    this.dispatchObservationEvent()
  }

  connectedCallback() {
    super.connectedCallback()
    this.fetchMore()
  }

  private async fetchMore() {
    try {
      const params = new URLSearchParams({
        limit: '20',
        offset: this.offset.toString(),
      })
      if (this._search.length > 0) {
        params.set('search', this._search)
      }
      const res = await fetch(`/api/v1/species/recommended?${params}`)
      this._speciesList.push(...(await res.json()).results)
      this.offset += 20
      this.requestUpdate()
    } catch (error) {
      console.error(error)
    }
  }

  private async search(e: TextInputEvent) {
    this._search = e.target.value
    this.offset = 0
    this._speciesList = []
    this.fetchMore()
  }

  render() {
    const renderSpecies = ({ species, common_name, images }: Species) => html`
      <li class="card overflow-clip relative p-0 transition-colors hover:bg-blue-200 data-selected:bg-blue-200 data-selected:border-sky-400 data-selected:shadow-md data-selected:shadow-sky-500/50 cursor-pointer" 
        @click=${() => this.selectSpecies(species)}
        ?data-selected=${this.observation.species === species}
        >
        <div title=${images?.at(0)?.attribution} class="aspect-square rounded-full w-5 flex items-center justify-center text-sm bg-slate-950/75 text-white absolute top-2 right-2">&copy;</div>
        <img src=${images?.at(0)?.src ?? 'placeholder.svg'} class="w-full object-cover aspect-[4/3]"/>
        <div class="px-2 font-bold mt-1">${common_name}</div>
        <div class="px-2 pb-2 text-sm">${species}</div>
      </li>
  `

    return html`
    <header class="bg-blue-600 px-3 py-2 border-b border-b-blue-500 shadow-sm">
      <input 
        type="text"
        class="bg-blue-700 border border-blue-800 text-blue-50 w-full rounded-md placeholder:text-blue-300"
        placeholder="Search species" 
        @input=${this.search} 
      />
    </header>
    <div class="relative overflow-y-scroll h-[calc(100%-6rem)]">
      <div class="absolute inset-0">
        <ul class="grid overflow-y-scroll gap-4 p-4 grid-cols-2 auto-rows-auto">
          ${repeat(this._speciesList, renderSpecies)}
        </ul>
        <footer class="w-full pb-4">
          <button class="mx-auto btn variant-slate" @click=${this.fetchMore}>Load more</button>
        </footer>
      </div>
    </div>
    `
  }

  static styles = TailwindElement.styles
}

type TextInputEvent = InputEvent & { target: HTMLInputElement }

@customElement('advanced-editor')
export class AdvancedEditor extends AbstractObservationEditor {
  // Mutating an array doesn't trigger an update.
  // https://lit.dev/docs/components/properties/#mutating-properties
  addRow() {
    this._keys.push('')
    this._vals.push('')
    this.requestUpdate()
    this.dispatchObservationEvent()
  }

  updateKey(idx: number, key: string) {
    this._keys[idx] = key
    this.requestUpdate()
    this.dispatchObservationEvent()
  }

  updateVal(idx: number, val: string) {
    this._vals[idx] = val
    this.requestUpdate()
    this.dispatchObservationEvent()
  }

  render() {
    const rows = repeat(
      zip(this._keys, this._vals),
      ([key, val], idx) => html`
      <tr>
      <td><input type="text" .value=${key} @input=${(ev: TextInputEvent) =>
        this.updateKey(idx, ev.target.value)}></input></td>
      <td><input type="text" .value=${val} @input=${(ev: TextInputEvent) =>
        this.updateVal(idx, ev.target.value)}></input></td>
      </tr>
    `
    )

    return html`
    <div class="root">
      <div class="card">
      <table>
      <thead>
          <tr>
              <th>Property</th>
              <th>Value</th>
          </tr>
      </thead>
      <tbody>
      ${rows}
      <tr>
      <td colspan="2"><button class="btn-sm btn-orange btn-fullwidth" @click=${this.addRow}>+ Add information</button></td>
      </tr>
      </tbody>
      </table>
      </div>
    </div>
    `
  }

  static styles = css`
    ${resetcss}
    ${btncss}
    .root {
      padding: 1rem;
    }
    .card {
      background-color: var(--gray-50);
      padding: 1rem;
      border-radius: .5rem;
      box-shadow:  var(--shadow-sm);
      width: 100%;
    }

    table {
      font-size: 0.8rem;
      width: 100%;
      border-spacing: 0;
    }
    table th:nth-child(1) {
      width: 40%;
    }
    table th {
      text-align: left;
      border-bottom: 1px solid var(--gray-200);
    }
    table th, table td {
      padding: 0.25rem;
    }
    table tbody tr:hover {
      background-color: var(--gray-50)
    }

    tbody:last-child td {

      padding: 0.5rem;
    }
    .btn-fullwidth {
      text-align: center;
      width: 100%;  
    }

    input {
      width: 100%;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'advanced-editor': AdvancedEditor
    'species-selection': SpeciesSelection
  }
}
