# Capture Sources
**Authors:** Scott Barnard

Capture sources are cameras that produces video streams. Each capture source has a name, location, camera hardware information and an optional unique identifier - site ID. OpenFish provides APIs to create, retrieve, update and delete capture sources, and features to query captures sources by their name and location.


## Fetching a single capture source
::: code-group
```http [Request]
GET /api/v1/capturesources/<capture source ID>
```

```http [Response]
HTTP/1.1 200
content-type: application/json

{
    "id": <capture source ID>,
    "name": "Stony Point Cuttle Cam",
    "location": "-32.12345,139.12345",
    "camera_hardware": "pi cam v2 (wide angle lens)",
    "site_id": 246813579,
}
```
:::

## Querying capture sources
This will filter only those capture sources that have name=`Stony Point Cuttle Cam` and location=`-32.12345,139.12345`

::: code-group
```http [Request]
GET /api/v1/capturesources?name=Stony Point Cuttle Cam&location=-32.12345,139.12345
```

```http [Response]
HTTP/1.1 200
content-type: application/json

{
  "results": [
    {
      "id": 123456789,
      "name": "Stony Point Cuttle Cam",
      "location": "-32.12345,139.12345",
      "camera_hardware": "pi cam v2 (wide angle lens)",
      "site_id": 246813579,
    }
  ],
  "offset": 0,
  "limit": 20,
  "total": 1
}
```
:::


## Creating capture sources
::: code-group
```http [Request]
POST /api/v1/capturesources HTTP/1.1
content-type: application/json

{
  "name": "Camera 1",
  "location": "-37.12345678,140.12345678",
  "camera_hardware": "rpi cam v2 - wide angle lens",
  "site_id": 10192840284
}
```

```http [Response]
HTTP/1.1 200
content-type: application/json

{
  "id": <ID of newly created capture source>
}
```
:::


## Updating a capture source
Only include the data you wish to update. Successful update operation will return 200 OK.
::: code-group
```http [Request]
PATCH /api/v1/capturesources/<capture source ID>
content-type: application/json

{
  "name": <new name here>
}
```

```http [Response]
HTTP/1.1 200
```
:::


## Deleting a capture source 
Successful delete will return 200 OK.


::: code-group
```http [Request]
DELETE /api/v1/capturesources/<capture source ID>
```

```http [Response]
HTTP/1.1 200
```
:::
