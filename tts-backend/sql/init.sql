-- Init schema for tts_db
CREATE DATABASE IF NOT EXISTS tts_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE tts_db;

-- Users
CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(100),
  `role` VARCHAR(20) DEFAULT 'user',
  `status` VARCHAR(20) DEFAULT 'active',
  `balance` DECIMAL(10, 2) DEFAULT 0.00,
  `character_count` BIGINT DEFAULT 0,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Admin operation logs
CREATE TABLE IF NOT EXISTS `admin_log` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `actor_user_id` BIGINT NOT NULL,
  `action` TEXT NOT NULL,
  `ip` VARCHAR(64),
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_admin_log_created_at` (`created_at`),
  INDEX `idx_admin_log_actor` (`actor_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Voices
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

-- User default voice (per user)
CREATE TABLE IF NOT EXISTS `user_voice_default` (
  `user_id` BIGINT PRIMARY KEY,
  `voice_id` BIGINT NOT NULL,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_voice_default_voice_id` (`voice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- TTS tasks
CREATE TABLE IF NOT EXISTS `tts_task` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `task_id` VARCHAR(64) NOT NULL UNIQUE,
  `title` VARCHAR(255) DEFAULT '',
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

-- TTS segments
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

-- Custom voice requests (simplified workflow)
CREATE TABLE IF NOT EXISTS `custom_voice_request` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `name` VARCHAR(50) NOT NULL,
  `tone` VARCHAR(50),
  `gender` VARCHAR(10),
  `sample_text` TEXT,
  `sample_urls` TEXT,
  `status` VARCHAR(20) DEFAULT 'pending',
  `result_voice_id` BIGINT,
  `error_msg` TEXT,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_custom_user_id` (`user_id`),
  INDEX `idx_custom_status` (`status`),
  INDEX `idx_custom_result_voice_id` (`result_voice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Feedback / suggestions
CREATE TABLE IF NOT EXISTS `feedback` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `category` VARCHAR(50),
  `content` TEXT NOT NULL,
  `contact` VARCHAR(100),
  `status` VARCHAR(20) DEFAULT 'open',
  `reply` TEXT,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_feedback_user_id` (`user_id`),
  INDEX `idx_feedback_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Seed demo users (for local development / admin pages)
-- Note: login supports plain-text passwords for backward compatibility in this demo project.
INSERT IGNORE INTO `user` (`username`, `password`, `email`, `role`, `status`, `balance`, `character_count`) VALUES
('admin', 'admin123', 'admin@example.com', 'admin', 'active', 1000, 0),
('demo', 'demo123', 'demo@example.com', 'user', 'active', 200, 0),
('alice', 'alice123', 'alice@example.com', 'engineer', 'active', 300, 0);

-- Seed default voices (idempotent)
INSERT INTO `voice` (`name`, `tone`, `gender`, `preview_url`, `is_default`)
SELECT 'XiaoXiao', 'young', 'female', 'https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav', 1
WHERE NOT EXISTS (SELECT 1 FROM `voice` v WHERE v.name = 'XiaoXiao')
LIMIT 1;

INSERT INTO `voice` (`name`, `tone`, `gender`, `preview_url`, `is_default`)
SELECT 'YunFei', 'young', 'male', 'https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav', 0
WHERE NOT EXISTS (SELECT 1 FROM `voice` v WHERE v.name = 'YunFei')
LIMIT 1;

INSERT INTO `voice` (`name`, `tone`, `gender`, `preview_url`, `is_default`)
SELECT 'XiaoMei', 'child', 'female', 'https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav', 0
WHERE NOT EXISTS (SELECT 1 FROM `voice` v WHERE v.name = 'XiaoMei')
LIMIT 1;

INSERT INTO `voice` (`name`, `tone`, `gender`, `preview_url`, `is_default`)
SELECT 'A-Qiang', 'middle', 'male', 'https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav', 0
WHERE NOT EXISTS (SELECT 1 FROM `voice` v WHERE v.name = 'A-Qiang')
LIMIT 1;

-- Seed an approved custom voice (for preview)
INSERT INTO `voice` (`name`, `tone`, `gender`, `preview_url`, `is_default`)
SELECT 'Radio Host - Dawn', 'warm', 'female', 'https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav', 0
WHERE NOT EXISTS (SELECT 1 FROM `voice` v WHERE v.name = 'Radio Host - Dawn')
LIMIT 1;

-- Seed user default voice mapping (idempotent)
INSERT INTO `user_voice_default` (`user_id`, `voice_id`)
SELECT u.id, v.id
FROM `user` u
JOIN `voice` v ON v.name = 'XiaoXiao'
WHERE u.username IN ('admin', 'demo', 'alice')
ON DUPLICATE KEY UPDATE `voice_id` = VALUES(`voice_id`);

-- Seed feedback mock data
INSERT INTO `feedback` (`user_id`, `category`, `content`, `contact`, `status`, `reply`)
SELECT u.id, 'Bug', 'Sometimes progress stays at 0% after clicking Generate All; refreshing fixes it.', 'demo@example.com', 'open', ''
FROM `user` u
WHERE u.username = 'demo'
  AND NOT EXISTS (
    SELECT 1 FROM `feedback` f WHERE f.user_id = u.id AND f.content = 'Sometimes progress stays at 0% after clicking Generate All; refreshing fixes it.'
  )
LIMIT 1;

INSERT INTO `feedback` (`user_id`, `category`, `content`, `contact`, `status`, `reply`)
SELECT u.id, 'Feature', 'Please add batch download and tags filter on Works page.', 'alice@example.com', 'closed', 'Noted. We will add batch actions and filters in a future update.'
FROM `user` u
WHERE u.username = 'alice'
  AND NOT EXISTS (
    SELECT 1 FROM `feedback` f WHERE f.user_id = u.id AND f.content = 'Please add batch download and tags filter on Works page.'
  )
LIMIT 1;

-- Seed custom voice requests mock data (sample_urls stores a JSON array string)
INSERT INTO `custom_voice_request`
  (`user_id`, `name`, `tone`, `gender`, `sample_text`, `sample_urls`, `status`, `result_voice_id`, `error_msg`)
SELECT
  u.id,
  'Radio Host - Dawn',
  'warm',
  'female',
  'Read 1-2 minutes of clean speech in a quiet environment.',
  '["https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav"]',
  'success',
  (SELECT v.id FROM `voice` v WHERE v.name = 'Radio Host - Dawn' ORDER BY v.id DESC LIMIT 1),
  NULL
FROM `user` u
WHERE u.username = 'demo'
  AND NOT EXISTS (
    SELECT 1 FROM `custom_voice_request` r WHERE r.user_id = u.id AND r.name = 'Radio Host - Dawn'
  )
LIMIT 1;

INSERT INTO `custom_voice_request`
  (`user_id`, `name`, `tone`, `gender`, `sample_text`, `sample_urls`, `status`, `result_voice_id`, `error_msg`)
SELECT
  u.id,
  'Game Character - Demo',
  'lively',
  'male',
  'Provide 5-10 minutes of recordings. Dialogue / narration are both ok.',
  '["https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav","https://www2.cs.uic.edu/~i101/SoundFiles/ImperialMarch60.wav"]',
  'pending',
  NULL,
  NULL
FROM `user` u
WHERE u.username = 'alice'
  AND NOT EXISTS (
    SELECT 1 FROM `custom_voice_request` r WHERE r.user_id = u.id AND r.name = 'Game Character - Demo'
  )
LIMIT 1;

INSERT INTO `custom_voice_request`
  (`user_id`, `name`, `tone`, `gender`, `sample_text`, `sample_urls`, `status`, `result_voice_id`, `error_msg`)
SELECT
  u.id,
  'Support Agent - Sweet',
  'sweet',
  'female',
  'Please provide clearer recordings with less noise.',
  '["https://www2.cs.uic.edu/~i101/SoundFiles/BabyElephantWalk60.wav"]',
  'failed',
  NULL,
  'Low quality samples: too much background noise. Please re-record and submit again.'
FROM `user` u
WHERE u.username = 'demo'
  AND NOT EXISTS (
    SELECT 1 FROM `custom_voice_request` r WHERE r.user_id = u.id AND r.name = 'Support Agent - Sweet'
  )
LIMIT 1;
