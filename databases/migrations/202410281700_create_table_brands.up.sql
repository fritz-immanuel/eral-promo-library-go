CREATE TABLE brands (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) DEFAULT '',
  code VARCHAR(255) DEFAULT '',
  logo_img_url LONGTEXT NOT NULL,
  business_id VARCHAR(255) NOT NULL,

  status_id VARCHAR(5) DEFAULT "1",
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_brand_business_id (business_id),
  INDEX idx_brand_status_id (status_id)
);