CREATE TABLE employee_roles (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) DEFAULT '',
  is_supervisor INT DEFAULT 0,

  status_id VARCHAR(5) DEFAULT "1",
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_status_id (status_id)
);