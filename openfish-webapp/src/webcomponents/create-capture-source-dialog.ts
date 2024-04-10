import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement } from 'lit/decorators.js'
import { type Ref, createRef, ref } from 'lit/directives/ref.js'
import resetcss from '../styles/reset.css?raw'
import btncss from '../styles/buttons.css?raw'
import './location-picker.ts'
import type { CaptureSource } from '../utils/api.types.ts'

@customElement('create-capture-source-dialog')
export class CreateCaptureSourceDialog extends LitElement {
  dialogRef: Ref<HTMLDialogElement> = createRef()

  show() {
    this.dialogRef.value?.showModal()
  }

  async submit(e: SubmitEvent) {
    e.preventDefault()

    const formdata = new FormData(e.target as HTMLFormElement)
    const payload: Partial<CaptureSource> = {
      name: formdata.get('name')?.toString(),
      camera_hardware: formdata.get('camera_hardware')?.toString(),
      location: formdata.get('location')?.toString() as `${number},${number}`,
    }

    if (formdata.get('site_id')) {
      payload.site_id = Number(formdata.get('site_id'))
    }

    await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/capturesources`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    this.dialogRef.value?.close()
    this.dispatchEvent(new Event('createitem'))
  }

  cancel() {
    this.dialogRef.value?.close()
  }

  render() {
    return html`
      <dialog ${ref(this.dialogRef)}>
      <h3>Create a new capture source</h3>
      <form @submit=${this.submit}>
        <section>
          <label for="name">Name</label>
          <input type="text" id="name" name="name" placeholder="Enter name of the capture source" required />
        </section>

        <section>
          <label for="camera_hardware">Camera Hardware</label>
          <input type="text" id="camera_hardware" name="camera_hardware" placeholder="Enter description of camera hardware" required />
        </section>

        <section>
          <label for="site_id">Site ID</label>
          <input type="text" id="site_id" name="site_id" placeholder="Enter site ID (optional)"  />
        </section>

        <section>
          <label>Location</label>
          <location-picker id="location" name="location" ></location-picker>
        </section>


      <footer>
        <input class="btn-orange btn-sm" type="submit" value="Create Capture Source" />
        
        <button class="btn-secondary btn-sm" @click=${this.cancel}>Cancel</button>
      </footer>
      </form>
      </dialog>
    `
  }

  static styles = css`
  ${unsafeCSS(resetcss)}
  ${unsafeCSS(btncss)}

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
    gap: 1.5rem;
  }

  input, label {
    display: block;
    width: 100%;
  }


  label {
    margin-top: -0.5rem;
    margin-bottom: 0.25rem;
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
    'create-capturesource-dialog': CreateCaptureSourceDialog
  }
}
