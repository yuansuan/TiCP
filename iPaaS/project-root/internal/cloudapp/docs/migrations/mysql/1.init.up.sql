CREATE TABLE `cloudapp_software`
(
  `id`          bigint(20) unsigned NOT NULL,
  `zone`        varchar(32)  NOT NULL DEFAULT '' COMMENT '可用区ID',
  `name`        varchar(64)  NOT NULL DEFAULT '' COMMENT '软件方案名字',
  `desc`        varchar(255) NOT NULL DEFAULT '' COMMENT '软件方案描述',
  `icon`        varchar(255) NOT NULL DEFAULT '' COMMENT '软件方案图标',
  `platform`    varchar(32)  NOT NULL DEFAULT '' COMMENT '软件平台：DESKTOP, APPLICATION',
  `image_id`    varchar(64)  NOT NULL DEFAULT '' COMMENT '镜像Id',
  `init_script` text COMMENT '初始化脚本内容',
  `gpu_desired` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否需要GPU支持',
  `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY           `idx_software_imageid` (`image_id`),
  KEY           `idx_software_zone` (`zone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公有云可视化软件参数';

CREATE TABLE `cloudapp_hardware`
(
  `id`              bigint(20) unsigned NOT NULL,
  `zone`            varchar(32)  NOT NULL DEFAULT '' COMMENT '可用区ID',
  `name`            varchar(64)  NOT NULL DEFAULT '' COMMENT '硬件方案名字',
  `desc`            varchar(255) NOT NULL DEFAULT '' COMMENT '硬件方案描述',
  `instance_type`   varchar(64)  NOT NULL DEFAULT '' COMMENT '实例类型名字',
  `instance_family` varchar(32)  NOT NULL DEFAULT '' COMMENT '实例机型系列',
  `network`         int(11) NOT NULL DEFAULT '0' COMMENT '实例的最大内网带宽',
  `cpu`             int(11) NOT NULL DEFAULT '0' COMMENT 'CPU核数',
  `cpu_model`       varchar(256) NOT NULL DEFAULT '' COMMENT 'CPU型号',
  `mem`             int(11) NOT NULL DEFAULT '0' COMMENT '内存容量，单位G',
  `gpu`             int(11) NOT NULL DEFAULT '0' COMMENT 'GPU数量',
  `gpu_model`       varchar(256) NOT NULL DEFAULT '' COMMENT 'GPU型号',
  `create_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY               `idx_hardware_instancetype` (`instance_type`),
  KEY               `idx_hardware_instancefamily` (`instance_family`),
  KEY               `idx_hardware_zone` (`zone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公有云可视化硬件配置';

CREATE TABLE `cloudapp_session`
(
  `id`                bigint(20) unsigned NOT NULL,
  `zone`              varchar(32)   NOT NULL DEFAULT '' COMMENT '可用区ID',
  `user_id`           bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建会话用户Id',
  `instance_id`       bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '会话所运行的实例Id',
  `status`            varchar(64)   NOT NULL DEFAULT '' COMMENT '会话状态',
  `desktop_url`       varchar(1024) NOT NULL DEFAULT '' COMMENT '桌面接入地址',
  `start_time`        datetime               DEFAULT NULL COMMENT '会话开始时间',
  `end_time`          datetime               DEFAULT NULL COMMENT '会话结束时间',
  `close_signal`      tinyint(1) NOT NULL DEFAULT '0' COMMENT '用户关闭会话信号',
  `user_close_time`   datetime               DEFAULT NULL COMMENT '用户关闭会话时间',
  `exit_reason`       text COMMENT '会话退出原因',
  `deleted`           tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已删除',
  `create_time`       datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time`       datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY                 `idx_session_userid` (`user_id`),
  KEY                 `idx_session_instanceid` (`instance_id`),
  KEY                 `idx_session_closesignal` (`close_signal`),
  KEY                 `idx_session_zone` (`zone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公有云可视化会话列表';

CREATE TABLE `cloudapp_instance`
(
  `id`              bigint(20) unsigned NOT NULL,
  `zone`            varchar(32) NOT NULL DEFAULT '' COMMENT '可用区ID',
  `hardware_id`     bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '绑定的硬件Id',
  `software_id`     bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '绑定的软件Id',
  `init_script`     text COMMENT '实例初始化脚本模板',
  `user_params`     text COMMENT '实例初始化用户传入参数',
  `user_script`     text COMMENT '实例实际执行的脚本内容',
  `instance_id`     varchar(64) NOT NULL DEFAULT '' COMMENT '实例Id',
  `instance_data`   text COMMENT '腾讯云实例数据',
  `ssh_password`    varchar(64) NOT NULL DEFAULT '' COMMENT '登录密码',
  `instance_status` varchar(32) NOT NULL DEFAULT '' COMMENT '腾讯云实例状态',
  `start_time`      datetime             DEFAULT NULL COMMENT '实例开始时间',
  `end_time`        datetime             DEFAULT NULL COMMENT '实例结束时间',
  `create_time`     datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time`     datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY               `idx_instance_hardwareid` (`hardware_id`),
  KEY               `idx_instance_softwareid` (`software_id`),
  KEY               `idx_instance_instanceid` (`instance_id`),
  KEY               `idx_instance_zone` (`zone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公有云可视化实例列表';

CREATE TABLE `cloudapp_remote_app`
(
  `id`          bigint(20) unsigned NOT NULL,
  `software_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '搭配的软件Id',
  `desc`        text COMMENT '描述',
  `name`        varchar(256)  NOT NULL DEFAULT '' COMMENT 'RemoteApp名称',
  `dir`         varchar(256)  NOT NULL DEFAULT '' COMMENT 'RemoteApp目录',
  `args`        text COMMENT 'RemoteApp启动参数',
  `logo`        varchar(1024) NOT NULL DEFAULT '' COMMENT 'RemoteApp logo 前端显示',
  `disable_gfx` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否使用h264视频渲染桌面 true 不使用视频 false 使用视频',
  `create_time` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_idx_remote_app` (`software_id`,`name`),
  KEY           `idx_remote_app_software_id` (`software_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公有云可视化RemoteApp';
