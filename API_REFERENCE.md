# REST API Reference

## TOC

- [REST API Reference](#rest-api-reference)
  * [TOC](#toc)
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
  * [Change Password](#change-password)
    + [Request](#request-2)
      - [Endpoint](#endpoint-2)
      - [Request Header](#request-header-2)
      - [Request Parameter](#request-parameter-2)
    + [Response](#response-2)
      - [Response if change password success](#response-if-change-password-success)
      - [Response if change password failed (no request body, old password incorrect)](#response-if-change-password-failed-no-request-body-old-password-incorrect)
      
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

## Change Password

### Request

#### Endpoint

```
POST /user/change_password
```

#### Request Header

```
Authorization: your_jwt_token_obtained_from_login
Content-type: application/json
```

#### Request Parameter

Parameter | Type | Required? | Remark
----------|------|-----------|--------
`old_password` | `string` | yes | -
`new_password` | `string` | yes | -

### Response

Response type : `Content-type: application/json`

#### Response if change password success

HTTP Status Code 200

Parameter | Type | Remark
----------|------|--------
`message` | `string` | -

#### Response if change password failed (no request body, old password incorrect)

HTTP Status Code 400, 500

Parameter | Type | Remark
----------|------|--------
`error_message` | `string` | -
`error_code` | `int` | -

## Change Profile

### Request

#### Endpoint

```
POST /user/change_profile
```

#### Request Header

```
Authorization: your_jwt_token_obtained_from_login
Content-type: application/json
```

#### Request Parameter

Parameter | Type | Required? | Remark
----------|------|-----------|--------
`username` | `string` | yes | -
`full_name` | `string` | yes | -
`gender` | `string` | yes | -

### Response

Response type : `Content-type: application/json`

#### Response if change password success

HTTP Status Code 200

Parameter | Type | Remark
----------|------|--------
`message` | `string` | -

#### Response if change password failed (no request body, and some server error)

HTTP Status Code 400, 500

Parameter | Type | Remark
----------|------|--------
`error_message` | `string` | -
`error_code` | `int` | -