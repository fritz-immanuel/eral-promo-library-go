// Code generated by rice embed-go; DO NOT EDIT.
package databases

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "202410251400_create_table_status.up.sql",
		FileModTime: time.Unix(1729840094, 0),

		Content: string("CREATE TABLE status (\r\n  id VARCHAR(5) NOT NULL,\r\n  name VARCHAR(255) DEFAULT \"\",\r\n  PRIMARY KEY (id)\r\n);"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "202410251401_insert_status_data.up.sql",
		FileModTime: time.Unix(1729840102, 0),

		Content: string("INSERT INTO\r\n  status (id, name)\r\nVALUES\r\n  ('0', 'Inactive'),\r\n  ('1', 'Active');"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "202410251402_create_table_days.up.sql",
		FileModTime: time.Unix(1729840127, 0),

		Content: string("CREATE TABLE days (\r\n  id VARCHAR(5) NOT NULL,\r\n  name VARCHAR(255) DEFAULT \"\",\r\n  name_en VARCHAR(255) DEFAULT \"\",\r\n  PRIMARY KEY (id)\r\n);"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "202410251403_insert_days_data.up.sql",
		FileModTime: time.Unix(1729840181, 0),

		Content: string("INSERT INTO\r\n  days (id, name, name_en)\r\nVALUES\r\n  ('1', 'Senin', 'Monday'),\r\n  ('2', 'Selasa', 'Tuesday'),\r\n  ('3', 'Rabu', 'Wednesday'),\r\n  ('4', 'Kamis', 'Thursday'),\r\n  ('5', 'Jumat', 'Friday'),\r\n  ('6', 'Sabtu', 'Saturday'),\r\n  ('7', 'Minggu', 'Sunday');"),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "202410251404_create_table_users.up.sql",
		FileModTime: time.Unix(1729855932, 0),

		Content: string("CREATE TABLE users (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  email VARCHAR(255) DEFAULT '',\r\n  username VARCHAR(255) NOT NULL,\r\n  password VARCHAR(255) NOT NULL,\r\n  \r\n  status_id VARCHAR(5) DEFAULT '1',\r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n  PRIMARY KEY (id),\r\n  INDEX idx_username (username)\r\n);"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "202410251405_insert_users_init_data.up.sql",
		FileModTime: time.Unix(1729839900, 0),

		Content: string("INSERT INTO\r\n  users (id, name, email, username, password, business_id, status_id)\r\nVALUES\r\n  (UUID(), 'Fritz', 'fritz@erajaya.com', 'fritz', 'e10adc3949ba59abbe56e057f20f883e', 1, 1);"),
	}
	file8 := &embedded.EmbeddedFile{
		Filename:    "202410251406_create_table_business.up.sql",
		FileModTime: time.Unix(1729934340, 0),

		Content: string("CREATE TABLE business (\r\n  id INT AUTO_INCREMENT,\r\n  name VARCHAR(255) DEFAULT '',\r\n  longitude DECIMAL(15,10) DEFAULT 0.0,\r\n  latitude DECIMAL(15,10) DEFAULT 0.0,\r\n  address VARCHAR(255) DEFAULT '',\r\n  phone_number VARCHAR(255) DEFAULT '',\r\n  wa_number VARCHAR(255) DEFAULT '',\r\n  code VARCHAR(255) DEFAULT '',\r\n  logo_img_url VARCHAR(255) DEFAULT '',\r\n  app_url VARCHAR(255) DEFAULT '',\r\n  on_prem_queue_url VARCHAR(255) DEFAULT '',\r\n  access_token VARCHAR(255) DEFAULT '',\r\n  token VARCHAR(255) DEFAULT '',\r\n\r\n  status_id INT DEFAULT 1,\r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n  PRIMARY KEY (id),\r\n  INDEX idx_business_status_id (status_id)\r\n);"),
	}
	file9 := &embedded.EmbeddedFile{
		Filename:    "202410251407_create_table_business_configs.up.sql",
		FileModTime: time.Unix(1729841194, 0),

		Content: string("CREATE TABLE business_configs (\r\n  id VARCHAR(255) NOT NULL,\r\n  business_id INT DEFAULT 0,\r\n  sub_url_name VARCHAR(255) DEFAULT '',\r\n  config LONGTEXT NOT NULL,\r\n  \r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n  PRIMARY KEY (id),\r\n  INDEX idx_business_id (business_id),\r\n  INDEX idx_sub_url_name (sub_url_name)\r\n);"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "202410251408_create_table_user_actions.up.sql",
		FileModTime: time.Unix(1729840060, 0),

		Content: string("CREATE TABLE user_actions (\r\n  id VARCHAR(255) NOT NULL,\r\n  user_id INT DEFAULT 0,\r\n  table_name VARCHAR(200) DEFAULT NULL,\r\n  action VARCHAR(100) DEFAULT NULL,\r\n  action_value INT DEFAULT 0,\r\n  ref_id VARCHAR(255) DEFAULT '0',\r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  PRIMARY KEY (id),\r\n  INDEX idx_user_id (user_id),\r\n  INDEX idx_ref_id (ref_id),\r\n  INDEX idx_table_name (table_name)\r\n);"),
	}
	fileb := &embedded.EmbeddedFile{
		Filename:    "202410251409_create_table_user_permissions.sql",
		FileModTime: time.Unix(1729842674, 0),

		Content: string("CREATE TABLE user_permissions (\r\n  id VARCHAR(255) NOT NULL,\r\n  user_id VARCHAR(255) DEFAULT '',\r\n  permission_id INT DEFAULT 0,\r\n\r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n\r\n  PRIMARY KEY (id),\r\n  INDEX idx_user_id (user_id),\r\n  INDEX idx_permission_id (permission_id)\r\n);"),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "202410251410_create_table_permissions.up.sql",
		FileModTime: time.Unix(1729842008, 0),

		Content: string("CREATE TABLE permissions (\r\n  id VARCHAR(255) NOT NULL,\r\n  package VARCHAR(255) DEFAULT '',\r\n  module_name VARCHAR(255) DEFAULT '',\r\n  action_name VARCHAR(255) DEFAULT '',\r\n  display_module_name VARCHAR(255) DEFAULT '',\r\n  display_action_name VARCHAR(255) DEFAULT '',\r\n  http_method VARCHAR(255) DEFAULT '',\r\n  route VARCHAR(255) DEFAULT '',\r\n  table_name VARCHAR(255) DEFAULT '',\r\n  is_hidden INT DEFAULT 0,\r\n\r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n\r\n  PRIMARY KEY (id),\r\n  INDEX idx_package (package),\r\n  INDEX idx_module_name (module_name)\r\n);"),
	}
	filed := &embedded.EmbeddedFile{
		Filename:    "202410251411_insert_permission_data_for_users.up.sql",
		FileModTime: time.Unix(1729862777, 0),

		Content: string("INSERT INTO\r\n  `permission`(\r\n    package,\r\n    module_name,\r\n    action_name,\r\n    display_module_name,\r\n    display_action_name,\r\n    http_method,\r\n    route,\r\n    table_name,\r\n    created_at,\r\n    created_by,\r\n    updated_at,\r\n    updated_by\r\n  )\r\nVALUES\r\n  (\r\n    'WebsiteAdmin',\r\n    'User',\r\n    'List',\r\n    'User',\r\n    'List',\r\n    'GET',\r\n    '/web/v1/users',\r\n    'users',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'User',\r\n    'View',\r\n    'User',\r\n    'View',\r\n    'GET',\r\n    '/web/v1/users/:id',\r\n    'users',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'User',\r\n    'Create',\r\n    'User',\r\n    'Create',\r\n    'POST',\r\n    '/web/v1/users',\r\n    'users',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'User',\r\n    'Edit',\r\n    'User',\r\n    'Edit',\r\n    'PUT',\r\n    '/web/v1/users/:id',\r\n    'users',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'User',\r\n    'UpdateStatus',\r\n    'User',\r\n    'Update Status',\r\n    'PUT',\r\n    '/web/v1/users/:id/status',\r\n    'users',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'User',\r\n    'UpdatePassword',\r\n    'User',,\r\n    'Update Password',\r\n    'PUT',\r\n    '/web/v1/users/:id/password',\r\n    'users',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'User',\r\n    'ResetPassword',\r\n    'User',,\r\n    'Reset Password',\r\n    'PUT',\r\n    '/web/v1/users/:id/reset-password',\r\n    'users',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  );"),
	}
	filee := &embedded.EmbeddedFile{
		Filename:    "202410251412_insert_permission_data_for_business.up.sql",
		FileModTime: time.Unix(1729862777, 0),

		Content: string("INSERT INTO\r\n  `permission`(\r\n    package,\r\n    module_name,\r\n    action_name,\r\n    display_module_name,\r\n    display_action_name,\r\n    http_method,\r\n    route,\r\n    table_name,\r\n    created_at,\r\n    created_by,\r\n    updated_at,\r\n    updated_by\r\n  )\r\nVALUES\r\n  (\r\n    'WebsiteAdmin',\r\n    'Business',\r\n    'List',\r\n    'Business',\r\n    'List',\r\n    'GET',\r\n    '/web/v1/business',\r\n    'business',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'Business',\r\n    'View',\r\n    'Business',\r\n    'View',\r\n    'GET',\r\n    '/web/v1/business/:id',\r\n    'business',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'Business',\r\n    'Create',\r\n    'Business',\r\n    'Create',\r\n    'POST',\r\n    '/web/v1/business',\r\n    'business',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'Business',\r\n    'Edit',\r\n    'Business',\r\n    'Edit',\r\n    'PUT',\r\n    '/web/v1/business/:id',\r\n    'business',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  ), (\r\n    'WebsiteAdmin',\r\n    'Business',\r\n    'UpdateStatus',\r\n    'Business',\r\n    'Update Status',\r\n    'PUT',\r\n    '/web/v1/business/:id/status',\r\n    'business',\r\n    CURRENT_TIMESTAMP,\r\n    '0',\r\n    CURRENT_TIMESTAMP,\r\n    '0'\r\n  );"),
	}
	filef := &embedded.EmbeddedFile{
		Filename:    "202410251800_create_table_employees.up.sql",
		FileModTime: time.Unix(1729855961, 0),

		Content: string("CREATE TABLE employees (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  email VARCHAR(255) DEFAULT '',\r\n  username VARCHAR(255) NOT NULL,\r\n  password VARCHAR(255) NOT NULL,\r\n  \r\n  status_id VARCHAR(5) DEFAULT '1',\r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n  PRIMARY KEY (id),\r\n  INDEX idx_username (username)\r\n);"),
	}
	fileg := &embedded.EmbeddedFile{
		Filename:    "202410251801_insert_employee_init_data.up.sql",
		FileModTime: time.Unix(1729858544, 0),

		Content: string("INSERT INTO\r\n  employees (id, name, email, username, password, business_id, status_id)\r\nVALUES\r\n  (UUID(), 'Fritz', 'fritz@erajaya.com', 'fritz', 'e10adc3949ba59abbe56e057f20f883e', 1, 1);"),
	}
	fileh := &embedded.EmbeddedFile{
		Filename:    "202410251802_create_table_employee_permissions.up.sql",
		FileModTime: time.Unix(1729856223, 0),

		Content: string("CREATE TABLE employee_permissions (\r\n  id VARCHAR(255) NOT NULL,\r\n  employee_id VARCHAR(255) DEFAULT '',\r\n  permission_id INT DEFAULT 0,\r\n\r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n\r\n  PRIMARY KEY (id),\r\n  INDEX idx_employee_id (employee_id),\r\n  INDEX idx_permission_id (permission_id)\r\n);"),
	}
	filei := &embedded.EmbeddedFile{
		Filename:    "202410261700_create_table_promos.up.sql",
		FileModTime: time.Unix(1729937187, 0),

		Content: string("CREATE TABLE promos (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) DEFAULT '',\r\n  code VARCHAR(255) DEFAULT '',\r\n  promo_type_id VARCHAR(255) NOT NULL,\r\n  img_url VARCHAR(255) DEFAULT '',\r\n  company_id VARCHAR(255) DEFAULT '',\r\n  business_id VARCHAR(255) DEFAULT '',\r\n  total_promo_budget DECIMAL(25,2) DEFAULT 0.0,\r\n  principle_support DECIMAL(3,2) DEFAULT 0.0,\r\n  internal_support DECIMAL(3,2) DEFAULT 0.0,\r\n\r\n  status_id VARCHAR(5) DEFAULT '1', \r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n\r\n  PRIMARY KEY (id),\r\n  INDEX idx_promo_type_id (promo_type_id),\r\n  INDEX idx_company_id (company_id),\r\n  INDEX idx_business_id (business_id),\r\n  INDEX idx_status_id (status_id)\r\n);"),
	}
	filej := &embedded.EmbeddedFile{
		Filename:    "202410261701_create_table_promo_documents.up.sql",
		FileModTime: time.Unix(1729940157, 0),

		Content: string("CREATE TABLE promo_documents (\r\n  id VARCHAR(255) NOT NULL,\r\n  promo_id VARCHAR(255) DEFAULT '',\r\n  document_url VARCHAR(255) DEFAULT '',\r\n\r\n  status_id VARCHAR(5) DEFAULT '1', \r\n  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  created_by VARCHAR(255) DEFAULT '',\r\n  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,\r\n  updated_by VARCHAR(255) DEFAULT '',\r\n\r\n  PRIMARY KEY (id),\r\n  INDEX idx_promo_id (promo_id),\r\n  INDEX idx_status_id (status_id)\r\n);"),
	}
	filek := &embedded.EmbeddedFile{
		Filename:    "202410261702_create_table_promo_status.up.sql",
		FileModTime: time.Unix(1729942247, 0),

		Content: string("CREATE TABLE promo_status (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) DEFAULT '',\r\n  PRIMARY KEY (id)\r\n);"),
	}
	filel := &embedded.EmbeddedFile{
		Filename:    "202410261703_insert_promo_status_data.up.sql",
		FileModTime: time.Unix(1729942323, 0),

		Content: string("INSERT INTO\r\n  promo_status (id, name)\r\nVALUES\r\n  ('0', 'Inactive'),\r\n  ('1', 'Active'),\r\n  ('2', 'Submitted');"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1729942270, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "202410251400_create_table_status.up.sql"
			file3, // "202410251401_insert_status_data.up.sql"
			file4, // "202410251402_create_table_days.up.sql"
			file5, // "202410251403_insert_days_data.up.sql"
			file6, // "202410251404_create_table_users.up.sql"
			file7, // "202410251405_insert_users_init_data.up.sql"
			file8, // "202410251406_create_table_business.up.sql"
			file9, // "202410251407_create_table_business_configs.up.sql"
			filea, // "202410251408_create_table_user_actions.up.sql"
			fileb, // "202410251409_create_table_user_permissions.sql"
			filec, // "202410251410_create_table_permissions.up.sql"
			filed, // "202410251411_insert_permission_data_for_users.up.sql"
			filee, // "202410251412_insert_permission_data_for_business.up.sql"
			filef, // "202410251800_create_table_employees.up.sql"
			fileg, // "202410251801_insert_employee_init_data.up.sql"
			fileh, // "202410251802_create_table_employee_permissions.up.sql"
			filei, // "202410261700_create_table_promos.up.sql"
			filej, // "202410261701_create_table_promo_documents.up.sql"
			filek, // "202410261702_create_table_promo_status.up.sql"
			filel, // "202410261703_insert_promo_status_data.up.sql"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`./migrations`, &embedded.EmbeddedBox{
		Name: `./migrations`,
		Time: time.Unix(1729942270, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"202410251400_create_table_status.up.sql":                 file2,
			"202410251401_insert_status_data.up.sql":                  file3,
			"202410251402_create_table_days.up.sql":                   file4,
			"202410251403_insert_days_data.up.sql":                    file5,
			"202410251404_create_table_users.up.sql":                  file6,
			"202410251405_insert_users_init_data.up.sql":              file7,
			"202410251406_create_table_business.up.sql":               file8,
			"202410251407_create_table_business_configs.up.sql":       file9,
			"202410251408_create_table_user_actions.up.sql":           filea,
			"202410251409_create_table_user_permissions.sql":          fileb,
			"202410251410_create_table_permissions.up.sql":            filec,
			"202410251411_insert_permission_data_for_users.up.sql":    filed,
			"202410251412_insert_permission_data_for_business.up.sql": filee,
			"202410251800_create_table_employees.up.sql":              filef,
			"202410251801_insert_employee_init_data.up.sql":           fileg,
			"202410251802_create_table_employee_permissions.up.sql":   fileh,
			"202410261700_create_table_promos.up.sql":                 filei,
			"202410261701_create_table_promo_documents.up.sql":        filej,
			"202410261702_create_table_promo_status.up.sql":           filek,
			"202410261703_insert_promo_status_data.up.sql":            filel,
		},
	})
}
