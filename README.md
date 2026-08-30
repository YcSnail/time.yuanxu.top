# ⏳ time.yuanxu.top · 我的倒计时

移动端优先的倒计时展示站:用户用「自己的密码」作为身份,自定义精确到秒的目标时间与标题,
主页面按**最近优先**展示,毫秒级滚动倒计时。

## 功能

- **密码即身份**:进入时输入自己的密码;同一密码 = 同一账号,首次输入自动创建,之后直接登录
- **密码复杂度校验**:必须同时包含大小写字母和数字、长度 ≥ 6,否则拒绝创建(前后端双重校验)
- **倒计时管理**:自定义标题 + 精确到秒的目标时间,可删除
- **最近优先**:未到期的倒计时按目标时间升序(最近的排最前),已结束的沉底置灰
- **毫秒滚动**:秒以下三位毫秒实时滚动(60fps requestAnimationFrame)
- **手机端预览**:430px 手机壳布局,桌面居中、真机全宽

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go + Gin + GORM,MySQL 5.6,BCrypt + JWT |
| 前端 | Vue 3 + Vite + Vue Router,pnpm |
| 部署 | GitHub Actions → 阿里云 ACR 镜像 → 服务器 docker compose + nginx,HTTPS(acme.sh) |

## 目录结构

```
time.yuanxu.top/
├── backend/            # Gin 后端
│   ├── main.go         # 入口:路由 / CORS / 端口
│   ├── config/         # 环境配置(DSN、JWT)
│   ├── models/         # User / Countdown
│   ├── handlers/       # 进入(注册/登录)、倒计时 CRUD
│   ├── middleware/     # JWT 鉴权
│   ├── utils/          # 密码复杂度校验
│   └── Dockerfile.prod
├── frontend/           # Vue 3 前端
│   ├── src/views/      # EnterView / HomeView / CreateView
│   ├── src/components/ # CountdownCard(毫秒滚动)
│   └── Dockerfile.prod
├── deploy/             # docker-compose.prod.yml、nginx vhost
└── .github/workflows/  # Deploy:构建镜像 → 推 ACR → 服务器部署
```

## 本地开发

```bash
# 后端(需要 .env,连服务器 MySQL 5.6)
cd backend && go run .

# 前端
cd frontend && pnpm install && pnpm dev
# 开发服务器 http://localhost:5173,/api 代理到 :8080
```

## API

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/enter` | 密码进入(注册或登录),弱密码拒绝创建 |
| GET | `/api/me` | 当前用户 |
| GET | `/api/countdowns` | 我的倒计时,按目标时间升序 |
| POST | `/api/countdowns` | 创建 `{title, target_time}`(RFC3339) |
| DELETE | `/api/countdowns/:id` | 删除自己的倒计时 |

## 生产部署

推送 `main` 分支即触发 GitHub Actions(需配置仓库 Secrets/Vars):

- `vars.ACR_NAMESPACE` = `yc_snail`
- `secrets.ACR_USERNAME` / `secrets.ACR_PASSWORD` — 阿里云 ACR 登录
- `secrets.SSH_HOST` / `secrets.SSH_USER` / `secrets.SSH_PRIVATE_KEY` — 部署服务器

服务器约定:MySQL 5.6 在宿主机 `127.0.0.1:3306`,后端容器 host 网络占 **8095**,
前端容器占 **8082**,宿主机 nginx 统一反代(证书在 `/usr/local/nginx/conf/ssl/time.yuanxu.top/`)。
