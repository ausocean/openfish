# Users
**Authors:** Scott Barnard

A user is identified by their email and has a role that gives them permissions. A user is created
when they first login to OpenFish. There are APIs for updating user's role, listing users and
deleting a user account.


## Fetching a user's details.
::: code-group
```http [Request]
GET /api/v1/users/<user email address>
```

```http [Response]
HTTP/1.1 200
content-type: application/json

{
  "email": "user@example.com",
  "role": "annotator"
}
```
:::

## Listing users
This will show a list of all users. This is restricted to admins only.

::: code-group
```http [Request]
GET /api/v1/users
```

```http [Response]
HTTP/1.1 200
content-type: application/json

{
  "results": [
    {
      "email": "user@example.com",
      "role": "annotator"
    }
  ],
  "offset": 0,
  "limit": 20,
  "total": 1
}
```
:::


## Updating a user's role
User role can be updated using PATCH. Successful update operation will return 200 OK.
::: code-group
```http [Request]
PATCH /api/v1/users/<user email address>
content-type: application/json

{
  "role": "admin"
}
```

```http [Response]
HTTP/1.1 200
```
:::


## Deleting a user
Successful delete will return 200 OK.


::: code-group
```http [Request]
DELETE /api/v1/users/<user email address>
```

```http [Response]
HTTP/1.1 200
```
:::
