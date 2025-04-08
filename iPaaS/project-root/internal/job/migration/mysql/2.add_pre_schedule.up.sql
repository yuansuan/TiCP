CREATE TABLE
    IF NOT EXISTS `pre_schedule` (
        `id` BIGINT (20) UNSIGNED NOT NULL COMMENT '预调度ID',
        `params` TEXT COMMENT '用户参数',
        `expected_min_cpus` BIGINT (20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '期望的最小核数',
        `expected_max_cpus` BIGINT (20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '期望的最大核数',
        `expected_memory` BIGINT (20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '期望的内存数',
        `zone` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '预调度的分区',
        `command` TEXT NOT NULL COMMENT '作业实际执行命令',
        `work_dir` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '工作目录',
        `app_id` BIGINT (20) UNSIGNED NOT NULL COMMENT '计算应用ID',
        `app_name` VARCHAR(255) NOT NULL COMMENT '计算应用名',
        `envs` TEXT COMMENT '环境变量',
        `used` TINYINT (1) NOT NULL DEFAULT 0 COMMENT '是否已经被使用',
        `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '计算预调度';

ALTER TABLE `job`
ADD COLUMN `pre_schedule_id` VARCHAR(64) COMMENT '预调度ID' AFTER `no_round`;