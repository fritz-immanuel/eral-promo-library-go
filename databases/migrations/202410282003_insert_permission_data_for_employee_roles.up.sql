INSERT INTO
  `permission`(
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
    updated_by
  )
VALUES
  (
    'WebsiteAdmin',
    'EmployeeRole',
    'List',
    'Employee Role',
    'List',
    'GET',
    '/admin/v1/employee-roles',
    'employee_roles',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'WebsiteAdmin',
    'EmployeeRole',
    'View',
    'Employee Role',
    'View',
    'GET',
    '/admin/v1/employee-roles/:id',
    'employee_roles',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'WebsiteAdmin',
    'Employee Role',
    'Create',
    'EmployeeRole',
    'Create',
    'POST',
    '/admin/v1/employee-roles',
    'employee_roles',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'WebsiteAdmin',
    'EmployeeRole',
    'Edit',
    'Employee Role',
    'Edit',
    'PUT',
    '/admin/v1/employee-roles/:id',
    'employee_roles',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'WebsiteAdmin',
    'EmployeeRole',
    'UpdateStatus',
    'Employee Role',
    'Update Status',
    'PUT',
    '/admin/v1/employee-roles/:id/status',
    'employee_roles',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  );