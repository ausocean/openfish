import { html } from 'lit'
import { customElement, property } from 'lit/decorators.js'

import { arrow, computePosition, offset, shift } from '@floating-ui/dom'
import { TailwindElement } from '@openfish/ui/components/tailwind-element'

@customElement('tooltip-elem')
export class Tooltip extends TailwindElement {
  @property({ type: String })
  accessor placement: 'top' | 'left' | 'bottom' | 'right'

  @property({ type: String })
  accessor for = ''

  @property({ type: String })
  accessor trigger: 'click' | 'hover' = 'hover'

  target: Element | null = null

  setTarget() {
    let showEvents: string[] = []
    let hideEvents: string[] = []

    if (this.trigger === 'hover') {
      showEvents = ['pointerenter', 'focus']
      hideEvents = ['pointerleave', 'blur', 'keydown', 'click']
    } else {
      showEvents = ['click']
      hideEvents = ['keydown', 'blur']
    }

    const newTarget = this.parentElement?.querySelector(`#${this.for}`) ?? null
    for (const name of showEvents) {
      this.target?.removeEventListener(name, this.show)
      newTarget?.addEventListener(name, this.show)
    }
    for (const name of hideEvents) {
      this.target?.removeEventListener(name, this.hide)
      newTarget?.addEventListener(name, this.hide)
    }
    this.target = newTarget
  }

  connectedCallback() {
    super.connectedCallback()
    this.hide()
    this.setTarget()
  }

  render() {
    return html`
      <div
        id="arrow"
        class="w-3 h-3 absolute bg-inherit z-0"
        style="clip-path: polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%);"
      ></div>
      <div class="overflow-clip z-10 p-2 relative">
        <slot></slot>
      </div>
    `
  }

  show = async () => {
    const arrowEl = this.shadowRoot!.querySelector('#arrow')! as HTMLDivElement
    if (this.target !== null) {
      this.style.cssText = ''
      const { x, y, middlewareData } = await computePosition(this.target, this, {
        strategy: 'absolute',
        placement: this.placement,
        middleware: [offset(6), shift(), arrow({ element: arrowEl, padding: 6 })],
      })

      Object.assign(this.style, {
        left: `${x}px`,
        top: `${y}px`,
      })

      const { x: ax, y: ay } = middlewareData!.arrow!

      const staticSide = {
        top: 'bottom',
        right: 'left',
        bottom: 'top',
        left: 'right',
      }[this.placement]
      Object.assign(arrowEl.style, {
        left: ax != null ? `${ax}px` : '',
        top: ay != null ? `${ay}px` : '',
        right: '',
        bottom: '',
        [staticSide]: '-6px',
      })
    }
  }

  hide = () => {
    this.style.display = 'none'
  }

  static styles = [TailwindElement.styles!]
}
