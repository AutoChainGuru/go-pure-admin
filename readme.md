# go-pure-admin

一款纯粹的后管权限框架。

## 技术栈

| 端 | 技术 |
|----|------|
| 前端 | Vite 8、Vue 3、JavaScript、Element Plus、UnoCSS、Pinia、Vue Router |
| 后端 | Go 1.25+、Gin v1.12、GORM、Viper、lo、JWT、Zap + Lumberjack 按天滚动日志 |

## 快速开始

### 1. 数据库

```bash
createdb go_pure_admin
psql -U postgres -d go_pure_admin -f pg.sql
```

连接信息见 `server/config.yaml`（与 readme 中数据库说明一致）。

### 2. 后端

```bash
cd server
go run .
# 默认 http://127.0.0.1:9999
```

首次启动会自动迁移表结构并写入种子数据。

默认账号：**admin / admin123**

### 3. 前端

```bash
cd web
npm install
npm run dev
# 默认 http://127.0.0.1:8000 ，API 代理至后端
```

## 目录结构

```
├── docs          # 数据库 DDL
├── server/          # Gin 后端
└── web/             # Vue 前端
```

## 权限说明

- **前端菜单/按钮**：`sys_menus` + `sys_role_menus`，按钮使用 `menu_type=2` 的 `perm_code`，指令 `v-auth`
- **接口**：登录（JWT）即可访问已挂载的业务接口，不再做 Casbin 细粒度校验
- **数据**：`sys_roles.data_scope` + `sys_role_depts`（预留）
- **超级管理员**：`super_admin=true` 拥有全部菜单与按钮权限码

## 页面风格

现代简洁布局：侧栏 + 顶栏 + 卡片内容区，样式优先使用 UnoCSS 工具类。
