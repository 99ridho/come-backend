# REST API Reference

## TOC

- [REST API Reference](#rest-api-reference)
  * [Login](#login)
    + [Request](#request)
      - [Endpoint](#endpoint)
      - [Request Header](#request-header)
      - [Request Parameter](#request-parameter)
    + [Response](#response)
      - [Response if logged in successfully](#response-if-logged-in-successfully)
      - [Response if login failed (no request body, password incorrect, not registered, etc)](#response-if-login-failed-no-request-body-password-incorrect-not-registered-etc)
  * [Register](#register)
    + [Request](#request-1)
      - [Endpoint](#endpoint-1)
      - [Request Header](#request-header-1)
      - [Request Parameter](#request-parameter-1)
    + [Response](#response-1)
      - [Response if registered successfully](#response-if-registered-successfully)
      - [Response if register failed (no request body, account already registered)](#response-if-register-failed-no-request-body-account-already-registered)

## Login

### Request

#### Endpoint

```
POST /login
```

#### Request Header

```
Content-type: application/json
```

#### Request Parameter

Parameter | Type | Required? | Remark
----------|------|-----------|--------
`email` | `string` | yes | -
`password` | `string` | yes | -

### Response

Response type : `Content-type: application/json`

#### Response if logged in successfully

HTTP Status Code 200

Parameter | Type | Remark
----------|------|--------
`message` | `string` | -
`token` | `string` | Use this token at your `Authorization` header for requesting a resource from protected endpoint

#### Response if login failed (no request body, password incorrect, not registered, etc)

HTTP Status Code 400, 500

Parameter | Type | Remark
----------|------|--------
`error_message` | `string` | -
`error_code` | `int` | -

## Register

### Request

#### Endpoint

```
POST /register
```

#### Request Header

```
Content-type: application/json
```

#### Request Parameter

Parameter | Type | Required? | Remark
----------|------|-----------|--------
`username` | `string` | yes | -
`email` | `string` | yes | -
`password` | `string` | yes | -
`full_name` | `string` | yes | -
`gender` | `string` | yes | value is `male` or `female`
`fcm_token` | `string` | yes | FCM token that obtained after initializing firebase cloud messaging



### Response

Response type : `Content-type: application/json`

#### Response if registered successfully

HTTP Status Code 200

Parameter | Type | Remark
----------|------|--------
`message` | `string` | -

#### Response if register failed (no request body, account already registered)

HTTP Status Code 400, 500

Parameter | Type | Remark
----------|------|--------
`error_message` | `string` | -
`error_code` | `int` | -