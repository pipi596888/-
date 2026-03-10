# TTS 语音合成系统

一个前后端分离的文本转语音（TTS）项目，包含生成音频、作品管理、音色管理、定制声音、意见反馈与系统管理等功能。

## 技术栈

- 前端：Vue 3 + TypeScript + Vite + Element Plus + Pinia + Vue Router
- 后端：Go + go-zero（多服务）
- 基础设施：MySQL、Redis（建议用 Docker 容器）

## 目录结构

```
tts-front/           # 前端
tts-backend/         # 后端（user-api / voice-api / tts-api / tts-worker）
start.bat            # Windows 一键启动（本机开发）
启动指南.md           # 更详细的启动/初始化/排错
docker-compose.yaml  # DBHub（可选，用于 Database MCP）
```

## 快速开始（Windows 推荐）

1) 安装依赖：Go（建议 1.21+）、Node.js（建议 18+）、Docker Desktop。  
2) 确保本机存在并运行名为 `mysql`、`redis` 的容器（脚本会执行 `docker start mysql redis`）。  
3) 首次运行先初始化数据库（见下节）。  
4) 双击运行 `start.bat`，打开：`http://localhost:3000`

## 初始化数据库（仅首次）

Windows（PowerShell/CMD）：

```bash
type tts-backend\sql\init.sql | docker exec -i mysql mysql -uroot -proot
```

macOS/Linux：

```bash
cat tts-backend/sql/init.sql | docker exec -i mysql mysql -uroot -proot
```

## 手动启动（便于调试）

后端（4 个终端）：

```bash
cd tts-backend/user-api   && go run . -f etc/user-api.yaml
cd tts-backend/voice-api  && go run . -f etc/voice-api.yaml
cd tts-backend/tts-api    && go run . -f etc/tts-api.yaml
cd tts-backend/tts-worker && go run . -f etc/tts-worker.yaml
```

前端：

```bash
cd tts-front
npm install
npm run dev
```

## 端口

- 前端：`http://localhost:3000`
- user-api：`http://localhost:8081`
- voice-api：`http://localhost:8082`
- tts-api：`http://localhost:8083`
- MySQL：`localhost:3306`（默认 root/root）
- Redis：`localhost:6379`

## 安全提示（重要）

- **生产环境必须使用 HTTPS**：如果用 `http://` 访问前端/接口，账号密码会以明文出现在浏览器 Network 请求里，并且传输层也无法防止被中间人窃听。
- 浏览器里“能看到自己输入/提交的密码”属于正常现象（客户端必须把凭证提交给服务端验证）。要保护传输安全，请在部署时使用 Nginx/网关为站点配置 HTTPS，再反向代理到 `user-api/voice-api/tts-api`。

## 常见问题

- 端口被占用（Windows）：

```bash
netstat -ano | findstr :3000
taskkill /F /PID <PID>
```

- `npm run build` 报 `spawn EPERM`：多见于安全软件/权限拦截 `esbuild` 子进程；将项目目录加入白名单或以管理员权限运行终端再试。

更多见 `启动指南.md`。
