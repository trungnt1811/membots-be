-- Active: 1692779176081@@mysql-reward-do-user-1468358-0.b.db.ondigitalocean.com@25060@reward
-- CAMPAIGN
CREATE TABLE aff_campaign (
  id INT NOT NULL AUTO_INCREMENT,
  accesstrade_id NVARCHAR(100) NOT NULL,
  active_status SMALLINT,
  logo NVARCHAR(256),
  max_com NVARCHAR(256),
  created_at DATETIME,
  updated_at DATETIME,
  merchant NVARCHAR(256),
  name NVARCHAR(256),
  scope NVARCHAR(256),
  approval NVARCHAR(256),
  status SMALLINT,
  type SMALLINT,
  url NVARCHAR(256),
  category NVARCHAR(256),
  sub_category NVARCHAR(256),
  cookie_duration INT,
  cookie_policy NVARCHAR(256),
  start_time DATETIME,
  end_time DATETIME,
  brand_id INT,
  PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE INDEX aff_campaign_created ON aff_campaign (created_at);
CREATE INDEX aff_campaign_updated ON aff_campaign (updated_at);
CREATE INDEX aff_campaign_accesstrade_id ON aff_campaign (accesstrade_id);
CREATE INDEX aff_campaign_merchant ON aff_campaign (merchant);
CREATE INDEX aff_campaign_category ON aff_campaign (category);

CREATE TABLE aff_campaign_description (
  id INT NOT NULL AUTO_INCREMENT,
  campaign_id INT NOT NULL,
  action_point MEDIUMTEXT,
  commission_policy MEDIUMTEXT,
  cookie_policy MEDIUMTEXT,
  introduction MEDIUMTEXT,
  other_notice MEDIUMTEXT,
  rejected_reason MEDIUMTEXT,
  traffic_building_policy MEDIUMTEXT,
  created_at DATETIME,
  updated_at DATETIME,
  PRIMARY KEY (id),
  FOREIGN KEY (campaign_id) REFERENCES aff_campaign(id)
) ENGINE = InnoDB;
CREATE INDEX aff_campaign_description_created ON aff_campaign_description (created_at);
CREATE INDEX aff_campaign_description_updated ON aff_campaign_description (updated_at);

-- LINK
CREATE TABLE aff_link (
  id INT NOT NULL AUTO_INCREMENT,
  campaign_id INT NOT NULL,
  aff_link NVARCHAR(1024),
  first_link NVARCHAR(1024),
  short_link NVARCHAR(1024),
  url_origin NVARCHAR(1024),
  created_at DATETIME,
  updated_at DATETIME,
  active_status SMALLINT DEFAULT 1,
  PRIMARY KEY (id),
  FOREIGN KEY (campaign_id) REFERENCES aff_campaign(id)
) ENGINE = InnoDB;

CREATE INDEX aff_link_created ON aff_link (created_at);
CREATE INDEX aff_link_updated ON aff_link (updated_at);

-- BRAND
ALTER TABLE aff_campaign ADD COLUMN brand_id INT;
ALTER TABLE aff_campaign
ADD FOREIGN KEY (brand_id) REFERENCES brand(id);
CREATE INDEX aff_campaign_brand ON aff_campaign (brand_id);
