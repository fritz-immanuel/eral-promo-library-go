CREATE TABLE permissions (
  id VARCHAR(255) NOT NULL,
  package VARCHAR(255) DEFAULT '',
  module_name VARCHAR(255) DEFAULT '',
  action_name VARCHAR(255) DEFAULT '',
  display_module_name VARCHAR(255) DEFAULT '',
  display_action_name VARCHAR(255) DEFAULT '',
  http_method VARCHAR(255) DEFAULT '',
  route VARCHAR(255) DEFAULT '',
  table_name VARCHAR(255) DEFAULT '',
  is_hidden INT DEFAULT 0,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',

  PRIMARY KEY (id),
  INDEX idx_package (package),
  INDEX idx_module_name (module_name)
);