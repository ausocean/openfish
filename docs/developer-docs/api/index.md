# OpenFish API documentation

OpenFish provides an API to access stored marine footage, and video annotations / labels, allowing clients to retrieve and filter the data.

Clients can download segments of footage or video annotations by querying
by location, time, and other parameters.

## Authentication
OpenFish has optional support for requiring user authentication.

User authentication is provided using Google Cloud's Identity Aware Proxy (IAP). By default it is disabled, to use it you need to pass the command line flag `--iap` or set the environmental variable `IAP=\"true\"` to enable it.

## Roles and permissions
If user authentication is enabled, the following roles and permissions apply:

| Role               | Permissions                                                                       |
| ------------------ | --------------------------------------------------------------------------------- |
| Admin              | Can add and remove annotations, videostreams, capturesources, users, and species. |
| Curator            | Can select streams for classification.                                            |
| Annotator          | Can add annotations, and delete their own annotations                             |
| Readonly (default) | A readonly user is only be able to look at annotations, not make any              |
