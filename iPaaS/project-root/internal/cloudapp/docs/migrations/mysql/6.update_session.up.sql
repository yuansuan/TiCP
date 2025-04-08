alter table cloudapp_session
    add pay_by_account_id bigint unsigned default 0 null comment '代支付账户' after account_id;
