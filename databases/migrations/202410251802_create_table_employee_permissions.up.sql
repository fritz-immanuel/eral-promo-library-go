CREATE TABLE employee_permissions (
  id VARCHAR(255) NOT NULL,
  employee_id VARCHAR(255) DEFAULT '',
  permission_id INT DEFAULT 0,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',

  PRIMARY KEY (id),
  INDEX idx_employee_id (employee_id),
  INDEX idx_permission_id (permission_id)
);