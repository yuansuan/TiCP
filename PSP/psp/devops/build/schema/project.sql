create table project
(
    id            bigint unsigned         not null comment 'ID' primary key,
    project_name  varchar(256) default '' not null comment '项目名称',
    project_owner bigint       default 0  not null comment '项目管理员',
    state         varchar(64)  default '' not null comment '项目状态:  Init:初始化，Running:进行中，Terminated:已终止，Completed:已结束',
    start_time    datetime                null comment '开始时间',
    end_time      datetime                null comment '结束时间',
    comment       varchar(512)            null comment '项目描述',
    file_path     varchar(512) default '' null comment '项目路径',
    is_delete     tinyint      default 0  not null,
    create_time   datetime                null comment '创建时间',
    update_time   datetime                null comment '更新时间'
) ENGINE = InnoDB DEFAULT CHARSET = utf8 comment '项目表';

create table project_member
(
    id          bigint unsigned             not null comment 'ID' primary key,
    project_id  bigint unsigned default '0' not null comment '项目id',
    user_id     bigint unsigned default '0' not null comment '用户id',
    is_delete   tinyint         default 0   not null,
    link_path   varchar(512)    default ''  null comment '链接文件路径',
    create_time datetime                    not null comment '创建时间',
    update_time datetime                    not null comment '更新时间'
) ENGINE = InnoDB DEFAULT CHARSET = utf8 comment '项目成员表';