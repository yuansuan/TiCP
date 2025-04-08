CREATE TABLE IF NOT EXISTS `job` (
  `id` bigint(20) unsigned NOT NULL COMMENT 'ID',
  `app_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '应用ID',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户ID',
  `job_set_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作业集ID',
  `project_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '项目ID',
  `out_job_id` varchar(255) NOT NULL DEFAULT '' COMMENT '外部接口作业ID',
  `real_job_id` varchar(255) NOT NULL DEFAULT '' COMMENT '调度器作业ID',
  `upload_task_id` varchar(255) NOT NULL DEFAULT '' COMMENT '上传任务ID',
  `type` varchar(64) NOT NULL DEFAULT '' COMMENT '作业类型',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '作业名称',
  `queue` varchar(255) NOT NULL DEFAULT '' COMMENT '作业队列',
  `state` varchar(64) NOT NULL DEFAULT '' COMMENT '作业状态',
  `raw_state` varchar(64) NOT NULL DEFAULT '' COMMENT '作业原始状态',
  `data_state` varchar(64) NOT NULL DEFAULT '' COMMENT '作业数据状态',
  `exit_code` varchar(64) NOT NULL DEFAULT '' COMMENT '作业退出码',
  `app_name` varchar(255) NOT NULL DEFAULT '' COMMENT '应用名称',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名称',
  `job_set_name` varchar(255) NOT NULL DEFAULT '' COMMENT '作业集名称',
  `project_name` varchar(255) NOT NULL DEFAULT '' COMMENT '项目名称',
  `cluster_name` varchar(255) NOT NULL DEFAULT '' COMMENT '集群名称(Zone)',
  `priority` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作业优先级',
  `cpus_alloc` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作业已分配核数',
  `mem_alloc` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作业已分配内存(MB)',
  `exec_duration` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作业实际运行时长(秒)',
  `exec_host_num` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作业执行节点数量',
  `burst_num` int(11) NOT NULL DEFAULT 0 COMMENT '爆发次数',
  `vis_analysis` text comment '可视化分析',
  `reason` text COMMENT '作业状态原因',
  `work_dir` text COMMENT '作业工作目录',
  `exec_hosts` text COMMENT '作业执行节点名称',
  `submit_time` datetime COMMENT '作业提交时间',
  `pend_time` datetime COMMENT '作业等待时间',
  `start_time` datetime COMMENT '作业开始时间',
  `end_time` datetime COMMENT '作业结束时间',
  `terminate_time` datetime COMMENT '作业终止时间',
  `suspend_time` datetime COMMENT '作业暂停时间',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '作业表';

CREATE INDEX idx_queue ON job (queue);
CREATE INDEX idx_state ON job (state);
CREATE INDEX idx_app_name ON job (app_name);
CREATE INDEX idx_user_name ON job (user_name);
CREATE INDEX idx_submit_time ON job (submit_time);
CREATE INDEX idx_type_app_name ON job (type, app_name);

CREATE TABLE IF NOT EXISTS `job_attr` (
  `job_id` bigint(20) unsigned NOT NULL COMMENT '作业业务ID',
  `key` varchar(255) NOT NULL COMMENT '属性名称',
  `value` MEDIUMTEXT COMMENT '属性值',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`job_id`, `key`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '作业属性表';

CREATE TABLE IF NOT EXISTS `job_timeline` (
  `job_id` bigint(20) unsigned NOT NULL COMMENT '作业业务ID',
  `event_name` varchar(255) NOT NULL COMMENT '事件名称',
  `event_time` datetime COMMENT '事件时间',
  PRIMARY KEY (`job_id`, `event_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '作业时间线';