INSERT INTO
  `permissions`(
    package,
    module_name,
    action_name,
    display_module_name,
    display_action_name,
    http_method,
    route,
    table_name,
    created_at,
    created_by,
    updated_at,
    updated_by,
    is_hidden
  )
VALUES
  (
    'WebsiteApp',
    'Employee',
    'ViewpProfile',
    'Employee',
    'View Profile',
    'GET',
    '/web/v1/employees/profile',
    'employees',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0',
    1
  ), (
    'WebsiteApp',
    'Employee',
    'EditPassword',
    'Employee',
    'Edit Password',
    'PUT',
    '/web/v1/employees/profile/password',
    'employees',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0',
    1
  );