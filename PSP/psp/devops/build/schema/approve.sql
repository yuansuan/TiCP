-- audit_log
CREATE TABLE IF NOT EXISTS `audit_log`
(
    `id`                bigint(20)   NOT NULL,
    `user_id`           bigint(20)   NOT NULL,
    `user_name`         varchar(128) NOT NULL COMMENT '操作用户',
    `ip_address`        varchar(15)  NOT NULL COMMENT '操作用户',
    `operate_type`      varchar(20)  NOT NULL COMMENT '操作类型',
    `operate_user_type` tinyint(1)   NOT NULL COMMENT '操作用户类型 1:普通用户 2:系统管理员 3:安全管理员',
    `operate_content`   text         NOT NULL COMMENT '操作内容',
    `operate_time`      datetime DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '操作日志';

-- approve_record
CREATE TABLE IF NOT EXISTS `approve_record`
(
    `id`              bigint(20)   NOT NULL,
    `type`            tinyint(4)   NOT NULL COMMENT '审批类型 新增用户等等',
    `approve_info`    text         NOT NULL COMMENT '审批信息 json格式，包含请求入参等',
    `status`          tinyint(1)   NOT NULL COMMENT '审批状态 1:等待审批 2:审批通过 3:审批拒绝 4:审批撤销 5:审批失败',
    `apply_user_id`   bigint(20)   NOT NULL COMMENT '审批发起人ID',
    `apply_user_name` varchar(128) NOT NULL COMMENT '审批发起人名称',
    `sign`            varchar(50)  NOT NULL COMMENT '审批唯一签名，如修改/删除同一用户的，状态为审批中的审批只能有一份',
    `content`         text COMMENT '显示文案',
    `create_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`     datetime              DEFAULT NULL COMMENT '修改时间',
    `approve_time`    datetime              DEFAULT NULL COMMENT '审批时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '审批记录表';

-- approve_user
CREATE TABLE IF NOT EXISTS `approve_user`
(
    `id`                bigint(20)   NOT NULL,
    `approve_record_id` bigint(20)   NOT NULL COMMENT '审批记录id',
    `approve_user_id`   bigint(20)   NOT NULL COMMENT '审批人(安全管理员)id',
    `approve_user_name` varchar(128) NOT NULL COMMENT '审批人名称',
    `result`            tinyint(1)   NOT NULL COMMENT '审批结果 0:默认状态 1:审批通过 2:审批拒绝',
    `suggest`           text COMMENT '审批意见(备注)',
    `create_time`       datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`       datetime              DEFAULT NULL COMMENT '修改时间',
    `approve_time`      datetime              DEFAULT NULL COMMENT '审批时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '审批人关联审批表';