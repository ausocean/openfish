Video streams are captured videos and have a start time, end time, stream URL and linked capture source. 

A stream URL specifies where the video data is stored, so that clients can play back that video. 

Examples:
- `http://vidgrind.ausocean.org/get?id=1`
- `https://www.youtube.com/watch?v=abcdefghijk`


**Page Contents**
- [Fetching a single video stream](#fetching-a-single-video-stream)
- [Querying video streams](#querying-video-streams)
- [Creating a video stream](#creating-a-video-stream)
- [Live streams](#live-streams)
- [Updating a video stream](#updating-a-video-stream)
- [Deleting a video stream](#deleting-a-video-stream)

---


### Fetching a single video stream
**Request**
```
GET /api/v1/videostreams/<video stream ID>
```
**Response (JSON, HTTP 200)**
```js
{
  "id": 4586454965551104,
  "startTime": "2023-06-07T08:00:00Z",
  "endTime": "2023-06-07T16:30:00Z",
  "stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
  "capturesource": 6636835711221760
}
```

### Querying video streams
Video streams can be filtered by a start and an end time, and also by the capture source that produced it.

**Request**
```
GET /api/v1/videostreams?timespan[start]=2023-05-24T00:00:00Z&timespan[end]=2023-06-01T00:00:00Z&capturesource=<capture source ID>
```
**Time spans explained**

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


### Creating a video stream
**Request**
```
POST http://localhost:3000/api/v1/videostreams
Content-Type: application/json

{
  "stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
  "capturesource": 5661542255165440,
  "startTime": "2023-06-07T08:00:00.00Z",
  "endTime": "2023-06-07T16:30:00.00Z"
}
```
**Response (JSON, HTTP 200)**
```js
{
  "id": /* ID of newly created videostream*/
}
```

### Live streams
Live streams are different to uploading an existing video. This is because we don't know the end time when we start it. The API has the `/api/v1/videostreams/live` endpoint for handling these scenarios.

To start a stream use POST. It takes the current time as the start time.

**Request**
```
POST /api/v1/videostreams/live
Content-Type: application/json

{
  "stream_url": "https://www.youtube.com/watch?v=abcdefghijk",
  "capturesource": 5970800033136640
}
```
**Response (JSON, HTTP 200)**
```js
{
  "id": /* ID of newly created videostream*/
}
```

To finish a stream use PATCH. It uses the current time as the end time.

**Request**
```
PATCH  /api/v1/videostreams/<video stream id>/live
```
**Response (HTTP 200)**

### Updating a video stream
Updating a video stream is used if it was to move to a new URL.

**Request**
```
PATCH /api/v1/videostreams/<video stream ID>
Content-Type: application/json

{
  "stream_url": "http://newurl.com.au/hello"
}
```
**Response (HTTP 200)**


### Deleting a video stream
Successful delete will return 200 OK.

**Request**
```
DELETE /api/v1/videostreams/<video stream ID>
```