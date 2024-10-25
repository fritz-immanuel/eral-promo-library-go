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
    'Business',
    'List',
    'Business',
    'List',
    'GET',
    '/web/v1/business',
    'business',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'Business',
    'View',
    'Business',
    'View',
    'GET',
    '/web/v1/business/:id',
    'business',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'Business',
    'Create',
    'Business',
    'Create',
    'POST',
    '/web/v1/business',
    'business',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'Business',
    'Edit',
    'Business',
    'Edit',
    'PUT',
    '/web/v1/business/:id',
    'business',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  ), (
    'Website',
    'Business',
    'UpdateStatus',
    'Business',
    'Update Status',
    'PUT',
    '/web/v1/business/:id/status',
    'business',
    CURRENT_TIMESTAMP,
    '0',
    CURRENT_TIMESTAMP,
    '0'
  );