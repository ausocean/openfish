// Subtracts datetimes and returns difference as a number in seconds.
export function datetimeDifference(a: DateLike, b: DateLike): number {
  return (new Date(a).getTime() - new Date(b).getTime()) / 1000
}

// toDatetimeLocal converts a Date to a format <input type="datetime-local"> can accept as a value.
export function toDatetimeLocal(dt: DateLike): string {
  const date = new Date(dt)
  date.setMinutes(date.getMinutes() - date.getTimezoneOffset())
  return date.toISOString().slice(0, 16)
}

// A Date or RFC3339 string representation of a date.
export type DateLike = Date | string

export function formatAsDate(dt: DateLike): string {
  return new Intl.DateTimeFormat('en-AU', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: '2-digit',
  }).format(new Date(dt))
}

export function formatAsTime(dt: DateLike): string {
  return new Intl.DateTimeFormat('en-AU', {
    hour: 'numeric',
    minute: 'numeric',
    second: 'numeric',
  }).format(new Date(dt))
}

export function formatAsDatetime(dt: DateLike): string {
  return `${new Intl.DateTimeFormat('en-AU', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: 'numeric',
    minute: 'numeric',
    second: 'numeric',
  }).format(new Date(dt))} [${formatAsTimeZone(new Date(dt))}]`
}

export function formatAsTimeZone(dt: DateLike): string {
  return new Intl.DateTimeFormat('en-AU', {
    day: '2-digit',
    timeZoneName: 'short',
  })
    .format(new Date(dt))
    .slice(4)
}

export function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds - h * 3600) / 60)
  const s = seconds - h * 3600 - m * 60
  return `${h}h ${m}m ${s}s `
}

export function formatVideoTime(seconds: number, ms = false): string {
  const date = new Date(0)
  date.setMilliseconds(seconds * 1000)
  return date.toISOString().substring(11, ms ? 23 : 19)
}

export function parseVideoTime(str: string): number {
  const [h, m, s, ms] = str.split(/[:\.]/)

  return Number(h) * 60 * 60 + Number(m) * 60 + Number(s) + Number(ms) / 1000
}
