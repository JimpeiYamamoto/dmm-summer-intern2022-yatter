CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(255),
  `avatar` text,
  `header` text,
  `note` text,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

CREATE TABLE `status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account_id` bigint(20) NOT NULL,
  `content` text NOT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_account_id` (`account_id`),
  CONSTRAINT `fk_status_account_id` FOREIGN KEY (`account_id`) REFERENCES  `account` (`id`)
);

CREATE TABLE `relation` (
  `follower_id` bigint(20),
  `followee_id` bigint(20),
  PRIMARY KEY (`follower_id`, `followee_id`),
  CONSTRAINT `fk_relation_follower_id` FOREIGN KEY (`follower_id`) REFERENCES  `account` (`id`),
  CONSTRAINT `fk_relation_followee_id` FOREIGN KEY (`followee_id`) REFERENCES  `account` (`id`)
);