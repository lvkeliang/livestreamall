# **直播平台**

本项目提供了一个简易的直播平台实现，旨在学习流服务的搭建。

---

## **目录**

1. [项目简介](#项目简介)  
2. [快速开始](#快速开始)  
   - [运行环境](#运行环境)  
   - [MySQL部署与配置](#MySQL部署与配置)  
   - [FFmpeg的安装](#FFmpeg的安装)  
   - [Nginx部署与配置](#Nginx部署与配置)  
   - [启动项目](#启动项目)  
3. [接口概览](#接口概览)  
4. [其他说明](#其他说明)  

---

## **项目简介**

本系统旨在提供一个基础功能完整的直播平台服务，包括：
- **用户管理**：支持用户注册、登录及获取用户信息。
- **直播功能**：支持主播推流、观看直播和直播间管理。
- **推流认证**：对推流请求进行权限验证。
- **流媒体支持**：支持RTMP推流和HLS拉流协议。

---

## **快速开始**

### **运行环境**
- **开发语言**：Go (Golang)
- **依赖框架**：Gin (HTTP框架)
- **音视频编解码**：FFmpeg (流转换)
- **流媒体服务器**：Nginx + RTMP 模块
- **数据库**：MySQL 或其他兼容关系型数据库，本项目使用GORM进行数据库操作
- **环境依赖**：
  - 安装 `Go` (1.13 或更高版本)
  - 配置并启动 Nginx + RTMP 服务
  - 数据库初始化并导入必要表结构

---

### **MySQL部署与配置**

本处不赘述MySQL的部署和配置，请自行查阅，并记住你设置的账号和密码，用于配置数据库连接。

### **FFmpeg的安装**

FFmpeg 是一个处理视频、音频等多媒体的开发工具包，在下面Nginx的RTMP模块会依赖FFmpeg将输入的 RTMP 流转为适合分发的格式 (如 HLS 或 DASH)。

本处不赘述FFmpeg的安装过程，请自行查阅FFmpeg的安装，并记住安装地址，用于配置nginx.conf文件。

### **Nginx部署与配置**

#### **1. 安装必要的软件包**
在 Ubuntu 官方 Nginx 中并不默认包含 nginx-rtmp-module，需要从源码安装 Nginx 并编译 RTMP 模块。

运行以下命令安装必要的依赖：
```bash
sudo apt update
sudo apt install -y build-essential libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev
```

#### **2. 下载 Nginx 和 RTMP 模块源码**
下载 Nginx 和 nginx-rtmp-module 的源码：
```bash
# 下载 Nginx 源码
cd /usr/local/src
wget http://nginx.org/download/nginx-1.25.2.tar.gz
tar -zxvf nginx-1.25.2.tar.gz

# 下载 nginx-rtmp-module
git clone https://github.com/arut/nginx-rtmp-module.git
```

#### **3. 编译 Nginx 并启用 RTMP 模块**
进入 Nginx 源码目录并开始编译：
```bash
cd /usr/local/src/nginx-1.25.2

# 配置 Nginx 编译选项，启用 RTMP 模块
./configure --with-http_ssl_module --add-module=/usr/local/src/nginx-rtmp-module

# 编译并安装
make
sudo make install
```
安装完成后，Nginx 默认会安装到 `/usr/local/nginx/`。

#### **4. 配置 Nginx 添加 RTMP 功能**
编辑 Nginx 配置文件 `/usr/local/nginx/conf/nginx.conf`，添加 RTMP 支持：

**如果你的FFmpeg安装路径与本文件中的 `/usr/bin/ffmpeg` 不同，请将 `exec_push  /usr/bin/ffmpeg -i` 这一行的FFmpeg路径替换为你自己的。**

```bash
sudo nano /usr/local/nginx/conf/nginx.conf
```
修改文件为以下配置(你也可以在本项目的nginxConf文件夹下找到该配置文件)：
```nginx
worker_processes  auto;
error_log /var/log/nginx/error.log debug;
#error_log  logs/error.log;

events {
    worker_connections  1024;
}

# RTMP configuration
rtmp {
    server {
		listen 1935; # Listen on standard RTMP port
		chunk_size 4000;
		# ping 30s;
		# notify_method get;

		# This application is to accept incoming stream
		application live {
			live on; # Allows live input

			# 启用访问日志
    		access_log /var/log/nginx/rtmp_access.log;

			# 推流认证
			on_publish http://127.0.0.1:9089/auth/publish;

			# 推流停止时触发 on_done
    		on_done http://127.0.0.1:9089/auth/stop_publish;

			# for each received stream, transcode for adaptive streaming
			# This single ffmpeg command takes the input and transforms
			# the source into 4 different streams with different bitrates
			# and qualities. # these settings respect the aspect ratio.
			exec_push  /usr/bin/ffmpeg -i rtmp://localhost:1935/$app/$name -async 1 -vsync -1
						-c:v libx264 -c:a aac -b:v 256k  -b:a 64k  -vf "scale=480:trunc(ow/a/2)*2"  -tune zerolatency -preset superfast -crf 23 -f flv rtmp://localhost:1935/show/$name_low
						-c:v libx264 -c:a aac -b:v 768k  -b:a 128k -vf "scale=720:trunc(ow/a/2)*2"  -tune zerolatency -preset superfast -crf 23 -f flv rtmp://localhost:1935/show/$name_mid
						-c:v libx264 -c:a aac -b:v 1024k -b:a 128k -vf "scale=960:trunc(ow/a/2)*2"  -tune zerolatency -preset superfast -crf 23 -f flv rtmp://localhost:1935/show/$name_high
						-c:v libx264 -c:a aac -b:v 1920k -b:a 128k -vf "scale=1280:trunc(ow/a/2)*2" -tune zerolatency -preset superfast -crf 23 -f flv rtmp://localhost:1935/show/$name_hd720
						-c copy -f flv rtmp://localhost:1935/show/$name_src;
			drop_idle_publisher 10s;
		}

		# This is the HLS application
		application show {
			live on; # Allows live input from above application
			deny play all; # disable consuming the stream from nginx as rtmp

			hls on; # Enable HTTP Live Streaming
			hls_fragment 3;
			hls_playlist_length 20;
			hls_path /mnt/hls/;  # hls fragments path
			# Instruct clients to adjust resolution according to bandwidth
			hls_variant _src BANDWIDTH=4096000; # Source bitrate, source resolution
			hls_variant _hd720 BANDWIDTH=2048000; # High bitrate, HD 720p resolution
			hls_variant _high BANDWIDTH=1152000; # High bitrate, higher-than-SD resolution
			hls_variant _mid BANDWIDTH=448000; # Medium bitrate, SD resolution
			hls_variant _low BANDWIDTH=288000; # Low bitrate, sub-SD resolution

			# MPEG-DASH
            dash on;
            dash_path /mnt/dash/;  # dash fragments path
			dash_fragment 3;
			dash_playlist_length 20;
		}
	}
}


http {
	sendfile off;
	tcp_nopush on;
	directio 512;
	# aio on;

	# HTTP server required to serve the player and HLS fragments
	server {
		listen 8080;

		# Serve HLS fragments
		location /hls {
			types {
				application/vnd.apple.mpegurl m3u8;
				video/mp2t ts;
			}

			root /mnt;

            add_header Cache-Control no-cache; # Disable cache

			# CORS setup
			add_header 'Access-Control-Allow-Origin' '*' always;
			add_header 'Access-Control-Expose-Headers' 'Content-Length';

			# allow CORS preflight requests
			if ($request_method = 'OPTIONS') {
				add_header 'Access-Control-Allow-Origin' '*';
				add_header 'Access-Control-Max-Age' 1728000;
				add_header 'Content-Type' 'text/plain charset=UTF-8';
				add_header 'Content-Length' 0;
				return 204;
			}
		}

        # Serve DASH fragments
        location /dash {
            types {
                application/dash+xml mpd;
                video/mp4 mp4;
            }

			root /mnt;

			add_header Cache-Control no-cache; # Disable cache


            # CORS setup
            add_header 'Access-Control-Allow-Origin' '*' always;
            add_header 'Access-Control-Expose-Headers' 'Content-Length';

            # Allow CORS preflight requests
            if ($request_method = 'OPTIONS') {
                add_header 'Access-Control-Allow-Origin' '*';
                add_header 'Access-Control-Max-Age' 1728000;
                add_header 'Content-Type' 'text/plain charset=UTF-8';
                add_header 'Content-Length' 0;
                return 204;
            }
        }

		# This URL provides RTMP statistics in XML
		location /stat {
			rtmp_stat all;
			rtmp_stat_stylesheet stat.xsl; # Use stat.xsl stylesheet
		}

		location /stat.xsl {
			# XML stylesheet to view RTMP stats.
			root /usr/local/nginx/html;
		}

	}
}
```

#### **5. 创建 HLS 目录**
为 HLS 配置创建目录，用于存放 HLS 切片文件：
```bash
sudo mkdir -p /mnt/hls
sudo mkdir -p /mnt/dash

sudo chmod -R 777 /mnt/hls
sudo chmod -R 777 /mnt/dash
```

#### **6. 启动 Nginx**
启动 Nginx 并检查是否运行正常：
```bash
sudo /usr/local/nginx/sbin/nginx -c /usr/local/nginx/conf/nginx.conf

# 检查 Nginx 是否启动
sudo /usr/local/nginx/sbin/nginx -t

# 查看 1935 端口是否被监听
netstat -tuln | grep 1935
```

如果配置无误，Nginx 应该成功启动，并且 1935 端口在监听状态。

---

### **启动项目**

1. **克隆项目代码**：
   
   ```bash
   git clone https://github.com/lvkeliang/livestreamall.git
   cd livestreamall
   ```
   
2. **安装依赖**：
   ```bash
   go mod tidy
   ```
   
3. **配置数据库连接**：

   请将本项目下的 `./dao/dao.go` 中的 `InitDB` 函数中的 `dsn` 变量替换为你自己的MySQL连接字符串
   
4. **配置IP设置**：

   请将本项目下的 `./config/config.go` 中的IP地址按照你个人的IP进行配置。
   
   - `YOUR_IP_ADDRESS` 是该项目运行机器的IP
   - `YOUR_STREAM_IP_ADDRESS` 是你的Nginx流服务器运行机器的IP
   
   如果你的Nginx和项目都运行在同一机器，则YOUR_IP_ADDRESS和YOUR_STREAM_IP_ADDRESS都使用本机的IP。

5. **运行项目**：
   
   ```bash
   go run main.go
   ```
   项目默认运行在 `http://localhost:9089`。

---

## **接口概览**

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
| `/stream/{stream_name}` | GET | 获取拉流相关地址 |
| `rtmp://{流服务器地址}/live` | POST | 推流地址 |
| `http://{流服务器地址}/hls/{流名称}.m3u8`   | GET | 拉流地址 |

详细接口信息请参考[接口文档](https://github.com/lvkeliang/livestreamall/blob/main/livestreamallAPI.md)。

---

## **其他说明**

这是一个个人实验小项目，如果能帮助到你，不胜荣幸！