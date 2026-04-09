CREATE TABLE IF NOT EXISTS `ads`
(
    `id`          bigint                                                        NOT NULL AUTO_INCREMENT,
    `title`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Ads title',
    `type`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Ads type',
    `content`     text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'Ads content',
    `target_url`  varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT '' COMMENT 'Ads target url',
    `start_time`  datetime                                                               DEFAULT NULL COMMENT 'Ads start time',
    `end_time`    datetime                                                               DEFAULT NULL COMMENT 'Ads end time',
    `status`      tinyint(1)                                                             DEFAULT '0' COMMENT 'Ads status,0 disable,1 enable',
    `created_at`  datetime(3)                                                            DEFAULT NULL COMMENT 'Create Time',
    `updated_at`  datetime(3)                                                            DEFAULT NULL COMMENT 'Update Time',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Description',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

UPDATE `auth_method`
SET `config` = CASE
                   WHEN JSON_VALID(`config`) AND JSON_EXTRACT(CAST(`config` AS JSON), '$.show_ads') IS NULL
                       THEN CAST(JSON_SET(CAST(`config` AS JSON), '$.show_ads', CAST('false' AS JSON)) AS CHAR)
                   ELSE `config`
    END,
    `updated_at` = CURRENT_TIMESTAMP(3)
WHERE `method` = 'device';
