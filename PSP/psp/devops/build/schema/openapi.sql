CREATE TABLE IF NOT EXISTS `openapi_user_certificate` (
    `id` bigint(20) unsigned NOT NULL COMMENT 'ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `certificate` varchar(36) NOT NULL COMMENT '用户凭证',
    `created_at` datetime DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT = 'openapi用户凭证';