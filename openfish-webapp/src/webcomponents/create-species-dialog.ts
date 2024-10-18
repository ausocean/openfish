import { LitElement, css, html, unsafeCSS } from 'lit'
import { customElement } from 'lit/decorators.js'
import { type Ref, createRef, ref } from 'lit/directives/ref.js'
import resetcss from '../styles/reset.css?raw'
import btncss from '../styles/buttons.css?raw'
import './location-picker.ts'
import type { Species } from '../utils/api.types.ts'

@customElement('create-species-dialog')
export class CreateSpeciesDialog extends LitElement {
  dialogRef: Ref<HTMLDialogElement> = createRef()

  show() {
    this.dialogRef.value?.showModal()
  }

  async submit(e: SubmitEvent) {
    e.preventDefault()

    const formdata = new FormData(e.target as HTMLFormElement)
    const payload: Partial<Species> = {
      common_name: formdata.get('common_name')?.toString(),
      species: formdata.get('species')?.toString(),
    }

    await fetch(`${import.meta.env.VITE_API_HOST}/api/v1/species`, {
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
      <h3>Create a new species</h3>
      <form @submit=${this.submit}>
        <section>
          <label for="common_name">Common Name</label>
          <input type="text" id="common_name" name="common_name" placeholder="Enter common name of the species" required />
        </section>

        <section>
          <label for="species">Scientific Name</label>
          <input type="text" id="species" name="species" placeholder="Enter the scientific name of the species" required />
        </section>

        <!-- TODO: images -->
      <footer>
        <input class="btn-orange btn-sm" type="submit" value="Create Species" />
        
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
    'create-species-dialog': CreateSpeciesDialog
  }
}
