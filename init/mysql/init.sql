create database if not exists biz_envoy;

use biz_envoy;
create table if not exists `control_version`(
    `id` bigint(20) not null auto_increment comment '主键',
    `version` int default 0 not null comment '时间戳',
    `service_id` varchar(32) default '' not null comment '服务名',
    primary key (`id`),
    unique key uk_service(`service_id`)
)engine=InnoDB default charset =utf8mb4 comment='表d';
