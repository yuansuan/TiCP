CREATE TABLE `account`
(
    `id`               bigint(20) unsigned NOT NULL,
    `customer_id`      bigint(20) DEFAULT NULL,
    `real_customer_id` bigint(20) DEFAULT NULL,
    `name`             varchar(255) DEFAULT NULL,
    `currency`         varchar(8)   DEFAULT NULL,
    `account_balance`  bigint(20) DEFAULT NULL,
    `freezed_amount`   bigint(20) DEFAULT NULL,
    `normal_balance`   bigint(20) DEFAULT NULL,
    `award_balance`    bigint(20) DEFAULT NULL,
    `withdraw_enabled` tinyint(4) DEFAULT NULL,
    `credit_quota`     bigint(20) DEFAULT NULL,
    `status`           tinyint(4) DEFAULT NULL,
    `is_freeze`        tinyint(4) DEFAULT NULL,
    `account_type`     int(11) DEFAULT NULL,
    `create_time`      datetime     DEFAULT NULL,
    `update_time`      datetime     DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `account_bill`
(
    `id`                    bigint(20) unsigned NOT NULL,
    `account_id`            bigint(20) unsigned DEFAULT NULL,
    `sign`                  tinyint(4) DEFAULT NULL,
    `amount`                bigint(20) DEFAULT NULL,
    `trade_type`            int(11) DEFAULT NULL,
    `trade_id`              varchar(128) NOT NULL DEFAULT '' COMMENT '交易单ID',
    `idempotent_id`         varchar(255) NOT NULL DEFAULT '' COMMENT '幂等ID',
    `account_balance`       bigint(20) DEFAULT NULL,
    `voucher_balance`       bigint(20) NOT NULL DEFAULT '0' COMMENT '代金券余额',
    `freezed_amount`        bigint(20) DEFAULT NULL,
    `delta_normal_balance`  bigint(20) DEFAULT NULL,
    `delta_award_balance`   bigint(20) DEFAULT NULL,
    `delta_voucher_balance` bigint(20) NOT NULL DEFAULT '0' COMMENT '代金券消费金额',
    `comment`               varchar(255)          DEFAULT NULL,
    `create_time`           datetime              DEFAULT NULL,
    `update_time`           datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `out_trade_id`          bigint(20) unsigned DEFAULT NULL,
    `account_voucher_ids`   varchar(512)          DEFAULT '' COMMENT '账户代金券关联ids',
    `merchandise_id`        varchar(128) NOT NULL DEFAULT '' COMMENT '商品id',
    `unit_price`            bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '单价',
    `merchandise_name`      varchar(255) NOT NULL DEFAULT '' COMMENT '商品名称',
    `product_name`          varchar(64)  NOT NULL DEFAULT '' COMMENT '产品类型 pass求解作业: CloudCompute，云应用: CloudApp',
    `price_des`             varchar(255) NOT NULL DEFAULT '' COMMENT '单价描述',
    `quantity` double NOT NULL DEFAULT '0' COMMENT '消耗数量',
    `quantity_unit`         varchar(255) NOT NULL DEFAULT '' COMMENT '消耗数量单位描述',
    `resource_id`           varchar(128) NOT NULL DEFAULT '' COMMENT '资源id',
    `start_time`            datetime              DEFAULT '1970-01-01 00:00:00' COMMENT '扣费周期开始时间，按量付费使用',
    `end_time`              datetime              DEFAULT '1970-01-01 00:00:00' COMMENT '扣费周期结束时间，按量付费使用',
    PRIMARY KEY (`id`),
    UNIQUE KEY `account_bill__uindex_idempotent` (`account_id`,`idempotent_id`) COMMENT '幂等性'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `account_cash_voucher_relation`
(
    `id`                  bigint(20) NOT NULL COMMENT '主键',
    `account_id`          bigint(20) DEFAULT NULL COMMENT '账户ID',
    `cash_voucher_id`     bigint(20) unsigned NOT NULL COMMENT '优惠券id',
    `cash_voucher_amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '代金券原始总金额',
    `used_amount`         bigint(20) NOT NULL DEFAULT '0' COMMENT '已使用金额',
    `remaining_amount`    bigint(20) NOT NULL DEFAULT '0' COMMENT '剩余金额',
    `status`              int(11) NOT NULL DEFAULT '0' COMMENT '账户代金券状态: 0:正常，1:禁用',
    `expired_time`        datetime          DEFAULT NULL,
    `is_expired`          int(11) DEFAULT '0' COMMENT '是否过期 0:正常 1:过期',
    `is_deleted`          tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标记 0:正常 1:删除 ',
    `opt_user_id`         bigint(20) NOT NULL DEFAULT '0' COMMENT '操作用户id',
    `create_time`         datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`         datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY                   `account_id_index` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代金券账户关系表';

CREATE TABLE `account_log`
(
    `id`           bigint(20) unsigned DEFAULT NULL,
    `account_id`   bigint(20) unsigned DEFAULT NULL,
    `operator_uid` bigint(20) unsigned DEFAULT NULL,
    `params`       mediumtext,
    `old`          mediumtext,
    `updated`      mediumtext,
    `create_time`  datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `cash_voucher`
(
    `id`                  bigint(20) NOT NULL COMMENT '主键',
    `name`                varchar(255) NOT NULL COMMENT '代金券名称',
    `amount`              bigint(20) NOT NULL DEFAULT '0' COMMENT '代金券金额',
    `availability_status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '上下架状态 0:下架 1:上架',
    `opt_user_id`         bigint(20) DEFAULT NULL COMMENT '业务平台操作用户id',
    `is_expired`          tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否过期 0:正常 1:过期',
    `abs_expired_time`    datetime              DEFAULT NULL COMMENT '绝对过期时间',
    `rel_expired_time`    bigint(20) DEFAULT '0' COMMENT '相对过期时间，以秒来计算',
    `expired_type`        tinyint(4) NOT NULL DEFAULT '0' COMMENT '过期类型 1:绝对 2:相对',
    `comment`             varchar(255)          DEFAULT NULL COMMENT '备注',
    `is_deleted`          tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标记',
    `create_time`         datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`         datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代金券基础信息表';

CREATE TABLE `account_cash_voucher_log`
(
    `id`                      bigint(20) NOT NULL,
    `account_id`              bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '账户id',
    `cash_voucher_id`         bigint(20) NOT NULL DEFAULT '0' COMMENT '代金券id',
    `account_cash_voucher_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户代金券id',
    `sign_type`               tinyint(4) NOT NULL COMMENT '使用标记 1:消费 2:过期 ',
    `opt_user_id`             bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '平台用户id',
    `account_bill_id`         bigint(20) unsigned DEFAULT '0' COMMENT '消费账单记录',
    `source_info`             text COMMENT '账户代金券修改前信息',
    `target_info`             text COMMENT '账户代金券修改后信息',
    `comment`                 varchar(255)      DEFAULT NULL COMMENT '使用备注',
    `create_time`             datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`             datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代金券使用记录表';
