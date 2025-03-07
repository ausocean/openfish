export type Species = {
  id: number
  species: string
  common_name: string
  images?: Image[]
}

export type Image = { src: string; attribution: string }
