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
    CONSTRAINT "aff_reward_fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE,
    CONSTRAINT "aff_reward_fk_accesstrade_order_id" FOREIGN KEY ("accesstrade_order_id") REFERENCES "aff_order" ("accesstrade_order_id") ON UPDATE CASCADE
) ENGINE = InnoDB;
CREATE TABLE aff_reward_withdraw (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    shipping_request_id VARCHAR(255),
    amount DECIMAL(30, 18) DEFAULT 0,
    fee DECIMAL(30, 18) DEFAULT 0,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT "aff_reward_withdraw_fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE,
    KEY "shipping_request_id" ("shipping_request_id")
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