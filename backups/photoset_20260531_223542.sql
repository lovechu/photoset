-- MySQL dump 10.13  Distrib 8.0.45, for Linux (x86_64)
--
-- Host: localhost    Database: photoset
-- ------------------------------------------------------
-- Server version	8.0.45-0ubuntu0.24.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `achievements`
--

DROP TABLE IF EXISTS `achievements`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `achievements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `icon` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `badge_image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `condition_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `condition_value` bigint NOT NULL DEFAULT '0',
  `reward_points` bigint DEFAULT '0',
  `reward_title` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `sort_order` bigint DEFAULT '0',
  `is_hidden` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_achievements_name` (`name`),
  KEY `idx_achievements_deleted_at` (`deleted_at`),
  KEY `idx_achievements_type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `achievements`
--

LOCK TABLES `achievements` WRITE;
/*!40000 ALTER TABLE `achievements` DISABLE KEYS */;
INSERT INTO `achievements` VALUES (1,'2026-05-31 22:32:58.456','2026-05-31 22:32:58.456',NULL,'first_post','初次发帖','发布第一篇帖子','📝','','post','post_count',1,20,'初出茅庐',1,0),(2,'2026-05-31 22:32:58.635','2026-05-31 22:32:58.635',NULL,'post_10','小有成就','累计发布10篇帖子','📚','','post','post_count',10,50,'',2,0),(3,'2026-05-31 22:32:58.794','2026-05-31 22:32:58.794',NULL,'post_50','创作达人','累计发布50篇帖子','✍️','','post','post_count',50,200,'',3,0),(4,'2026-05-31 22:32:58.980','2026-05-31 22:32:58.980',NULL,'post_100','百帖大神','累计发布100篇帖子','🎯','','post','post_count',100,500,'',4,0),(5,'2026-05-31 22:32:59.262','2026-05-31 22:32:59.262',NULL,'first_reply','初次回复','发布第一篇回复','💬','','reply','reply_count',1,10,'热心网友',10,0),(6,'2026-05-31 22:32:59.426','2026-05-31 22:32:59.426',NULL,'reply_100','评论达人','累计发布100篇回复','🗣️','','reply','reply_count',100,100,'',11,0),(7,'2026-05-31 22:32:59.554','2026-05-31 22:32:59.554',NULL,'first_like','初次获赞','收到第一个点赞','❤️','','like','like_received',1,10,'',20,0),(8,'2026-05-31 22:32:59.691','2026-05-31 22:32:59.691',NULL,'like_100','万人迷','累计收到100个点赞','💖','','like','like_received',100,100,'',21,0),(9,'2026-05-31 22:32:59.867','2026-05-31 22:32:59.867',NULL,'like_1000','超级网红','累计收到1000个点赞','🌟','','like','like_received',1000,500,'',22,0),(10,'2026-05-31 22:33:00.029','2026-05-31 22:33:00.029',NULL,'first_follow','初次关注','关注第一个用户','🤝','','follow','following_count',1,10,'',30,0),(11,'2026-05-31 22:33:00.247','2026-05-31 22:33:00.247',NULL,'follower_10','小有人气','拥有10个粉丝','👥','','follow','follower_count',10,50,'',31,0),(12,'2026-05-31 22:33:00.416','2026-05-31 22:33:00.416',NULL,'follower_100','人气之星','拥有100个粉丝','⭐','','follow','follower_count',100,200,'',32,0),(13,'2026-05-31 22:33:00.530','2026-05-31 22:33:00.530',NULL,'level_5','钻石之路','达到5级','💎','','level','level_reached',5,300,'',40,0),(14,'2026-05-31 22:33:00.766','2026-05-31 22:33:00.766',NULL,'level_10','荣耀巅峰','达到10级','🏆','','level','level_reached',10,1000,'',41,0),(15,'2026-05-31 22:33:00.988','2026-05-31 22:33:00.988',NULL,'early_bird','早起的鸟儿','在早上6点前发帖','🐦','','special','special',1,50,'',50,1),(16,'2026-05-31 22:33:01.157','2026-05-31 22:33:01.157',NULL,'night_owl','夜猫子','在凌晨2点后发帖','🦉','','special','special',2,50,'',51,1);
/*!40000 ALTER TABLE `achievements` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `admin_logs`
--

DROP TABLE IF EXISTS `admin_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `admin_id` bigint unsigned DEFAULT NULL,
  `admin_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `action` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `target` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `detail` text COLLATE utf8mb4_unicode_ci,
  `ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_admin_logs_admin_id` (`admin_id`),
  KEY `idx_admin_logs_action` (`action`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin_logs`
--

LOCK TABLES `admin_logs` WRITE;
/*!40000 ALTER TABLE `admin_logs` DISABLE KEYS */;
INSERT INTO `admin_logs` VALUES (1,3,'','ban_user','用户#10','状态改为 0','127.0.0.1',1776257167),(2,3,'','unban_user','用户#10','状态改为 1','127.0.0.1',1776257171),(3,3,'','role_change','用户#10','角色改为 member','127.0.0.1',1776257175);
/*!40000 ALTER TABLE `admin_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `api_keys`
--

DROP TABLE IF EXISTS `api_keys`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `api_keys` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `key` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `secret` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `created_by` bigint unsigned DEFAULT NULL,
  `last_used` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_api_keys_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `api_keys`
--

LOCK TABLES `api_keys` WRITE;
/*!40000 ALTER TABLE `api_keys` DISABLE KEYS */;
/*!40000 ALTER TABLE `api_keys` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sort_order` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_categories_name` (`name`),
  UNIQUE KEY `idx_categories_slug` (`slug`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'2026-04-09 21:01:02.000','2026-04-09 21:01:02.000','自然风光','nature','大自然的美景',1),(2,'2026-04-09 21:01:02.000','2026-04-09 21:01:02.000','人像摄影','portrait','人物肖像摄影',2),(3,'2026-04-09 21:01:02.000','2026-04-09 21:01:02.000','城市街拍','street','城市街头摄影',3);
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comment_likes`
--

DROP TABLE IF EXISTS `comment_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comment_likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `comment_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_comment` (`user_id`,`comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comment_likes`
--

LOCK TABLES `comment_likes` WRITE;
/*!40000 ALTER TABLE `comment_likes` DISABLE KEYS */;
/*!40000 ALTER TABLE `comment_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `photoset_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `image_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `parent_id` bigint unsigned DEFAULT NULL,
  `like_count` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_comments_deleted_at` (`deleted_at`),
  KEY `idx_comments_photo_set_id` (`photoset_id`),
  KEY `idx_comments_user_id` (`user_id`),
  KEY `idx_comments_parent_id` (`parent_id`),
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `drafts`
--

DROP TABLE IF EXISTS `drafts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `drafts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `category` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'discussion',
  `post_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'dynamic',
  `visibility` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'public',
  PRIMARY KEY (`id`),
  KEY `idx_drafts_deleted_at` (`deleted_at`),
  KEY `idx_drafts_user_id` (`user_id`),
  CONSTRAINT `fk_drafts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `drafts`
--

LOCK TABLES `drafts` WRITE;
/*!40000 ALTER TABLE `drafts` DISABLE KEYS */;
/*!40000 ALTER TABLE `drafts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `favorites`
--

DROP TABLE IF EXISTS `favorites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `favorites` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `photoset_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_photoset` (`user_id`,`photoset_id`),
  UNIQUE KEY `uk_user_photoset` (`user_id`,`photoset_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_photoset_id` (`photoset_id`),
  CONSTRAINT `fk_favorites_photo_set` FOREIGN KEY (`photoset_id`) REFERENCES `photosets` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `favorites`
--

LOCK TABLES `favorites` WRITE;
/*!40000 ALTER TABLE `favorites` DISABLE KEYS */;
INSERT INTO `favorites` VALUES (3,'2026-04-08 21:15:25.600',NULL,NULL,10,2),(4,'2026-04-08 21:31:10.927',NULL,NULL,3,4);
/*!40000 ALTER TABLE `favorites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `follows`
--

DROP TABLE IF EXISTS `follows`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `follows` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `following_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_follow` (`user_id`,`following_id`),
  KEY `idx_following_id` (`following_id`),
  CONSTRAINT `fk_follows_following` FOREIGN KEY (`following_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_follows_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `follows`
--

LOCK TABLES `follows` WRITE;
/*!40000 ALTER TABLE `follows` DISABLE KEYS */;
/*!40000 ALTER TABLE `follows` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `memberships`
--

DROP TABLE IF EXISTS `memberships`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `memberships` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `duration` bigint NOT NULL COMMENT '会员天数',
  `price` decimal(10,2) NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `status` tinyint DEFAULT '1' COMMENT '1-on,0-off',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `memberships`
--

LOCK TABLES `memberships` WRITE;
/*!40000 ALTER TABLE `memberships` DISABLE KEYS */;
/*!40000 ALTER TABLE `memberships` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `from_user_id` bigint unsigned NOT NULL,
  `to_user_id` bigint unsigned NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_read` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_messages_from_user_id` (`from_user_id`),
  KEY `idx_messages_to_user_id` (`to_user_id`),
  KEY `idx_messages_is_read` (`is_read`),
  KEY `idx_messages_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_messages_from_user` FOREIGN KEY (`from_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_messages_to_user` FOREIGN KEY (`to_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `is_read` tinyint(1) NOT NULL DEFAULT '0',
  `sender_id` bigint unsigned DEFAULT NULL,
  `target_id` bigint unsigned DEFAULT NULL,
  `target_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_deleted_at` (`deleted_at`),
  KEY `idx_notifications_user_id` (`user_id`),
  KEY `idx_notifications_type` (`type`),
  KEY `idx_notifications_is_read` (`is_read`),
  KEY `idx_notifications_sender_id` (`sender_id`),
  CONSTRAINT `fk_notifications_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifications`
--

LOCK TABLES `notifications` WRITE;
/*!40000 ALTER TABLE `notifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  `membership_id` bigint unsigned DEFAULT NULL,
  `photo_set_id` bigint unsigned DEFAULT NULL,
  `paid_at` datetime(3) DEFAULT NULL,
  `expire_seconds` bigint DEFAULT '1800' COMMENT '支付过期秒数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_orders_order_no` (`order_no`),
  KEY `idx_orders_user_id` (`user_id`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `fk_orders_membership` (`membership_id`),
  KEY `fk_orders_photo_set` (`photo_set_id`),
  CONSTRAINT `fk_orders_membership` FOREIGN KEY (`membership_id`) REFERENCES `memberships` (`id`),
  CONSTRAINT `fk_orders_photo_set` FOREIGN KEY (`photo_set_id`) REFERENCES `photosets` (`id`),
  CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pages`
--

DROP TABLE IF EXISTS `pages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pages` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `slug` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content_md` text COLLATE utf8mb4_unicode_ci,
  `user_id` bigint unsigned NOT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'published',
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_pages_slug` (`slug`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_pages_user_id` (`user_id`),
  KEY `idx_pages_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pages`
--

LOCK TABLES `pages` WRITE;
/*!40000 ALTER TABLE `pages` DISABLE KEYS */;
INSERT INTO `pages` VALUES (1,'about','关于我们','# 关于我们\n\n我们是摄影爱好者社区，致力于分享美丽瞬间。',1,'published',20260415176,20260415176),(2,'terms','使用协议','# 使用协议\n\n请遵守平台规则。',1,'published',20260415176,20260415176),(3,'privacy','隐私政策','# 隐私政策\n\n我们非常重视您的隐私。',1,'published',20260415176,20260415176),(4,'faq','常见问题','# 常见问题\n\n### 如何上传作品？\n请注册并进入创作者后台。',1,'published',20260415176,20260415176),(5,'contact','联系我们','# 联系我们\n\n邮箱：support@photoset.io',1,'published',20260415176,20260415176),(6,'ceshi','测试','测试',3,'published',1776248690,1776248690);
/*!40000 ALTER TABLE `pages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `password_reset_tokens`
--

DROP TABLE IF EXISTS `password_reset_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `password_reset_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `used` tinyint(1) DEFAULT '0',
  `expire` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_password_reset_tokens_token` (`token`),
  KEY `idx_password_reset_tokens_user_id` (`user_id`),
  KEY `idx_password_reset_tokens_email` (`email`),
  KEY `idx_password_reset_tokens_expire` (`expire`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `password_reset_tokens`
--

LOCK TABLES `password_reset_tokens` WRITE;
/*!40000 ALTER TABLE `password_reset_tokens` DISABLE KEYS */;
INSERT INTO `password_reset_tokens` VALUES (1,'2026-04-16 19:17:03.349','2026-04-16 19:17:03.349',3,'root@honker.org','b1fec09f7fa91c9a5deadebdc6f363fe7bfc626ac749fc87504650edde29c9d6',0,'2026-04-16 19:47:03.349');
/*!40000 ALTER TABLE `password_reset_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `photos`
--

DROP TABLE IF EXISTS `photos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `photos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sort_order` bigint DEFAULT '0',
  `photoset_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_photos_deleted_at` (`deleted_at`),
  KEY `fk_photosets_photos` (`photoset_id`),
  CONSTRAINT `fk_photosets_photos` FOREIGN KEY (`photoset_id`) REFERENCES `photosets` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=460 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `photos`
--

LOCK TABLES `photos` WRITE;
/*!40000 ALTER TABLE `photos` DISABLE KEYS */;
INSERT INTO `photos` VALUES (1,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/0b66dd95-dc78-4f7a-8a93-9bc3e5847ce2.jpg',0,5),(2,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/eb698099-ae66-4cb1-80df-695019004570.jpg',1,5),(3,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/30eac5e5-d6d6-4a02-ac43-461fd8e8ac5a.jpg',2,5),(4,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/f30a068c-f129-45f5-8aa7-50553066d0eb.jpg',3,5),(5,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/d0962b78-5863-4f60-9882-a6f97ddafc30.jpg',4,5),(6,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/fef4774e-8e8a-4d27-baa7-59fa65cbaaa4.jpg',5,5),(7,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/b002fb99-2acf-452c-8d30-58863f3afc75.jpg',6,5),(8,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/c22cea88-41f3-477c-903b-7d4426ff06d3.jpg',7,5),(9,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/8e891929-3db2-48a0-86d0-0169a13f2a21.jpg',8,5),(10,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/6b9dd4f7-c2e1-413d-b7d6-1abd85a7e5c8.jpg',9,5),(11,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/f9c24417-17ed-4616-84a1-54f5341f8698.jpg',10,5),(12,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/60d364f5-aca2-4574-9081-b76b2ca9a658.jpg',11,5),(13,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/22a63b75-7613-49a4-ae84-7ec12f2082bb.jpg',12,5),(14,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4582174d-7df8-4052-a2d6-c261583db259.jpg',13,5),(15,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/b2a3202b-7383-48e3-8db3-7b84e00ae86c.jpg',14,5),(16,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/6b236b85-1fd1-44dc-88d7-274f1d74c17c.jpg',15,5),(17,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/68897063-ab67-4a3f-9c7a-8d733a784c66.jpg',16,5),(18,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/1f62c691-3215-4eba-b515-6fe2b022dc64.jpg',17,5),(19,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/68263bad-dcff-48d4-b224-aa310fca7846.jpg',18,5),(20,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/fdac087c-c6fa-4d75-99eb-7afb2cd6cfdc.jpg',19,5),(21,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/41a3aeba-ce9b-420d-8608-5d33de36e8dc.jpg',20,5),(22,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/cf03cba7-38f3-4788-bde6-c651edbc3649.jpg',21,5),(23,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/6a0154a7-f1d4-4488-869b-a5ac519ec3b5.jpg',22,5),(24,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4d1034e5-56c7-4aa4-879d-9e67d9f62675.jpg',23,5),(25,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/6e2ec9b2-6ef5-4507-a66c-5b8a856c6dd9.jpg',24,5),(26,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/5108a08a-b25d-4e4d-b9e7-23c84bab3b62.jpg',25,5),(27,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/369849f7-a9ec-4d9c-a1f5-5e9bd17155ef.jpg',26,5),(28,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/f6034cb9-d311-4a9c-a836-41a1a11e2208.jpg',27,5),(29,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4db7c36a-e107-4b19-89ff-90cff9038ae3.jpg',28,5),(30,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/76923a88-561f-4c18-b476-fe66701c1e2f.jpg',29,5),(31,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/da30a73f-d70a-409d-a12e-6734e9109aa4.jpg',30,5),(32,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4c979498-b039-49e5-8f15-a03d61b1a7fb.jpg',31,5),(33,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/63b991ad-2356-4ae2-8b91-25fefce433f9.jpg',32,5),(34,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/c25ce6ec-93e6-491f-bea7-b21255e61f11.jpg',33,5),(35,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4a03f5ba-7c0e-4694-9120-4d890660ce44.jpg',34,5),(36,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/0bd30f2e-e965-4afc-a663-2f2933ae6612.jpg',35,5),(37,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/cbaf175e-a843-4122-9b15-10cbfb8df5e1.jpg',36,5),(38,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/ef5a086a-977f-4d88-8495-e72c052fadcd.jpg',37,5),(39,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/c8b37912-7219-4c32-a9be-371d48b416b1.jpg',38,5),(40,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/d32183f8-5e2f-4dfe-b2b8-940f1954a232.jpg',39,5),(41,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4c4473f8-0c7a-4761-88f6-0e32b6962c39.jpg',40,5),(42,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/6b918462-2d4e-4b34-8569-a5999db9dd1a.jpg',41,5),(43,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/433429c4-ca9c-4770-846d-9d6ba8ea2df0.jpg',42,5),(44,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/bb23ed4b-dd3b-45c4-81b3-e97d64104647.jpg',43,5),(45,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4d8bb9a1-3955-4df1-b003-84fa7a82e6bc.jpg',44,5),(46,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/9b0a3fe2-5720-4cb6-b217-cb3a5bf3a7d2.jpg',45,5),(47,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/b0a7fe2b-05f4-4f33-abd3-4a0770abfd48.jpg',46,5),(48,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/961da0ed-8ce4-4eaa-af00-bd3c9e79391d.jpg',47,5),(49,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4f183ebe-11f1-4337-bb46-3c055b5a8cf4.jpg',48,5),(50,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/e9c27cbb-2247-49ff-8986-6d7ac8016f46.jpg',49,5),(51,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/fa7f881f-71f1-491e-b915-f6482ea91737.jpg',50,5),(52,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/8859e0f0-8b13-4e2a-832a-e76746e09348.jpg',51,5),(53,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/db92f1e6-4630-4b08-9d39-e7eec6f6887d.jpg',52,5),(54,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/ecc5059c-ca28-4b21-9ac3-f68d8d56dda8.jpg',53,5),(55,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/71579db6-f27b-46db-915e-a1917a8dc0ff.jpg',54,5),(56,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/f6b99137-2f24-4b59-a09c-e51ae29de230.jpg',55,5),(57,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/1e971dbe-50b7-43dd-8ad5-2a870f7ce051.jpg',56,5),(58,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/a1fb96ff-0c2c-43f3-8902-fd25a6f4a81b.jpg',57,5),(59,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/f86def03-170b-4005-b4ce-179552a6e82d.jpg',58,5),(60,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/a72c2b9a-b40d-411d-8d06-7893a16d9ca3.jpg',59,5),(61,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/6c5426a6-25ec-4d12-9af3-057a24a980c0.jpg',60,5),(62,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/e175c1d3-48f9-4f56-b331-19e447024b54.jpg',61,5),(63,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/9ce1a3ce-a1bd-4676-a8f2-729a26a71163.jpg',62,5),(64,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/5d027d97-7700-450d-b7a4-610555cb4b50.jpg',63,5),(65,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/f0e48c00-33e8-442c-8507-ec5d09b3e330.jpg',64,5),(66,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/594b3837-404e-4bf6-9744-faf65e5f0cea.jpg',65,5),(67,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/9ae93aa2-f4cd-4c3d-a423-1a35715736ae.jpg',66,5),(68,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/6286b609-49c9-49e1-869c-7558ad9a715e.jpg',67,5),(69,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/53566119-d442-4f95-86da-6e0a44076229.jpg',68,5),(70,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/45f30f28-9284-49f9-9084-9f06d937f6af.jpg',69,5),(71,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/4839a4df-3991-4956-ba27-0647f9c21a91.jpg',70,5),(72,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/eb7a3247-2fbc-4224-8bdf-3688051dcede.jpg',71,5),(73,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/c1dc0924-9fd4-452d-b4e7-ba0435df40d0.jpg',72,5),(74,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/795b3d9f-b6e7-4420-9514-1ec0b541ff7f.jpg',73,5),(75,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/5fcfbd76-ef3d-48be-a211-df392e3f6607.jpg',74,5),(76,'2026-04-09 23:02:41.486',NULL,'/uploads/images/2026/04/09/a8fbc1d1-c650-46db-a98f-6662376b5738.jpg',75,5),(193,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/e911b175-f167-4142-90c0-ed8b2869a547.jpg',0,7),(194,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/bc33baa0-b7d5-4f7d-960c-c469606ca1d8.jpg',1,7),(195,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/d2a53737-801d-44ed-a11d-283aae4df943.jpg',2,7),(196,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/2cc18ca5-97e9-4aa1-ae2f-d893db3d5ef7.jpg',3,7),(197,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/7934f5bd-f126-4ed4-a7c9-bab8ebf3dc7d.jpg',4,7),(198,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/88695d2c-982a-4766-ac4f-bd0a12108e93.jpg',5,7),(199,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/698addad-f8fa-4b1d-9707-c036217a0c32.jpg',6,7),(200,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/257f1950-5f4a-47fd-968a-f5f62925cdcd.jpg',7,7),(201,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/73618eb2-6dd7-4b23-8bab-146705c88623.jpg',8,7),(202,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/2cd67b5f-1f1e-40a2-88f6-d64c96d94d22.jpg',9,7),(203,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/f36328c5-4028-45b2-93bf-359241779411.jpg',10,7),(204,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/1bce3cd3-eb1c-4f68-ace6-4c2ed7e8979a.jpg',11,7),(205,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/152e90a7-559b-4966-9bda-7787e07e150e.jpg',12,7),(206,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/e06dc3f1-c103-454e-9562-1b65a4e242e9.jpg',13,7),(207,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/7778c06f-f3e6-4b2e-8c60-8e7b19d8d0f4.jpg',14,7),(208,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/54806435-2467-416f-95df-ca90f14f091e.jpg',15,7),(209,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/bcfcf253-cfc5-4083-bd7c-929c33d73c42.jpg',16,7),(210,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/9eea5255-c074-4ff8-aea8-e5b6f8be833b.jpg',17,7),(211,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/985de4ad-5873-4a35-9215-18ac83449a73.jpg',18,7),(212,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/39efe868-4f91-4e18-8e70-ae9cf241a63b.jpg',19,7),(213,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/fc582b69-ac28-4124-a584-dce310dcc26b.jpg',20,7),(214,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/f19aed33-3807-4056-ba6e-95e4eb2753d5.jpg',21,7),(215,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/9d774ea0-5954-47a5-bebf-553f4d36ef0d.jpg',22,7),(216,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/5340c4fa-406d-4d60-80d5-8062d3ba891b.jpg',23,7),(217,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/1b05ed38-d17d-4be7-91e9-f5b430593954.jpg',24,7),(218,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/d80027de-7d30-456f-ab7d-2e69a618b718.jpg',25,7),(219,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/f2d914ef-4a5f-48d9-88e1-c854b7ea497d.jpg',26,7),(220,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/2d5479a8-4573-4415-90b7-335acf4a7cbe.jpg',27,7),(221,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/d09988de-a32e-4ec9-88d4-4d41f07fed90.jpg',28,7),(222,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/946914d1-413a-4838-80ac-f121adcb3fa9.jpg',29,7),(223,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/4e27dadd-9fba-40d7-9bc0-5ff0538eb269.jpg',30,7),(224,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/f374b8ae-bdaf-4fa6-b1cb-a5b47c4324fe.jpg',31,7),(225,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/eb1b6dee-0c13-485e-ab8d-44ddfb0d09f8.jpg',32,7),(226,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/950a782b-c7f8-4492-aebe-4251ca88824d.jpg',33,7),(227,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/16581660-63ce-451e-8a6e-1904d9f7d464.jpg',34,7),(228,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/ba0b1078-00ba-451e-bae9-61a111260c3e.jpg',35,7),(229,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/a1d0faea-f4dd-4530-af11-26fffd301604.jpg',36,7),(230,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/e292de11-e949-47b7-928f-278a7501d820.jpg',37,7),(231,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/b208bc42-6609-4c92-948a-09446e36ce8a.jpg',38,7),(232,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/3b681c47-3e85-487b-b2e8-7f83004dd369.jpg',39,7),(233,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/21184dcb-81d9-4695-a654-b6a670fbde0b.jpg',40,7),(234,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/c8b3e3f6-9da6-43fc-91de-54cd7f9848cf.jpg',41,7),(235,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/f2fbc485-e590-4efa-8379-c6be318eb906.jpg',42,7),(236,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/c100762d-7d76-4772-a650-4f90f8f8e85e.jpg',43,7),(237,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/9f7f91d7-9608-483e-af4f-c4351ab44fe9.jpg',44,7),(238,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/cf5eb494-a10a-4b21-89c1-4614bd2c9bc8.jpg',45,7),(239,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/11e9e0f6-02a6-4f60-af04-0a0c14af8a62.jpg',46,7),(240,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/a7aabf0d-6937-47bd-a91d-39aa86616557.jpg',47,7),(241,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/dbb4afd5-687f-4eaf-95be-25bf6b8f8bb1.jpg',48,7),(242,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/ac270801-4000-4e78-b0eb-a9a4ecde0b72.jpg',49,7),(243,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/13535dec-9d64-4f63-87c4-7b20de83dc6b.jpg',50,7),(244,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/95cde964-5682-486b-8019-f4150922e372.jpg',51,7),(245,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/1611cd9b-5049-4872-a92c-ec6e30f0f1d3.jpg',52,7),(246,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/70f7035b-a512-4983-ba1b-3039283cf2c5.jpg',53,7),(247,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/19f00629-c663-44e3-8614-16002fe04530.jpg',54,7),(248,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/4e7c6c05-cc90-4e48-a346-42bdfeddc325.jpg',55,7),(249,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/3327b36d-3c6a-4875-b421-5d7cd7e3258a.jpg',56,7),(250,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/f2455089-0f45-4d0b-91fd-e1e90d074bea.jpg',57,7),(251,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/9e51021d-2b7e-4aea-af94-f945210ff184.jpg',58,7),(252,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/6131c736-d8a5-46ed-a874-a12c4dc14cf5.jpg',59,7),(253,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/5c62e211-4e45-4a42-84b6-cce3a0a849c6.jpg',60,7),(254,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/eb335eac-c5d3-4e13-b2e6-49fe01bb7adf.jpg',61,7),(255,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/2b13b47f-3431-406a-bf6e-735f13e3dd99.jpg',62,7),(256,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/a7ccfa2e-0735-4797-8bc5-1f9f301aeca1.jpg',63,7),(257,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/fc29bb8b-9cee-4962-ab36-4d6fde4eb16f.jpg',64,7),(258,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/ac81a455-48b8-4b9a-9211-a2881fb20ec1.jpg',65,7),(259,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/9fd1020e-858c-4e19-9cc1-cdc3d439e479.jpg',66,7),(260,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/f36b63bc-3a04-419e-8806-3f6f5a81ccfd.jpg',67,7),(261,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/91a0e800-a4a1-4c67-978c-506b9798ac06.jpg',68,7),(262,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/14745e3e-d159-4984-92c6-a8a1399215ea.jpg',69,7),(263,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/be53d1b8-fe58-4d9b-ad2f-c57cc67b9694.jpg',70,7),(264,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/aec02e5d-c11d-4a7d-958d-164a16dd11d3.jpg',71,7),(265,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/6d4a6ae5-b1b9-4748-b9f9-dbb736aebec4.jpg',72,7),(266,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/efd56d06-e15f-4dc7-88b4-fbbf24928788.jpg',73,7),(267,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/c7346c44-c1e1-4a85-aba5-46fff9532fbf.jpg',74,7),(268,'2026-04-09 23:40:10.546',NULL,'/uploads/images/2026/04/09/aeb0e9a6-c982-4346-b9dd-07d33273e8f1.jpg',75,7),(273,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160326/04/7cc8e7cc-cdd1-4417-88ee-00b07d350664.jpg',0,6),(274,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160326/04/aedc2084-f1bb-4fd1-b9f3-683e3cdeb15a.jpg',1,6),(275,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160326/04/1467d9d3-222d-4f02-b8f2-70a065af144f.jpg',2,6),(276,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160326/04/1ab727e6-9a3e-4637-9c36-bff3ea6ba8e1.jpg',3,6),(277,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160326/04/c0af1dc9-1cfe-48b1-ada5-69e4c24f0b88.jpg',4,6),(278,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160326/04/de4b4493-3271-4483-a008-ccc80eaeee7d.jpg',5,6),(279,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160326/04/d277e403-f614-4e3f-8c33-08641e162f19.jpg',6,6),(280,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/659aef80-fd13-4490-a49f-546e999bd7c1.jpg',7,6),(281,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/b4ac1260-260d-4866-8bc1-6cbd787777e3.jpg',8,6),(282,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/8430877c-0cec-46ed-bd82-5bd1655c0147.jpg',9,6),(283,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/8d6099e6-db80-48a6-bce2-26827206c7f7.jpg',10,6),(284,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/456f5f3a-c3ec-4d46-8b89-8fc6cae843b5.jpg',11,6),(285,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/9a6e8971-dc40-4ac6-922a-a1083b7d09ef.jpg',12,6),(286,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/13a6212b-3741-44d3-806f-38ef7af6ecb2.jpg',13,6),(287,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/c744b024-8c96-4440-8cff-e19d4db3210c.jpg',14,6),(288,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/4aa1ae70-2b20-4d7a-bd8f-e681b0ed3a8c.jpg',15,6),(289,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/7cad2ae5-23ab-439b-8693-68f44783eed1.jpg',16,6),(290,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/aff51ef1-45c6-47f8-8b7c-9552b6c95f65.jpg',17,6),(291,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/34a1919b-f61f-44e4-8989-c2cd51c9719f.jpg',18,6),(292,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/44a8d33a-9c7e-4fb7-9f83-40c83ac619e7.jpg',19,6),(293,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/eb10dfc9-11e7-4968-a1aa-1cd72010773e.jpg',20,6),(294,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/6e8b23ec-d991-43c3-87a3-fad704361cd1.jpg',21,6),(295,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/a5877b3d-18e8-4375-890a-fa9d6e344e99.jpg',22,6),(296,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160327/04/51e8466a-7655-4cea-91e1-02df8c66510b.jpg',23,6),(297,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/d27e639a-ec74-4338-b817-8f24870ef90f.jpg',24,6),(298,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/6cc341b8-491c-40a2-b71e-6295767c9bd4.jpg',25,6),(299,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/32486420-3cdc-4eb3-875b-4a165bb9d65e.jpg',26,6),(300,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/62b96623-4f2e-4af5-abc7-c7b53ae11dde.jpg',27,6),(301,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/75d0efee-bd50-4062-b8ef-2086d96bc999.jpg',28,6),(302,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/f97aaa49-d971-42d1-b940-601305112fe6.jpg',29,6),(303,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/b0a77562-c158-4413-af7f-4b8525788ba8.jpg',30,6),(304,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/a98c055f-3bee-4ff1-93ae-22f796511d33.jpg',31,6),(305,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/57375821-e8d0-4e54-8f99-ef8757580f33.jpg',32,6),(306,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/4ce267db-6852-43e2-8820-0bfed261eebd.jpg',33,6),(307,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/adb0bd59-58cd-4323-bc9b-ee8c1155a559.jpg',34,6),(308,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/ad8b9f1b-4f60-4869-b336-c3c51d9e0c98.jpg',35,6),(309,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/5777c0f2-6ce8-4c45-9141-92f8106f43f0.jpg',36,6),(310,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/5ca7d562-e518-4d0e-837b-c0dec159703f.jpg',37,6),(311,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/60441807-98cf-46c9-afb9-398ca3d6bc59.jpg',38,6),(312,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/eac8aa02-0991-44e4-b9b0-f09a50eff0c1.jpg',39,6),(313,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/4b8fd733-189a-4ad0-92e3-72d3675b6903.jpg',40,6),(314,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160328/04/28b90126-46be-4be2-88bf-7f77af4c65fa.jpg',41,6),(315,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/b9a6677e-4144-4622-a7fb-461e9cdfef96.jpg',42,6),(316,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/abe400d8-4b78-44cf-84a5-659122f9e43b.jpg',43,6),(317,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/1ae92da4-315a-4d9d-aa01-23e115fd4aee.jpg',44,6),(318,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/dce06628-5e10-47ee-a250-820e33d9c707.jpg',45,6),(319,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/4e32d72a-70dc-44e8-bfd5-a525fb4522b8.jpg',46,6),(320,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/0ee63186-a688-43a9-8ab2-a781a1f663bf.jpg',47,6),(321,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/677dbe11-7b25-4c21-b36f-62821f637f42.jpg',48,6),(322,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/726db1f2-045a-4d1d-aa92-d44d5090a32d.jpg',49,6),(323,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/70e8b3dd-33d0-465f-a2b5-cd44070d5d74.jpg',50,6),(324,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/17834d1d-de81-4fe6-8d3c-9ff4f1a4841a.jpg',51,6),(325,'2026-04-16 16:03:38.618',NULL,'/uploads/photos/20260416160329/04/7b34d56f-a3cc-4c9f-abf8-92bd20ca1305.jpg',52,6),(326,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162515/04/a265c5e3-c2df-4cfa-8cca-c03009305678.jpg',0,3),(327,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/72970d74-6f81-4856-b2bf-e89e5cb960bb.jpg',1,3),(328,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/6abff22f-6b7c-43f8-be3a-f1088ce7fa27.jpg',2,3),(329,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/6431b160-af04-4906-b1e6-511f00ed6f58.jpg',3,3),(330,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/e9cc58bb-28c1-40be-b601-94c02dd89e23.jpg',4,3),(331,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/7753ec33-7b44-45ea-80da-bea9eff9d95f.jpg',5,3),(332,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/3c7362b4-70a5-43f6-8a47-56b668204641.jpg',6,3),(333,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/6630c176-0a5a-4c60-88c8-5978d1c6be25.jpg',7,3),(334,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/7b56616a-37ef-4a09-8bbd-ea1d9a2c8a39.jpg',8,3),(335,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/8114ce0a-695f-421e-8868-d8d6b0a08adc.jpg',9,3),(336,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/f1de7f48-8bed-4529-9494-930e06f01381.jpg',10,3),(337,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/0b8fe380-bcb3-4d56-9821-c8a5854381a5.jpg',11,3),(338,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/0bd95752-06ae-486a-9d3b-22b50f2bd84d.jpg',12,3),(339,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162516/04/6c5cd729-b496-4c12-a3c3-6f20eccd4e53.jpg',13,3),(340,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/244eecc4-a32b-44d4-a118-2bd0a09de297.jpg',14,3),(341,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/ed9a3e90-5c56-46b6-baad-e37803d8b47d.jpg',15,3),(342,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/765c39e0-07c5-4e44-a981-f0eb4c0fdf8c.jpg',16,3),(343,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/0fae9d8b-1416-4407-af6e-8e075efe1915.jpg',17,3),(344,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/1148bd21-0ec7-4b7e-913f-e7ffe4a5c986.jpg',18,3),(345,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/9228827f-ef7e-456a-aeb9-a120fc1285ce.jpg',19,3),(346,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/e30d5661-d8df-4c3c-980b-e35fe33e9cc5.jpg',20,3),(347,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/2d67663d-7d6d-419e-a541-ac2f2997b450.jpg',21,3),(348,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/cd02eda0-a93d-493c-b544-6cf783b9e60a.jpg',22,3),(349,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/77f367b6-3377-43f0-813e-a0530b49bf39.jpg',23,3),(350,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/77735ddd-faa9-4272-815a-03b23bf577db.jpg',24,3),(351,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/368827f2-b5e9-4adf-aaba-34877ca8bdc5.jpg',25,3),(352,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/d205aaba-009f-4911-8fc5-035b876fbf61.jpg',26,3),(353,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/1e25f618-05ab-4c89-af5a-5608ed7c2b1c.jpg',27,3),(354,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/5a412b27-ab7e-4207-95a1-3bf79674538f.jpg',28,3),(355,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/968678cb-8ed1-4edd-b497-13e3fdc8b58b.jpg',29,3),(356,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162517/04/7cf34fe3-a020-4a32-bcf8-1569efb91a0c.jpg',30,3),(357,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/7000b7e0-c49f-46e4-9bc9-ac024add2f4a.jpg',31,3),(358,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/37b85b43-c820-4c18-8007-0c0e715e01d5.jpg',32,3),(359,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/a28aac01-1f89-4695-ad60-1f261ed1cabc.jpg',33,3),(360,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/26ea239b-aa02-45dc-8ff0-c34c7c14b226.jpg',34,3),(361,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/ebfe273e-0082-4586-a790-2e2c59e3be5e.jpg',35,3),(362,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/75bd8adc-23fd-4a6d-935a-abf2df7a6f8c.jpg',36,3),(363,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/bd8b3380-1ba7-499c-816b-ffc98a7a37af.jpg',37,3),(364,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/6609f876-5de3-437a-9a88-2d0c6ff2ec94.jpg',38,3),(365,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/37814634-7473-4a76-b6f2-5abaf8bd1382.jpg',39,3),(366,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/5ada5890-974b-4bf5-8117-942e61b3a6d8.jpg',40,3),(367,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/3046da94-0f4b-479a-8d63-ea8672a36783.jpg',41,3),(368,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162518/04/00554938-185b-4588-a281-f273e1c924c5.jpg',42,3),(369,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/ae1058e9-7bbc-48e8-9dde-c912e654e9fd.jpg',43,3),(370,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/65faf6e3-87c2-4687-adbd-33fc117516f7.jpg',44,3),(371,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/8072ccb5-bc0a-471f-8e29-ae30cd2e71aa.jpg',45,3),(372,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/5acadf71-e7ba-4fbe-bc54-cab8269840b7.jpg',46,3),(373,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/7e84787a-4070-4cf2-8d21-45f5dd73104a.jpg',47,3),(374,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/ff62d1ff-86d8-4b71-a252-dd23ddb58acc.jpg',48,3),(375,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/b0d08051-7acb-4433-bb62-38383ac43e30.jpg',49,3),(376,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/05ce2fbd-8225-42bb-b5fd-9ab2cd4989e1.jpg',50,3),(377,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/8c7f3146-2be6-4134-9f35-c2b9bb188625.jpg',51,3),(378,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/81ddf354-ef53-4423-b0eb-b5df70c8f8a1.jpg',52,3),(379,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/4425403b-3fbc-4cf9-a03d-ac98a7a2085d.jpg',53,3),(380,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/5ec391a3-8317-4ce1-ab1d-1a6d3b1c7d66.jpg',54,3),(381,'2026-04-16 16:25:36.400',NULL,'/uploads/photos/20260416162519/04/05c38850-cd2e-4e81-83ab-f36f34d5a1aa.jpg',55,3),(382,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174339/04/8eb3b3f4-791a-4343-903b-60e746196116.jpg',0,8),(383,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174339/04/21c3b0f6-ab97-42dd-ab34-28ff66012f71.jpg',1,8),(384,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174339/04/1140e270-a690-42fc-984b-740c2b928ae6.jpg',2,8),(385,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174339/04/0a932436-9b32-43f7-804b-af9046e997a7.jpg',3,8),(386,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174339/04/9dd18de4-89a2-490f-8a8e-9db3022038e8.jpg',4,8),(387,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174339/04/bccf45e1-4252-4f39-9b6e-fc989ed45029.jpg',5,8),(388,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174339/04/b8e4efb6-f12b-4205-a83f-a6dee01b1f0d.jpg',6,8),(389,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/99f9e899-e672-438c-97ed-5ec8c1a09cc9.jpg',7,8),(390,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/77e2c1cb-2980-4ddb-9d64-1df819f9f169.jpg',8,8),(391,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/a88524da-2711-41d5-a6ff-a198cf8fd9a0.jpg',9,8),(392,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/4aaa1ff4-02aa-4b30-a0d7-845709692a1c.jpg',10,8),(393,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/f72ee517-f5b7-4630-9575-eba230c57cec.jpg',11,8),(394,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/26afb31e-4786-4edc-9d0d-0cc4834398bb.jpg',12,8),(395,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/862df91d-9f6d-43ac-a071-c065461e4b5c.jpg',13,8),(396,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/b12e0338-623a-4215-bc82-e055f594182e.jpg',14,8),(397,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/b789d39b-d72b-4875-a65e-615965be3122.jpg',15,8),(398,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/0424a5a5-fde7-4e90-8a8b-9599555b0dbe.jpg',16,8),(399,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/0c63c9a3-4f09-46f5-a493-28e29ac6636d.jpg',17,8),(400,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/a3ff98f9-28ca-4f20-b750-1a6453157a4f.jpg',18,8),(401,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174340/04/1f4a7f42-1d80-4fab-a901-979e2cad8c3e.jpg',19,8),(402,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/2b46cc13-51c7-4349-af9c-0f0ed7d1cb57.jpg',20,8),(403,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/4bbd62c9-e6e7-4593-8a5e-9cf130abcbef.jpg',21,8),(404,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/795747e9-8ae6-4623-b772-9d4c8b0e28ba.jpg',22,8),(405,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/a91cd524-f798-4f3e-b077-4f24984300b1.jpg',23,8),(406,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/ca64fed1-e7df-4cb7-9b80-6c663398c520.jpg',24,8),(407,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/835ed691-08e4-4a47-9869-2845e8d0728d.jpg',25,8),(408,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/26fa1527-d0e1-4150-9450-3bf9dfca72ac.jpg',26,8),(409,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/cec1ac01-78fe-49bf-8ed2-1e588ef1308b.jpg',27,8),(410,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/28c1f31a-da98-4ed0-ae89-5b3a5e39937c.jpg',28,8),(411,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/8e146a43-abec-4e14-819d-05eb7545d411.jpg',29,8),(412,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/1b35f3ec-0f6e-4179-bcca-34ce7f35f18d.jpg',30,8),(413,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/5b0064b1-4517-4903-954d-fcdbd8577517.jpg',31,8),(414,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/52b9c5a7-1aed-4d0f-816b-88c1e4e92bda.jpg',32,8),(415,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174341/04/8d662613-dc3a-4ff6-89e2-846d207f0dcf.jpg',33,8),(416,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/fa68ef38-6f0f-4442-ba34-194e3dbb3bcc.jpg',34,8),(417,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/2c89d9f3-725d-4d9a-a268-121e1422bba3.jpg',35,8),(418,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/b3c4e33f-2c0f-4993-a1b4-2821342f44b7.jpg',36,8),(419,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/01e6e9ae-efc4-4234-99b3-379aa421aeea.jpg',37,8),(420,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/eadeeda3-fd43-4d4a-9588-58d8ee22acb4.jpg',38,8),(421,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/0c2915f3-0890-4a90-a71f-e7a6602aab21.jpg',39,8),(422,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/bcba6f07-44e9-4d4b-a304-1d29ea444493.jpg',40,8),(423,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/2fb0a49b-c78e-4eef-ab3a-8ac009c00189.jpg',41,8),(424,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/fe821708-4364-4ba4-9877-932c7d11a5c6.jpg',42,8),(425,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/1fbb4fd5-c6cb-4e0d-9a3f-1575ba56d583.jpg',43,8),(426,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/70258865-64f5-42ba-96dc-3c723fb588b3.jpg',44,8),(427,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174342/04/9e95f2f0-85db-4169-b6c3-23e29f27cfa3.jpg',45,8),(428,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174344/04/7219e0ef-c352-4843-9934-9009cffe2b16.jpg',46,8),(429,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174344/04/f861986c-c5f4-4388-927b-ae2c1bee5aa3.jpg',47,8),(430,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174344/04/3afc1f66-0a7e-4229-ad43-1a1aad07da8d.jpg',48,8),(431,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174344/04/05e6b8ed-fea7-4ef7-80ae-48b03548e48b.jpg',49,8),(432,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174344/04/cf5f6511-0a9f-403d-876b-495bb8544172.jpg',50,8),(433,'2026-04-16 17:44:40.733',NULL,'/uploads/photos/20260416174344/04/3472fe35-4602-45f8-a696-4853af06fbc0.jpg',51,8),(447,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174620/04/4939f4b4-a68f-461e-b06e-009b78f7a5bf.jpg',0,9),(448,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174620/04/af8fab70-5a1a-4b09-8cb3-45b0f493a593.jpg',1,9),(449,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174621/04/607d945f-c698-43b2-be7f-c6880b751b4d.jpg',2,9),(450,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174621/04/e5c41964-d71c-4e40-9b25-06f431148bc8.jpg',3,9),(451,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174621/04/83855b2b-5826-4f9a-b9b1-9962f5cba3f9.jpg',4,9),(452,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174621/04/2b12c0d1-e60b-4f36-94ca-fadceb183cd1.jpg',5,9),(453,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174625/04/678ef891-6669-4ef3-9260-bb9f49f45a0a.jpg',6,9),(454,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174627/04/768f6808-3648-4bbc-927a-89d7af2a38b0.jpg',7,9),(455,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174628/04/2b15c6e6-0bdf-440f-b062-0f741cac63bb.jpg',8,9),(456,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174628/04/6c9d937f-6aa7-462c-b102-cd873835ed01.jpg',9,9),(457,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174628/04/ad15c67a-7d12-40f1-a57a-fcb991effb09.jpg',10,9),(458,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174629/04/db539b9a-d907-4b0a-9bbe-1a33d6a45abd.jpg',11,9),(459,'2026-04-16 17:49:00.056',NULL,'/uploads/photos/20260416174629/04/605f826f-5fd9-44f9-8705-c56850674799.jpg',12,9);
/*!40000 ALTER TABLE `photos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `photoset_tags`
--

DROP TABLE IF EXISTS `photoset_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `photoset_tags` (
  `tag_id` bigint unsigned NOT NULL,
  `photoset_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`tag_id`,`photoset_id`),
  KEY `fk_photoset_tags_photo_set` (`photoset_id`),
  CONSTRAINT `fk_photoset_tags_photo_set` FOREIGN KEY (`photoset_id`) REFERENCES `photosets` (`id`),
  CONSTRAINT `fk_photoset_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `photoset_tags`
--

LOCK TABLES `photoset_tags` WRITE;
/*!40000 ALTER TABLE `photoset_tags` DISABLE KEYS */;
INSERT INTO `photoset_tags` VALUES (3,4),(4,4),(1,5),(5,7),(5,8),(5,9);
/*!40000 ALTER TABLE `photoset_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `photosets`
--

DROP TABLE IF EXISTS `photosets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `photosets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `cover` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `photo_count` int unsigned NOT NULL DEFAULT '0',
  `is_free` tinyint DEFAULT '1' COMMENT '1-free,0-paid',
  `price` decimal(10,2) DEFAULT '0.00',
  `user_id` bigint unsigned NOT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'draft' COMMENT 'draft,published,pending',
  `category` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_photosets_deleted_at` (`deleted_at`),
  KEY `fk_photosets_user` (`user_id`),
  KEY `idx_photosets_category` (`category`),
  FULLTEXT KEY `ft_title_description` (`title`,`description`) /*!50100 WITH PARSER `ngram` */ ,
  CONSTRAINT `fk_photosets_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `photosets`
--

LOCK TABLES `photosets` WRITE;
/*!40000 ALTER TABLE `photosets` DISABLE KEYS */;
INSERT INTO `photosets` VALUES (1,'2026-04-08 19:29:56.206','2026-04-08 19:29:56.206',NULL,'Nature Collection','https://picsum.photos/800/600','Beautiful nature photos',0,1,0.00,2,'published',''),(2,'2026-04-08 19:29:57.161','2026-04-08 19:29:57.161',NULL,'Premium Portraits','https://picsum.photos/800/600','Professional portrait photos',0,1,9.99,2,'published',''),(3,'2026-04-08 19:35:11.857','2026-04-16 16:25:35.926',NULL,'测试数据12','/uploads/photos/20260416162505/04/70ac0e5f-db02-4838-8acc-cb344f935d1f.jpg','黑丝',0,0,1.01,2,'published','portrait'),(4,'2026-04-08 19:35:12.367','2026-04-08 19:35:12.367',NULL,'Premium Portraits','https://picsum.photos/800/600','Professional portrait photos',0,1,9.99,2,'published',''),(5,'2026-04-08 19:47:32.507','2026-04-09 23:02:40.785',NULL,'Test Photoset 1775648852','/uploads/images/2026/04/09/b359dc88-dbe4-481a-8c7c-fc0a1ae8a1bd.jpg','Integration test',0,1,0.00,9,'published','portrait'),(6,'2026-04-09 23:38:17.086','2026-04-16 16:03:38.200',NULL,'林测试','/uploads/images/2026/04/10/ca575230-ead4-4db1-9ecf-7fe5bd0be990.jpg','',0,1,0.00,3,'published','portrait'),(7,'2026-04-09 23:40:10.190','2026-04-09 23:40:10.190',NULL,'踩踩踩','/uploads/images/2026/04/09/2f6b02b7-02f1-40ce-b81a-8f9b12432ba0.jpg','测试',0,1,0.00,3,'published','portrait'),(8,'2026-04-16 17:44:40.384','2026-04-16 17:44:40.384',NULL,'允爾','/uploads/photos/20260416174327/04/772c3bc4-0163-46fb-8f07-ac7ecdd5477e.jpg','允爾',0,1,0.00,3,'published','portrait'),(9,'2026-04-16 17:47:00.522','2026-04-16 17:48:58.711',NULL,'扬','/uploads/photos/20260416174600/04/ea818218-a727-4798-b551-a3f9632b9357.jpg','',0,0,0.01,3,'published','portrait');
/*!40000 ALTER TABLE `photosets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `points_mall_items`
--

DROP TABLE IF EXISTS `points_mall_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `points_mall_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `category` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `points_cost` bigint NOT NULL DEFAULT '0',
  `item_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `item_value` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `total_stock` bigint DEFAULT '-1',
  `used_stock` bigint DEFAULT '0',
  `is_unlimited` tinyint(1) DEFAULT '1',
  `min_level` bigint DEFAULT '1',
  `is_active` tinyint(1) DEFAULT '1',
  `sort_order` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_points_mall_items_deleted_at` (`deleted_at`),
  KEY `idx_points_mall_items_category` (`category`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `points_mall_items`
--

LOCK TABLES `points_mall_items` WRITE;
/*!40000 ALTER TABLE `points_mall_items` DISABLE KEYS */;
INSERT INTO `points_mall_items` VALUES (1,'2026-05-31 22:33:01.411','2026-05-31 22:33:01.411',NULL,'专属徽章-新手上路','新手专属徽章，展示你的社区起点','','badge',100,'badge','badge_newbie',-1,0,1,1,1,1),(2,'2026-05-31 22:33:01.562','2026-05-31 22:33:01.562',NULL,'专属徽章-活跃达人','活跃达人徽章，展示你的社区活力','','badge',500,'badge','badge_active',-1,0,1,2,1,2),(3,'2026-05-31 22:33:01.728','2026-05-31 22:33:01.728',NULL,'专属徽章-创作大师','创作大师徽章，展示你的创作才华','','badge',2000,'badge','badge_creator',-1,0,1,4,1,3),(4,'2026-05-31 22:33:01.923','2026-05-31 22:33:01.923',NULL,'称号-社区之星','闪耀的社区之星称号','','title',300,'title','社区之星',-1,0,1,2,1,10),(5,'2026-05-31 22:33:02.069','2026-05-31 22:33:02.069',NULL,'称号-创作先锋','勇于创作的先锋称号','','title',800,'title','创作先锋',-1,0,1,3,1,11),(6,'2026-05-31 22:33:02.228','2026-05-31 22:33:02.228',NULL,'称号-意见领袖','社区意见领袖称号','','title',1500,'title','意见领袖',-1,0,1,4,1,12),(7,'2026-05-31 22:33:02.472','2026-05-31 22:33:02.472',NULL,'VIP体验卡-7天','7天VIP体验，享受会员特权','','privilege',1000,'vip_days','7',-1,0,1,3,1,20),(8,'2026-05-31 22:33:02.585','2026-05-31 22:33:02.585',NULL,'VIP体验卡-30天','30天VIP体验，享受会员特权','','privilege',3000,'vip_days','30',-1,0,1,5,1,21),(9,'2026-05-31 22:33:02.706','2026-05-31 22:33:02.706',NULL,'自定义头像框','解锁专属头像框装饰','','privilege',2000,'custom','custom_avatar_frame',-1,0,1,4,1,22);
/*!40000 ALTER TABLE `points_mall_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_categories`
--

DROP TABLE IF EXISTS `post_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `color` varchar(7) COLLATE utf8mb4_unicode_ci DEFAULT '#409EFF',
  `icon` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `sort_order` int DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_post_categories_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_categories`
--

LOCK TABLES `post_categories` WRITE;
/*!40000 ALTER TABLE `post_categories` DISABLE KEYS */;
INSERT INTO `post_categories` VALUES (1,'discussion','讨论','','#409EFF','',1,'2026-05-24 20:36:36.216','2026-05-24 20:36:36.216'),(2,'qa','问答','','#67C23A','',2,'2026-05-24 20:36:36.216','2026-05-24 20:36:36.216'),(3,'showcase','作品展示','','#E6A23C','',3,'2026-05-24 20:36:36.216','2026-05-24 20:36:36.216'),(4,'suggestion','建议','','#F56C6C','',4,'2026-05-24 20:36:36.216','2026-05-24 20:36:36.216');
/*!40000 ALTER TABLE `post_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_likes`
--

DROP TABLE IF EXISTS `post_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `post_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_post` (`user_id`,`post_id`),
  KEY `idx_post_id` (`post_id`),
  CONSTRAINT `fk_post_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_posts_likes` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_likes`
--

LOCK TABLES `post_likes` WRITE;
/*!40000 ALTER TABLE `post_likes` DISABLE KEYS */;
/*!40000 ALTER TABLE `post_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_replies`
--

DROP TABLE IF EXISTS `post_replies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_replies` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `post_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_reply_id` bigint unsigned DEFAULT NULL,
  `like_count` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_post_replies_post_id` (`post_id`),
  KEY `idx_post_replies_user_id` (`user_id`),
  KEY `idx_post_replies_parent_reply_id` (`parent_reply_id`),
  CONSTRAINT `fk_post_replies_children` FOREIGN KEY (`parent_reply_id`) REFERENCES `post_replies` (`id`),
  CONSTRAINT `fk_post_replies_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_posts_replies` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_replies`
--

LOCK TABLES `post_replies` WRITE;
/*!40000 ALTER TABLE `post_replies` DISABLE KEYS */;
/*!40000 ALTER TABLE `post_replies` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_reply_likes`
--

DROP TABLE IF EXISTS `post_reply_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_reply_likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `reply_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_reply` (`user_id`,`reply_id`),
  KEY `idx_reply_id` (`reply_id`),
  CONSTRAINT `fk_post_replies_likes` FOREIGN KEY (`reply_id`) REFERENCES `post_replies` (`id`),
  CONSTRAINT `fk_post_reply_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_reply_likes`
--

LOCK TABLES `post_reply_likes` WRITE;
/*!40000 ALTER TABLE `post_reply_likes` DISABLE KEYS */;
/*!40000 ALTER TABLE `post_reply_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_reports`
--

DROP TABLE IF EXISTS `post_reports`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_reports` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `post_id` bigint unsigned DEFAULT NULL,
  `reply_id` bigint unsigned DEFAULT NULL,
  `reporter_id` bigint unsigned NOT NULL,
  `reason` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `handler_id` bigint unsigned DEFAULT NULL,
  `handled_at` datetime(3) DEFAULT NULL,
  `handle_note` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_post_reports_post_id` (`post_id`),
  KEY `idx_post_reports_reply_id` (`reply_id`),
  KEY `idx_post_reports_reporter_id` (`reporter_id`),
  KEY `idx_post_reports_status` (`status`),
  KEY `fk_post_reports_handler` (`handler_id`),
  CONSTRAINT `fk_post_reports_handler` FOREIGN KEY (`handler_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_post_reports_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
  CONSTRAINT `fk_post_reports_reply` FOREIGN KEY (`reply_id`) REFERENCES `post_replies` (`id`),
  CONSTRAINT `fk_post_reports_reporter` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_reports`
--

LOCK TABLES `post_reports` WRITE;
/*!40000 ALTER TABLE `post_reports` DISABLE KEYS */;
/*!40000 ALTER TABLE `post_reports` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_shares`
--

DROP TABLE IF EXISTS `post_shares`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_shares` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `post_id` bigint unsigned NOT NULL,
  `platform` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'other',
  PRIMARY KEY (`id`),
  KEY `idx_share_user` (`user_id`),
  KEY `idx_share_post` (`post_id`),
  CONSTRAINT `fk_post_shares_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_posts_shares` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_shares`
--

LOCK TABLES `post_shares` WRITE;
/*!40000 ALTER TABLE `post_shares` DISABLE KEYS */;
/*!40000 ALTER TABLE `post_shares` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_tags`
--

DROP TABLE IF EXISTS `post_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_tags` (
  `post_id` bigint unsigned NOT NULL,
  `tag_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`post_id`,`tag_id`),
  KEY `fk_post_tags_tag` (`tag_id`),
  CONSTRAINT `fk_post_tags_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
  CONSTRAINT `fk_post_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_tags`
--

LOCK TABLES `post_tags` WRITE;
/*!40000 ALTER TABLE `post_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `post_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post_topics`
--

DROP TABLE IF EXISTS `post_topics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_topics` (
  `topic_id` bigint unsigned NOT NULL,
  `post_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`topic_id`,`post_id`),
  KEY `fk_post_topics_post` (`post_id`),
  CONSTRAINT `fk_post_topics_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
  CONSTRAINT `fk_post_topics_topic` FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_topics`
--

LOCK TABLES `post_topics` WRITE;
/*!40000 ALTER TABLE `post_topics` DISABLE KEYS */;
/*!40000 ALTER TABLE `post_topics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `photoset_id` bigint unsigned DEFAULT NULL,
  `category` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'discussion',
  `visibility` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'public',
  `is_pinned` tinyint(1) NOT NULL DEFAULT '0',
  `view_count` bigint NOT NULL DEFAULT '0',
  `reply_count` bigint NOT NULL DEFAULT '0',
  `like_count` bigint NOT NULL DEFAULT '0',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'approved',
  `post_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'dynamic',
  `is_essence` tinyint(1) NOT NULL DEFAULT '0',
  `share_count` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_posts_deleted_at` (`deleted_at`),
  KEY `idx_posts_user_id` (`user_id`),
  KEY `fk_posts_photoset` (`photoset_id`),
  CONSTRAINT `fk_posts_photoset` FOREIGN KEY (`photoset_id`) REFERENCES `photosets` (`id`),
  CONSTRAINT `fk_posts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sensitive_words`
--

DROP TABLE IF EXISTS `sensitive_words`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sensitive_words` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `replacement` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '***',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sensitive_words_word` (`word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sensitive_words`
--

LOCK TABLES `sensitive_words` WRITE;
/*!40000 ALTER TABLE `sensitive_words` DISABLE KEYS */;
/*!40000 ALTER TABLE `sensitive_words` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `site_settings`
--

DROP TABLE IF EXISTS `site_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `site_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci,
  `group` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'general',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_site_settings_key` (`key`),
  KEY `idx_site_settings_group` (`group`)
) ENGINE=InnoDB AUTO_INCREMENT=134 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `site_settings`
--

LOCK TABLES `site_settings` WRITE;
/*!40000 ALTER TABLE `site_settings` DISABLE KEYS */;
INSERT INTO `site_settings` VALUES (1,'2026-04-11 23:29:53.012','2026-04-11 23:29:53.012','site_name','测试','general'),(2,'2026-04-11 23:29:53.270','2026-04-15 17:33:45.000','site_description','调试站','general'),(3,'2026-04-11 23:29:53.425','2026-04-15 17:33:45.000','site_icp','京ICP备04000001号','general'),(4,'2026-04-11 23:29:53.582','2026-04-15 17:33:45.000','register_enabled','true','general'),(5,'2026-04-11 23:29:53.729','2026-04-15 17:33:45.000','email_verify_required','false','general'),(23,'2026-04-12 00:13:53.729','2026-04-15 17:33:45.000','site_title','调试站','general'),(80,'2026-04-15 15:44:19.000','2026-04-15 15:44:19.000','site_keywords','丝袜室，黑丝','general'),(81,'2026-04-15 15:44:38.000','2026-04-15 15:44:38.000','watermark_enabled','true','general'),(82,'2026-04-15 15:44:38.000','2026-04-15 15:44:38.000','watermark_text','siwashi','general'),(83,'2026-04-15 15:44:38.000','2026-04-15 15:44:38.000','watermark_opacity','30','general'),(89,'2026-04-15 22:37:28.000','2026-04-16 17:59:14.000','smtp_from_name','本地调试测试发送邮件','general'),(90,'2026-04-15 22:37:28.000','2026-04-16 17:59:14.000','smtp_host','smtp.126.com','general'),(91,'2026-04-15 22:37:28.000','2026-04-16 17:59:14.000','smtp_port','465','general'),(92,'2026-04-15 22:37:29.000','2026-04-16 17:59:15.000','smtp_user','siwashi@126.com','general'),(93,'2026-04-15 22:37:29.000','2026-04-16 17:59:15.000','smtp_password','UBvr6T6f9tUGVwPW','general'),(129,'2026-04-16 21:09:41.000','2026-04-16 21:09:41.000','site_url','','general'),(130,'2026-04-16 21:09:41.000','2026-04-16 21:09:41.000','api_url','','general'),(131,'2026-04-16 21:09:41.000','2026-04-16 21:09:41.000','dev_api_url','http://localhost:3000','general'),(132,'2026-05-16 23:04:38.000','2026-05-16 23:13:50.000','nav_menu','[{\"slug\":\"street\"},{\"slug\":\"portrait\"},{\"slug\":\"nature\"}]','general');
/*!40000 ALTER TABLE `site_settings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tags`
--

DROP TABLE IF EXISTS `tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tags_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tags`
--

LOCK TABLES `tags` WRITE;
/*!40000 ALTER TABLE `tags` DISABLE KEYS */;
INSERT INTO `tags` VALUES (1,'2026-04-08 19:29:56.410','nature'),(2,'2026-04-08 19:29:56.648','landscape'),(3,'2026-04-08 19:29:57.401','portrait'),(4,'2026-04-08 19:29:57.559','studio'),(5,'2026-04-09 23:06:01.071','黑丝');
/*!40000 ALTER TABLE `tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `topics`
--

DROP TABLE IF EXISTS `topics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `topics` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `cover` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `post_count` bigint NOT NULL DEFAULT '0',
  `is_hot` tinyint(1) NOT NULL DEFAULT '0',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_topics_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `topics`
--

LOCK TABLES `topics` WRITE;
/*!40000 ALTER TABLE `topics` DISABLE KEYS */;
/*!40000 ALTER TABLE `topics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_achievements`
--

DROP TABLE IF EXISTS `user_achievements`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_achievements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `achievement_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_achievements_user_id` (`user_id`),
  KEY `idx_user_achievements_achievement_id` (`achievement_id`),
  CONSTRAINT `fk_user_achievements_achievement` FOREIGN KEY (`achievement_id`) REFERENCES `achievements` (`id`),
  CONSTRAINT `fk_user_achievements_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_achievements`
--

LOCK TABLES `user_achievements` WRITE;
/*!40000 ALTER TABLE `user_achievements` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_achievements` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_blocks`
--

DROP TABLE IF EXISTS `user_blocks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_blocks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `blocked_id` bigint unsigned NOT NULL,
  `block_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'block',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_block` (`user_id`,`blocked_id`),
  KEY `idx_blocked_id` (`blocked_id`),
  CONSTRAINT `fk_user_blocks_blocked` FOREIGN KEY (`blocked_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_user_blocks_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_blocks`
--

LOCK TABLES `user_blocks` WRITE;
/*!40000 ALTER TABLE `user_blocks` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_blocks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_level_configs`
--

DROP TABLE IF EXISTS `user_level_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_level_configs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `level` bigint NOT NULL,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `icon` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `color` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '#FFD700',
  `min_points` bigint NOT NULL DEFAULT '0',
  `max_points` bigint NOT NULL DEFAULT '0',
  `description` text COLLATE utf8mb4_unicode_ci,
  `can_create_post` tinyint(1) DEFAULT '1',
  `can_create_reply` tinyint(1) DEFAULT '1',
  `can_upload_image` tinyint(1) DEFAULT '1',
  `can_upload_video` tinyint(1) DEFAULT '0',
  `can_create_topic` tinyint(1) DEFAULT '0',
  `can_pin_post` tinyint(1) DEFAULT '0',
  `can_delete_reply` tinyint(1) DEFAULT '0',
  `max_post_per_day` bigint DEFAULT '5',
  `max_reply_per_day` bigint DEFAULT '10',
  `max_image_per_post` bigint DEFAULT '9',
  `max_video_per_post` bigint DEFAULT '0',
  `max_post_length` bigint DEFAULT '5000',
  `reward_points` bigint DEFAULT '0',
  `reward_badge` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `reward_title` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_level_configs_level` (`level`),
  KEY `idx_user_level_configs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_level_configs`
--

LOCK TABLES `user_level_configs` WRITE;
/*!40000 ALTER TABLE `user_level_configs` DISABLE KEYS */;
INSERT INTO `user_level_configs` VALUES (1,'2026-05-31 22:32:56.441','2026-05-31 22:32:56.441',NULL,1,'新手','🌱','#8BC34A',0,99,'刚加入社区的新手',1,1,1,0,0,0,0,5,10,9,0,5000,0,'newbie','社区新人'),(2,'2026-05-31 22:32:56.578','2026-05-31 22:32:56.578',NULL,2,'活跃成员','⭐','#2196F3',100,499,'积极参与社区讨论的成员',1,1,1,0,0,0,0,10,20,9,0,8000,50,'active','活跃达人'),(3,'2026-05-31 22:32:56.756','2026-05-31 22:32:56.756',NULL,3,'资深成员','🌟','#9C27B0',500,1999,'社区的资深贡献者',1,1,1,0,1,0,0,15,30,12,0,10000,100,'senior','资深社区成员'),(4,'2026-05-31 22:32:57.200','2026-05-31 22:32:57.200',NULL,4,'金牌成员','🏅','#FF9800',2000,4999,'社区的金牌贡献者',1,1,1,1,1,0,0,20,50,15,1,15000,200,'gold','金牌创作者'),(5,'2026-05-31 22:32:57.314','2026-05-31 22:32:57.314',NULL,5,'钻石成员','💎','#00BCD4',5000,9999,'社区的钻石级成员',1,1,1,1,1,1,1,30,100,20,3,20000,500,'diamond','钻石创作者'),(6,'2026-05-31 22:32:57.492','2026-05-31 22:32:57.492',NULL,6,'至尊成员','👑','#F44336',10000,19999,'社区的至尊级成员',1,1,1,1,1,1,1,50,200,30,5,30000,1000,'supreme','至尊创作者'),(7,'2026-05-31 22:32:57.662','2026-05-31 22:32:57.662',NULL,7,'荣耀L7','🏆','#E91E63',20000,29999,'荣耀等级第七级',1,1,1,1,1,1,1,100,500,50,10,50000,2000,'glory7','荣耀守护者'),(8,'2026-05-31 22:32:57.835','2026-05-31 22:32:57.835',NULL,8,'荣耀L8','🏆','#673AB7',30000,39999,'荣耀等级第八级',1,1,1,1,1,1,1,100,500,50,10,50000,3000,'glory8','荣耀精英'),(9,'2026-05-31 22:32:58.147','2026-05-31 22:32:58.147',NULL,9,'荣耀L9','🏆','#3F51B5',40000,49999,'荣耀等级第九级',1,1,1,1,1,1,1,100,500,50,10,50000,5000,'glory9','荣耀大师'),(10,'2026-05-31 22:32:58.319','2026-05-31 22:32:58.319',NULL,10,'荣耀L10','🏆','#FFD700',50000,999999,'荣耀等级最高级',1,1,1,1,1,1,1,999,999,99,20,99999,10000,'glory10','荣耀传说');
/*!40000 ALTER TABLE `user_level_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_points`
--

DROP TABLE IF EXISTS `user_points`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_points` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `points` bigint NOT NULL DEFAULT '0',
  `level` bigint NOT NULL DEFAULT '1',
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_user_points_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_points`
--

LOCK TABLES `user_points` WRITE;
/*!40000 ALTER TABLE `user_points` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_points` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_points_exchanges`
--

DROP TABLE IF EXISTS `user_points_exchanges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_points_exchanges` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `item_id` bigint unsigned NOT NULL,
  `points` bigint NOT NULL,
  `item_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `item_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `item_value` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'completed',
  PRIMARY KEY (`id`),
  KEY `idx_user_points_exchanges_user_id` (`user_id`),
  KEY `fk_user_points_exchanges_item` (`item_id`),
  CONSTRAINT `fk_user_points_exchanges_item` FOREIGN KEY (`item_id`) REFERENCES `points_mall_items` (`id`),
  CONSTRAINT `fk_user_points_exchanges_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_points_exchanges`
--

LOCK TABLES `user_points_exchanges` WRITE;
/*!40000 ALTER TABLE `user_points_exchanges` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_points_exchanges` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'user',
  `status` tinyint DEFAULT '1' COMMENT '1-active,0-inactive',
  `last_login_at` datetime(3) DEFAULT NULL,
  `membership_expires` datetime(3) DEFAULT NULL,
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `bio` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `ip_location` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `level` tinyint DEFAULT '1',
  `exp` bigint DEFAULT '0',
  `circle_count` bigint DEFAULT '0',
  `following_count` bigint DEFAULT '0',
  `follower_count` bigint DEFAULT '0',
  `like_count` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_users_membership_expires` (`membership_expires`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2026-04-08 19:27:26.241','2026-04-08 19:35:13.530',NULL,'TestUser','test@photoset.dev','$2a$10$gUt10KB3N/dxdRtZKGoPpubW2Bh44T0Wa0offDcnSGyjrkddkK0pO','user',1,'2026-04-08 19:35:13.530',NULL,'','','',1,0,0,0,0,0),(2,'2026-04-08 19:27:26.861','2026-04-08 19:35:11.212',NULL,'Creator1','creator@photoset.dev','$2a$10$vXPJQQrKjT3HHAtbvvmvJ.kYqeHyBwXSTXjEP9uOXx3JGRgQCS3/O','creator',1,'2026-04-08 19:35:11.202',NULL,'','','',1,0,0,0,0,0),(3,'2026-04-08 19:41:37.726','2026-05-31 22:35:13.141',NULL,'dulang','root@honker.org','$2a$10$0Gp2lxqS8lOjmO9Fc.CaIufA4MTAOmiNdY4g9ZviBBcAUBwA/hUBe','admin',1,'2026-05-31 22:35:13.140',NULL,'','','',1,0,0,0,0,0),(4,'2026-04-08 19:45:39.277','2026-04-08 19:45:39.277',NULL,'DirectTest','direct@test.dev','$2a$10$LfzOdJ2.XFEIOyMiKspW7Ogtb5ZCMp9FoRXk.dnMc0hdbFyHQWzZW','user',1,NULL,NULL,'','','',1,0,0,0,0,0),(5,'2026-04-08 19:45:39.580','2026-04-08 19:45:39.580',NULL,'ProxyTest','proxy@test.dev','$2a$10$pspf8WHW0eMM0JP8ZHc1bOIIVhqjPxIqkprP6n8pukIu/wR4Fw1ta','user',1,NULL,NULL,'','','',1,0,0,0,0,0),(6,'2026-04-08 19:45:39.966','2026-04-08 19:45:39.966',NULL,'ProxyTest2','proxy2@test.dev','$2a$10$46osFQnvs2DDBl.wxP/Ewu6wSMUIgFpS93ngItfqnl0T.p9bHf3Ma','user',1,NULL,NULL,'','','',1,0,0,0,0,0),(7,'2026-04-08 19:46:11.485','2026-04-08 19:46:11.826',NULL,'IntegTest','integtest@photoset.dev','$2a$10$mO52ZPIhmDNzPIaXXIoc1.7g0oLDyIDndhTu9CHn2Ne3chmAw4EGa','user',1,'2026-04-08 19:46:11.826',NULL,'','','',1,0,0,0,0,0),(8,'2026-04-08 19:47:31.035','2026-04-08 19:47:31.304',NULL,'TestUser1775648850','test1775648850@photoset.dev','$2a$10$BMgrW8XRIM/xT.OGPoMY1OX5WaamNox6ZuQTTBW2KThEJDEW/q2Q.','user',1,'2026-04-08 19:47:31.303',NULL,'','','',1,0,0,0,0,0),(9,'2026-04-08 19:47:31.702','2026-04-08 19:47:32.270',NULL,'Creator1775648851','creator1775648851@photoset.dev','$2a$10$EBe8VSG0.4Alv0FG3u7TF.PszkImZIZE7Skj3NLpjwUWvbT2ETVRa','creator',1,'2026-04-08 19:47:32.269',NULL,'','','',1,0,0,0,0,0),(10,'2026-04-08 21:11:50.549','2026-04-15 20:46:15.009',NULL,'testuser','testuser@example.com','$2a$10$wN/t/8vkJOkmKxx42H99e.XJE1.6StqZvfGOnNs.LU2DuthUGhEH6','member',1,'2026-04-08 21:20:47.497',NULL,'','','',1,0,0,0,0,0),(11,'2026-04-09 00:31:29.000','2026-04-09 00:35:35.418',NULL,'admin','admin@photoset.dev','$2a$10$6HTkZ8Cz..1LFbO8nxBmquouXGrgeRwjalJFWOI8uzw5RHlaiJx2C','admin',1,'2026-04-09 00:35:35.417',NULL,'','','',1,0,0,0,0,0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'photoset'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-05-31 22:35:43
