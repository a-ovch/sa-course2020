CREATE TABLE `user`
(
    `id`         BINARY(16)   NOT NULL,
    `username`   VARCHAR(255) NOT NULL,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name`  VARCHAR(255) NOT NULL,
    `email`      VARCHAR(255) NOT NULL,
    `phone`      VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;