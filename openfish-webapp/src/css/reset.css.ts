import { css } from 'lit'

// Lit creates custom elements / web components that use the shadow DOM.
// Shadow DOM elements are isolated and do not inherit all styles from
// the page. This fragment contains default styles that each component
// should have, but is refusing to inherit from the page because of its
// isolation.
export const resetcss = css`
  *, *:before, *:after {
    box-sizing: border-box;
  }`
