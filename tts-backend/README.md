# tts-backend（go-zero）

后端由多个服务组成：

- `user-api`：用户/登录/权限/余额等
- `voice-api`：音色管理、定制声音等
- `tts-api`：TTS 任务提交、任务状态查询等
- `tts-worker`：异步执行 TTS 任务

## 配置

各服务配置文件位于：

- `user-api/etc/user-api.yaml`
- `voice-api/etc/voice-api.yaml`
- `tts-api/etc/tts-api.yaml`
- `tts-worker/etc/tts-worker.yaml`

请根据本机 MySQL/Redis 地址与端口进行调整。

## 启动（本机开发）

建议先按根目录 `启动指南.md` 初始化数据库后启动：

```bash
cd tts-backend/user-api   && go run . -f etc/user-api.yaml
cd tts-backend/voice-api  && go run . -f etc/voice-api.yaml
cd tts-backend/tts-api    && go run . -f etc/tts-api.yaml
cd tts-backend/tts-worker && go run . -f etc/tts-worker.yaml
```

## 端口（默认）

- `user-api`：8081
- `voice-api`：8082
- `tts-api`：8083

具体以各自 `etc/*.yaml` 为准。
