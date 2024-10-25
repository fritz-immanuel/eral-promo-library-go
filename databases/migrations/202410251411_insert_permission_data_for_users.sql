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
    'Website',
    'User',
    'List',
    'User',
    'List',
    'GET',
    '/web/v1/users',
    'users',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'User',
    'View',
    'User',
    'View',
    'GET',
    '/web/v1/users/:id',
    'users',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'User',
    'Create',
    'User',
    'Create',
    'POST',
    '/web/v1/users',
    'users',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'User',
    'Edit',
    'User',
    'Edit',
    'PUT',
    '/web/v1/users/:id',
    'users',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'User',
    'UpdateStatus',
    'User',
    'Update Status',
    'PUT',
    '/web/v1/users/:id/status',
    'users',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'User',
    'UpdatePassword',
    'User',,
    'Update Password',
    'PUT',
    '/web/v1/users/:id/password',
    'users',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'User',
    'ResetPassword',
    'User',,
    'Reset Password',
    'PUT',
    '/web/v1/users/:id/reset-password',
    'users',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  );