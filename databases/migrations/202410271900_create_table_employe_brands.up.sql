CREATE TABLE employee_brands (
  id VARCHAR(255) NOT NULL,
  employee_id VARCHAR(255) DEFAULT '',
  brand_id VARCHAR(255) DEFAULT '',

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_employee_id (employee_id)
);