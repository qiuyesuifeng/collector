-- tidb 脚本
CREATE DATABASE IF NOT EXISTS hackathon;
USE hackathon;

-- 用户信息表
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '用户 ID',
  `brand` varchar(64) NOT NULL COMMENT '手机品牌',
  `model` varchar(64) NOT NULL COMMENT '手机型号',
  `name` varchar(32) DEFAULT NULL COMMENT '用户名称',
  `province` varchar(32) DEFAULT NULL COMMENT '用户所在省会',
  `city` varchar(32) DEFAULT NULL COMMENT '用户所在城市',
  `age` int(4) DEFAULT '0' COMMENT '用户年龄',
  `create_time` timestamp NULL DEFAULT NULL COMMENT '用户创建时间',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '用户最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_key` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

Insert Into user values(1, 101, "Apple", "5", "Test User 101", "北京", "北京", 18, now(), now());
Insert Into user values(2, 102, "Apple", "6", "Test User 102", "北京", "北京", 19, now(), now());
Insert Into user values(3, 103, "Apple", "6S", "Test User 103", "北京", "北京", 20, now(), now());
Insert Into user values(4, 104, "Apple", "7", "Test User 104", "北京", "北京", 20, now(), now());
Insert Into user values(5, 105, "Apple", "8", "Test User 105", "北京", "北京", 30, now(), now());
Insert Into user values(6, 106, "Apple", "7", "Test User 106", "北京", "北京", 32, now(), now());
Insert Into user values(7, 107, "Apple", "8", "Test User 107", "北京", "北京", 20, now(), now());
Insert Into user values(8, 108, "Apple", "6S", "Test User 108", "北京", "北京", 35, now(), now());
Insert Into user values(9, 109, "Apple", "5S", "Test User 109", "北京", "北京", 40, now(), now());
Insert Into user values(10, 110, "Apple", "5", "Test User 110", "北京", "北京", 42, now(), now());
Insert Into user values(11, 111, "Apple", "5", "Test User 111", "北京", "北京", 18, now(), now());
Insert Into user values(12, 112, "Apple", "8", "Test User 112", "北京", "北京", 19, now(), now());
Insert Into user values(13, 113, "Apple", "6S", "Test User 113", "北京", "北京", 20, now(), now());
Insert Into user values(14, 114, "Apple", "6", "Test User 114", "北京", "北京", 20, now(), now());
Insert Into user values(15, 115, "Apple", "6S", "Test User 115", "北京", "北京", 30, now(), now());
Insert Into user values(16, 116, "Apple", "7", "Test User 116", "北京", "北京", 32, now(), now());
Insert Into user values(17, 117, "Apple", "8", "Test User 117", "北京", "北京", 20, now(), now());
Insert Into user values(18, 118, "Apple", "7", "Test User 118", "北京", "北京", 35, now(), now());
Insert Into user values(19, 119, "Apple", "5S", "Test User 119", "北京", "北京", 40, now(), now());
Insert Into user values(20, 120, "Apple", "6", "Test User 120", "北京", "北京", 42, now(), now());

-- 广告信息表
CREATE TABLE IF NOT EXISTS `ads` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ads_id` bigint(20) NOT NULL COMMENT '广告 ID',
  `title` varchar(64) NOT NULL COMMENT '广告 Title',
  `content` varchar(256) NOT NULL COMMENT '广告内容',
  `status` varchar(16) NOT NULL COMMENT '用户状态',
  `create_time` timestamp NULL DEFAULT NULL COMMENT '用户创建时间',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '用户最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ads_id_key` (`ads_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

Insert Into ads values(1, 1001, "Test ADS 1", "Hello TiDB Hackathon 1", "normal", now(), now());
Insert Into ads values(2, 1002, "Test ADS 2", "Hello TiDB Hackathon 2", "normal", now(), now());
Insert Into ads values(3, 1003, "Test ADS 3", "Hello TiDB Hackathon 3", "abnormal", now(), now());
Insert Into ads values(4, 1004, "Test ADS 4", "Hello TiDB Hackathon 4", "normal", now(), now());
Insert Into ads values(5, 1005, "Test ADS 5", "Hello TiDB Hackathon 5", "normal", now(), now());
Insert Into ads values(6, 1006, "Test ADS 6", "Hello TiDB Hackathon 6", "abnormal", now(), now());
Insert Into ads values(7, 1007, "Test ADS 7", "Hello TiDB Hackathon 7", "normal", now(), now());
Insert Into ads values(8, 1008, "Test ADS 8", "Hello TiDB Hackathon 8", "normal", now(), now());
Insert Into ads values(9, 1009, "Test ADS 9", "Hello TiDB Hackathon 9", "normal", now(), now());
Insert Into ads values(10, 1010, "Test ADS 10", "Hello TiDB Hackathon 10", "abnormal", now(), now());

-- 广告点击信息流表
CREATE STREAM tidb_kafka_stream_table_click(id bigint(20), user_id bigint(20), ads_id bigint(20), create_time timestamp) with ('type' =  'kafka', 'topic' = 'click');

