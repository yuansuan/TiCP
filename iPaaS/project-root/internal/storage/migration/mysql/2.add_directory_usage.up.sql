CREATE TABLE IF NOT EXISTS  `directory_usage`
(
    `id`          VARCHAR(255) NOT NULL COMMENT '目录用量计算任务id',
    `user_id`     VARCHAR(255) NOT NULL COMMENT '用户id',
    `path`        VARCHAR(255) NOT NULL DEFAULT '' COMMENT '目录路径',
    `size`        BIGINT(20) NOT NULL DEFAULT 0 COMMENT '目录大小 单位为字节',
    `status`      TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态 -1: 失败 0: 计算中 1: 成功 2: 已取消',
    `err_msg`     TEXT  COMMENT '错误信息',
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='目录用量表';

ALTER TABLE `compress_info`
    MODIFY COLUMN `status` tinyint(1) NOT NULL DEFAULT 0
    COMMENT '状态 -1: 失败 0: 压缩中 1: 成功 2: 已取消';