# API Usage
**Authors:** Scott Barnard

OpenFish's API uses URLs with a resource type and resource ID for fetching single resources:

#### Request
```http
GET /api/v1/<resource>/<resource ID>
```
For a list of resources, use the following where offset is the number of items to skip and limit is the total number of items to fetch:

#### Request
```http
GET /api/v1/<resource>?offset=0&limit=20
```
#### Response (JSON, HTTP 200)
```json
{
  "results": [
    /* result 1 */
    /* result 2 */
    /* result 3 */
  ],
  "offset": 0,
  "limit": 20,
  "total": 3
}
```

To select only the data you need, use format to specify what keys from the JSON you need:

#### Request
```http
GET /api/v1/<resource>/<resource ID>?format=key 1, key 2
```
**or**
```http
GET /api/v1/<resource>?offset=0&limit=20&format=key 1, key 2
```

Many APIs can filter data - see specific examples.