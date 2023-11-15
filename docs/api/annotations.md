# Annotations
**Authors:** Scott Barnard

Annotations are used for labeling interesting things in videos. They store a linked video stream, a bounding box, a start and end time, the observer's name, and the observations themselves. Observations have a flexible format, they use key-value pairs so you can add all sorts of different information. Most commonly used is `species=<species name>`.

OpenFish provides APIs to create, retrieve, update (under development) and delete (under development) annotations, and features to query annotations by the person who made the observation, what kind of observations were made (presence of a key), what was observed (presence of a key and given value), and by the location (under development), video stream or capture source (under development).


## Fetching a single annotation
::: code-group
```http [Request]
GET /api/v1/annotations/<annotation ID>
```

```http [Response]
HTTP/1.1 200

{
  "id": <annotation ID>,
  "videostreamId": 4586454965551104,
  "timespan": {
    "Start": "0001-01-01T00:00:00Z",
    "End": "0001-01-01T00:00:00Z"
  },
  "observer": "scott@ausocean.org",
  "observation": {
    "common_name": "Zebrafish",
    "species": "Girella Zebra"
  }
}
```
:::

## Querying annotations
Annotations can be queried by the presence of an observation key `observation[<key>]=*` or a key-value pair `observation[<key>]=<value>`
Annotations can be filtered to show only those made by a given user. They can also be filtered to return only those for a given video stream or for a given location (Work in progress).

**Supported requests**
```http
GET /api/v1/annotations?observation[common_name]=Giant Cuttlefish&observation[behaviour]=*
GET /api/v1/annotations?observer=scott@ausocean.org
GET /api/v1/annotations?videostream=<video stream id>
GET /api/v1/annotations?capturesource=<capture source id>
GET /api/v1/annotations?location=-37.12345678,140.12345678
```

## Creating annotations
Annotations have a flexible format which gives you options when creating one. An annotation must have a video stream ID, timespan and observer, but a bounding box is optional, and the observation can be any amount of key-value pairs.

::: code-group
```http [Request]
POST /api/v1/annotations
Content-Type: application/json

{
  "videostreamId": 4586454965551104,
  "timespan": { "start_time": "2023-06-07T16:30:00.00", "end_time": "2023-06-07T16:30:02.00" },
  "boundingBox": {"x1": 84, "y1": 160, "x2": 205, "y2": 295},
  "observer": "scott@ausocean.com",
  "observation": {
    "common_name": "Giant Cuttlefish",
    "species": "Sepia Apama",
    "behaviour": "Eating"
  }
}
```

```http [Response]
HTTP/1.1 200

{
  "id": <annotation ID>
}
```
:::

## Updating an annotation
::: warning
This feature is under development
:::

## Deleting an annotation
::: warning
This feature is under development
:::
