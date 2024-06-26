-- meme table
CREATE TABLE meme (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    total_supply VARCHAR(50) NOT NULL,
    decimals SMALLINT NOT NULL,
    logo_url VARCHAR(256) NOT NULL,
    banner_url VARCHAR(256) NOT NULL,
    creator_address VARCHAR(50) NOT NULL,
    contract_address VARCHAR(50) NOT NULL,
    swap_fee_bps SMALLINT NOT NULL,
    vesting_alloc_bps SMALLINT NOT NULL,
    meta VARCHAR(10) NOT NULL,
    live BOOLEAN NOT NULL,
    network_id INT NOT NULL,
    website VARCHAR(50) NOT NULL,
    salt VARCHAR(66) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB;
CREATE INDEX meme_contract_address_idx ON meme (contract_address);
CREATE INDEX meme_symbol_idx ON meme (symbol);
-- memception table
CREATE TABLE memeception (
    id INT NOT NULL AUTO_INCREMENT,
    meme_id INT NOT NULL,
    start_at INT NOT NULL,
    status SMALLINT NOT NULL,
    ama BOOLEAN NOT NULL,
    contract_address VARCHAR(50) NOT NULL,
    target_eth FLOAT NOT NULL,
    collected_eth FLOAT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT 1,
    updated_at_epoch INT NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (meme_id) REFERENCES meme(id)
) ENGINE = InnoDB;
CREATE INDEX memeception_start_at_idx ON memeception (start_at);
-- social table
CREATE TABLE social (
    id INT NOT NULL AUTO_INCREMENT,
    meme_id INT NOT NULL,
    provider VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    photo_url VARCHAR(256),
    url VARCHAR(256),
    PRIMARY KEY (id),
    FOREIGN KEY (meme_id) REFERENCES meme(id)
) ENGINE = InnoDB;