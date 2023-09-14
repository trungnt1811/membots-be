
CREATE TABLE aff_merchant_brand (
  id INT NOT NULL AUTO_INCREMENT,
  merchant NVARCHAR(256),
  brand_id INT NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (brand_id) REFERENCES brand(id)
) ENGINE = InnoDB;

CREATE INDEX aff_merchant_brand_idx ON aff_merchant_brand (merchant);

-- AFF TRANSACTIONS AND ORDERS
CREATE TABLE aff_transaction (
  id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  -- ACCESS TRADE DATA
  accesstrade_id NVARCHAR(256),
  accesstrade_order_id NVARCHAR(256),
  accesstrade_conversion_id BIGINT,
  merchant NVARCHAR(256),
  status SMALLINT,
  click_time DATETIME,
  transaction_time DATETIME,
  transaction_value DOUBLE(15, 2),
  update_time DATETIME,
  confirmed_time DATETIME,
  is_confirmed SMALLINT,
  commission DOUBLE(15,2),
  product_id NVARCHAR(256),
  product_name NVARCHAR(256),
  product_price DOUBLE(15,2),
  product_quantity INT,
  product_image NVARCHAR(1024),
  product_category NVARCHAR(256),
  extra JSON,
  category_name NVARCHAR(256),
  conversion_platform NVARCHAR(256),
  click_url NVARCHAR(1024),
  utm_term NVARCHAR(256),
  utm_source NVARCHAR(256),
  utm_campaign NVARCHAR(256),
  utm_medium NVARCHAR(256),
  utm_content NVARCHAR(256),
  reason_rejected TEXT,
  customer_type NVARCHAR(256),
  -- INDEXING
  PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE INDEX aff_transaction_created ON aff_transaction (created_at);
CREATE INDEX aff_transaction_updated ON aff_transaction (updated_at);
CREATE INDEX aff_transaction_accesstrade_order_id ON aff_transaction (accesstrade_order_id);

CREATE TABLE aff_order (
  id INT NOT NULL AUTO_INCREMENT,
  aff_link NVARCHAR(1024),
  created_at DATETIME,
  updated_at DATETIME,
  user_id INT NOT NULL,
  order_status ENUM ('initial','pending','approved', 'rejected', 'rewarding', 'success'),
  -- ACCESS TRADE DATA
  at_product_link NVARCHAR(1024),
  billing DOUBLE(15,2),
  browser NVARCHAR(1024),
  category_name NVARCHAR(256),
  client_platform NVARCHAR(256),
  click_time DATETIME,
  confirmed_time DATETIME,
  conversion_platform NVARCHAR(256),
  customer_type NVARCHAR(256),
  is_confirmed SMALLINT,
  landing_page NVARCHAR(1024),
  merchant NVARCHAR(256),
  accesstrade_order_id NVARCHAR(256),
  order_pending SMALLINT,
  order_reject SMALLINT,
  order_approved SMALLINT,
  product_category NVARCHAR(256),
  products_count INT,
  pub_commission DOUBLE(15,2),
  sales_time DATETIME,
  update_time DATETIME,
  website NVARCHAR(1024),
  website_url NVARCHAR(1024),
  utm_term NVARCHAR(256),
  utm_source NVARCHAR(256),
  utm_campaign NVARCHAR(256),
  utm_medium NVARCHAR(256),
  utm_content NVARCHAR(256),
  -- INDEXING
  PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE INDEX aff_order_created ON aff_order (created_at);
CREATE INDEX aff_order_updated ON aff_order (updated_at);
CREATE INDEX aff_order_accesstrade_order_id ON aff_order (accesstrade_order_id);


CREATE TABLE aff_postback_log (
  id INT NOT NULL AUTO_INCREMENT,
  order_id NVARCHAR(256),
  created_at DATETIME,
  updated_at DATETIME,
  data JSON,
  PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE INDEX aff_postback_log_order_id ON aff_postback_log (order_id);
CREATE INDEX aff_postback_log_created ON aff_postback_log (created_at);
CREATE INDEX aff_postback_log_updated ON aff_postback_log (updated_at);

CREATE TABLE aff_tracked_click (
  id INT NOT NULL AUTO_INCREMENT,
  campaign_id INT,
  user_id INT,
  link_id INT,
  created_at DATETIME,
  updated_at DATETIME,
  aff_link NVARCHAR(1024),
  short_link NVARCHAR(1024),
  url_origin NVARCHAR(1024),
  order_id NVARCHAR(256),
  PRIMARY KEY (id),
  FOREIGN KEY (campaign_id) REFERENCES aff_campaign(id),
  FOREIGN KEY (link_id) REFERENCES aff_link(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE = InnoDB;

CREATE INDEX aff_tracked_click_created ON aff_tracked_click (created_at);
CREATE INDEX aff_tracked_click_order_id ON aff_tracked_click (order_id);
CREATE INDEX aff_tracked_click_updated ON aff_tracked_click (updated_at);
