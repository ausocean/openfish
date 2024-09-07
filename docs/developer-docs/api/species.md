# Species
**Authors:** Scott Barnard

Species are used for providing suggestions to our users when annotating videos. They store the latin and common name, and an array of images. Images have a source and attribution - we use this to give the author credit and to abide by the rules of the license.

OpenFish provides APIs to create, get and delete species. The `/api/v1/species/recommended` API returns the most relevant species first, depending on the supplied video stream, and capture source.


## Fetching the most relevant species
So OpenFish can provide a list of species most relevant, you should provide the video stream and capture source for context. These are optional, if either is omitted, they will not be used to determine the most relevant species, if both are omitted then the list will be unsorted.

::: code-group
```http [Request]
GET /api/v1/species/recommended?videostream=12345&capturesource=67890
```

```http [Response]
HTTP/1.1 200
content-type: application/json

{
  "results": [
    {
      "id": 5701666712059904,
      "species": "Sepioteuthis australis",
      "common_name": "Southern Reef Squid",
      "images": [
        {
          "src": "https://inaturalist-open-data.s3.amazonaws.com/photos/340064435/medium.jpg",
          "attribution": "Tiffany Kosch, CC BY-NC-SA 4.0"
        }
      ]
    }
  ],
  "offset": 0,
  "limit": 20,
  "total": 1
}
```
:::

## Fetching a single species
::: code-group
```http [Request]
GET /api/v1/species/<species ID>
```

```http [Response]
HTTP/1.1 200

{
  "id": <species ID>,
  "species": "Sepioteuthis australis",
  "common_name": "Southern Reef Squid",
  "images": [
    {
      "src": "https://inaturalist-open-data.s3.amazonaws.com/photos/340064435/medium.jpg",
      "attribution": "Tiffany Kosch, CC BY-NC-SA 4.0"
    }
  ]
}
```
:::


## Creating species
Species have a latin and common name, and a list of images. All fields are mandatory.
::: code-group
```http [Request]
POST /api/v1/species
content-type: application/json

{
    "species": "Sepioteuthis australis",
    "common_name": "Southern Reef Squid",
    "images": [
        {
        "src": "https://inaturalist-open-data.s3.amazonaws.com/photos/340064435/medium.jpg",
        "attribution": "Tiffany Kosch, CC BY-NC-SA 4.0"
        }
    ]
}
```

```http [Response]
HTTP/1.1 200

{
  "id": <species ID>
}
```
:::

## Deleting a species
Successful delete will return 200 OK.

::: code-group
```http [Request]
DELETE /api/v1/species/<species ID>
```

```http [Response]
HTTP/1.1 200
```
:::