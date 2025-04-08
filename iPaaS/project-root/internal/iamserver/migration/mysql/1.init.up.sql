CREATE TABLE `policy`
(
    `id`              bigint(20) NOT NULL AUTO_INCREMENT,
    `userId`          varchar(255) DEFAULT NULL,
    `policyName`      varchar(255) DEFAULT NULL,
    `statementShadow` longtext,
    `version`         longtext,
    `createdAt`       datetime(3) DEFAULT NULL,
    `updatedAt`       datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_user_policy` (`userId`,`policyName`)
) ENGINE=InnoDB AUTO_INCREMENT=1901539978882584578 DEFAULT CHARSET=utf8;

CREATE TABLE `policy_audit`
(
    `id`           bigint(20) NOT NULL AUTO_INCREMENT,
    `subject`      varchar(255) DEFAULT NULL,
    `policyShadow` longtext,
    `createdAt`    timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updatedAt`    timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=253446124 DEFAULT CHARSET=utf8;

CREATE TABLE `role`
(
    `id`                bigint(20) NOT NULL AUTO_INCREMENT,
    `userId`            varchar(255) DEFAULT NULL,
    `roleName`          varchar(255) DEFAULT NULL,
    `description`       longtext,
    `trustPolicyShadow` longtext,
    `createdAt`         datetime(3) DEFAULT NULL,
    `updatedAt`         datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_user_role` (`userId`,`roleName`)
) ENGINE=InnoDB AUTO_INCREMENT=1901539978870001665 DEFAULT CHARSET=utf8;

CREATE TABLE `role_policy_relation`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT,
    `roleId`    bigint(20) DEFAULT NULL,
    `policyId`  bigint(20) DEFAULT NULL,
    `createdAt` datetime(3) DEFAULT NULL,
    `updatedAt` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_role_policy` (`roleId`,`policyId`)
) ENGINE=InnoDB AUTO_INCREMENT=1901539978890973186 DEFAULT CHARSET=utf8;

CREATE TABLE `secret`
(
    `accessKeyId`     varchar(191) NOT NULL,
    `accessKeySecret` longtext,
    `sessionToken`    longtext,
    `expiration`      datetime(3) DEFAULT NULL,
    `parentUser`      longtext,
    `claims_shadows`  longtext,
    `description`     longtext,
    `status`          tinyint(1) DEFAULT NULL,
    `tag`             longtext,
    `createdAt`       datetime(3) DEFAULT NULL,
    `updatedAt`       datetime(3) DEFAULT NULL,
    PRIMARY KEY (`accessKeyId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
