CREATE TABLE IF NOT EXISTS `sys_config` (
  `id` bigint(20) COMMENT 'ID',
  `key` varchar(255) COMMENT '配置键',
  `value` text COMMENT '配置值',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '系统配置';

INSERT INTO `sys_config` (`id`, `key`, `value`, `create_time`, `update_time`)
VALUES (1699622219002417152, 'DefaultRoleId', '2', '2023-09-07 11:14:49', '2023-09-07 11:14:49');

CREATE TABLE `alert_notification` (
  `id` bigint NOT NULL COMMENT 'ID',
  `key` varchar(255) default '' not null COMMENT '配置键',
  `value` VARCHAR(255) default '' not null COMMENT '配置值',
  `type` varchar(16) default '' not null COMMENT '类型',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='告警通知配置';