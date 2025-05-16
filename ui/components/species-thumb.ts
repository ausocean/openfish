import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import type { Species } from '@openfish/client'
import './image-carousel'
import './tooltip'

@customElement('species-thumb')
export class SpeciesThumb extends TailwindElement {
  @property({ type: Object })
  species: Species

  @state()
  private _idx = 0

  render() {
    const carousel = () => html`
    <image-carousel
        class="aspect-[4/3] w-full block rounded-t-sm overflow-clip"
        @update=${(e: CustomEvent<number>) => {
          this._idx = e.detail
        }}
    >
        ${repeat(
          this.species.images,
          (image) => html`
            <img src=${image.src} class="object-cover" />
        `
        )}
    </image-carousel>
    `
    const singleImage = () => html`
    <img
        src=${this.species.images[0].src}
        class="aspect-[4/3] w-full block rounded-t-sm overflow-clip object-cover"
    />
    `

    return html`
      <li
        class="card p-0 relative cursor-pointer h-full"
      >
        ${this.species.images.length > 1 ? carousel() : singleImage()}

        <div
            class="aspect-square rounded-full w-5 flex items-center justify-center text-sm bg-slate-950/75 text-white absolute top-2 right-2"
            id="attribution-icon-${this.species.id}"
        >
            &copy;
        </div>
        <tooltip-elem
            for="attribution-icon-${this.species.id}"
            trigger="hover"
            placement="bottom"
            class="z-10 max-w-256 text-xs"
        >
            ${this.species.images[this._idx].attribution}
        </tooltip-elem>

        <div class="px-2 font-bold mt-1 text-sm text-blue-900">${this.species.common_name}</div>
        <div class="px-2 pb-2 text-xs text-blue-900">${this.species.scientific_name}</div>
    </li>
    `
  }

  static styles = TailwindElement.styles
}

declare global {
  interface HTMLElementTagNameMap {
    'species-thumb': SpeciesThumb
  }
}
