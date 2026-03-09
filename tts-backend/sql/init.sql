-- 创建数据库
CREATE DATABASE IF NOT EXISTS tts_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE tts_db;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS tts_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE tts_db;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(100),
  `balance` DECIMAL(10, 2) DEFAULT 0.00,
  `character_count` BIGINT DEFAULT 0,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 音色表
CREATE TABLE IF NOT EXISTS `voice` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL,
  `tone` VARCHAR(50),
  `gender` VARCHAR(10),
  `preview_url` VARCHAR(255),
  `is_default` TINYINT(1) DEFAULT 0,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- TTS 任务表
CREATE TABLE IF NOT EXISTS `tts_task` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `task_id` VARCHAR(64) NOT NULL UNIQUE,
  `user_id` BIGINT,
  `status` VARCHAR(20) DEFAULT 'pending',
  `progress` INT DEFAULT 0,
  `audio_url` VARCHAR(255),
  `format` VARCHAR(10) DEFAULT 'mp3',
  `channel` VARCHAR(10) DEFAULT 'mono',
  `error_msg` TEXT,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_task_id` (`task_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- TTS 片段表
CREATE TABLE IF NOT EXISTS `tts_segment` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `task_id` VARCHAR(64) NOT NULL,
  `voice_id` BIGINT NOT NULL,
  `emotion` VARCHAR(50),
  `text` TEXT NOT NULL,
  `sort` INT DEFAULT 0,
  INDEX `idx_task_id` (`task_id`),
  FOREIGN KEY (`task_id`) REFERENCES `tts_task`(`task_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 初始化默认音色数据
INSERT INTO `voice` (`name`, `tone`, `gender`, `is_default`) VALUES
('晓晓', '女-青年', 'female', 1),
('云飞', '男-青年', 'male', 0),
('小美', '女-童声', 'female', 0),
('阿强', '男-中年', 'male', 0);
