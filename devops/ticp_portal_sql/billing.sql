-- auto-generated definition
create table IF NOT EXISTS billing_cloud (
    id bigint unsigned default 0 not null comment '账单ID',
    user_id bigint unsigned default 0 not null comment '用户ID',
    user_name varchar(100) default '' not null comment '用户名',
    account_id varchar(128) default '' not null comment '账户ID',
    freeze_amount bigint unsigned default 0 not null comment '账户冻结金额',
    out_biz_id varchar(128) default '' not null comment '外部订单号',
    out_resource_id varchar(128) default '' not null comment '外部资源id',
    out_product_name varchar(128) default '' not null comment '资源类型 CloudCompute: pass求解作业, CloudApp: 云应用 ',
    real_resource_id bigint unsigned default 0 not null comment 'TiCP平台对应的业务id',
    project_id bigint unsigned default '1687026933658816512' not null comment '项目id',
    project_name varchar(255) default 'default' not null comment '项目名称',
    name varchar(64) default '' not null comment '账单业务名称',
    comment text null comment '说明',
    merchandise_id varchar(128) default '' not null comment '商品ID',
    merchandise_name varchar(255) default '' not null comment '商品名称',
    merchandise_price_unit bigint unsigned default 0 not null comment '商品单价',
    merchandise_price_desc varchar(100) default '' not null comment '商品价格描述',
    merchandise_quantity double default 0 not null comment '商品数量',
    merchandise_quantity_unit varchar(32) null comment '消耗数量单位',
    amount bigint unsigned default 0 not null comment '总价：商品数量 * 单价',
    discount_amount bigint unsigned default 0 null comment '折扣金额',
    deduction_amount bigint unsigned default 0 null comment '抵扣金额',
    refund_amount bigint default 0 null comment '退款金额',
    start_time datetime null comment '扣费周期开始时间，按量付费使用',
    end_time datetime null comment '扣费周期结束时间，按量付费使用',
    job_submit_time datetime null comment '作业提交时间',
    job_start_time datetime null comment '作业开始时间',
    job_end_time datetime null comment '作业结束时间',
    job_app_id   bigint unsigned default '0' null comment '作业app id',
    job_app_name varchar(256) default '' null comment '作业对应app名称',
    latest_trade_time datetime null comment '最后出账扣费时间',
    is_deleted tinyint default 0 not null comment '删除标记 ',
    create_time datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    update_time datetime default CURRENT_TIMESTAMP not null comment '更新时间',
    PRIMARY KEY (`id`)
) engine = InnoDB default charset = utf8 comment '账单表';

create index billing_idx_user_name on billing_cloud (user_name);

create index billing_index_real_resource_id on billing_cloud (real_resource_id);