CREATE TABLE
  IF NOT EXISTS `application_allow` (
    `id` bigint (20) unsigned NOT NULL,
    `application_id` bigint (20) unsigned NOT NULL COMMENT '应用id',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_application_allow_unique` (`application_id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '计算应用白名单';
