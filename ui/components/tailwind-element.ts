// See Also: https://github.com/butopen/web-components-tailwind-starter-kit

import { LitElement, unsafeCSS } from 'lit'

import globalStyles from '../theme.css?inline'

export class TailwindElement extends LitElement {
  static styles: typeof LitElement.styles = unsafeCSS(globalStyles)
}
