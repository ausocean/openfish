// Convert a video time in seconds to a datetime.
export function videotimeToDatetime(streamStart: string, time: number): Date {
  const playbackDatetime = new Date(streamStart)
  playbackDatetime.setSeconds(playbackDatetime.getSeconds() + time)
  return playbackDatetime
}

// Convert a datetime to a video time in seconds.
export function datetimeToVideoTime(streamStart: DateLike, datetime: DateLike): number {
  return (new Date(datetime).getTime() - new Date(streamStart).getTime()) / 1000
}

// Subtracts datetimes and returns difference as a number in seconds.
export function datetimeDifference(a: DateLike, b: DateLike): number {
  return (new Date(a).getTime() - new Date(b).getTime()) / 1000
}

// A Date or ISO string representation of a date.
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
  return (
    new Intl.DateTimeFormat('en-AU', {
      weekday: 'short',
      year: 'numeric',
      month: 'short',
      day: '2-digit',
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
    }).format(new Date(dt)) + ` [${formatAsTimeZone(new Date(dt))}]`
  )
}

export function formatAsTimeZone(dt: DateLike): string {
  return new Intl.DateTimeFormat('en-AU', { day: '2-digit', timeZoneName: 'short' })
    .format(new Date(dt))
    .slice(4)
}

export function formatDuration(duration: number): string {
  const h = Math.floor(duration / 3600)
  const m = Math.floor((duration - h * 3600) / 60)
  const s = duration - h * 3600 - m * 60
  return `${h}h ${m}m ${s}s `
}
