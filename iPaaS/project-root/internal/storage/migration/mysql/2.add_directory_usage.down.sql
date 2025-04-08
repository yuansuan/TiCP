DROP TABLE IF EXISTS `directory_usage`;

ALTER TABLE `compress_info`
    MODIFY COLUMN `status` tinyint(1) NOT NULL DEFAULT 0
    COMMENT '状态 0: 未完成 1: 已完成';
