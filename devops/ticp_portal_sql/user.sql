-- user
CREATE TABLE IF NOT EXISTS `user` (
    `id` bigint(20) NOT NULL COMMENT '用户snowflake id' primary key,
    `name` varchar(128) NOT NULL COMMENT '用户名',
    `real_name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户姓名',
    `approve_status` tinyint NOT NULL DEFAULT -1 COMMENT '审批状态: 0-审批中; 1-通过; 2-拒绝',
    `mobile` varchar(16) DEFAULT NULL,
    `email` varchar(128) DEFAULT NULL,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `password` varchar(64) DEFAULT NULL COMMENT '加密密码',
    `is_internal` tinyint(4) DEFAULT NULL COMMENT '是否内部用户 1-是 0-否',
    `enabled` tinyint(4) NOT NULL COMMENT '是否启用 1-启用 0-禁用',
    `enable_openapi` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否允许调用openapi 1-启用 0-禁用',
    `is_deleted` datetime DEFAULT NULL COMMENT '逻辑删字段'
) ENGINE = InnoDB DEFAULT CHARSET = utf8 comment '用户表';

-- org_structure
CREATE TABLE IF NOT EXISTS `org_structure` (
    `id` bigint(20) not null comment '组织架构id' primary key,
    `name` varchar(128) not null comment '名称',
    `type` tinyint(1) null comment '类型 1-部门',
    `remark` varchar(255) null comment '描述',
    `parent_id` bigint(20) null comment '上级组织架构id: 默认为0'
) ENGINE = InnoDB DEFAULT CHARSET = utf8 comment '组织架构表';

-- org_user_relation
CREATE TABLE IF NOT EXISTS `org_user_relation` (
    `id` bigint(20) not null comment '主键id' primary key,
    `org_id` bigint(20) not null comment '组织架构id',
    `user_id` bigint(20) not null comment '用户id'
) ENGINE = InnoDB DEFAULT CHARSET = utf8 comment '组织架构与用户关系表';