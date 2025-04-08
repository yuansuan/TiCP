ALTER TABLE upload_info MODIFY COLUMN path VARCHAR(255);
ALTER TABLE upload_info MODIFY COLUMN tmp_path VARCHAR(255);

ALTER TABLE compress_info MODIFY COLUMN tmp_path VARCHAR(255);
ALTER TABLE compress_info MODIFY COLUMN paths VARCHAR(255);
ALTER TABLE compress_info MODIFY COLUMN target_path VARCHAR(255);
ALTER TABLE compress_info MODIFY COLUMN base_path VARCHAR(255);

ALTER TABLE storage_operation_log MODIFY COLUMN src_path VARCHAR(255);
ALTER TABLE storage_operation_log MODIFY COLUMN dest_path VARCHAR(255);

ALTER TABLE shared_directory MODIFY COLUMN path VARCHAR(255);

ALTER TABLE directory_usage MODIFY COLUMN path VARCHAR(255);


