CREATE TABLE IF NOT EXISTS `shared_directory` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` VARCHAR(255) NOT NULL COMMENT '用户id',
  `path` VARCHAR(255) NOT NULL COMMENT '指定路径',
  `shared_user_name` VARCHAR(255) NOT NULL COMMENT '用户名',
  `shared_password` VARCHAR(255) NOT NULL COMMENT '密码',
  `shared_host` VARCHAR(255) NOT NULL COMMENT '共享主机地址',
  `shared_src` VARCHAR(255) NOT NULL COMMENT '共享目录路径',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
