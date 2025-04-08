ALTER TABLE cloudapp_session
  DROP COLUMN `account_id`;
ALTER TABLE cloudapp_session
  DROP COLUMN `charge_type`;
ALTER TABLE cloudapp_session
  DROP COLUMN `is_paid_finished`;

DROP TABLE `cloudapp_bill`;

DROP TABLE `cloudapp_hardware_user`;

DROP TABLE `cloudapp_software_user`;
