ALTER TABLE `sc_job`
    ADD COLUMN exec_hosts text COMMENT '执行作业使用的节点' AFTER completed_time,
    ADD COLUMN exec_host_num bigint(20) NOT NULL DEFAULT 0 COMMENT '作业实际使用的节点数量' AFTER exec_hosts;