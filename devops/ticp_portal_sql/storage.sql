-- share_file_record
CREATE TABLE IF NOT EXISTS `share_file_record`
(
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime DEFAULT NULL COMMENT '修改时间',
    `id`          bigint(20)   NOT NULL COMMENT '主键id',
    `file_path`   varchar(255) NOT NULL COMMENT '分享文件路径',
    `owner`       varchar(128) NOT NULL COMMENT '文件持有人名称',
    `type`        tinyint(1)   NOT NULL COMMENT '分享方式,1-复制 2-硬链接',
    `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '分享记录表';

-- 分享文件记录
CREATE TABLE IF NOT EXISTS `share_file_user`
(
    `share_record_id` bigint(20) NOT NULL COMMENT '分享文件记录id',
    `user_id`         bigint(20) NOT NULL COMMENT '用户id',
    `state`           tinyint(1)   NOT NULL COMMENT '状态,1-未处理 2-已处理'
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '分享文件用户关联表';

