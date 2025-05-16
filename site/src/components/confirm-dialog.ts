import { TailwindElement } from '@openfish/ui/components/tailwind-element'
import { html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { type Ref, createRef, ref } from 'lit/directives/ref.js'

@customElement('confirm-dialog')
export class ConfirmDialog extends TailwindElement {
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
      <dialog ${ref(this.dialogRef)} class="dialog">
        <slot></slot>

        <footer class="flex mt-2 justify-end gap-2 w-full">
          <button class="btn variant-orange" @click=${this.confirm}>${this.confirmMessage}</button>
          <button class="btn variant-slate" @click=${this.cancel}>${this.cancelMessage}</button>
        </footer>
      </dialog>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'confirm-dialog': ConfirmDialog
  }
}
