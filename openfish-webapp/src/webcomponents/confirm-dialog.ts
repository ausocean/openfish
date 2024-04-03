import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { type Ref, createRef, ref } from 'lit/directives/ref.js'
import resetcss from '../styles/reset.css?raw'
import btncss from '../styles/buttons.css?raw'

@customElement('confirm-dialog')
export class ConfirmDialog extends LitElement {
  @property({ type: String })
  confirmMessage = 'Yes'

  @property({ type: String })
  cancelMessage = 'No'

  callback: () => void = () => {}

  dialogRef: Ref<HTMLDialogElement> = createRef()

  show(callback: () => void) {
    this.dialogRef.value?.showModal()
    this.callback = callback
  }

  confirm() {
    this.dialogRef.value?.close()
    this.callback()
  }

  cancel() {
    this.dialogRef.value?.close()
  }

  render() {
    return html`
      <dialog ${ref(this.dialogRef)}>
      <slot></slot>
      <footer>
        <button class="btn-orange btn-sm" @click=${this.confirm}>${this.confirmMessage}</button>
        <button class="btn-secondary btn-sm" @click=${this.cancel}>${this.cancelMessage}</button>
      </footer>
      </dialog>
    `
  }

  static styles = css`
  ${unsafeCSS(resetcss)}
  ${unsafeCSS(btncss)}

  dialog {
    border: none;
    border-radius: 0.5rem;
  }

  ::backdrop {
    background-color: var(--gray-950);
    opacity: 0.5;
  }

  footer {
    margin-top: 0.5rem;
    display: flex;
    width: 100%;
    justify-content: flex-end;
    gap: 1rem;
  }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'confirm-dialog': ConfirmDialog
  }
}
