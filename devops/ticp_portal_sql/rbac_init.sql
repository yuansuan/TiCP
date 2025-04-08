INSERT INTO `resource` (`id`, `name`, `action`, `type`, `display_name`, `custom`, `external_id`, `parent_id`)
VALUES (1, 'sys_manager', 'NONE', 'system', '系统管理', 1, 0, 0);

INSERT INTO `resource` (`id`, `name`, `action`, `type`, `display_name`, `custom`, `external_id`, `parent_id`)
VALUES (2, 'file_manager', 'NONE', 'system', '文件管理', 1, 0, 0);

INSERT INTO `resource` (`id`, `name`, `action`, `type`, `display_name`, `custom`, `external_id`, `parent_id`)
VALUES (3, 'job_manager', 'NONE', 'system', '作业管理(管理员)', 1, 0, 0);

INSERT INTO `resource` (`id`, `name`, `action`, `type`, `display_name`, `custom`, `external_id`, `parent_id`)
VALUES (4, 'personal_job_manager', 'NONE', 'system', '作业管理(个人)', 1, 0, 0);

INSERT INTO `resource` (`id`, `name`, `action`, `type`, `display_name`, `custom`, `external_id`, `parent_id`)
VALUES (5, 'cluster_monitor', 'NONE', 'system', '集群监控', 1, 0, 0);

INSERT INTO `resource` (id, name, action, type, display_name, custom, external_id, parent_id)
VALUES (6, 'project_manager', 'NONE', 'system', '项目管理(项目管理员)', 1, 0, 0);

INSERT INTO `resource` (id, name, action, type, display_name, custom, external_id, parent_id)
VALUES (7, 'personal_project_manager', 'NONE', 'system', '项目管理(个人)', 1, 0, 0);

INSERT INTO `role` (`id`, `name`, `comment`, `type`)
VALUES (1, '超级管理员', '超级管理员', 1);
INSERT INTO `role` (`id`, `name`, `comment`, `type`)
VALUES (2, '系统管理员', '系统管理员', 0);
INSERT INTO `role` (`id`, `name`, `comment`, `type`)
VALUES (3, '普通用户', '普通用户', 0);
INSERT INTO `role` (id, name, comment, type)
VALUES (4, '项目管理员', '项目管理员', 0);



INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r1', '1', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r1', '2', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r1', '3', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r1', '5', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r1', '6', 'NONE');


INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r2', '1', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r2', '2', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r2', '3', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r2', '5', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r2', '6', 'NONE');

INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r3', '2', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r3', '4', 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r3', '7', 'NONE');

INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r4', '6', 'NONE');

-- 3d云应用权限
INSERT INTO `resource` (`name`, `action`, `type`, `display_name`, `custom`, `external_id`, `parent_id`)
VALUES ('visual', 'NONE', 'system', '3D云应用', 1, 0, 0);

SET @max_id = (SELECT MAX(id) FROM resource);

INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r1', @max_id, 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r2', @max_id, 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r3', @max_id, 'NONE');
INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r4', @max_id, 'NONE');

-- 内置计算应用 StarCCM 权限
INSERT INTO `resource` (`name`, `action`, `type`, `display_name`, `custom`, `external_id`, `parent_id`)
VALUES ('STAR-CCM+ 12.02-R8', 'NONE', 'local_app', 'STAR-CCM+ 12.02-R8', 1, 1689929831401132032, 0);

SET @max_id = (SELECT MAX(id) FROM resource);

INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`)
VALUES ('p', 'r1', @max_id, 'NONE');