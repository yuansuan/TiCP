CREATE TABLE IF NOT EXISTS `storage_quota` (
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `storage_usage` float(10,2) NOT NULL DEFAULT '0.00' COMMENT '存储空间用量',
    `storage_limit` float(10,2) NOT NULL DEFAULT '0.00' COMMENT '存储上限',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存储配额表';

CREATE TABLE IF NOT EXISTS `shared_directory` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `user_id` varchar(255) NOT NULL COMMENT '用户id',
    `path` varchar(255) NOT NULL COMMENT '指定路径',
    `shared_user_name` varchar(255) NOT NULL COMMENT '用户名',
    `shared_password` varchar(255) NOT NULL COMMENT '密码',
    `shared_host` varchar(255) NOT NULL COMMENT '共享主机地址',
    `shared_src` varchar(255) NOT NULL COMMENT '共享目录路径',
    `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='共享目录表';


CREATE TABLE IF NOT EXISTS `storage_operation_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `user_id` varchar(255) NOT NULL COMMENT '用户id',
    `file_name` varchar(255) NOT NULL COMMENT '文件名',
    `src_path` varchar(255) NOT NULL COMMENT '源路径',
    `dest_path` varchar(255) NOT NULL COMMENT '目标路径',
    `file_type` varchar(20) NOT NULL COMMENT '文件类型, 可选值: file-普通文件, folder-文件夹, batch-批量操作',
    `operation_type` varchar(20) NOT NULL COMMENT '操作类型, 可选值: upload-上传, download-下载, delete-删除, move-移动, mkdir-添加文件夹, compress-压缩, copy-拷贝, copy_range-指定范围拷贝, create-创建, link-链接, read_at-读, write_at-写',
    `size` varchar(20) NOT NULL COMMENT '文件/文件夹大小',
    `is_deleted` int(2) NOT NULL COMMENT '是否删除',
    `create_time` datetime NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';


CREATE TABLE IF NOT EXISTS `upload_info` (
    `id` varchar(255) NOT NULL COMMENT '主键id',
    `user_id` varchar(255) DEFAULT NULL COMMENT '用户id',
    `tmp_path` varchar(255) NOT NULL DEFAULT '' COMMENT '临时文件路径',
    `path` varchar(255) NOT NULL DEFAULT '' COMMENT '文件路径',
    `size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小',
    `overwrite` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否覆盖',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='上传任务信息表';


CREATE TABLE IF NOT EXISTS `compress_info` (
    `id` varchar(255) NOT NULL COMMENT '上传 ID',
    `user_id` varchar(32) NOT NULL COMMENT '请求者用户 ID',
    `tmp_path` varchar(255) NOT NULL DEFAULT '' COMMENT '临时文件路径',
    `paths` varchar(255) NOT NULL DEFAULT '' COMMENT '源文件/文件夹路径',
    `target_path` varchar(255) NOT NULL DEFAULT '' COMMENT '目标文件路径',
    `base_path` varchar(255) DEFAULT '' COMMENT '压缩包起始文件夹路径',
    `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态 0: 未完成 1: 已完成',
    `error_msg` text COMMENT '错误信息',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='压缩任务信息表';