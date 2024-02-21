# Video Streams
**Authors:** Scott Barnard

The video stream API provides the metadata for video streams. A video stream has a start time, end time, stream URL and linked capture source. 

Using the video stream API we can register our video streams with OpenFish so it can annotate and play back that stream. The stream URL specifies where the video data is stored. 

Examples:
- `http://vidgrind.ausocean.org/get?id=1`
- `https://www.youtube.com/watch?v=abcdefghijk`


## Fetching a single video stream
::: code-group
```http [Request]
GET /api/v1/videostreams/<video stream ID>
```

```http [Response]
HTTP/1.1 200

{
  "id": <video stream ID>,
  "startTime": "2023-06-07T08:00:00Z",
  "endTime": "2023-06-07T16:30:00Z",
  "stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
  "capturesource": 6636835711221760
}
```
:::

## Querying video streams
Video streams can be filtered by a start and an end time, and also by the capture source that produced it.

::: code-group
```http [Request]
GET /api/v1/videostreams?timespan[start]=2023-05-24T00:00:00Z&timespan[end]=2023-06-01T00:00:00Z&capturesource=123456
```

```http [Response]
HTTP/1.1 200
content-type: application/json

{
  "results": [
    {
      "id": 4586454965551104,
      "startTime": "2023-05-25T08:00:00Z",
      "endTime": "2023-05-25T16:30:00Z",
      "stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
      "capturesource": 123456
    },
    {
      "id": 1231104458645496,
      "startTime": "2023-05-25T08:00:00Z",
      "endTime": "2023-05-25T16:30:00Z",
      "stream_url": "https://www.youtube.com/watch?v=lmnopqrstuv",
      "capturesource": 123456
    }
  ],
  "offset": 0,
  "limit": 20,
  "total": 2
}
```
:::


::: tip Tip: Time spans explained
Requests will return all video streams that overlap with the specified time span, not just those contained within it.
In this scenario we are querying for video streams from 4pm til 7:30pm. This request will therefore return video streams 1, 2, 3 & 4 but not 5.
```
          |         :<-- Query  -->:
          |         :              :
Stream 1: |   •-----:-------•      :
Stream 2: |         :  •--------•  :
Stream 3: |         :       •------:-------•
Stream 4: |  •------:--------------:-------•
Stream 5: | •---•   :              :
          |         :              :
          '---------'---------'---------'---------'     (Time)
          2pm       4pm       6pm       8pm       10pm
```
:::

## Registering a video stream
::: code-group
```http [Request]
POST http://localhost:8080/api/v1/videostreams
Content-Type: application/json

{
  "stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
  "capturesource": 5661542255165440,
  "startTime": "2023-06-07T08:00:00.00Z",
  "endTime": "2023-06-07T16:30:00.00Z"
}
```

```http [Response]
{
  "id": <ID of newly created videostream>
}
```
:::

## Working with live streams
Live streams are different to registering an existing video. This is because we don't know the end time when we start it. The API has the `/api/v1/videostreams/live` endpoint for handling these scenarios.

To register a stream when it starts use POST. It takes the current time as the start time.

::: code-group
```http [Request]
POST /api/v1/videostreams/live
Content-Type: application/json

{
  "stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
  "capturesource": 5970800033136640
}
```

```http [Response]
HTTP/1.1 200

{
  "id": <ID of newly created videostream>
}
```
:::


To finish a stream use PATCH. It uses the current time as the end time.

::: code-group
```http [Request]
PATCH  /api/v1/videostreams/<video stream ID>/live
```

```http [Response]
HTTP/1.1 200
```
:::

## Updating a video stream
Updating a video stream is used if it was to move to a new URL.

::: code-group
```http [Request]
PATCH /api/v1/videostreams/<video stream ID>
Content-Type: application/json

{
  "stream_url": "http://newurl.com.au/hello"
}
```

```http [Response]
HTTP/1.1 200
```
:::


## Deleting a video stream
Successful delete will return 200 OK.

::: code-group
```http [Request]
DELETE /api/v1/videostreams/<video stream ID>
```

```http [Response]
HTTP/1.1 200
```
:::
