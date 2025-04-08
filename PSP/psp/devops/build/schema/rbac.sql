-- resource
CREATE TABLE IF NOT EXISTS `resource` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键' primary key,
    `name` varchar(128) CHARACTER SET utf8 NOT NULL COMMENT '资源名: 如果是接口则为url',
    `action` varchar(8) CHARACTER SET utf8 NOT NULL COMMENT '资源操作方式 GET、POST、PUT、DELETE、NONE 默认NONE',
    `type` varchar(64) CHARACTER SET utf8 NOT NULL COMMENT 'system-菜单 job_sub_app-求解应用 remote_app-远程应用 api-接口 internal ',
    `display_name` varchar(128) CHARACTER SET utf8 DEFAULT NULL COMMENT '前端展示用',
    `custom` tinyint(1) NOT NULL COMMENT '是否可自定义 -1-false 1-true',
    `external_id` bigint(20) DEFAULT NULL COMMENT '外部id: 部分权限需要关联到其他表的主键id',
    `parent_id` bigint(20) DEFAULT NULL COMMENT '权限父级id'
) ENGINE = InnoDB DEFAULT CHARSET = utf8 comment '资源表' AUTO_INCREMENT = 1000;

-- role auto-generated definition
CREATE TABLE IF NOT EXISTS role (
    id bigint auto_increment primary key COMMENT '主键',
    name varchar(64) not null COMMENT '名称',
    comment varchar(256) not null COMMENT '描述',
    type tinyint null COMMENT '类型 1-超级管理员 0-自定义角色',
    constraint uniq_name unique (name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 comment '角色表' AUTO_INCREMENT = 100;

-- casbin_rule auto-generated definition
CREATE TABLE IF NOT EXISTS `casbin_rule` (
    `p_type` varchar(100) NOT NULL DEFAULT '',
    `v0` varchar(100) NOT NULL DEFAULT '',
    `v1` varchar(100) NOT NULL DEFAULT '',
    `v2` varchar(100) NOT NULL DEFAULT '',
    `v3` varchar(100) NOT NULL DEFAULT '',
    `v4` varchar(100) NOT NULL DEFAULT '',
    `v5` varchar(100) NOT NULL DEFAULT '',
    KEY `IDX_casbin_rule_v2` (`v2`),
    KEY `IDX_casbin_rule_v3` (`v3`),
    KEY `IDX_casbin_rule_v4` (`v4`),
    KEY `IDX_casbin_rule_v5` (`v5`),
    KEY `IDX_casbin_rule_p_type` (`p_type`),
    KEY `IDX_casbin_rule_v0` (`v0`),
    KEY `IDX_casbin_rule_v1` (`v1`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;