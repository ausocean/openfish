import { TailwindElement } from '@openfish/ui/components/tailwind-element'
import { css, html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { unsafeSVG } from 'lit/directives/unsafe-svg.js'
import caretLeft from '../icons/caret-left.svg?raw'
import caretRight from '../icons/caret-right.svg?raw'

@customElement('image-carousel')
export class ImageCarousel extends TailwindElement {
  private _idx = 0

  goto(direction: 1 | -1) {
    // Get slotted elements.
    const slot = this.shadowRoot?.querySelector('slot')!
    const slottedElems = slot.assignedElements({ flatten: true })

    // Scroll next or previous into view.
    this._idx = (this._idx + direction + slottedElems.length) % slottedElems.length
    slottedElems[this._idx].scrollIntoView({ behavior: 'smooth', block: 'nearest' })

    this.dispatchEvent(new CustomEvent('update', { detail: this._idx }))
  }

  render() {
    return html`
        <div class="relative group w-full h-full">
            <div class="absolute inset-0 flex w-full overflow-clip  overflow-x-scroll snap-x snap-mandatory snap-scroll-smooth">
                <slot></slot>
            </div>

            <button class="absolute left-2 bottom-2 rounded-full btn variant-transparent p-0 group-hover:opacity-100 opacity-0 transition-all text-slate-50 aspect-square h-8 w-8 *:w-4 *:h-4 *:fill-current" @click=${() => this.goto(-1)}>
                ${unsafeSVG(caretLeft)}
            </button>
            <button class="absolute right-2 bottom-2 rounded-full btn variant-transparent p-0 group-hover:opacity-100 opacity-0 transition-all text-slate-50 aspect-square h-8 w-8 *:w-4 *:h-4 *:fill-current" @click=${() => this.goto(+1)}>
                ${unsafeSVG(caretRight)}
            </button>
        </div>
    `
  }

  static styles = [
    TailwindElement.styles!,
    css`
        ::slotted(*) {
            scroll-snap-align: center;
        }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'image-carousel': ImageCarousel
  }
}
