export function extractVideoID(url: string | undefined | null): string | null {
  if (url) {
    const parsed = new URL(url)
    return parsed.searchParams.get('v')
  }
  return null
}
