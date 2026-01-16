# OpenFish API documentation

OpenFish provides an API to access stored marine footage, and video annotations / labels, allowing clients to retrieve and filter the data.

Clients can download segments of footage or video annotations by querying
by location, time, and other parameters.

## Authentication
OpenFish has optional support for requiring user authentication.

User authentication is provided using Google Cloud's Identity Aware Proxy (IAP) or JSON Web Tokens (JWT). By default both are disabled.
Once authenticated, the user must create an account using the POST /api/v1/auth/me endpoint if they do not have one already.

### Using Identity Aware Proxy (IAP)
To enable IAP, pass the command line flag `--iap` or set the environment variable `IAP="true"`. You must also provide the audience using the command line flag `--jwt-audience` or the environment variable `JWT_AUDIENCE`.

### Using JWT Authentication
To enable JWT authentication, pass the command line flag `--jwt` or set the environment variable `JWT="true"`. You must provide an audience using the command line flag `--jwt-audience` or the environment variable `JWT_AUDIENCE`. This can be any value that isn't being used by any of your other services, for example "openfish". You must provide an issuer for the JWT tokens using the command line flag `--jwt-issuer` or the environment variable `JWT_ISSUER`. This can be some sort of identifer for the service issuing the JWT such as a service account email address. You must provide a key for the JWT in a secrets file in the format `jwtSecret:{64 hexadecimal characters}` and provide the location in the environment variable `OPENFISH_SECRETS`. The key can be generated using the following openssl command `openssl rand -hex 32` or your preferred way of generating 32 byte long keys. When issuing a JWT it must have an audience and issuer matching the provided values, an expiry date, and a subject which should be the email address of the user.

## Roles and permissions
If user authentication is enabled, the following roles and permissions apply:

| Role                | Permissions                                                                       |
| ------------------- | --------------------------------------------------------------------------------- |
| Admin               | Can add and remove annotations, videostreams, capturesources, users, and species. |
| Curator             | Can select streams for classification.                                            |
| Annotator (default) | Can add annotations, and delete their own annotations                             |
| Readonly            | A readonly user is only be able to look at annotations, not make any              |
