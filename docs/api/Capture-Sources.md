Capture sources are cameras that produces video streams, and each has a name, location, camera hardware information and an optional unique identifier - site ID. OpenFish provides APIs to create, retrieve, update and delete capture sources, and features to query captures sources by their name and location.

**Page Contents**
- [Fetching a single capture source](#fetching-a-single-capture-source)
- [Querying capture sources](#querying-capture-sources)
- [Creating capture sources](#creating-capture-sources)
- [Updating a capture source](#updating-a-capture-source)
- [Deleting a capture source](#deleting-a-capture-source)

---

### Fetching a single capture source
**Request**
```
GET /api/v1/capturesources/<capture source ID>
```
**Response (JSON, HTTP 200)**
```js
{
    "id": /* capture source ID */,
    "name": "Stony Point Cuttle Cam",
    "location": "-32.12345,139.12345",
    "camera_hardware": "pi cam v2 (wide angle lens)",
    "site_id": 246813579,
}
```

### Querying capture sources
This will filter only those capture sources that have name=`Stony Point Cuttle Cam` and location=`-32.12345,139.12345`

**Request**
```
GET /api/v1/capturesources?name=Stony Point Cuttle Cam&location=-32.12345,139.12345
```


### Creating capture sources
**Request**
```
POST /api/v1/capturesources HTTP/1.1
Content-Type: application/json

{
  "name": "Camera 1",
  "location": "-37.12345678,140.12345678",
  "camera_hardware": "rpi cam v2 - wide angle lens",
  "site_id": 10192840284
}
```

**Response (JSON, HTTP 200)**
```js
{
  "id": /* ID of newly created capture source */
}
```

### Updating a capture source
Only include the data you wish to update. Successful update operation will return 200 OK.

**Request**
```
PATCH /api/v1/capturesources/<capture source ID>
Content-Type: application/json

{
  "name": <new name here>
}
```
**Response (HTTP 200)**

### Deleting a capture source 
Successful delete will return 200 OK.

**Request**
```
DELETE /api/v1/capturesources/<capture source ID>
```
**Response (HTTP 200)**