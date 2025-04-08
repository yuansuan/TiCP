ALTER TABLE `cloudapp_instance`
  ADD COLUMN `boot_volume_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '启动卷ID' AFTER `instance_status`;
