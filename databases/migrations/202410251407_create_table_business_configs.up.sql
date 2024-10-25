CREATE TABLE business_configs (
  id VARCHAR(255) NOT NULL,
  business_id INT DEFAULT 0,
  sub_url_name VARCHAR(255) DEFAULT '',
  config LONGTEXT NOT NULL,
  
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_business_id (business_id),
  INDEX idx_sub_url_name (sub_url_name)
);