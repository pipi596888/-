# TTS 多角色系统后端 (go-zero)

## 项目结构

```
tts-backend/
├── user-api/           # 用户与余额服务
├── voice-api/          # 音色管理服务
├── tts-api/           # TTS 生成 API
├── tts-worker/        # TTS 任务执行器
├── common/            # 公共模块
└── sql/               # 数据库脚本
```

## 服务说明

| 服务 | 端口 | 说明 |
|------|------|------|
| user-api | 8081 | 用户注册登录、余额管理 |
| voice-api | 8082 | 音色 CRUD |
| tts-api | 8080 | 提交生成任务、查询状态 |
| tts-worker | - | 异步执行 TTS 任务 |

## 快速开始

### 1. 初始化数据库

```bash
mysql -u root -p < sql/init.sql
```

### 2. 修改配置

各服务配置文件在 `etc/*.yaml`，请根据实际情况修改：
- MySQL 连接信息
- Redis 连接信息
- OSS 配置（用于存储生成的音频）

### 3. 启动服务

```bash
# 启动 user-api
cd user-api
go run user.go -f etc/user-api.yaml

# 启动 voice-api
cd voice-api
go run voice.go -f etc/voice-api.yaml

# 启动 tts-api
cd tts-api
go run tts.go -f etc/tts-api.yaml

# 启动 tts-worker
cd tts-worker
go run tts-worker.go -f etc/tts-worker.yaml
```

### 4. 使用 go-zero 工具生成代码

```bash
# 安装 go-zero 工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 生成 user-api 代码
goctl api go -api api/user.api -dir . -style gozero

# 生成 voice-api 代码
goctl api go -api api/voice.api -dir . -style gozero

# 生成 tts-api 代码
goctl api go -api api/tts.api -dir . -style gozero
```

## API 接口

### 用户服务 (user-api:8081)

- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `GET /api/user/info` - 获取用户信息

### 音色服务 (voice-api:8082)

- `GET /api/voice/list` - 获取音色列表
- `POST /api/voice/create` - 创建音色
- `DELETE /api/voice/:id` - 删除音色
- `PUT /api/voice/default/:id` - 设置默认音色

### TTS 服务 (tts-api:8080)

- `POST /api/tts/generate` - 提交生成任务
- `GET /api/tts/task/:taskId` - 查询任务状态

## 技术栈

- **Go**: 1.21+
- **go-zero**: 微服务框架
- **MySQL**: 持久化存储
- **Redis**: 缓存、限流
- **RabbitMQ**: 消息队列（可选）
- **OSS**: 音频文件存储
- **ffmpeg**: 音频处理

## TTS 引擎

支持多种 TTS 引擎：
- VITS（本地模型）
- 阿里云 TTS
- 腾讯云 TTS

引擎接口定义见 `tts-worker/internal/engine/interface.go`

## 扩展能力

- 情绪合成：emotion 参数透传
- 双声道对话：左右声道分离
- 批量生成：批处理任务
- AI 自动分角色：接入 LLM
