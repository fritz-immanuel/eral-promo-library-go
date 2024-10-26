CREATE TABLE promos (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) DEFAULT '',
  code VARCHAR(255) DEFAULT '',
  promo_type_id VARCHAR(255) NOT NULL,
  img_url VARCHAR(255) DEFAULT '',
  company_id VARCHAR(255) DEFAULT '',
  business_id VARCHAR(255) DEFAULT '',
  total_promo_budget DECIMAL(25,2) DEFAULT 0.0,
  principle_support DECIMAL(3,2) DEFAULT 0.0,
  internal_support DECIMAL(3,2) DEFAULT 0.0,

  status_id VARCHAR(5) DEFAULT '1', 
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',

  PRIMARY KEY (id),
  INDEX idx_promo_type_id (promo_type_id),
  INDEX idx_company_id (company_id),
  INDEX idx_business_id (business_id),
  INDEX idx_status_id (status_id)
);