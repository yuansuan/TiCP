CREATE TABLE IF NOT EXISTS `notice_message` (
  `id` bigint(20) unsigned NOT NULL COMMENT 'ID',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户ID',
  `state` int(11) NOT NULL DEFAULT 0 COMMENT '消息状态: 1-未读; 2-已读',
  `type` varchar(255) NOT NULL DEFAULT '' COMMENT '消息类型',
  `content` text COMMENT '消息内容',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '通知信息表';