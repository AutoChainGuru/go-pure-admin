# go-pure-admin

一款**轻量、可二次开发**的后台权限管理系统：后端用 Go 提供 REST API，前端用 Vue 3 实现动态菜单与页面权限，适合作为企业内部系统、SaaS 管理端或学习「RBAC + 前后端分离」的脚手架。

设计取向：**权限模型清晰、依赖克制、代码易读**——菜单与按钮权限在前端统一管控，接口层以 JWT 登录为准；不引入 Casbin 等额外授权引擎，降低上手与维护成本。

---

## 功能概览

| 模块 | 说明 |
|------|------|
| 认证 | 登录 / 登出、JWT、个人信息与修改密码 |
| 用户管理 | 用户、部门、角色关联；超级管理员标识 |
| 角色管理 | 角色 CRUD、菜单授权、数据范围（`data_scope`，预留行级过滤） |
| 菜单管理 | 目录 / 页面 / 按钮三级；驱动动态路由与 `v-auth` 按钮显隐 |
| 部门管理 | 树形组织架构 |
| 字典管理 | 字典类型 + 字典项 |
| 审计 | 登录日志、操作日志（写操作自动记录） |
| 前端体验 | 侧栏 + 顶栏 + 多标签页、页面缓存（`keep_alive`）、浅色/深色主题、页面刷新 |

---

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | [Vite](https://vite.dev/) 8、[Vue 3](https://vuejs.org/)、JavaScript、[Element Plus](https://element-plus.org/)、[UnoCSS](https://unocss.dev/)、[Pinia](https://pinia.vuejs.org/)、[Vue Router](https://router.vuejs.org/)、Axios |
| 后端 | [Go](https://go.dev/) 1.25+、[Gin](https://gin-gonic.com/) v1.12、[GORM](https://gorm.io/)、[Viper](https://github.com/spf13/viper)、[lo](https://github.com/samber/lo)、JWT、[Zap](https://github.com/uber-go/zap) + Lumberjack 日志滚动 |
| 数据库 | PostgreSQL |

---

## 环境要求

- Go **1.25+**
- Node.js **18+**（推荐 20+）
- PostgreSQL **12+**

---

## 快速开始

### 1. 准备数据库

创建空库（库名可与配置一致）：

```bash
createdb go_pure_admin
```

> 表结构由后端 **GORM AutoMigrate** 自动创建；`docs/pg.sql` 为权限模型说明与参考 DDL，可按需手工执行，非必须。

### 2. 配置后端

编辑 `server/config.yaml`（或通过环境变量 `CONFIG_PATH` 指定其它配置文件路径）：

```yaml
database:
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: your_password
  dbname: go_pure_admin
  sslmode: disable

server:
  port: 9999

jwt:
  secret: change-me-in-production
  expire_hours: 72

cors:
  allow_origins:
    - http://127.0.0.1:8000
```

### 3. 启动后端

在仓库根目录或 `server` 目录下执行：

```bash
cd server
go mod download
go run .
# 监听 http://127.0.0.1:9999
```

首次启动会：

1. 自动迁移表结构；
2. 在**无用户数据**时写入种子数据（默认管理员、菜单、示例字典等）。

**默认账号**：`admin` / `admin123`（生产环境请立即修改）

### 4. 启动前端

```bash
cd web
npm install   # 或使用 pnpm install
npm run dev
# 默认 http://127.0.0.1:8000
```

开发环境下，`/api` 由 Vite 代理到 `http://127.0.0.1:9999`（见 `web/vite.config.js`）。

浏览器访问前端地址，使用默认账号登录即可。

---

## 生产部署（简要）

**后端**

```bash
cd server
go build -o bin/server .
CONFIG_PATH=./config.yaml ./bin/server
```

**前端**

```bash
cd web
npm run build
# 将 web/dist 交由 Nginx 等静态托管，并将 /api 反向代理到后端
```

部署时请修改 `jwt.secret`、数据库密码，并在 `cors.allow_origins` 中加入实际前端域名。

---

## 目录结构

```
go-pure-admin/
├── docs/              # 数据库参考脚本与权限模型说明（pg.sql）
├── server/            # Go 后端
│   ├── config.yaml    # 默认配置
│   ├── main.go
│   └── internal/      # config、handler、service、model、middleware 等
├── web/               # Vue 前端
│   └── src/
│       ├── api/       # 接口封装
│       ├── views/     # 页面
│       ├── layouts/   # 布局
│       ├── stores/    # Pinia 状态
│       └── router/    # 路由与动态路由
└── readme.md
```

---

## 权限模型

权限按维度拆分，运行时各问一处，避免重复配置：

| 维度 | 存储 | 说明 |
|------|------|------|
| 身份 | `sys_users` + `sys_user_roles` | 用户绑定角色 |
| 菜单 / 按钮 | `sys_menus` + `sys_role_menus` | `menu_type`：0 目录 · 1 页面 · 2 按钮；按钮的 `perm_code` 供前端 `v-auth` |
| 接口 | JWT | 登录后可访问已挂载的业务接口，不做 Casbin 细粒度校验 |
| 数据范围 | `sys_roles.data_scope` + `sys_role_depts` | 1 全部 · 2 自定义部门 · 3 本部门 · 4 本部门及下级 · 5 仅本人（业务查询可按需接入） |
| 审计 | `sys_login_log` / `sys_oper_log` | 登录与写操作留痕 |

`super_admin = true` 的用户拥有全部菜单与按钮权限码。

更完整的表结构与设计说明见 [`docs/pg.sql`](docs/pg.sql)。

---

## 前端约定

- **动态路由**：登录后根据菜单树生成路由；目录项 `component = Layout`，页面项指向 `views` 下组件路径。
- **页面缓存**：菜单开启 `keep_alive` 后，组件 `name` 需与路由 `name` 一致（如 `SystemUser`）。
- **按钮权限**：`v-auth="'system:user:add'"`，无权限时隐藏（不删 DOM，避免表格异常）。

---

## 常见问题

**Q：修改菜单或 `keep_alive` 后路由未更新？**  
A：重新登录，以重新拉取菜单并注册动态路由。

**Q：前端请求 401？**  
A：确认后端已启动、代理目标端口为 `9999`，且 token 未过期。

**Q：只想重置业务数据？**  
A：清空相关表后重启后端；若 `sys_users` 为空会再次执行种子。生产环境请用备份与迁移流程，勿直接删库。

---

## 开源协议

见仓库根目录 [LICENSE](LICENSE)。
