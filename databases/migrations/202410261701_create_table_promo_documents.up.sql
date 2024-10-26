CREATE TABLE promo_documents (
  id VARCHAR(255) NOT NULL,
  promo_id VARCHAR(255) DEFAULT '',
  document_url VARCHAR(255) DEFAULT '',

  status_id VARCHAR(5) DEFAULT '1', 
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',

  PRIMARY KEY (id),
  INDEX idx_promo_id (promo_id),
  INDEX idx_status_id (status_id)
);