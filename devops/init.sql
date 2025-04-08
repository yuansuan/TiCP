
CREATE DATABASE IF NOT EXISTS ticp_portal DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE DATABASE IF NOT EXISTS ticp DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- CREATE USER IF NOT EXISTS 'ticp_user'@'%' IDENTIFIED BY 'ticp6655';

GRANT ALL PRIVILEGES ON ticp.* TO 'ticp_user'@'%';
GRANT ALL PRIVILEGES ON ticp_portal.* TO 'ticp_user'@'%';
FLUSH PRIVILEGES;

USE ticp;

CREATE TABLE IF NOT EXISTS `account_bill_version`
(
    `current` bigint NOT NULL,
    PRIMARY KEY (`current`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;
INSERT INTO account_bill_version (current)
VALUES (0);

CREATE TABLE IF NOT EXISTS `iamserver_version`
(
    `current` bigint NOT NULL,
    PRIMARY KEY (`current`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;
INSERT INTO `iamserver_version` (current)
VALUES (0);

CREATE TABLE IF NOT EXISTS `job_version`
(
    `current` bigint NOT NULL,
    PRIMARY KEY (`current`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;
INSERT INTO job_version (current)
VALUES (0);

CREATE TABLE IF NOT EXISTS `license_version`
(
    `current` bigint NOT NULL,
    PRIMARY KEY (`current`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;
INSERT INTO license_version (current)
VALUES (0);

CREATE TABLE IF NOT EXISTS `storage_version`
(
    `current` bigint NOT NULL,
    PRIMARY KEY (`current`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;
INSERT INTO storage_version (current)
VALUES (0);

CREATE TABLE IF NOT EXISTS `hydra_lcp_version`
(
    `current` bigint NOT NULL,
    PRIMARY KEY (`current`)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4;
INSERT INTO hydra_lcp_version (current)
VALUES (0);
