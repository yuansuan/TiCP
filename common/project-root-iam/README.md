
# API DOC

## Add Policy

Adds a new policy to the system.

### Request

`POST /iam/v1/policies`

#### Request Body

The request body should be a JSON object with the following fields:

| Field      | Type     | Required | Description                                                                                     |
|------------|----------|----------|-------------------------------------------------------------------------------------------------|
| PolicyName | string   | Yes      | 权限策略名称.                                                           |
| Version    | string   | Yes      | 版本.                                                        |
| Effect     | string   | Yes      | 授权效果包括两种：允许（allow）和拒绝（deny）.                                 |
| Resources  | []string | Yes      | 资源是指被授权的具体对象.                                                       |
| Actions    | []string | Yes      | 操作是指对具体资源的操作.                                                   |

#### Example Request Body

```json
{
  "PolicyName": "example-policy",
  "Version": "1.0",
  "Effect": "allow",
  "Resources": [
    "yrn:ys:cs::4TiSxuPtJEm:path/4T4ZZvA2tVb/<.*>"
  ],
  "Actions": [
    "<.*>"
  ]
}
```

### Response

#### Success Response

If the policy is successfully created, the response will have a `200 OK` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is missing a required field or contains an invalid value. |
| 500 Internal Server Error | An error occurred while creating the policy. |

### Example

#### Request

```
POST /iam/v1/policies HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json

{
  "PolicyName": "example-policy",
  "Version": "1.0",
  "Effect": "allow",
  "Resources": [
    "yrn:ys:cs::4TiSxuPtJEm:path/4T4ZZvA2tVb/<.*>"
  ],
  "Actions": [
    "<.*>"
  ]
}
```

#### Response

```
HTTP/1.1 200 OK
```


## Get Policy

Retrieves the details of a policy.

### Request

`GET  /iam/v1/policies/{policyName}`

#### Path Parameters

| Parameter  | Type   | Required | Description         |
|------------|--------|----------|---------------------|
| policyName | string | Yes      | The name of the policy to retrieve. |

### Response

#### Success Response

If the policy is found, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field      | Type     | Description                                                                                     |
|------------|----------|-------------------------------------------------------------------------------------------------|
| PolicyName | string   | The name of the policy.                                                                         |
| Effect     | string   | 同上.                                 |
| Resources  | []string | 同上.                                                       |
| Actions    | []string | 同上.                                                   |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified policy does not exist. |
| 400 Bad Request | The request body is missing a required field or contains an invalid value. |
| 500 Internal Server Error | An error occurred while retrieving the policy. |

### Example

#### Request

```
GET  /iam/v1/policies/example-policy HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "PolicyName": "example-policy",
  "Effect": "Allow",
  "Resources": [
    "arn:aws:s3:::example-bucket/*"
  ],
  "Actions": [
    "s3:GetObject"
  ]
}
```



## List Policies

Retrieves a list of all policies.

### Request

`GET /iam/v1/policies`

### Response

#### Success Response

If the policies are found, the response will have a `200 OK` status code and a JSON object with the following fields:

同上

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 500 Internal Server Error | An error occurred while retrieving the policies. |

### Example

#### Request

```
GET /iam/v1/policies HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "Policies": [
    {
      "PolicyName": "example-policy-1",
      "Effect": "Allow",
      "Resources": [
        "arn:aws:s3:::example-bucket-1/*"
      ],
      "Actions": [
        "s3:GetObject"
      ]
    },
    {
      "PolicyName": "example-policy-2",
      "Effect": "Deny",
      "Resources": [
        "arn:aws:s3:::example-bucket-2/*"
      ],
      "Actions": [
        "s3:GetObject"
      ]
    }
  ]
}
```


## Delete Policy

Deletes a policy from the system.

### Request

`DELETE /iam/v1/policies/{policyName}`

#### Path Parameters

| Parameter  | Type   | Required | Description         |
|------------|--------|----------|---------------------|
| policyName | string | Yes      | The name of the policy to delete. |

### Response

#### Success Response

If the policy is successfully deleted, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified policy does not exist. |
| 500 Internal Server Error | An error occurred while deleting the policy. |

### Example

#### Request

```
DELETE /iam/v1/policies/example-policy HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 OK
```



## Update Policy

Updates an existing policy in the system.

### Request

`PUT /iam/v1/policies/{policyName}`

#### Path Parameters

| Parameter  | Type   | Required | Description         |
|------------|--------|----------|---------------------|
| policyName | string | Yes      | The name of the policy to update. |

#### Request Body

The request body must be a JSON object with the following fields:

| Field      | Type     | Required | Description                                                                                     |
|------------|----------|----------|-------------------------------------------------------------------------------------------------|
| PolicyName | string   | No       | 同上.               |
| Version    | string   | No       | 同上.         |
| Effect     | string   | No       | 同上. |
| Resources  | []string | No       | 同上. |
| Actions    | []string | No       | 同上. |

### Response

#### Success Response

If the policy is successfully updated, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified policy does not exist. |
| 500 Internal Server Error | An error occurred while updating the policy. |

### Example

#### Request

```
PUT /iam/v1/policies/example-policy HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json

{
  "Effect": "deny",
  "Resources": [
    "arn:aws:s3:::example-bucket/*"
  ],
  "Actions": [
    "s3:GetObject"
  ]
}
```

#### Response

```
HTTP/1.1 200 OK
```



## Get Secret

Retrieves the details of a secret.

### Request

`GET /iam/v1/secrets/{accessKeyId}`

#### Path Parameters

| Parameter   | Type   | Required | Description         |
|-------------|--------|----------|---------------------|
| accessKeyId | string | Yes      | The access key ID of the secret to retrieve. |

### Response

#### Success Response

If the secret is found, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field           | Type       | Description                                                                                     |
|-----------------|------------|-------------------------------------------------------------------------------------------------|
| AccessKeyId     | string     | The access key ID of the secret.                                                                |
| AccessKeySecret | string     | The access key secret of the secret.                                                            |
| YSId            | string     | The YS ID of the secret.                                                                         |
| Expire          | time.Time  | The expiration time of the secret.                                                               |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified secret does not exist. |
| 500 Internal Server Error | An error occurred while retrieving the secret. |

### Example

#### Request

```
GET /iam/v1/secrets/example-access-key-id HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "AccessKeyId": "example-access-key-id",
  "AccessKeySecret": "example-access-key-secret",
  "YSId": "example-ys-id",
  "Expire": "2022-01-01T00:00:00Z"
}
```



## List Secret

Retrieves a list of all secrets.

### Request

`GET /iam/v1/secrets`

### Response

#### Success Response

If the secrets are found, the response will have a `200 OK` status code and a JSON object with the following fields:

同上

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 500 Internal Server Error | An error occurred while retrieving the secrets. |

### Example

#### Request

```
GET /iam/v1/secrets HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "Secrets": [
    {
      "AccessKeyId": "example-access-key-id-1",
      "AccessKeySecret": "example-access-key-secret-1",
      "YSId": "example-ys-id-1",
      "Expire": "2022-01-01T00:00:00Z"
    },
    {
      "AccessKeyId": "example-access-key-id-2",
      "AccessKeySecret": "example-access-key-secret-2",
      "YSId": "example-ys-id-2",
      "Expire": "2022-01-01T00:00:00Z"
    }
  ]
}
```



## Delete Secret

Deletes a secret from the system.

### Request

`DELETE /secret/{accessKeyId}`

#### Path Parameters

| Parameter   | Type   | Required | Description         |
|-------------|--------|----------|---------------------|
| accessKeyId | string | Yes      | The access key ID of the secret to delete. |

### Response

#### Success Response

If the secret is successfully deleted, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified secret does not exist. |
| 500 Internal Server Error | An error occurred while deleting the secret. |

### Example

#### Request

```
DELETE /iam/v1/secrets/example-access-key-id HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 Ok
```



## Add Secret

Adds a new secret and returns temporary security credentials.

### Request

`POST /iam/v1/secrets`

#### Request Body

               |

### Response

#### Success Response

If the secret is successfully added, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field            | Type   | Description                                                                                     |
|------------------|--------|-------------------------------------------------------------------------------------------------|
| AccessKeyId      | string | The access key ID of the temporary security credentials.                                       |
| AccessKeySecret  | string | The secret access key of the temporary security credentials.                                   |
| YSId             | string | The YS ID of the secret.                                                                        |
| Expire           | string | The expiration time of the temporary security credentials, in RFC3339 format.                  |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 500 Internal Server Error | An error occurred while adding the secret. |

### Example

#### Request

```
POST /iam/v1/secrets HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json

```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "AccessKeyId": "EXAMPLE_ACCESS_KEY_ID",
  "AccessKeySecret": "EXAMPLE_SECRET_ACCESS_KEY",
  "YSId": "EXAMPLE_YS_ID",
  "Expire": "2022-01-01T00:00:00Z"
}
```




## Get Role

Retrieves the details of a role.

### Request

`GET /iam/v1/roles/{roleName}`

#### Path Parameters

| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| roleName  | string | Yes      | The name of the role to retrieve. |

### Response

#### Success Response

If the role is found, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field        | Type       | Description                                                                                     |
|--------------|------------|-------------------------------------------------------------------------------------------------|
| RoleName     | string     | 角色名.                                                                           |
| Description  | string     | 描述.                                                                    |
| TrustPolicy  | RolePolicy | 扮演角色的权限策略.                          |

The `RolePolicy` object has the following fields:

| Field      | Type     | Description                                                                                     |
|------------|----------|-------------------------------------------------------------------------------------------------|
| Actions    | []string | 例如 "sts:AssumeRole".                                                     |
| Effect     | string   | "allow" or "deny".                                   |
| Principals | []string | 授信实体.                                                        |
| Resources  | []string | 资源.                                                          |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified role does not exist. |
| 500 Internal Server Error | An error occurred while retrieving the role. |

### Example

#### Request

```
GET /iam/v1/roles/example-role HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "RoleName": "example-role",
  "Description": "An example role",
  "TrustPolicy": {
    "Actions": [
      "sts:AssumeRole"
    ],
    "Effect": "Allow",
    "Principals": [
      "arn:aws:iam::123456789012:root"
    ],
    "Resources": [
      "arn:aws:iam::123456789012:role/example-role"
    ]
  }
}
```



## List Role

Retrieves a list of all roles.

### Request

`GET /iam/v1/roles`

### Response

#### Success Response

If the roles are found, the response will have a `200 OK` status code and a JSON object with the following fields:

同上

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 500 Internal Server Error | An error occurred while retrieving the roles. |

### Example

#### Request

```
GET /iam/v1/roles HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "Roles": [
    {
      "RoleName": "example-role-1",
      "Description": "An example role",
      "TrustPolicy": {
        "Actions": [
          "sts:AssumeRole"
        ],
        "Effect": "Allow",
        "Principals": [
          "arn:aws:iam::123456789012:root"
        ],
        "Resources": [
          "arn:aws:iam::123456789012:role/example-role-1"
        ]
      }
    },
    {
      "RoleName": "example-role-2",
      "Description": "Another example role",
      "TrustPolicy": {
        "Actions": [
          "sts:AssumeRole"
        ],
        "Effect": "Allow",
        "Principals": [
          "arn:aws:iam::123456789012:root"
        ],
        "Resources": [
          "arn:aws:iam::123456789012:role/example-role-2"
        ]
      }
    }
  ]
}
```


## Add Role

Adds a new role to the system.

### Request

`POST /iam/v1/roles`

#### Request Body

The request body must be a JSON object with the following fields:

| Field        | Type       | Required | Description                                                                                     |
|--------------|------------|----------|-------------------------------------------------------------------------------------------------|
| RoleName     | string     | Yes      | 同上.                                                                           |
| Description  | string     | No       | 同上.                                                                    |
| TrustPolicy  | RolePolicy | Yes      | 同上.                          |

The `RolePolicy` object has the following fields:

| Field      | Type     | Required | Description                                                                                     |
|------------|----------|----------|-------------------------------------------------------------------------------------------------|
| Actions    | []string | Yes      | 同上.                                                     |
| Effect     | string   | Yes      | 同上.                                   |
| Principals | []string | Yes      | 同上.                                                        |
| Resources  | []string | Yes      | 同上.                                                          |

### Response

#### Success Response

If the role is successfully added, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 500 Internal Server Error | An error occurred while adding the role. |

### Example

#### Request

```
POST /iam/v1/roles HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json

{
  "RoleName": "example-role",
  "Description": "An example role",
  "TrustPolicy": {
    "Actions": [
      "sts:AssumeRole"
    ],
    "Effect": "Allow",
    "Principals": [
      "arn:aws:iam::123456789012:root"
    ],
    "Resources": [
      "arn:aws:iam::123456789012:role/example-role"
    ]
  }
}
```

#### Response

```
HTTP/1.1 200 Ok
```




## Delete Role

Deletes a role from the system.

### Request

`DELETE /iam/v1/roles/{roleName}`

#### Path Parameters

| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| roleName  | string | Yes      | The name of the role to delete. |

### Response

#### Success Response

If the role is successfully deleted, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified role does not exist. |
| 500 Internal Server Error | An error occurred while deleting the role. |

### Example

#### Request

```
DELETE /iam/v1/roles/example-role HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 Ok
```



## Update Role

Updates an existing role in the system.

### Request

`PUT /iam/v1/roles/{roleName}`

#### Path Parameters

| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| roleName  | string | Yes      | The name of the role to update. |

#### Request Body

The request body must be a JSON object with the following fields:

| Field        | Type       | Required | Description                                                                                     |
|--------------|------------|----------|-------------------------------------------------------------------------------------------------|
| RoleName     | string     | No       | 同上.                     |
| Description  | string     | No       | 同上.       |
| TrustPolicy  | RolePolicy | No       | 同上.     |

The `RolePolicy` object has the following fields:

| Field      | Type     | Required | Description                                                                                     |
|------------|----------|----------|-------------------------------------------------------------------------------------------------|
| Actions    | []string | Yes      | 同上.                                                     |
| Effect     | string   | Yes      | 同上.                                   |
| Principals | []string | Yes      | 同上.                                                        |
| Resources  | []string | Yes      | 同上.                                                          |

### Response

#### Success Response

If the role is successfully updated, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 404 Not Found | The specified role does not exist. |
| 500 Internal Server Error | An error occurred while updating the role. |

### Example

#### Request

```
PUT /iam/v1/roles/example-role HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json

{
  "Description": "An updated example role",
  "TrustPolicy": {
    "Actions": [
      "sts:AssumeRole"
    ],
    "Effect": "Allow",
    "Principals": [
      "arn:aws:iam::123456789012:root"
    ],
    "Resources": [
      "arn:aws:iam::123456789012:role/example-role"
    ]
  }
}
```

#### Response

```
HTTP/1.1 200 Ok
```



## Add Role-Policy Relation

Adds a relation between a role and a policy.

### Request

`PATCH /iam/v1/roles/{roleName}/policies/{policyName}`

#### Path Parameters

| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| roleName  | string | Yes      | The name of the role to add the relation to. |
| policyName  | string | Yes      | The name of the policy to add the relation to. |

#### Request Body



### Response

#### Success Response

If the relation is successfully added, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified role or policy does not exist. |
| 500 Internal Server Error | An error occurred while adding the relation. |

### Example

#### Request

```
PATCH /iam/v1/roles/example-role/policies/example-policy HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json


```

#### Response

```
HTTP/1.1 200 Ok
```


## Delete Role-Policy Relation

Deletes a relation between a role and a policy.

### Request

`DELETE /iam/v1/roles/{roleName}/policies/{policyName}`

#### Path Parameters

| Parameter | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| roleName  | string | Yes      | The name of the role to delete the relation from. |
| policyName  | string | Yes      | The name of the policy to delete the relation from. |

### Response

#### Success Response

If the relation is successfully deleted, the response will have a `200 Ok` status code and an empty response body.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified role or policy does not exist. |
| 500 Internal Server Error | An error occurred while deleting the relation. |

### Example

#### Request

```
DELETE /iam/v1/roles/example-role/policies/example-policy HTTP/1.1
Host: openapi4-test.yuansuan.com
```

#### Response

```
HTTP/1.1 200 Ok
```



## Assume Role

Assumes a role and returns temporary security credentials.

### Request

`POST /iam/v1/AssumeRole`

#### Request Body

The request body must be a JSON object with the following fields:

| Field            | Type   | Required | Description                                                                                     |
|------------------|--------|----------|-------------------------------------------------------------------------------------------------|
| RoleYrn          | string | Yes      | The YRN of the role to assume.                                                                  |
| RoleSessionName  | string | Yes      | An identifier for the assumed role session.                                                     |
| DurationSeconds  | int    | No       | The duration, in seconds, of the role session. If not specified, the default value is 3600 (1 hour). |

### Response

#### Success Response

If the role is successfully assumed, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field       | Type         | Description                                                                                     |
|-------------|--------------|-------------------------------------------------------------------------------------------------|
| Credentials | Credentials  | The temporary security credentials.                                                             |
| ExpireTime  | time.Time    | The expiration time of the temporary security credentials.                                      |

The `Credentials` object has the following fields:

| Field            | Type   | Description                                                                                     |
|------------------|--------|-------------------------------------------------------------------------------------------------|
| AccessKeyId      | string | The access key ID of the temporary security credentials.                                       |
| AccessKeySecret  | string | The secret access key of the temporary security credentials.                                   |
| SessionToken     | string | The session token of the temporary security credentials.                                       |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 403 Forbidden | The caller is not authorized to assume the specified role. |
| 500 Internal Server Error | An error occurred while assuming the role. |

### Example

#### Request

```
POST /iam/v1/AssumeRole HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json

{
  "RoleYrn": "yrn:ys:iam::123456:role/CloudComputeRole",
  "RoleSessionName": "example-session",
  "DurationSeconds": 3600
}
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "Credentials": {
    "AccessKeyId": "EXAMPLE_ACCESS_KEY_ID",
    "AccessKeySecret": "EXAMPLE_SECRET_ACCESS_KEY",
    "SessionToken": "EXAMPLE_SESSION_TOKEN"
  },
  "ExpireTime": "2022-01-01T00:00:00Z"
}
```



## IsAllow

Checks whether the specified action on the specified resource is allowed for the specified subject.

### Request

`POST /iam/v1/IsAllow`

#### Request Body

The request body must be a JSON object with the following fields:

| Field     | Type   | Required | Description                                                                                     |
|-----------|--------|----------|-------------------------------------------------------------------------------------------------|
| Action    | string | Yes      | The action to check.                                                                            |
| Resource  | string | Yes      | The resource to check.                                                                          |
| Subject   | string | Yes      | The access key ID of the subject to check.                                                      |

### Response

#### Success Response

If the action is allowed, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field    | Type   | Description                                                                                     |
|----------|--------|-------------------------------------------------------------------------------------------------|
| Allow    | bool   | `true` if the action is allowed, `false` otherwise.                                             |
| Message  | string | A message describing the result of the check.                                                   |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 500 Internal Server Error | An error occurred while checking the permission. |

### Example

#### Request

```
POST /iam/v1/IsAllow HTTP/1.1
Host: openapi4-test.yuansuan.com
Content-Type: application/json

{
  "Action": "s3:GetObject",
  "Resource": "arn:aws:s3:::example-bucket/example-object",
  "Subject": "EXAMPLE_ACCESS_KEY_ID"
}
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "Allow": true,
  "Message": "Access granted"
}
`````


# Admin API DOC


## Add Policy

请求体额外添加 `UserId`，其它与上面一致。

### Request

`POST /iam/admin/policies`

#### Example Request Body

```json
{
  "UserId": "example-user-id",
  "PolicyName": "example-policy",
  "Version": "1.0",
  "Effect": "allow",
  "Resources": [
    "yrn:ys:cs::4TiSxuPtJEm:path/4T4ZZvA2tVb/<.*>"
  ],
  "Actions": [
    "<.*>"
  ]
}
```

## Get Policy

Url额外添加 `UserId`，其它与上面一致。

### Request

`GET  /iam/admin/policies/{userId}/{policyName}`


## List Policies

Url额外添加 `UserId`，其它与上面一致。

### Request

`GET /iam/admin/policies/{userId}`


## Delete Policy

Url额外添加 `UserId`，其它与上面一致。

### Request

`DELETE /iam/admin/policies/{userId}/{policyName}`


## Update Policy

请求体及Url额外添加 `UserId`，其它与上面一致。

### Request

`PUT /iam/admin/policies/{userId}/{policyName}`


## Get Secret

与上面一致。

### Request

`GET /iam/admin/secrets/{accessKeyId}`


## List Secret

Url额外添加 `UserId`，其它与上面一致。

### Request

`GET /iam/admin/secrets/user/{userId}`


## Delete Secret

与上面一致。

### Request

`DELETE /iam/admin/secrets/{accessKeyId}`



## Add Secret

Adds a new secret for a user and returns temporary security credentials.

### Request

`POST /iam/admin/secrets`

#### Request Body

The request body must be a JSON object with the following fields:

| Field   | Type   | Required | Description                                                                                     |
|---------|--------|----------|-------------------------------------------------------------------------------------------------|
| UserId  | string | Yes      | The ID of the user to add the secret for.                                                       |
| Tag     | string | No       | A tag to associate with the secret.                                                             |

### Response

#### Success Response

If the secret is successfully added, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field            | Type   | Description                                                                                     |
|------------------|--------|-------------------------------------------------------------------------------------------------|
| AccessKeyId      | string | The access key ID of the temporary security credentials.                                       |
| AccessKeySecret  | string | The secret access key of the temporary security credentials.                                   |
| YSId             | string | The YS ID of the secret.                                                                        |
| Expire           | string | The expiration time of the temporary security credentials, in RFC3339 format.                  |
| Tag              | string | The tag associated with the secret, if any.                                                    |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 500 Internal Server Error | An error occurred while adding the secret. |

### Example

#### Request

```
POST /iam/admin/secrets HTTP/1.1
Host: example.com
Content-Type: application/json

{
  "UserId": "EXAMPLE_USER_ID",
  "Tag": "EXAMPLE_TAG"
}
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "AccessKeyId": "EXAMPLE_ACCESS_KEY_ID",
  "AccessKeySecret": "EXAMPLE_SECRET_ACCESS_KEY",
  "YSId": "EXAMPLE_YS_ID",
  "Expire": "2022-01-01T00:00:00Z",
  "Tag": "EXAMPLE_TAG"
}
```

## Get Role

Url额外添加 `UserId`，其它与上面一致。

### Request

`GET /iam/admin/roles/{userId}/{roleName}`


## List Role

Url额外添加 `UserId`，其它与上面一致。

### Request

`GET /iam/admin/roles/{userId}`

## Add Role

请求体额外添加 `UserId`，其它与上面一致。

### Request

`POST /iam/admin/roles`


## Delete Role

Url额外添加 `UserId`，其它与上面一致。

### Request

`DELETE /iam/admin/roles/{userId}/{roleName}`

## Patch Policy

请求体额外添加 `UserId`，其它与上面(Add Role-Policy Relation)一致。

### Request

`PATCH /iam/admin/roles/{roleName}`


## Detach Policy

请求体额外添加 `UserId`，其它与上面(Delete Role-Policy Relation)一致。

### Request

`POST /iam/admin/roles/{roleName}`



## Get User

Retrieves information about a user.

### Request

`GET /iam/admin/users/{userId}`

#### Path Parameters

| Parameter | Type   | Required | Description                                                                                     |
|-----------|--------|----------|-------------------------------------------------------------------------------------------------|
| userId    | string | Yes      | The ID of the user to retrieve.                                                                 |

### Response

#### Success Response

If the user is successfully retrieved, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field   | Type   | Description                                                                                     |
|---------|--------|-------------------------------------------------------------------------------------------------|
| userId  | string | The ID of the user.                                                                             |
| name    | string | The name of the user.                                                                           |
| phone   | string | The phone number of the user.                                                                   |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 404 Not Found | The specified user does not exist. |
| 500 Internal Server Error | An error occurred while retrieving the user. |

### Example

#### Request

```
GET /iam/admin/users/EXAMPLE_USER_ID HTTP/1.1
Host: example.com
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "userId": "EXAMPLE_USER_ID",
  "name": "John Doe",
  "phone": "18322385431"
}
```


## Add User

Adds a new user.

### Request

`POST /iam/admin/users`

#### Request Body

The request body must be a JSON object with the following fields:

| Field     | Type   | Required | Description                                                                                     |
|-----------|--------|----------|-------------------------------------------------------------------------------------------------|
| Phone    | string | Yes      | The phone number of the user.                                                                   |
| Password | string | Yes      | The password of the user.                                                                       |

### Response

#### Success Response

If the user is successfully added, the response will have a `200 OK` status code and a JSON object with the following fields:

| Field   | Type   | Description                                                                                     |
|---------|--------|-------------------------------------------------------------------------------------------------|
| userId  | string | The ID of the new user.                                                                         |

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 500 Internal Server Error | An error occurred while adding the user. |

### Example

#### Request

```
POST /iam/admin/users HTTP/1.1
Host: example.com
Content-Type: application/json

{
  "Phone": "+1-555-555-5555",
  "Password": "example_password"
}
```

#### Response

```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "userId": "EXAMPLE_USER_ID"
}
```


## Update User

Updates information about a user.

### Request

`PUT /iam/admin/users/{userId}`

#### Path Parameters

| Parameter | Type   | Required | Description                                                                                     |
|-----------|--------|----------|-------------------------------------------------------------------------------------------------|
| userId    | string | Yes      | The ID of the user to update.                                                                   |

#### Request Body

The request body must be a JSON object with the following fields:

| Field   | Type   | Required | Description                                                                                     |
|---------|--------|----------|-------------------------------------------------------------------------------------------------|
| name    | string | No       | The new name of the user. If not specified, the name will not be changed.                        |

### Response

#### Success Response

If the user is successfully updated, the response will have a `200 OK` status code.

#### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 Bad Request | The request body is invalid. |
| 404 Not Found | The specified user does not exist. |
| 500 Internal Server Error | An error occurred while updating the user. |

### Example

#### Request

```
PUT /iam/admin/users/EXAMPLE_USER_ID HTTP/1.1
Host: example.com
Content-Type: application/json

{
  "name": "New Name"
}
```

#### Response

```
HTTP/1.1 200 OK
```

## Latest Version

The latest version reference [here](http://wiki.yuansuan.com/pages/viewpage.action?pageId=13117491)