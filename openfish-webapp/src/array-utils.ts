export function zip<A, B>(a: A[], b: B[]): [A, B][] {
  return a.map((_, i) => [a[i], b[i]])
}
