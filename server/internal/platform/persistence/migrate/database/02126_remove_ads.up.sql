DROP TABLE IF EXISTS `ads`;

UPDATE `auth_method`
SET `config` = CASE
                   WHEN JSON_VALID(`config`) THEN CAST(JSON_REMOVE(CAST(`config` AS JSON), '$.show_ads') AS CHAR)
                   ELSE `config`
    END,
    `updated_at` = CURRENT_TIMESTAMP(3)
WHERE `method` = 'device';
