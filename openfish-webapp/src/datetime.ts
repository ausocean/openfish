// Convert a video time in seconds to a datetime.
export function videotimeToDatetime(streamStart: string, time: number): Date {
  let playbackDatetime = new Date(streamStart)
  playbackDatetime.setSeconds(playbackDatetime.getSeconds() + time)
  return playbackDatetime
}

// Convert a datetime to a video time in seconds.
export function datetimeToVideoTime(streamStart: string, datetime: string): number {
  return (new Date(datetime).getTime() - new Date(streamStart).getTime()) / 1000
}

// Subtracts datetimes and returns difference as a number in seconds.
export function datetimeDifference(a: string, b: string): number {
  return (new Date(a).getTime() - new Date(b).getTime()) / 1000
}
