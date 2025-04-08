CREATE TABLE
  IF NOT EXISTS `job` (
    -- 作业基本信息
    `id` bigint (20) unsigned NOT NULL COMMENT '作业ID',
    `name` varchar(255) NOT NULL DEFAULT '作业名称',
    `comment` text COMMENT '作业备注',
    `user_id` bigint (20) unsigned NOT NULL COMMENT '用户ID',
    `job_source` varchar(20) DEFAULT '' COMMENT '作业来源',
    -- 状态信息
    `state` int (11) unsigned NOT NULL COMMENT '作业状态',
    `sub_state` int (11) unsigned NOT NULL COMMENT '作业子状态',
    `state_reason` longtext COMMENT '作业等待或者中间其他原因',
    `exit_code` text COMMENT '作业退出码',
    `file_sync_state` varchar(32) NOT NULL DEFAULT '' COMMENT '文件同步状态',
    -- 参数信息
    `params` text COMMENT '用户原始参数',
    `user_zone` varchar(64) NOT NULL DEFAULT '' COMMENT '用户选择的分区 若空为未选择',
    `timeout` bigint (20) NOT NULL DEFAULT '-1' COMMENT '超时时间',
    `file_classifier` text COMMENT '文件分类器',
    `resource_usage_cpus` bigint (20) unsigned NOT NULL DEFAULT '0' COMMENT '用户选择使用核数',
    `resource_usage_memory` bigint (20) unsigned NOT NULL DEFAULT '0' COMMENT '用户选择使用内存',
    `custom_state_rule_key_statement` text COMMENT '自定义状态规则key语句',
    `custom_state_rule_result_state` varchar(20) DEFAULT '' COMMENT '自定义状态规则结果状态',
    `no_round` tinyint (4) NOT NULL DEFAULT '0' COMMENT '单节点是否不进行取整,仅限内部用户使用',
    -- 作业运行信息
    `hpc_job_id` varchar(64) NOT NULL COMMENT 'HPC作业ID',
    `zone` varchar(64) DEFAULT '' COMMENT '实际运行的分区',
    `resource_assign_cpus` bigint (20) unsigned NOT NULL DEFAULT '0' COMMENT '实际分配使用核数',
    `resource_assign_memory` bigint (20) unsigned DEFAULT '0' COMMENT '实际分配使用内存',
    `command` text NOT NULL COMMENT '作业实际执行命令',
    `work_dir` varchar(255) NOT NULL DEFAULT '' COMMENT '工作目录',
    `origin_job_id` varchar(64) NOT NULL COMMENT '调度器作业ID',
    `queue` varchar(255) NOT NULL COMMENT '作业实际运行的队列',
    `priority` bigint (20) NOT NULL DEFAULT '0' COMMENT '作业实际优先级',
    `exec_hosts` varchar(256) NOT NULL COMMENT '作业执行节点名称列表',
    `exec_host_num` varchar(64) NOT NULL COMMENT '作业执行节点总数',
    `execution_duration` int (11) NOT NULL DEFAULT '0' COMMENT '作业执行时间',
    -- 文件信息
    `input_type` varchar(20) NOT NULL DEFAULT '' COMMENT 'hpc_storage 或者cloud_storage，数据类型为超算存储或者远算云盒子',
    `input_dir` varchar(255) NOT NULL DEFAULT '' COMMENT '盒子上的作业输入文件目录',
    `destination` varchar(255) NOT NULL DEFAULT '' COMMENT '输入文件的目标路径',
    `output_type` varchar(20) NOT NULL DEFAULT '' COMMENT 'hpc_storage 或者cloud_storage，数据类型为超算存储或者远算云盒子',
    `output_dir` varchar(255) NOT NULL DEFAULT '' COMMENT '盒子上的作业输出文件目录',
    `no_needed_paths` text COMMENT '正则表达式,符合规则的文件路径将不会进行回传',
    `file_input_storage_zone` varchar(20) NOT NULL DEFAULT '' COMMENT '输入文件区域',
    `file_output_storage_zone` varchar(20) NOT NULL DEFAULT '' COMMENT '输出文件区域',
    `download_file_size_total` bigint (20) unsigned DEFAULT '0' COMMENT '下载文件总大小',
    `download_file_size_current` bigint (20) unsigned DEFAULT '0' COMMENT '下载文件当前下载大小',
    `upload_file_size_total` bigint (20) unsigned DEFAULT '0' COMMENT '上传文件总大小',
    `upload_file_size_current` bigint (20) unsigned DEFAULT '0' COMMENT '上传文件当前上传大小',
    -- 应用信息
    `app_id` bigint (20) unsigned NOT NULL COMMENT '计算应用 ID',
    `app_name` varchar(255) NOT NULL COMMENT '计算应用名',
    -- 标志信息
    `user_cancel` tinyint (4) NOT NULL DEFAULT '0' COMMENT '用户取消标记',
    `is_file_ready` tinyint (4) NOT NULL DEFAULT '0' COMMENT '作业文件是否准备完成',
    `download_finished` tinyint (4) NOT NULL DEFAULT '0' COMMENT '作业下载是否完成',
    `is_system_failed` tinyint (4) NOT NULL DEFAULT '0' COMMENT '是否系统失败',
    `is_deleted` tinyint (4) NOT NULL DEFAULT '0' COMMENT '标识作业是否已被删除, 0 - 未删除, 1 - 已删除',
    -- 计费信息
    `account_id` bigint (20) unsigned NOT NULL DEFAULT '0' COMMENT '账户ID',
    `charge_type` varchar(32) NOT NULL DEFAULT '' COMMENT '计费类型 PrePaid | PostPaid',
    `is_paid_finished` tinyint (1) NOT NULL DEFAULT '0' COMMENT '是否完成付费',
    -- 时间信息
    `upload_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作业计算文件上传完成时间 hpc完全获取到所有计算文件的时间',
    `download_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作业回传完成时间 hpc最后一次回传完成的时间, 即FileSyncState变成终态的时间',
    `pending_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态变成Pending的时间',
    `running_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态变成Running的时间',
    `terminating_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态变成Terminating的时间',
    `transmitting_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作业回传中时间 hpc最后一次回传开始的时间',
    `suspending_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态变成Suspending的时间',
    `suspended_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态变成Suspended的时间',
    `submit_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '提交给hpc的时间',
    `end_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作业结束时间,状态变成Completed、Terminated、Failed的时间',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '计算作业';

CREATE TABLE
  IF NOT EXISTS `application` (
    `id` bigint (20) NOT NULL,
    `name` varchar(255) NOT NULL COMMENT 'display of the application, such as: Abaqus 6.1.5',
    `type` varchar(255) NOT NULL COMMENT 'real name of the application, such as: Abaqus, used to classify applications without version',
    `version` varchar(32) NOT NULL COMMENT 'version of the application, such as: 6.1.5',
    `app_params_version` int (11) NOT NULL COMMENT 'app_params version',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `image` varchar(128) NOT NULL COMMENT 'image名称',
    `endpoint` varchar(255) NOT NULL COMMENT '超算中心endpoint',
    `command` text NOT NULL COMMENT '提交命令',
    `publish_status` varchar(32) NOT NULL DEFAULT 'unpublished' COMMENT '发布状态 published, unpublished',
    `description` text COMMENT '应用描述',
    `icon_url` varchar(128) NOT NULL COMMENT '应用图标',
    `cores_max_limit` bigint (20) NOT NULL DEFAULT '0',
    `cores_placeholder` varchar(256) NOT NULL DEFAULT '',
    `file_filter_rule` varchar(255) DEFAULT '{"result": "\\\\.dat$","model": "\\\\.(jou|cas)$","log": "\\\\.(sta|dat|msg|out|log)$","middle": "\\\\.(com|prt)$"}' COMMENT '文件过滤规则',
    `residual_enable` tinyint (1) DEFAULT '0' COMMENT '残差图是否开启',
    `residual_log_regexp` varchar(255) DEFAULT 'stdout.log' COMMENT '残差图文件',
    `residual_log_parser` varchar(255) DEFAULT '' COMMENT '残差图解析器',
    `monitor_chart_enable` tinyint (1) DEFAULT '0' COMMENT '监控图表是否开启',
    `monitor_chart_regexp` varchar(255) DEFAULT '.*\\.out' COMMENT '监控图表文件规则',
    `monitor_chart_parser` varchar(255) DEFAULT '' COMMENT '监控图表解析器',
    `snapshot_enable` tinyint (1) DEFAULT '0' COMMENT '云图是否开启',
    `bin_path` text COMMENT '应用bin路径',
    `extention_params` text COMMENT '扩展参数',
    `lic_manager_id` bigint (20) DEFAULT NULL COMMENT 'LicenseManager id, 如果为空代表免费软件',
    `need_limit_core` tinyint (1) DEFAULT '0' COMMENT '是否需要限制核数',
    `specify_queue` varchar(255) DEFAULT '指定队列',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_app_unique` (`version`, `type`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '计算应用';

CREATE TABLE
  IF NOT EXISTS `application_quota` (
    `id` bigint (20) unsigned NOT NULL,
    `application_id` bigint (20) unsigned NOT NULL COMMENT '应用id',
    `ys_id` bigint (20) unsigned NOT NULL COMMENT '用户id',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_application_quota_unique` (`application_id`, `ys_id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '计算应用配额';

CREATE TABLE
  IF NOT EXISTS `job_bill` (
    `id` BIGINT (20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `job_id` BIGINT (20) UNSIGNED NOT NULL COMMENT '作业ID',
    `order_id` BIGINT (20) UNSIGNED NOT NULL COMMENT '订单ID',
    `app_id` BIGINT (20) UNSIGNED NOT NULL COMMENT '软件ID',
    `billed_duration` BIGINT (20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '已经被收费的duration，单位s，对应job表中的ExecutionDuration',
    `bill_time` datetime DEFAULT NULL COMMENT '扣上一个BillDuration的时间点',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_idx_order_id` (`order_id`),
    UNIQUE KEY `uniq_idx_job_id` (`job_id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '付费信息';

CREATE TABLE
  IF NOT EXISTS `residual` (
    `id` BIGINT (20) NOT NULL AUTO_INCREMENT COMMENT '残差图ID',
    `job_id` BIGINT (20) NOT NULL COMMENT '作业ID',
    `content` LONGTEXT COMMENT '残差图内容,经过base64',
    `finished` TINYINT (1) NOT NULL DEFAULT '0' COMMENT '是否完成',
    `residual_log_regexp` VARCHAR(255) NOT NULL DEFAULT 'stdout.log' COMMENT '残差图文件',
    `residual_log_parser` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '残差图解析器类型',
    `failed_reason` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '失败原因',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_job_id` (`job_id`),
    KEY `idx_finished` (`finished`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '残差图';

CREATE TABLE
  IF NOT EXISTS `monitor_chart` (
    `id` BIGINT (20) NOT NULL AUTO_INCREMENT COMMENT '监控图表ID',
    `job_id` BIGINT (20) NOT NULL COMMENT '作业ID',
    `content` LONGTEXT COMMENT '监控图表内容',
    `finished` TINYINT (1) NOT NULL DEFAULT '0' COMMENT '是否完成',
    `monitor_chart_regexp` VARCHAR(255) NOT NULL DEFAULT '.*\\.out' COMMENT '监控图表文件规则',
    `monitor_chart_parser` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '监控图表解析器',
    `failed_reason` VARCHAR(512) DEFAULT NULL COMMENT '失败原因',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_job_id` (`job_id`),
    KEY `idx_finished` (`finished`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '监控图表';