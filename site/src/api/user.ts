export class User {
  email: string
  role: string
  display_name: string
  id: number

  get initial(): string {
    return this.email[0]
  }
}
