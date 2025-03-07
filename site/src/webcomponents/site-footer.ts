import { TailwindElement } from './tailwind-element'
import { html } from 'lit'
import { customElement } from 'lit/decorators.js'

@customElement('site-footer')
export class SiteFooter extends TailwindElement {
  render() {
    return html`
    <footer class="contain px-8 flex py-4 gap-2">
        <span class="flex-1">&copy; AusOcean 2025</span>
        <a href="https://ausocean.github.io/openfish/" class="link">Project information</a>
        <a href="https://ausocean.org" class="link">AusOcean.org</a>
        <a href="https://ausocean.tv" class="link">AusOcean.TV</a>
    </footer>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'site-footer': SiteFooter
  }
}
