import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { client } from '../api'
import type { Species } from '@openfish/client'

import './species-thumb'

export type SpeciesSelectionEvent = CustomEvent<number | null>
type TextInputEvent = InputEvent & { target: HTMLInputElement }

@customElement('species-selection')
export class SpeciesSelection extends TailwindElement {
  @state()
  protected _keys: string[] = []

  @state()
  protected _vals: string[] = []

  @property({ type: Object })
  selection: Species | null = null

  @state()
  private _speciesList: Species[] = []

  @state()
  offset = 0

  @state()
  _search = ''

  private selectSpecies(species: Species) {
    this.selection = species
    this.dispatchEvent(new CustomEvent('selection', { detail: this.selection?.id }))
  }

  connectedCallback() {
    super.connectedCallback()
    this.fetchMore()
  }

  private async fetchMore() {
    const { data, error } = await client.GET('/api/v1/species', {
      params: {
        query: {
          limit: 20,
          offset: this.offset,
          search: this._search.length > 0 ? this._search : undefined,
        },
      },
    })

    if (error !== undefined) {
      console.error(error)
    }

    if (data !== undefined) {
      this._speciesList.push(...data.results)
      this.offset += 20
      this.requestUpdate()
    }
  }

  private async search(e: TextInputEvent) {
    this._search = e.target.value
    this.offset = 0
    this._speciesList = []
    this.fetchMore()
  }

  render() {
    return html`
      <header
        class="bg-blue-600 px-3 py-2 border-b border-b-blue-500 shadow-sm"
      >
        <input
          type="text"
          class="bg-blue-700 border border-blue-800 text-blue-50 w-full rounded-md placeholder:text-blue-300"
          placeholder="Search species"
          @input=${this.search}
        />
      </header>
      <div class="relative overflow-y-scroll h-[calc(100%-3rem)]">
        <div class="absolute inset-0">
          <ul
            class="grid overflow-hidden gap-4 p-4 grid-cols-2 auto-rows-auto"
          >
            ${repeat(
              this._speciesList,
              (species) => species.id,
              (species) => html`
                <species-thumb
                    class="hover:ring-2 hover:ring-sky-400/50 data-selected:ring-2 data-selected:ring-sky-400 data-selected:ring-offset-4 ring-offset-1 ring-offset-blue-700 transition-shadow rounded-md"
                    .species=${species}
                    @click=${() => this.selectSpecies(species)}
                    ?data-selected=${this.selection?.id === species.id}
                ></species-thumb>
            `
            )}
          </ul>
          <footer class="w-full pb-4">
            <button class="mx-auto btn variant-slate" @click=${this.fetchMore}>
              Load more
            </button>
          </footer>
        </div>
      </div>
    `
  }

  static styles = TailwindElement.styles
}

declare global {
  interface HTMLElementTagNameMap {
    'species-selection': SpeciesSelection
  }
}
