ALTER TABLE `cloudapp_remote_app`
  ADD COLUMN `login_user` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '登陆的系统用户名' AFTER `disable_gfx`;

CREATE TABLE `cloudapp_remote_app_user_pass`
(
  `id`              BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id`      BIGINT(20) UNSIGNED NOT NULL COMMENT '会话ID',
  `remote_app_name` VARCHAR(64)         NOT NULL DEFAULT '' COMMENT 'RemoteApp名称',
  `username`        VARCHAR(64)         NOT NULL DEFAULT '' COMMENT '登陆系统用户',
  `password`        VARCHAR(64)         NOT NULL DEFAULT '' COMMENT '登陆用户密码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_session_id_remote_app_name` (`session_id`, `remote_app_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='远程应用用户密码表';
