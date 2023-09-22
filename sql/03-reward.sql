CREATE TABLE aff_reward (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    accesstrade_order_id NVARCHAR(256),
    amount DECIMAL(30, 18) DEFAULT 0,
    rewarded_amount DECIMAL(30, 18) DEFAULT 0,
    commission_fee FLOAT(5, 2) DEFAULT 0,
    ended_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY "aff_reward_user_id" ("user_id"),
    CONSTRAINT "aff_reward_fk_accesstrade_order_id" FOREIGN KEY ("accesstrade_order_id") REFERENCES "aff_order" ("accesstrade_order_id") ON UPDATE CASCADE,
    CONSTRAINT "aff_reward_unique_accesstrade_order_id" UNIQUE ("accesstrade_order_id")
) ENGINE = InnoDB;
CREATE TABLE aff_reward_withdraw (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    shipping_request_id VARCHAR(255),
    tx_hash VARCHAR(66) NOT NULL DEFAULT '',
    shipping_status enum('initial', 'sending', 'success', 'failed') NOT NULL DEFAULT 'initial',
    amount DECIMAL(30, 18) DEFAULT 0,
    fee DECIMAL(30, 18) DEFAULT 0,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY "aff_reward_withdraw_user_id" ("user_id"),
    KEY "aff_reward_withdraw_shipping_request_id" ("shipping_request_id"),
    KEY "aff_reward_withdraw_tx_hash" ("tx_hash")
) ENGINE = InnoDB;
CREATE TABLE aff_reward_order_history (
    id INT NOT NULL AUTO_INCREMENT,
    reward_id INT NOT NULL,
    reward_withdraw_id INT NOT NULL,
    amount DECIMAL(30, 18) DEFAULT 0,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT "aff_reward_order_history_fk_reward_id" FOREIGN KEY ("reward_id") REFERENCES "aff_reward" ("id") ON UPDATE CASCADE,
    CONSTRAINT "aff_reward_order_history_fk_reward_withdraw_id" FOREIGN KEY ("reward_withdraw_id") REFERENCES "aff_reward_withdraw" ("id") ON UPDATE CASCADE
) ENGINE = InnoDB;