CREATE TABLE IF NOT EXISTS `visual_hardware` (
  `id` bigint(20) COMMENT 'ID',
  `out_hardware_id` varchar(255) COMMENT '外部硬件ID',
  `name` varchar(64) COMMENT '名称',
  `desc` varchar(255) COMMENT '描述',
  `network` int(11) COMMENT '网络',
  `cpu` int(11) COMMENT 'CPU',
  `mem` int(11) COMMENT '内存',
  `gpu` int(11) COMMENT 'GPU',
  `cpu_model` varchar(64) COMMENT 'CPU型号',
  `gpu_model` varchar(64) COMMENT 'GPU型号',
  `instance_type` varchar(64) COMMENT '实例类型',
  `instance_family` varchar(64) COMMENT '实例族',
  `zone` varchar(64) COMMENT '区域',
  `deleted` DATETIME COMMENT '是否删除',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '硬件配置表';

CREATE TABLE IF NOT EXISTS `visual_session` (
  `id` bigint(20) COMMENT 'ID',
  `out_session_id` varchar(255) COMMENT '外部会话ID',
  `hardware_id` bigint(20) COMMENT '硬件ID',
  `out_hardware_id` varchar(255) COMMENT '硬件ID',
  `software_id` bigint(20) COMMENT '软件ID',
  `out_software_id` varchar(255) COMMENT '软件ID',
  `project_id` bigint(20) null COMMENT '项目 ID',
  `project_name` varchar(255) null COMMENT '项目名称',
  `user_id` bigint(20) COMMENT '用户ID',
  `user_name` varchar(64) COMMENT '用户名',
  `raw_status` varchar(64) COMMENT '外部会话状态',
  `status` varchar(64) COMMENT '状态',
  `stream_url` text COMMENT '流地址',
  `exit_reason` varchar(255) COMMENT '退出原因',
  `duration` bigint(20) COMMENT '时长',
  `zone` varchar(64) COMMENT '区域',
  `is_auto_close` tinyint(4) COMMENT '是否自动关闭',
  `deleted` DATETIME COMMENT '是否删除',
  `start_time` datetime COMMENT '开始时间',
  `end_time` datetime COMMENT '结束时间',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '会话记录表';

create index session_index_hardware_id on visual_session (hardware_id);

create index session_index_software_id on visual_session (software_id);

CREATE TABLE IF NOT EXISTS `visual_software` (
  `id` bigint(20) COMMENT 'ID',
  `out_software_id` varchar(255) COMMENT '外部软件ID',
  `name` varchar(64) COMMENT '名称',
  `desc` varchar(255) COMMENT '描述',
  `platform` varchar(64) COMMENT '平台',
  `image_id` varchar(255) COMMENT '镜像ID',
  `state` varchar(64) COMMENT '状态',
  `init_script` text COMMENT '初始化脚本',
  `icon` mediumtext COMMENT '图标',
  `gpu_desired` tinyint(4) COMMENT '是否需要GPU',
  `zone` varchar(64) COMMENT '区域',
  `deleted` DATETIME COMMENT '是否删除',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '软件配置表';

CREATE TABLE IF NOT EXISTS `visual_software_preset` (
  `id` bigint(20) COMMENT 'ID',
  `software_id` bigint(20) COMMENT '软件ID',
  `hardware_id` bigint(20) COMMENT '硬件ID',
  `defaulted` tinyint(4) COMMENT '是否默认',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '软件预设表';

CREATE TABLE IF NOT EXISTS `visual_remote_app` (
  `id` bigint(20) COMMENT 'ID',
  `out_remote_app_id` varchar(255) COMMENT '外部远程应用ID',
  `software_id` bigint(20) COMMENT '软件ID',
  `out_software_id` varchar(255) COMMENT '外部软件ID',
  `name` varchar(64) COMMENT '名称',
  `desc` varchar(255) COMMENT '描述',
  `dir` varchar(255) COMMENT '目录',
  `args` varchar(255) COMMENT '参数',
  `logo` varchar(255) COMMENT '图标',
  `disable_gfx` tinyint(4) COMMENT '是否禁用图形',
  `deleted` DATETIME COMMENT '是否删除',
  `create_time` datetime COMMENT '创建时间',
  `update_time` datetime COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '远程应用表';