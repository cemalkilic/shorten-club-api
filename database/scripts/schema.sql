
CREATE DATABASE `shorten`;
use shorten;


CREATE TABLE `records` (
    `id` INT NOT NULL AUTO_INCREMENT ,
    `username` VARCHAR(32) NOT NULL ,
    `slug` VARCHAR(32) NOT NULL ,
    `type` VARCHAR(20) NOT NULL,
    `content` TEXT NOT NULL ,
    `expire_at` DATETIME DEFAULT NULL ,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    `updated_at` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX slug ON records(slug);

CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(32) NOT NULL,
  `password` varchar(64) NOT NULL,
  `createdAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
