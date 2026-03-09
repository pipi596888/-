# TTS 语音合成系统

一个功能完善的文本转语音(Text-to-Speech)在线合成系统，支持多角色配音、音频下载、作品管理等功能。

## 项目简介

本系统采用前后端分离架构，前端使用 Vue 3 + TypeScript，后端使用 Go + go-zero 微服务框架。支持多种 TTS 引擎，包括 VITS（本地模型）、阿里云 TTS、腾讯云 TTS。

## 功能特性

- **多角色配音**：支持为不同文本片段指定不同音色
- **音色管理**：创建、编辑、删除音色，设置默认音色
- **作品管理**：查看和管理已生成的音频作品
- **音频下载**：一键下载合成的音频文件
- **任务队列**：异步处理 TTS 任务，支持任务状态查询

## 技术栈

### 前端
- Vue 3 + TypeScript
- Vite 构建工具
- Element Plus UI 组件库
- Pinia 状态管理
- Vue Router 路由管理
- Axios HTTP 客户端

### 后端
- Go 1.21+
- go-zero 微服务框架
- MySQL 数据库
- Redis 缓存
- RabbitMQ 消息队列（可选）

### 基础设施
- Docker / Docker Compose
- ffmpeg 音频处理

## 项目结构

```
文本转音频/
├── tts-front/                 # 前端项目
│   ├── src/
│   │   ├── api/              # API 接口层
│   │   ├── components/       # 公共组件
│   │   ├── views/            # 页面视图
│   │   ├── store/            # 状态管理
│   │   ├── router/           # 路由配置
│   │   ├── utils/            # 工具函数
│   │   ├── types/            # 类型定义
│   │   ├── assets/           # 静态资源
│   │   ├── App.vue           # 根组件
│   │   ├── main.ts           # 入口文件
│   │   └── style.css         # 全局样式
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
│
├── tts-backend/               # 后端项目
│   ├── user-api/             # 用户服务
│   ├── voice-api/            # 音色服务
│   ├── tts-api/              # TTS 生成服务
│   ├── tts-worker/           # TTS 任务执行器
│   ├── common/               # 公共模块
│   └── sql/                  # 数据库脚本
│
└── docker-compose.yaml       # Docker 配置
```

---

## 前端项目详解 (tts-front)

### 目录结构

```
tts-front/src/
├── api/                      # API 接口层
│   ├── user.ts              # 用户相关 API
│   ├── voice.ts             # 音色相关 API
│   ├── tts.ts               # TTS 生成相关 API
│   └── works.ts             # 作品相关 API
│
├── components/              # 公共组件
│   ├── VoiceSelector/       # 音色选择器组件
│   ├── AudioPlayer/         # 音频播放器组件
│   ├── GeneratePanel/       # 生成面板组件
│   ├── SegmentEditor/       # 文本片段编辑器组件
│   ├── RolePanel/           # 角色面板组件
│   └── HelloWorld.vue       # 示例组件
│
├── views/                   # 页面视图
│   ├── Login/               # 登录页面
│   ├── GenerateAudio/       # 音频生成页面
│   ├── Works/               # 我的作品页面
│   └── VoiceManage/         # 音色管理页面
│
├── store/                   # 状态管理 (Pinia)
│   ├── user.ts              # 用户状态管理
│   ├── voice.ts             # 音色状态管理
│   └── tts.ts               # TTS 状态管理
│
├── router/                  # 路由配置
│   └── index.ts             # 路由定义和守卫
│
├── utils/                   # 工具函数
│   ├── request.ts           # axios 封装
│   └── audio.ts             # 音频工具类
│
├── types/                   # TypeScript 类型定义
│   ├── api.d.ts             # API 类型
│   ├── tts.d.ts             # TTS 类型
│   └── voice.d.ts           # 音色类型
│
├── assets/                  # 静态资源
│   └── vue.svg              # Vue Logo
│
├── App.vue                  # 根组件
├── main.ts                  # 入口文件
└── style.css                # 全局样式
```

### 各模块详细功能

#### 1. API 层 (`src/api/`)

负责与后端服务进行 HTTP 通信，提供统一的数据请求接口。

| 文件 | 功能说明 |
|------|----------|
| `user.ts` | 用户登录、注册、获取用户信息 |
| `voice.ts` | 音色列表查询、创建音色、删除音色、设置默认音色 |
| `tts.ts` | 提交 TTS 生成任务、查询任务状态 |
| `works.ts` | 获取用户作品列表 |

#### 2. 组件层 (`src/components/`)

可复用的 Vue 组件，构成页面 UI 的基础。

| 组件 | 功能说明 |
|------|----------|
| `VoiceSelector` | 音色下拉选择器，用于选择文本片段使用的音色 |
| `AudioPlayer` | 音频播放器，支持播放、暂停、进度拖动、音量控制 |
| `GeneratePanel` | 生成控制面板，包含格式设置、声道选择、生成按钮 |
| `SegmentEditor` | 文本片段编辑器，支持添加、编辑、删除、拖拽排序文本片段 |
| `RolePanel` | 角色面板，展示和管理音色角色 |

#### 3. 视图层 (`src/views/`)

页面级组件，对应路由对应的页面视图。

| 页面 | 路由 | 功能说明 |
|------|------|----------|
| `Login` | `/login` | 用户登录/注册页面 |
| `GenerateAudio` | `/generate` | 音频生成主页面，包含文本编辑、音色选择、生成控制 |
| `Works` | `/works` | 我的作品页面，展示历史生成的音频作品 |
| `VoiceManage` | `/voice` | 音色管理页面，CRUD 音色 |

#### 4. 状态管理层 (`src/store/`)

使用 Pinia 进行全局状态管理。

| Store | 功能说明 |
|-------|----------|
| `user.ts` | 用户信息、Token 管理，提供登录、登出、获取用户信息功能 |
| `voice.ts` | 音色列表管理，提供获取、创建、删除、设置默认音色功能 |
| `tts.ts` | TTS 生成状态管理，管理文本片段、生成参数、任务状态、轮询任务 |

**TTS Store 核心功能：**
- `segments` - 文本片段数组，支持添加、修改、删除、排序
- `generateAudio()` - 提交生成任务
- `startPolling()` / `stopPolling()` - 任务状态轮询
- `totalCharacters` - 计算总字符数

#### 5. 路由层 (`src/router/`)

Vue Router 路由配置，包含路由守卫控制。

- 路由列表：Login、GenerateAudio、Works、VoiceManage
- 路由守卫：验证 Token，未登录自动跳转登录页

#### 6. 工具层 (`src/utils/`)

| 工具 | 功能说明 |
|------|----------|
| `request.ts` | 封装 Axios，添加请求/响应拦截器，自动携带 Token，处理错误信息 |
| `audio.ts` | `AudioPlayer` 类：音频播放控制；`downloadAudio()`：音频下载功能 |

#### 7. 类型定义层 (`src/types/`)

TypeScript 类型声明文件。

| 文件 | 类型说明 |
|------|----------|
| `api.d.ts` | 通用 API 类型 |
| `tts.d.ts` | TTS 相关类型：Segment、TTSGenerateParams、TTSTaskResponse 等 |
| `voice.d.ts` | 音色相关类型：Voice、VoiceCreateParams 等 |

---

## 后端项目详解 (tts-backend)

### 目录结构

```
tts-backend/
├── user-api/                 # 用户服务
│   ├── api/                 # API 定义文件
│   ├── etc/                 # 配置文件
│   ├── internal/
│   │   ├── config/         # 配置结构体
│   │   ├── handler/        # HTTP 处理器
│   │   ├── logic/          # 业务逻辑
│   │   ├── model/          # 数据模型
│   │   └── svc/            # 服务上下文
│   ├── user.go             # 入口文件
│   └── go.mod
│
├── voice-api/               # 音色服务
│   ├── api/
│   ├── etc/
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── logic/
│   │   ├── model/
│   │   ├── svc/
│   │   └── types/
│   ├── voice.go
│   └── go.mod
│
├── tts-api/                 # TTS 生成服务
│   ├── api/
│   ├── etc/
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── logic/
│   │   ├── model/
│   │   ├── svc/
│   │   └── types/
│   ├── tts.go
│   └── go.mod
│
├── tts-worker/              # TTS 任务执行器
│   ├── etc/
│   ├── internal/
│   │   ├── config/
│   │   ├── engine/          # TTS 引擎接口和实现
│   │   ├── model/
│   │   ├── utils/
│   │   └── worker/
│   ├── tts-worker.go
│   └── go.mod
│
├── common/                  # 公共模块
│
└── sql/                     # 数据库脚本
    └── init.sql
```

### 服务详解

#### 1. 用户服务 (user-api) - 端口 8081

提供用户认证和基本信息管理。

**API 接口：**

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | /api/user/register | 用户注册 |
| POST | /api/user/login | 用户登录 |
| GET | /api/user/info | 获取用户信息 |

**请求/响应类型：**

```go
// 登录请求
type LoginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// 登录响应
type LoginResp struct {
    Token string       `json:"token"`
    User  UserInfoResp `json:"user"`
}

// 用户信息
type UserInfoResp struct {
    Id             int64   `json:"id"`
    Username       string  `json:"username"`
    Balance        float64 `json:"balance"`
    CharacterCount int64   `json:"characterCount"`
}
```

**核心逻辑：**
- `register.go` - 处理用户注册，创建新用户记录
- `login.go` - 验证用户名密码，生成 JWT Token
- `getuserinfo.go` - 获取当前登录用户信息

#### 2. 音色服务 (voice-api) - 端口 8082

提供音色 CRUD 管理功能。

**API 接口：**

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | /api/voice/list | 获取音色列表 |
| POST | /api/voice/create | 创建音色 |
| DELETE | /api/voice/:id | 删除音色 |
| PUT | /api/voice/default/:id | 设置默认音色 |

**音色数据结构：**

```go
type Voice struct {
    Id          int64  `json:"id"`
    Name        string `json:"name"`
    Tone        string `json:"tone"`
    Gender      string `json:"gender"`
    PreviewUrl  string `json:"previewUrl,omitempty"`
    IsDefault   bool   `json:"isDefault"`
}
```

**核心逻辑：**
- `getvoicelist.go` - 获取所有音色列表
- `createvoice.go` - 创建新音色
- `voiceop.go` - 删除音色、设置默认音色

#### 3. TTS 生成服务 (tts-api) - 端口 8080

提供 TTS 任务提交和状态查询。

**API 接口：**

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | /api/tts/generate | 提交生成任务 |
| GET | /api/tts/task/:taskId | 查询任务状态 |

**请求/响应类型：**

```go
// 文本片段
type Segment struct {
    VoiceId int64  `json:"voiceId"`
    Emotion string `json:"emotion"`
    Text    string `json:"text"`
}

// 生成请求
type GenerateReq struct {
    Segments []Segment `json:"segments"`
    Format   string    `json:"format"`
    Channel  string    `json:"channel"`
}

// 生成响应
type GenerateResp struct {
    TaskId string `json:"taskId"`
}

// 任务状态响应
type TaskResp struct {
    TaskId    string `json:"taskId"`
    Status    string `json:"status"`
    Progress  int    `json:"progress"`
    AudioUrl  string `json:"audioUrl,omitempty"`
    Error     string `json:"error,omitempty"`
}
```

**核心逻辑：**
- `generate.go` - 接收生成请求，创建任务，返回 TaskId
- `querytask.go` - 根据 TaskId 查询任务状态

**任务状态流转：**
```
pending -> processing -> success / failed
```

#### 4. TTS 任务执行器 (tts-worker)

异步执行 TTS 任务的核心 worker，从队列获取任务并调用 TTS 引擎生成音频。

**核心模块：**

| 模块 | 功能说明 |
|------|----------|
| `engine/interface.go` | TTS 引擎接口定义 |
| `engine/vits.go` | VITS 引擎实现 |
| `worker/ttsworker.go` | Worker 主逻辑，任务消费和处理 |
| `model/worker.go` | 任务数据模型 |
| `utils/audio.go` | 音频处理工具 |

**TTS 引擎接口：**

```go
type TTSProvider interface {
    Generate(text string, voiceId int64, emotion string) ([]byte, error)
    GetVoiceName(voiceId int64) string
}
```

**支持的 TTS 引擎：**
- VITS（本地模型）
- 阿里云 TTS
- 腾讯云 TTS

**扩展能力：**
- 情绪合成：emotion 参数透传
- 双声道对话：左右声道分离
- 批量生成：批处理任务
- AI 自动分角色：接入 LLM

#### 5. 数据库 (sql/init.sql)

初始化数据库表结构，包括：
- 用户表
- 音色表
- TTS 任务表
- 作品表

---

## 快速开始

### 前端启动

```bash
cd tts-front

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

### 后端启动

1. 初始化数据库

```bash
mysql -u root -p < tts-backend/sql/init.sql
```

2. 修改配置文件

各服务配置文件在 `tts-backend/*/etc/*.yaml`，请根据实际情况修改：
- MySQL 连接信息
- Redis 连接信息
- OSS 配置（用于存储生成的音频）

3. 启动服务

```bash
# 启动 user-api
cd tts-backend/user-api
go run user.go -f etc/user-api.yaml

# 启动 voice-api
cd tts-backend/voice-api
go run voice.go -f etc/voice-api.yaml

# 启动 tts-api
cd tts-backend/tts-api
go run tts.go -f etc/tts-api.yaml

# 启动 tts-worker
cd tts-backend/tts-worker
go run tts-worker.go -f etc/tts-worker.yaml
```

### 使用 go-zero 工具生成代码

```bash
# 安装 go-zero 工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 生成各服务代码
goctl api go -api api/user.api -dir . -style gozero
goctl api go -api api/voice.api -dir . -style gozero
goctl api go -api api/tts.api -dir . -style gozero
```

---

## 页面功能说明

### 1. 登录页面 (`/login`)
- 用户名/密码登录
- 新用户注册
- Token 存储到 localStorage

### 2. 生成音频页面 (`/generate`)
- 三栏布局：角色面板 | 文本编辑 | 生成面板
- 支持添加多个文本片段
- 每个片段可指定不同音色和情绪
- 拖拽排序文本片段
- 实时显示生成进度
- 生成完成后自动播放

### 3. 我的作品页面 (`/works`)
- 展示历史生成的音频作品列表
- 支持在线播放和下载
- 显示生成时间、状态、格式

### 4. 音色管理页面 (`/voice`)
- 音色列表展示
- 创建新音色（名称、音调、性别、预览音频）
- 删除音色
- 设置默认音色

---

## 许可证

MIT License
