# API Documentation

gocron provides two types of API interfaces:
1. **Open API**: `/api/v1/*`, uses signature authentication, suitable for third-party system integration (currently supports limited features).
2. **Web API**: `/api/*`, uses Cookie/Session authentication, fully featured, suitable for deep customization or automation scripts (requires login simulation).

---

## Part 1: Open API (Signature Auth)

Designed for third-party integration, using API Key/Secret for signature verification.

### Authentication

::: warning Prerequisites
Configure `api.key` and `api.secret` in `conf/app.ini` before use.
:::

**Common Request Parameters**:

| Parameter | Type | Required | Description | Example |
|:---:|:---:|:---:|:---|:---|
| time | int64 | Yes | Current Unix timestamp (seconds) | `1495519926` |
| sign | string | Yes | Signature string (lowercase) | `418819b148815894e4abfa944d19f70f` |

**Signature Formula**:
```
sign = MD5(api.key + time + url_path + api.secret)
```

### Interface List

#### 1. Enable Task
- **Endpoint**: `/api/v1/task/enable/:id`
- **Method**: `POST`
- **Description**: Enable task by ID

#### 2. Disable Task
- **Endpoint**: `/api/v1/task/disable/:id`
- **Method**: `POST`
- **Description**: Disable task by ID

#### 3. Remove Task Log
- **Endpoint**: `/api/v1/tasklog/remove/:id`
- **Method**: `POST`
- **Description**: Remove log by ID

---

## Part 2: Web API (Cookie Auth)

Interfaces used by the management dashboard, providing full functionality. Requires login to obtain Cookie.

### Authentication

1. Call login interface `/api/user/login`
2. Get `Set-Cookie` from response Header
3. Include this Cookie in the Header of subsequent requests

### 1. User Login

- **Endpoint**: `/api/user/login`
- **Method**: `POST`

**Parameters**:

| Parameter | Type | Required | Description |
|:---:|:---:|:---:|:---|
| username | string | Yes | Username |
| password | string | Yes | Password |

**Response Example**:
```json
{
    "code": 0,
    "message": "Login successful",
    "data": {
        "token": "..." // Primarily for frontend state, backend relies on Cookie
    }
}
```

### 2. Task Management

#### Get Task List
- **Endpoint**: `/api/task`
- **Method**: `GET`
- **Parameters**:
  - `page`: Page number (default 1)
  - `page_size`: Items per page (default 20)
  - `id`: Task ID (optional)
  - `name`: Task name (optional)
  - `host_id`: Host ID (optional)
  - `status`: Status (1:Enabled 0:Disabled)

#### Get Task Detail
- **Endpoint**: `/api/task/:id`
- **Method**: `GET`

#### Create/Edit Task
- **Endpoint**: `/api/task/store`
- **Method**: `POST`
- **Parameters**:

| Parameter | Type | Required | Description |
|:---:|:---:|:---:|:---|
| id | int | No | Task ID (Required for edit, 0 for new) |
| name | string | Yes | Task Name |
| spec | string | Yes | Crontab Expression (Second-level) |
| command | string | Yes | Command to execute |
| protocol | int | Yes | Protocol (1:HTTP 2:Shell) |
| host_id | int | Yes | Host ID |
| timeout | int | No | Timeout (seconds) |
| multi | int | No | Single instance (1:Yes 2:No) |
| retry_times | int | No | Retry times |
| remark | string | No | Remark |

#### Manually Run Task
- **Endpoint**: `/api/task/run/:id`
- **Method**: `GET`

#### Remove Task
- **Endpoint**: `/api/task/remove/:id`
- **Method**: `POST`

### 3. Host Management

#### Get Host List
- **Endpoint**: `/api/host`
- **Method**: `GET`
- **Parameters**:
  - `page`: Page number
  - `page_size`: Items per page
  - `name`: Host name (optional)

#### Create/Edit Host
- **Endpoint**: `/api/host/store`
- **Method**: `POST`
- **Parameters**:

| Parameter | Type | Required | Description |
|:---:|:---:|:---:|:---|
| id | int | No | Host ID (Required for edit) |
| name | string | Yes | Host Name |
| alias | string | Yes | Alias |
| port | int | Yes | Port (Default 5921) |
| remark | string | No | Remark |

#### Remove Host
- **Endpoint**: `/api/host/remove/:id`
- **Method**: `POST`

### 4. Log Management

#### Get Task Logs
- **Endpoint**: `/api/task/log`
- **Method**: `GET`
- **Parameters**:
  - `page`: Page number
  - `page_size`: Items per page
  - `task_id`: Task ID (optional)
  - `status`: Status (1:Failed 2:Running 3:Success 4:Cancelled)

#### Stop Task Execution
- **Endpoint**: `/api/task/log/stop`
- **Method**: `POST`
- **Parameters**:
  - `id`: Log ID
  - `task_id`: Task ID