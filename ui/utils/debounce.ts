export function debounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  // Use ReturnType<typeof setTimeout> to handle both Node and Browser environments.
  let timer: ReturnType<typeof setTimeout>

  return function (this: any, ...args: Parameters<T>) {
    clearTimeout(timer)
    timer = setTimeout(() => {
      func.apply(this, args)
    }, delay)
  }
}
