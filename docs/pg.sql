-- =============================================================================
-- pure-go-admin · 权限模型（PostgreSQL）
-- =============================================================================
--
-- 设计原则：一种资源一种授权方式，运行时只问一处。
--
--   ┌─────────────┬──────────────────────┬─────────────────────────────┐
--   │ 维度        │ 资源表               │ 授权存储                     │
--   ├─────────────┼──────────────────────┼─────────────────────────────┤
--   │ 身份        │ sys_users            │ sys_user_roles              │
--   │ 前端路由/按钮│ sys_menus           │ sys_role_menus              │
--   │ 后端接口    │ —                    │ 登录（JWT）即可访问          │
--   │ 行级数据    │ sys_dept             │ sys_roles.data_scope        │
--   │             │                      │ + sys_role_depts（自定义）  │
--   │ 审计        │ sys_login_log        │ 登录成功/失败记录           │
--   │             │ sys_oper_log         │ 业务操作记录（写操作为主）  │
--   └─────────────┴──────────────────────┴─────────────────────────────┘
--
-- super_admin = true 时拥有全部菜单与按钮权限码。
--
-- menu_type：0 目录 · 1 页面 · 2 按钮（按钮节点的 perm_code 供前端 v-auth 使用）
-- data_scope：1 全部 · 2 自定义部门 · 3 本部门 · 4 本部门及下级 · 5 仅本人
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 组织
-- ---------------------------------------------------------------------------
create table sys_dept
(
    id         bigserial primary key,
    parent_id  bigint references sys_dept on delete set null,
    ancestors  varchar(512) not null default '',
    name       varchar(128) not null,
    sort       int          not null default 0,
    enabled    boolean      not null default true,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    deleted_at timestamptz
);

create index idx_sys_dept_parent_id on sys_dept (parent_id) where deleted_at is null;
create index idx_sys_dept_ancestors on sys_dept (ancestors) where deleted_at is null;

comment on table sys_dept is '部门树';
comment on column sys_dept.ancestors is '物化路径，如 0,1,3；应用层维护';

-- ---------------------------------------------------------------------------
-- 身份
-- ---------------------------------------------------------------------------
create table sys_users
(
    id            bigserial primary key,
    dept_id       bigint references sys_dept on delete set null,
    username      varchar(64)  not null,
    password      varchar(255) not null,
    nickname      varchar(64),
    phone         varchar(32),
    email         varchar(128),
    avatar        varchar(512),
    enabled       boolean      not null default true,
    super_admin   boolean      not null default false,
    last_login_at timestamptz,
    created_at    timestamptz  not null default now(),
    updated_at    timestamptz  not null default now(),
    deleted_at    timestamptz
);

create unique index uk_sys_users_username on sys_users (username) where deleted_at is null;
create index idx_sys_users_dept_id on sys_users (dept_id) where deleted_at is null;

comment on column sys_users.super_admin is '为 true 时拥有全部菜单与按钮权限';

create table sys_roles
(
    id          bigserial primary key,
    name        varchar(128) not null,
    code        varchar(64)  not null,
    sort        int          not null default 0,
    enabled     boolean      not null default true,
    data_scope  smallint     not null default 5,
    description varchar(512),
    created_at  timestamptz  not null default now(),
    updated_at  timestamptz  not null default now(),
    deleted_at  timestamptz,
    constraint chk_sys_roles_data_scope check (data_scope between 1 and 5)
);

create unique index uk_sys_roles_code on sys_roles (code) where deleted_at is null;

comment on column sys_roles.code is '角色唯一编码';
comment on column sys_roles.data_scope is '1 全部 2 自定义部门 3 本部门 4 本部门及下级 5 仅本人';

create table sys_user_roles
(
    user_id bigint not null references sys_users on delete cascade,
    role_id bigint not null references sys_roles on delete cascade,
    primary key (user_id, role_id)
);

create index idx_sys_user_roles_role_id on sys_user_roles (role_id);

create table sys_role_depts
(
    role_id bigint not null references sys_roles on delete cascade,
    dept_id bigint not null references sys_dept on delete cascade,
    primary key (role_id, dept_id)
);

create index idx_sys_role_depts_dept_id on sys_role_depts (dept_id);

comment on table sys_role_depts is 'data_scope = 2 时生效';

-- ---------------------------------------------------------------------------
-- 前端权限：菜单树（目录 / 页面 / 按钮）
-- ---------------------------------------------------------------------------
create table sys_menus
(
    id         bigserial primary key,
    parent_id  bigint references sys_menus on delete set null,
    menu_type  smallint     not null default 1,
    perm_code  varchar(128),                -- menu_type = 2 时必填，如 system:user:add
    path       varchar(255),
    name       varchar(128),
    component  varchar(255),
    redirect   varchar(255),
    title      varchar(128) not null,
    icon       varchar(128),
    sort       int          not null default 0,
    hidden     boolean      not null default false,
    keep_alive boolean      not null default false,
    enabled    boolean      not null default true,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    deleted_at timestamptz,
    constraint chk_sys_menus_menu_type check (menu_type in (0, 1, 2)),
    constraint chk_sys_menus_button_perm check (
        menu_type <> 2 or (perm_code is not null and perm_code <> '')
        )
);

create index idx_sys_menus_parent_id on sys_menus (parent_id) where deleted_at is null;
create unique index uk_sys_menus_perm_code on sys_menus (perm_code) where deleted_at is null and menu_type = 2;

comment on column sys_menus.menu_type is '0 目录 1 页面 2 按钮';
comment on column sys_menus.perm_code is '按钮权限标识；页面/目录可为空';

create table sys_role_menus
(
    role_id bigint not null references sys_roles on delete cascade,
    menu_id bigint not null references sys_menus on delete cascade,
    primary key (role_id, menu_id)
);

create index idx_sys_role_menus_menu_id on sys_role_menus (menu_id);

-- =============================================================================
-- 系统支撑表
-- =============================================================================

create table sys_dict
(
    id          bigserial primary key,
    name        varchar(128) not null,
    type        varchar(64)  not null,
    description varchar(512),
    enabled     boolean      not null default true,
    created_at  timestamptz  not null default now(),
    updated_at  timestamptz  not null default now(),
    deleted_at  timestamptz
);

create unique index uk_sys_dict_type on sys_dict (type) where deleted_at is null;

create table sys_dict_detail
(
    id         bigserial primary key,
    dict_id    bigint       not null references sys_dict on delete cascade,
    label      varchar(128) not null,
    value      varchar(128) not null,
    sort       int          not null default 0,
    enabled    boolean      not null default true,
    css_class  varchar(128),
    list_class varchar(128),
    is_default boolean      not null default false,
    remark     varchar(512),
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    deleted_at timestamptz
);

create unique index uk_sys_dict_detail_dict_value on sys_dict_detail (dict_id, value) where deleted_at is null;
create index idx_sys_dict_detail_dict_id on sys_dict_detail (dict_id) where deleted_at is null;

-- ---------------------------------------------------------------------------
-- 审计日志（只追加，不做软删）
-- ---------------------------------------------------------------------------

-- 登录日志：记录每次登录尝试（成功与失败）
create table sys_login_log
(
    id         bigserial primary key,
    user_id    bigint references sys_users on delete set null,
    username   varchar(64)  not null,
    nickname   varchar(64),
    status     smallint     not null default 0,
    message    varchar(512) not null default '',
    ip         varchar(64),
    location   varchar(255),
    user_agent varchar(512),
    created_at timestamptz  not null default now(),
    constraint chk_sys_login_log_status check (status in (0, 1))
);

create index idx_sys_login_log_user_id on sys_login_log (user_id);
create index idx_sys_login_log_username on sys_login_log (username);
create index idx_sys_login_log_status on sys_login_log (status);
create index idx_sys_login_log_ip on sys_login_log (ip);
create index idx_sys_login_log_created_at on sys_login_log (created_at desc);

comment on table sys_login_log is '登录日志';
comment on column sys_login_log.user_id is '登录成功且用户存在时写入；失败可为空';
comment on column sys_login_log.username is '登录时输入的用户名（快照）';
comment on column sys_login_log.status is '0 失败 · 1 成功';
comment on column sys_login_log.message is '失败原因或成功说明';
comment on column sys_login_log.location is 'IP 归属地，可由应用层解析后写入';

-- 操作日志：记录后台写操作与关键读操作（由应用层埋点写入）
create table sys_oper_log
(
    id            bigserial primary key,
    user_id       bigint references sys_users on delete set null,
    username      varchar(64),
    module        varchar(128) not null default '',
    action        varchar(128) not null default '',
    method        varchar(16),
    path          varchar(255),
    business_type varchar(64),
    business_id   bigint,
    status        smallint     not null default 1,
    error_msg     varchar(512) not null default '',
    request_body  text,
    response_body text,
    duration_ms   int          not null default 0,
    ip            varchar(64),
    user_agent    varchar(512),
    created_at    timestamptz  not null default now(),
    constraint chk_sys_oper_log_status check (status in (0, 1))
);

create index idx_sys_oper_log_user_id on sys_oper_log (user_id);
create index idx_sys_oper_log_module on sys_oper_log (module);
create index idx_sys_oper_log_status on sys_oper_log (status);
create index idx_sys_oper_log_business on sys_oper_log (business_type, business_id);
create index idx_sys_oper_log_created_at on sys_oper_log (created_at desc);

comment on table sys_oper_log is '操作日志';
comment on column sys_oper_log.module is '功能模块，如 用户管理、角色管理';
comment on column sys_oper_log.action is '操作名称，如 新增用户、分配菜单';
comment on column sys_oper_log.method is 'HTTP 方法';
comment on column sys_oper_log.path is '请求路径（如 /api/v1/system/users）';
comment on column sys_oper_log.business_type is '业务类型标识，如 user、role';
comment on column sys_oper_log.business_id is '业务主键，便于按单条数据追溯';
comment on column sys_oper_log.status is '0 失败 · 1 成功';
comment on column sys_oper_log.request_body is '请求参数 JSON 文本，应用层应截断脱敏';
comment on column sys_oper_log.response_body is '响应摘要，可选；大响应不建议完整入库';
comment on column sys_oper_log.duration_ms is '接口耗时（毫秒）';
