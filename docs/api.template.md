# APIs

All APIs accept body in JSON format

Headers

```
Content-Type: application/json
```

**Common Errors**

- Unknown Error

  Response code: 401

  ```json
  {
    "code": 1,
    "message": "Unknown Error"
  }
  ```

## Get ...

Endpoint: POST `/get_...`

### Headers

```
Content-Type: application/json
```

### Body Parameters

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
|          |      |          |         |             |

Example

```json
{}
```

### Responses

- Success

  Response Code: 200

  **Body (Data)**

  | Property | Type | Required | Default | Description |
  | -------- | ---- | -------- | ------- | ----------- |
  |          |      |          |         |             |

  Example response body (JSON)

  ```json
  {
    "code": 0,
    "message": "",
    "data": {}
  }
  ```

- Error

  Response Code: 400

  ```json
  {
    "code": 2,
    "message": ""
  }
  ```
