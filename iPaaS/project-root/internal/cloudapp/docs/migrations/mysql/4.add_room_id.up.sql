ALTER TABLE `cloudapp_session`
  ADD COLUMN `room_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'webrtc双端通信唯一标志' AFTER `is_paid_finished`;
