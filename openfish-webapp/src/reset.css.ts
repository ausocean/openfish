import { css } from 'lit'

// Lit creates custom elements / web components that use the shadow DOM.
// Shadow DOM elements are isolated and do not inherit all styles from
// the page. This fragment contains default styles that each component
// should have, but is refusing to inherit from the page because of its
// isolation.
export const resetcss = css`
  *, *:before, *:after {
    box-sizing: border-box;
  }
  `

export const buttonStyles = css`
.btn {
  width: min-content;
  height: 2.5rem;
  border-radius: 999px;
  font-size: 1rem;
  padding: 0 1rem;
  white-space: nowrap;
  border: 1px solid;
  cursor: pointer;

  &:disabled {
    opacity: 0.5;
  }
}

.btn-sm {
  width: min-content;
  height: 1.5rem;
  border-radius: 999px;
  font-size: 0.8rem;
  padding: 0 1rem;
  white-space: nowrap;
  border: none;
  border: 1px solid;
  cursor: pointer;

  &:disabled {
    opacity: 0.5;
  }
}

.btn-orange {    
  background-color: var(--orange-400);
  color: var(--orange-800);
  border-color: var(--orange-400);

  &:hover:enabled {
    background-color: var(--orange-500);
    border-color: var(--orange-500);
  }
}

.btn-secondary {
  background-color: var(--gray-200);
  color: var(--gray-900);
  border-color: var(--gray-200);

  &:hover:enabled {
    background-color: var(--gray-300);
    color: var(--gray-950);
    border-color: var(--gray-300);
  }
}

.btn-blue {
  background-color: var(--blue-800);
  color: var(--gray-100);
  border-color: var(--blue-800);

  &:hover:enabled {
    background-color: var(--blue-900);
    border-color: var(--blue-900);
  }

}

.btn-outline {
  background-color: transparent;
  color: var(--gray-50);
  border-color: currentColor;

  &:hover:enabled {
    color: var(--gray-300);
  }
}

.btn-transparent {
  background-color: rgba(0, 0, 0, 0.75);
  color: var(--gray-100);
  border-color: rgba(0, 0, 0, 0.75);

  &:hover:enabled {
    background-color: rgba(0, 0, 0, 0.9);
    border-color: rgba(0, 0, 0, 0.9);
  }
}
`
