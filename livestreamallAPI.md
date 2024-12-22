# **API 文档**

---

## **1. 总览**

| **接口**             | **请求方式** | **描述**                     |
| -------------------- | ------------ | ---------------------------- |
| `/home`        | GET          | 访问主页         |
| `/user/login`        | GET          | 访问登录和注册页面         |
| `/user/register`     | POST         | 注册新用户                   |
| `/user/login`        | POST         | 验证用户登录，返回Token      |
| `/user/info`         | GET          | 获取用户信息                 |
| `/auth/publish`      | POST         | 验证推流授权                 |
| `/auth/stop_publish` | POST         | 停止推流                     |
| `/live/play`         | GET          | 访问直播观看页面               |
| `/live/start`        | GET          | 访问开始直播页面           |
| `/live/start`        | POST         | 开始直播，返回推流地址和密钥 |
| `/live/live_rooms`   | GET          | 获取直播间列表               |
| `/live/live-room/{stream_name}`   | GET | 获取直播间信息 |
| `ws://{服务器地址}/live/ws/{stream_name}`   | GET | 聊天室收发 |
| `/stream/{stream_name}` | GET | 获取拉流相关地址 |
| `rtmp://{流服务器地址}/live` | POST | 推流地址 |
| `http://{流服务器地址}/hls/{流名称}.m3u8`   | GET | 拉流地址 |

---

## **2. 通用设置**

### **2.1 静态资源**
**路径：`/static`**

- 描述：提供静态资源文件支持，例如CSS、JS文件。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| 无参数       | -        | -        | -        |

---

### **2.2 HTML 模板**
**路径：`/html/*`**

- 描述：提供前端页面支持，用于渲染HTML页面。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| 无参数       | -        | -        | -        |

---

## **3. 主页**

**GET `/home`**

- 描述：访问主页。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| 无参数       | -        | -        | -        |

响应为对应的页面

---

## **4. 用户相关接口**

### **4.1 用户登录和注册页面**

**GET `/user/login`**

- 描述：访问登录和注册页面。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| 无参数       | -        | -        | -        |

响应为对应的页面

---

### **4.2 用户注册**

**POST `/user/register`**

- 描述：处理用户注册逻辑，保存用户信息。

携带表单: 

| **参数名称** | **类型** | **必填** | **描述**     |
| ------------ | -------- | -------- | ------------ |
| `nickname`   | String   | 否       | 用户昵称     |
| `password`   | String   | 是       | 用户密码     |
| `email`      | String   | 是       | 用户邮箱地址 |

响应示例：

```json
{
    "status": 10000,
    "info": "success",
    "data": null
}
```

---

### **4.3 用户登录**

**POST `/user/login`**

- 描述：验证用户的邮箱和密码，返回认证Token。

携带表单: 

| **参数名称** | **类型** | **必填** | **描述**     |
| ------------ | -------- | -------- | ------------ |
| `email`      | String   | 是       | 用户邮箱地址 |
| `password`   | String   | 是       | 用户密码     |

响应示例：

```json
{
    "status": 10000,
    "info": "success",
    "token": "your-token"
}
```

客户端收到正确响应后，需要将token存入本地存储，并在之后的一些需要登录状态的请求中，添加如下header: key为`Authorization`，value为你存储的token。

---

### **4.4 获取用户信息**

**GET `/user/info`**

- 描述：返回登录用户的详细信息，需携带认证Token。

携带Header: 

| **参数名称**    | **类型** | **必填** | **描述**              |
| --------------- | -------- | -------- | --------------------- |
| `Authorization` | Header   | 是       | 用户认证Token（必需） |

响应示例：

```json
{
    "user": {
        "created_at": "2024-12-21T22:04:40.192+08:00",
        "email": "abc@abc.com",
        "id": 4,
        "nickname": "abc"
    }
}
```

---

## **5. 推流认证相关接口**

### **5.1 推流认证**

**POST `/auth/publish`**

- 描述：验证推流是否被授权，后端将直播状态改为正在直播。
- 该接口由Nginx的on_publish钩子调用，不提供给客户端。

携带表单: 

| **参数名称**  | **类型** | **必填** | **描述** |
| ------------- | -------- | -------- | -------- |
| `name` | String   | 是       | 流名称(串流密钥) |
| `app` | String | 是       | 应用名称 |

响应200状态码即为允许推流，其他拒绝推流。

---

### **5.2 停止推流**

**POST `/auth/stop_publish`**

- 描述：停止推流，关闭直播流，后端将直播状态改为未开播。
- 该接口由Nginx的on_done钩子调用，不提供给客户端。

| **参数名称**  | **类型** | **必填** | **描述** |
| ------------- | -------- | -------- | -------- |
| `name` | String   | 是       | 流名称(串流密钥) |
| `app` | String | 是       | 应用名称 |

返回200状态码即可。

---

## **6. 直播相关接口**

### **6.1 观看直播页面**

**GET `/live/play?stream_id={stream-name}`**

- 描述：访问直播观看页面。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| `stream-name`       | url参数        | 是        | 流名称        |

响应为对应的页面

---

### **6.2 开始直播页面**

**GET `/live/start`**

- 描述：返回开始直播页面，用于主播推流。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| 无参数       | -        | -        | -        |

响应为对应的页面

---

### **6.3 开始直播**

**POST `/live/start`**

- 描述：开始直播，获取推流地址和推流密钥。

携带Header: 

| **参数名称**    | **类型** | **必填** | **描述**              |
| --------------- | -------- | -------- | --------------------- |
| `Authorization` | Header   | 是       | 用户认证Token（必需） |

携带表单: 

| **参数名称**    | **类型** | **必填** | **描述**              |
| --------------- | -------- | -------- | --------------------- |
| `title`         | String   | 是       | 直播标题              |
| `description` | String   | 是       | 直播简介 |

响应示例：

```json
{
    "stream_name": "4",
    "token": "your-stream-key",
    "user_id": "4"
}
```

---

### **6.4 获取直播间列表**

**GET `/live/live_rooms`**

- 描述：获取所有直播间的列表。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| 无参数       | -        | -        | -        |

响应示例：

```json
{
    "live_rooms": [
        {
            "ID": 20,
            "Title": "测试直播1",
            "StreamName": "stream-name",
            "Description": "测试测试",
            "IsLive": true,
            "CreatedAt": "2024-12-21T22:36:15.02+08:00",
            "UserID": "4",
            "Messages": null
        },
        {
            "ID": 20,
            "Title": "测试直播2",
            "StreamName": "stream-name2",
            "Description": "测试测试",
            "IsLive": true,
            "CreatedAt": "2024-12-21T22:38:15.02+08:00",
            "UserID": "3",
            "Messages": null
        }
    ]
}
```

---



### **6.5 获取直播间信息**

**GET `/live/live-room/{stream_name}`**

- 描述：根据推流名称获取指定直播间的信息。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| `stream-name`       | url参数        | 是        | 流名称        |

- 请求示例

**URL**: `/live/live-room/example-stream`

- 响应示例

成功响应:

```json
{
    "id": 1,
    "title": "Example Live Stream",
    "description": "This is an example live stream description.",
    "is_live": true
}
```

失败响应:

**404 Not Found**

```json
{
    "error": "live room not found"
}
```

**400 Bad Request**

```json
{
    "error": "stream_name parameter is required"
}
```

---

### **6.6 WebSocket 连接处理**

**GET `/live/ws/{stream_name}`**

- 描述：处理针对指定直播间的WebSocket连接，允许用户实时发送和接收消息。

- 请求示例

**URL**: `/live/ws/example-stream`

- 响应示例

成功连接：

*此请求不会返回传统的JSON响应，而是建立WebSocket连接。连接成功后，客户端可以发送消息。*

- **发送的消息示例**：

```json
{
    "token": "your-auth-token",
    "Content": "Hello, World!"
}
```

失败响应:

**400 Bad Request**

```json
{
    "error": "Missing stream_name"
}
```

**404 Not Found**

```json
{
    "error": "live room not found"
}
```

**500 Internal Server Error**

```json
{
    "error": "WebSocket upgrade failed"
}
```

- **接收消息示例**：

```json
{
    "liveRoomID": 18,
    "username": "消息发布者昵称",
    "content": "消息内容"
}
```

---

### 详细说明

1. **获取直播间信息**：`GET /live/live-room/:stream_name` 允许客户端根据推流名称查询直播间的详细信息，例如标题、简介和在线状态。这对于查看某个直播间的基本信息非常有用。

2. **WebSocket 连接处理**：`GET /live/ws/:stream_name` 将建立与指定直播间的WebSocket连接，用户可以通过这个连接发送消息并接收实时更新。 

希望这两个API文档符合您的要求，如果需要进一步修改或者添加更多细节，请告诉我！

---

## **7. 流相关接口**

### **7.1 获取拉流相关地址**

**GET `/stream/{stream_name}`**

- 描述：获取hls拉流地址、播放网页地址。

| **参数名称** | **类型** | **必填** | **描述** |
| ------------ | -------- | -------- | -------- |
| `stream-name`       | url参数        | 是        | 流名称        |

响应示例：

```json
{
    "dash_url": "http://47.92.137.133:8080/dash/stream_name.mpd",
    "hls_url": "http://47.92.137.133:8080/hls/stream_name.m3u8",
    "play_url": "http://47.92.137.133:9089/live/play?stream_id=stream_name"
}
```

---

**推流地址和拉流地址由流服务器提供，故IP和端口可能与其他接口不同。**

---

### **7.2 推流地址**
**路径：`rtmp://{流服务器地址}/live`**

- 描述：作为OBS等推流软件的推流地址(服务器)，同时需要填写由`POST /live/start`接口获取到的串流密钥。
- rtmp推流的默认端口是 1935 。

---

### **7.3 拉流地址**

**路径：`http://{流服务器地址}/hls/{流名称}.m3u8`**

- 描述：提供hls拉流，在`GET /live/play?stream_id={stream-name}`接口中的stream-name即是该拉流地址的流名称，在访问它的时候由前端JS提取`{stream-name}`并填入`{流名称}`来拉流。
