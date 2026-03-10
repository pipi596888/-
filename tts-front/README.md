# tts-front（Vue 3 + Vite）

前端包含以下页面：

- `/login`：登录/注册
- `/generate`：生成音频（支持多角色/片段）
- `/works`：我的作品（播放/下载/重命名等）
- `/voice`：音色管理（试听/设默认/删除等）
- `/custom-voice`：定制声音
- `/feedback`：意见反馈
- `/system`：系统管理（管理员）

## 开发启动

```bash
npm install
npm run dev
```

默认访问：`http://localhost:3000`

## 构建

```bash
npm run build
```

若遇到 `spawn EPERM`（Windows），通常是安全软件拦截 `esbuild` 子进程或权限问题；建议将项目目录加入白名单或以管理员权限运行终端。
