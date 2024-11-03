-- User / Admin Permissions
DELETE FROM
	user_permissions
WHERE
	user_id = "<user_id>";

INSERT INTO
	user_permissions (id, user_id, permission_id)
SELECT
	UUID(),
	"<user_id>",
	id
FROM
	permissions
WHERE permissions.package = "WebsiteAdmin";

-- Employee Role Permissions
DELETE FROM
	employee_role_permissions
WHERE
	employee_role_id = "<role_id>";

INSERT INTO
	employee_role_permissions (id, employee_role_id, permission_id)
SELECT
	UUID(),
	"<role_id>",
	id
FROM
	permissions
WHERE permissions.package = "WebsiteApp";