ALTER TABLE cloudapp_session
  ADD COLUMN `account_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '账户ID' AFTER `deleted`;
ALTER TABLE cloudapp_session
  ADD COLUMN `charge_type` varchar(32) NOT NULL DEFAULT '' AFTER `account_id`;
ALTER TABLE cloudapp_session
  ADD COLUMN `is_paid_finished` tinyint(0) NOT NULL DEFAULT '0' AFTER `charge_type`;

CREATE TABLE `cloudapp_bill`
(
  `id`            BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id`    BIGINT(20) UNSIGNED NOT NULL COMMENT '会话ID',
  `order_id`      BIGINT(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `resource_id`   BIGINT(20) UNSIGNED NOT NULL COMMENT '资源ID [软件|硬件]',
  `resource_type` VARCHAR(32) COMMENT '资源类型',
  `bill_time`     datetime                     DEFAULT NULL COMMENT '会话开始时间',
  `create_time`   datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time`   datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_idx_order_id` (`order_id`),
  UNIQUE KEY `uniq_idx_session_id_resource_id` (`session_id`, `resource_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='付费信息';

CREATE TABLE `cloudapp_hardware_user`
(
  `id`          BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `hardware_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '硬件ID',
  `user_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_idx_hardware_id_user_id` (`hardware_id`, `user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='硬件用户关联信息';

CREATE TABLE `cloudapp_software_user`
(
  `id`          BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `software_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '软件ID',
  `user_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_idx_software_id_user_id` (`software_id`, `user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='软件用户关联信息';


