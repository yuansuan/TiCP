ALTER TABLE job
ADD COLUMN needed_paths TEXT COMMENT '正则表达式,符合规则的文件路径将会进行回传';