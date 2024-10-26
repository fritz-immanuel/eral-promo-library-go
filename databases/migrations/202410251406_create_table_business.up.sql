CREATE TABLE business (
  id INT AUTO_INCREMENT,
  name VARCHAR(255) DEFAULT '',
  longitude DECIMAL(15,10) DEFAULT 0.0,
  latitude DECIMAL(15,10) DEFAULT 0.0,
  address VARCHAR(255) DEFAULT '',
  phone_number VARCHAR(255) DEFAULT '',
  wa_number VARCHAR(255) DEFAULT '',
  code VARCHAR(255) DEFAULT '',
  logo_img_url VARCHAR(255) DEFAULT '',
  app_url VARCHAR(255) DEFAULT '',
  on_prem_queue_url VARCHAR(255) DEFAULT '',
  access_token VARCHAR(255) DEFAULT '',
  token VARCHAR(255) DEFAULT '',

  status_id INT DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_business_status_id (status_id)
);