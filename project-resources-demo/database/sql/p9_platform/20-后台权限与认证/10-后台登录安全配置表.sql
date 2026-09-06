-- 文件用途：
-- 用于保存总网后台的登录安全配置。
-- 总网后台固定只有一条配置记录。
BEGIN;
  CREATE TABLE PUBLIC.sys_login_security_config (
    -- 总网固定使用ID 1
    id BIGINT NOT NULL DEFAULT 1,
    -- 是否强制后台用户绑定TOTP
    totp_required BOOLEAN NOT NULL DEFAULT FALSE,
    -- 是否启用登录IP白名单
    ip_whitelist_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    -- 是否允许同一个用户同时存在多个有效会话
    allow_concurrent_login BOOLEAN NOT NULL DEFAULT FALSE,
    -- 连续密码错误次数上限
    password_error_limit INTEGER NOT NULL DEFAULT 5,
    -- 连续验证码错误次数上限
    captcha_error_limit INTEGER NOT NULL DEFAULT 5,
    -- 达到错误次数上限后的锁定分钟数
    login_lock_minutes INTEGER NOT NULL DEFAULT 30,
    -- 创建时间
    created_at TIMESTAMPTZ (3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间
    updated_at TIMESTAMPTZ (3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_sys_login_security_config PRIMARY KEY (id)
  );
  INSERT INTO PUBLIC.sys_login_security_config (id, totp_required, ip_whitelist_enabled, allow_concurrent_login, password_error_limit, captcha_error_limit, login_lock_minutes)
  VALUES
  (1, FALSE, FALSE, FALSE, 5, 5, 30) ON CONFLICT (id) DO
    NOTHING;
    COMMENT ON TABLE PUBLIC.sys_login_security_config IS '总网后台登录安全配置表';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.id IS '固定配置主键，值为1';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.totp_required IS '是否强制绑定TOTP';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.ip_whitelist_enabled IS '是否启用登录IP白名单';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.allow_concurrent_login IS '是否允许同一用户同时存在多个有效会话';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.password_error_limit IS '连续密码错误次数上限';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.captcha_error_limit IS '连续验证码错误次数上限';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.login_lock_minutes IS '达到错误次数上限后的锁定分钟数';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.created_at IS '创建时间';
    COMMENT ON COLUMN PUBLIC.sys_login_security_config.updated_at IS '更新时间';
    COMMIT;