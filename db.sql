CREATE TABLE `pix_auth` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(50) NOT NULL
);

CREATE TABLE `pix_tag` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `state` tinyint(3) NOT NULL DEFAULT 1,
  `created_at` datetime,
  `created_by` varchar(100),
  `updated_at` datetime,
  `updated_by` varchar(100),
  `deleted_at` datetime,
  `deleted_by` varchar(100)
);

CREATE TABLE `pix_article` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `tag_id` int NOT NULL,
  `title` varchar(100),
  `desc` varchar(255),
  `content` text,
  `state` tinyint(3) NOT NULL DEFAULT 1,
  `created_at` datetime,
  `created_by` varchar(100),
  `updated_at` datetime,
  `updated_by` varchar(100),
  `deleted_at` datetime,
  `deleted_by` varchar(100)
);

ALTER TABLE `pix_article` ADD FOREIGN KEY (`tag_id`) REFERENCES `bbwl_tag` (`id`);
