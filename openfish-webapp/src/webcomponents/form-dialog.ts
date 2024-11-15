import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { type Ref, createRef, ref } from 'lit/directives/ref.js'
import resetcss from '../styles/reset.css?lit'
import btncss from '../styles/buttons.css?lit'

type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

@customElement('form-dialog')
export class FormDialog extends LitElement {
  dialogRef: Ref<HTMLDialogElement> = createRef()

  @property()
  action = ''

  @property()
  method: HTTPMethod = 'POST'

  @property()
  title = ''

  @property()
  btntext = ''

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
      <dialog ${ref(this.dialogRef)}>
        <h3>${this.title}</h3>
        <form @submit=${this.submit}>
          <slot @slotchange=${this.onSlotChange}></slot>
          <footer>
            <input class="btn-orange btn-sm" type="submit" value=${this.btntext} >
            <button class="btn-secondary btn-sm" @click=${this.cancel}>Cancel</button>
          </footer>
        </form>
      </dialog>
    `
  }

  static styles = css`
  ${resetcss}
  ${btncss}

  dialog {
    border: none;
    border-radius: 0.5rem;
    min-width: 30rem;
  }

  ::backdrop {
    background-color: var(--gray-950);
    opacity: 0.5;
  }

  form {
    margin-top: 1rem;
    padding: 1.5rem 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  ::slotted(input) {
    display: block;
    width: 100%;
  }

  ::slotted(label) {
    display: block;
    width: 100%;
  }

  ::slotted(label:not(:first-child)) {
    margin-top: 1rem;
  }

  footer {
    margin-top: 1rem;
    display: flex;
    width: 100%;
    justify-content: flex-end;
    gap: 1rem;
  }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'form-dialog': FormDialog
  }
}
