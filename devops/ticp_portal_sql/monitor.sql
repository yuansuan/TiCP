CREATE TABLE IF NOT EXISTS `monitor_node` (
    `id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'id',
    `node_name` varchar(64) NOT NULL DEFAULT '' COMMENT '节点名称',
    `scheduler_status` varchar(64) NOT NULL DEFAULT '' COMMENT '调度器状态（原始的）',
    `status` varchar(64) NOT NULL DEFAULT '' COMMENT '调度器状态（加工后的）',
    `node_type` varchar(16) NOT NULL DEFAULT '' COMMENT '节点类型',
    `queue_name` varchar(16)  NOT NULL DEFAULT '' COMMENT '所属队列',
    `platform_name` varchar(64)  NOT NULL DEFAULT '' COMMENT '标识',
    `total_core_num` int(11) NOT NULL DEFAULT 0 COMMENT '总核数',
    `used_core_num` int(11) NOT NULL DEFAULT 0 COMMENT '已使用核数',
    `free_core_num` int(11) NOT NULL DEFAULT 0 COMMENT '剩余的核数',
    `total_mem` int(11) NOT NULL DEFAULT 0 COMMENT '总内存空间',
    `used_mem` int(11) NOT NULL DEFAULT 0 COMMENT '已使用的内存空间',
    `available_mem` int(11) NOT NULL DEFAULT 0 COMMENT '可以使用的内存空间',
    `free_mem` int(11) NOT NULL DEFAULT 0 COMMENT '剩余内存空间',
    `create_time` datetime COMMENT '创建时间',
    `update_time` datetime COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '监控节点信息';