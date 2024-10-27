CREATE TABLE employees (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) DEFAULT '',
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  business_id VARCHAR(255) NOT NULL,
  employee_role_id VARCHAR(255) DEFAULT '',
  
  status_id VARCHAR(5) DEFAULT '1',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_username (username),
  INDEX idx_business_id (business_id),
  INDEX idx_employee_role_id (employee_role_id)
);