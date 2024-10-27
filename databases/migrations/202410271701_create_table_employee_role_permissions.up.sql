CREATE TABLE employee_role_permissions (
  id VARCHAR(255) NOT NULL,
  employee_role_id VARCHAR(255) DEFAULT '',
  permission_id INT DEFAULT 0,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_employee_role_id (employee_role_id)
);