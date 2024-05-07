CREATE TABLE `post` (
    `id` int NOT NULL AUTO_INCREMENT,
    `content` varchar(255) NOT NULL,
    `title` varchar(255) NOT NULL,
    `ownerId` int NOT NULL,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `created_by` int NOT NULL,
    `updated_at` datetime NULL,
    `updated_by` int NULL,
    `deleted_at` datetime NULL,
    `deleted_by` int NULL,
    PRIMARY KEY(`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC;