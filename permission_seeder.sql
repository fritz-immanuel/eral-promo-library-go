DELETE FROM
	user_permissions
WHERE
	user_id = "";

INSERT INTO
	user_permissions (id, user_id, permission_id)
SELECT
	UUID(),
	"",
	id
FROM
	permissions;