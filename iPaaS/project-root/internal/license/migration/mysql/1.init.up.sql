CREATE TABLE `license_info`
(
    `id`                      bigint(20) NOT NULL,
    `manager_id`              bigint(20) NOT NULL COMMENT 'manager id',
    `provider`                varchar(255) NOT NULL COMMENT '提供者',
    `license_server`          varchar(255) NOT NULL COMMENT '许可证变量',
    `mac_addr`                varchar(255) NOT NULL COMMENT 'Mac地址',
    `license_url`             varchar(255) NOT NULL COMMENT '许可证服务器',
    `license_port`            int(11) NOT NULL DEFAULT '0' COMMENT '端口',
    `license_proxies`         varchar(255) NOT NULL DEFAULT '' COMMENT 'HpcEndpoint对应的许可证服务器地址',
    `license_num`             varchar(255) NOT NULL COMMENT 'licenses许可证序列号',
    `weight`                  int(11) NOT NULL DEFAULT '0' COMMENT '调度优先级',
    `begin_time`              datetime     NOT NULL COMMENT '使用有效期 开始',
    `end_time`                datetime     NOT NULL COMMENT '使用有效期 结束',
    `auth`                    tinyint(1) NOT NULL DEFAULT '2' COMMENT '是否被授权 1-授权 2-未授权',
    `license_type`            tinyint(1) NOT NULL DEFAULT '0' COMMENT '供应商类型: 1-自有，2-外部，3-寄售',
    `tool_path`               varchar(255) NOT NULL COMMENT 'lmutil所在路径',
    `collector_type`          varchar(255) NOT NULL DEFAULT '' COMMENT '收集器类型: flex, lsdyna, altair, dsli',
    `hpc_endpoint`            varchar(255) NOT NULL DEFAULT '' COMMENT '超算endpoint',
    `allowable_hpc_endpoints` varchar(255) NOT NULL DEFAULT '' COMMENT '支持的HpcEndpoint范围',
    `license_server_status`   varchar(64)  NOT NULL DEFAULT 'abnormal' COMMENT 'license服务状态: normal-正常，abnormal-异常',
    `create_time`             datetime              DEFAULT NULL,
    `update_time`             datetime              DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `license_job`
(
    `id`          bigint(20) NOT NULL,
    `module_id`   bigint(20) NOT NULL COMMENT 'module id',
    `job_id`      bigint(20) NOT NULL COMMENT '作业ID',
    `licenses`    bigint(20) NOT NULL COMMENT 'licenses',
    `used`        tinyint(1) NOT NULL DEFAULT '1' COMMENT '1-使用中 2-使用完成',
    `create_time` datetime DEFAULT NULL,
    `update_time` datetime DEFAULT NULL,
    `license_id`  bigint(20) NOT NULL COMMENT 'license id',
    PRIMARY KEY (`id`),
    KEY           `IDX_license_job_job_id` (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `license_manager`
(
    `id`           bigint(20) NOT NULL,
    `app_type`     varchar(255)  NOT NULL DEFAULT '' COMMENT '求解器软件类型',
    `os`           tinyint(1) NOT NULL DEFAULT '1' COMMENT '操作系统 1-linux 2-win',
    `status`       tinyint(1) NOT NULL DEFAULT '2' COMMENT '发布状态 1-已发布 2-未发布',
    `description`  varchar(1024) NOT NULL DEFAULT '' COMMENT '描述',
    `compute_rule` varchar(1024) NOT NULL DEFAULT '' COMMENT 'license使用计算规则',
    `publish_time` datetime               DEFAULT NULL COMMENT '发布时间',
    `create_time`  datetime               DEFAULT NULL,
    `update_time`  datetime               DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `module_config`
(
    `id`           bigint(20) NOT NULL,
    `license_id`   bigint(20) NOT NULL COMMENT 'license id',
    `module_name`  varchar(255) NOT NULL COMMENT '模块名称',
    `total`        int(11) NOT NULL DEFAULT '0' COMMENT 'licenses数量',
    `used`         int(11) NOT NULL DEFAULT '0' COMMENT 'licenses已使用数量',
    `actual_total` int(11) NOT NULL DEFAULT '0' COMMENT '实时总数量',
    `actual_used`  int(11) NOT NULL DEFAULT '0' COMMENT '实时已使用数量',
    `create_time`  datetime DEFAULT NULL,
    `update_time`  datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `UQE_module_config_license_module_name` (`license_id`,`module_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

