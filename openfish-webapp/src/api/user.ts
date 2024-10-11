export class User {
  email: string
  role: string

  get initial(): string {
    return this.email[0]
  }
}
