CREATE TABLE IF NOT EXISTS `sc_job` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `idempotent_id` varchar(1024) NOT NULL DEFAULT '' COMMENT '幂等ID 保证作业只提交一次',

  `state` varchar(64) NOT NULL DEFAULT '' COMMENT '状态',
  `sub_state` VARCHAR(45) NULL DEFAULT '' COMMENT '子状态',
  `file_sync_state` varchar(64) NOT NULL DEFAULT '' COMMENT '回传文件文件同步状态',

  `state_reason` text COMMENT '状态原因: 异常状态会显示状态原因',
  `download_current_size` BIGINT(20) NULL DEFAULT 0 COMMENT '下载当前大小',
  `download_total_size` BIGINT(20) NULL DEFAULT 0 COMMENT '下载总大小',
  `upload_current_size` BIGINT(20) NULL DEFAULT 0 COMMENT '上传当前大小',
  `upload_total_size` BIGINT(20) NULL DEFAULT 0 COMMENT '上传总大小',

  `queue` varchar(1024) NOT NULL DEFAULT '' COMMENT '调度器队列',
  `priority` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '优先级',
  `request_cores` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '请求核数',
  `request_memory` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '请求内存',

  `app_mode` varchar(1024) NOT NULL DEFAULT '' COMMENT '使用应用的方式 image | local',
  `app_path` varchar(1024) NOT NULL DEFAULT '' COMMENT '应用的绝对路径',
  `singularity_image` varchar(1024) NOT NULL DEFAULT '' COMMENT '使用的singularity镜像',

  `inputs` text COMMENT '输入文件 json',
  `output` text COMMENT '输出文件 json',
  `env_vars` text COMMENT '环境变量 json',
  `command` text COMMENT '输入命令行',

  `workspace` varchar(1024) NOT NULL DEFAULT '' COMMENT '工作目录',
  `script` varchar(1024) NOT NULL DEFAULT '' COMMENT '提交脚本',
  `stdout` varchar(1024) NOT NULL DEFAULT '' COMMENT '标准输出文件路径',
  `stderr` varchar(1024) NOT NULL DEFAULT '' COMMENT '错误输出文件路径',

  `is_override` BOOLEAN NOT NULL DEFAULT 0 COMMENT '是否覆写原目录（原地计算）',
  `work_dir` varchar(1024) NOT NULL DEFAULT '' COMMENT '原地计算时的目录路径',

  `origin_job_id` varchar(64) NOT NULL DEFAULT '' COMMENT '原始作业id',
  `alloc_cores` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '分配核数',
  `alloc_memory` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '分配内存',
  `origin_state` varchar(64) NOT NULL DEFAULT '' COMMENT '原始作业状态',
  `exit_code` varchar(64) NOT NULL DEFAULT '' COMMENT '退出码',
  `execution_duration` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作业执行时间 单位s',
  `pending_time` datetime comment '入队时间 进入pendding状态时间',
  `running_time` datetime comment '开始运行时间 进入running状态时间',
  `completing_time` datetime comment '计算完成时间 进入completing状态时间',
  `completed_time` datetime comment '完成时间 进入completed状态时间',

  `control_bit_terminate` tinyint(4) NOT NULL DEFAULT 0 COMMENT '作业控制码 取消',
  `timeout` bigint(20) COMMENT '任务超时时间(秒)',
  `is_timeout` tinyint(4) NOT NULL DEFAULT 0 COMMENT '作业控制码 超时取消',

  `webhook` text COMMENT 'webhook地址 为空则表示不调用webhook',

  `custom_state_rule` text COMMENT '自定义作业状态规则 json',

  `create_time` datetime not null default now() comment '创建时间',
  `update_time` datetime not null default now() on update now() comment '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '作业表';

