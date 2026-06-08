-- Mac Cluster Monitor 数据库表结构
-- 创建数据库: CREATE DATABASE mac_cluster_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE mac_cluster_monitor;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码哈希',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    role VARCHAR(20) DEFAULT 'user' COMMENT '角色: admin=超级管理员, user=普通用户',
    license_key VARCHAR(36) NOT NULL UNIQUE COMMENT '卡密(UUID), 即账号唯一标识',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_license_key (license_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 设备表
CREATE TABLE IF NOT EXISTS devices (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    mac_address VARCHAR(17) NOT NULL UNIQUE COMMENT 'MAC地址 (aa:bb:cc:dd:ee:ff)',
    hostname VARCHAR(100) NOT NULL COMMENT '设备主机名',
    ip_address VARCHAR(45) COMMENT 'IP地址 (支持IPv6)',
    user_id BIGINT UNSIGNED COMMENT '绑定用户ID, NULL表示未绑定',
    total_memory BIGINT UNSIGNED COMMENT '总内存 (MB)',
    cpu_cores INT UNSIGNED COMMENT 'CPU核心数',
    status TINYINT DEFAULT 0 COMMENT '状态: 0=离线, 1=在线',
    last_seen_at TIMESTAMP NULL COMMENT '最后在线时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '首次注册时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_status (status),
    INDEX idx_mac_address (mac_address),
    INDEX idx_user_id (user_id),
    INDEX idx_last_seen (last_seen_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备表';

-- 设备心跳历史表
CREATE TABLE IF NOT EXISTS device_heartbeats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    device_id BIGINT UNSIGNED NOT NULL COMMENT '设备ID',
    ip_address VARCHAR(45) COMMENT '上报时的IP地址',
    memory_used BIGINT UNSIGNED COMMENT '已用内存 (MB)',
    memory_total BIGINT UNSIGNED COMMENT '总内存 (MB)',
    cpu_usage FLOAT COMMENT 'CPU使用率 (%)',
    reported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '上报时间',
    INDEX idx_device_id (device_id),
    INDEX idx_reported_at (reported_at),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备心跳历史表';

-- 用户盈利表 (静态模拟数据, 由后端定时更新)
CREATE TABLE IF NOT EXISTS user_profits (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '用户ID',
    today_profit DECIMAL(12,2) DEFAULT 0.00 COMMENT '今日盈利',
    total_profit DECIMAL(12,2) DEFAULT 0.00 COMMENT '累计盈利',
    last_reset_date DATE COMMENT '上次重置日期(每日重置today_profit)',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户盈利表';

-- 插入默认管理员用户 (密码: admin123)
-- license_key 为管理员专用, 实际不参与设备绑定
INSERT INTO users (username, password, email, role, license_key) VALUES
('admin', '$2a$10$5SCi73i/zvxd3C.ix/J3nOxyi6WMXjw.KQ2VRCLsMcQfV/H1s1PDe', 'admin@example.com', 'admin', UUID())
ON DUPLICATE KEY UPDATE id=id;
