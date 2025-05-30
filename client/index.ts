import createClientFromSchema from 'openapi-fetch'
import type { components, paths } from './schema' // generated by openapi-typescript

export const createClient = createClientFromSchema<paths>

export type OpenfishClient = ReturnType<typeof createClient>

export type PaginatedPath =
  | '/api/v1/videostreams'
  | '/api/v1/annotations'
  | '/api/v1/capturesources'
  | '/api/v1/species'
  | '/api/v1/users'

export type Species = components['schemas']['services.Species']
export type AnnotationWithJoins = components['schemas']['services.AnnotationWithJoins']
export type VideoStreamWithJoins = components['schemas']['services.VideoStreamWithJoins']
export type User = components['schemas']['services.User']
export type PublicUser = components['schemas']['services.PublicUser']
export type CaptureSource = components['schemas']['services.CaptureSource']
export type Identification = components['schemas']['services.Identification']
export type SpeciesSummary = components['schemas']['services.Identification']['species']
