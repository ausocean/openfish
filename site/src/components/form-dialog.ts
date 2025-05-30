import { TailwindElement } from '@openfish/ui/components/tailwind-element'
import { html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { type Ref, createRef, ref } from 'lit/directives/ref.js'

type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

@customElement('form-dialog')
export class FormDialog extends TailwindElement {
  dialogRef: Ref<HTMLDialogElement> = createRef()

  @property()
  accessor action = ''

  @property()
  accessor method: HTTPMethod = 'POST'

  @property()
  accessor title = ''

  @property()
  accessor btntext = ''

  show() {
    this.dialogRef.value?.showModal()
  }

  set(data: Record<string, any>) {
    for (const key in data) {
      const element = this.shadowRoot!.querySelector(`[name="${key}"]`) as
        | (HTMLElement & { value: string })
        | null
      if (element !== null) {
        element.value = data[key] ?? ''
      }
    }
  }

  // NOTE: Slotted content does not live in the shadow DOM, so won't appear in the
  // FormData. Here we move the elements to the shadow DOM, whenever they change.
  // See: https://stackoverflow.com/questions/53676756
  onSlotChange(e: Event & { target: HTMLSlotElement }) {
    const form = this.shadowRoot?.querySelector('form')!

    for (const field of e.target.assignedNodes()) {
      form.insertBefore(field, e.target)
    }
  }

  async submit(e: SubmitEvent & { target: HTMLFormElement }) {
    e.preventDefault()

    const formdata = new FormData(e.target)
    const payload: Record<string, any> = {}
    for (const [key, val] of formdata.entries()) {
      if (key === 'capturesource' || key === 'videostream') {
        payload[key] = Number(val)
      } else if (val !== '' && val !== undefined && val !== null) {
        payload[key] = val
      }
    }

    e.target.reset()

    await fetch(this.action, {
      method: this.method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    this.dialogRef.value?.close()
    this.dispatchEvent(new Event('formsubmit'))
  }

  cancel(e: PointerEvent & { target: HTMLButtonElement }) {
    e.target.form?.reset()
    this.dialogRef.value?.close()
  }

  render() {
    return html`
      <dialog ${ref(this.dialogRef)} class="dialog min-w-120 text-slate-800">
        <h3 class="text-lg font-bold">${this.title}</h3>
        <form @submit=${this.submit} class="mt-4 flex flex-col gap-2">
          <slot @slotchange=${this.onSlotChange}></slot>
          <footer class="flex mt-2 justify-end gap-2 w-full">
            <input class="btn variant-orange" type="submit" value=${this.btntext} >
            <button class="btn variant-slate" @click=${this.cancel}>Cancel</button>
          </footer>
        </form>
      </dialog>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'form-dialog': FormDialog
  }
}
