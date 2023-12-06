-- name: GetUserById :one
SELECT id, username, role, email, "displayName", "phoneNumber", "mainWallet", dob, "profileImage", "isEmailVerified" FROM betonz."User" WHERE id = $1;

-- name: GetExtendedPlayerByUsername :one
SELECT * FROM betonz."User" WHERE username = $1 AND role = 'PLAYER'::betonz."Role";

-- name: GetExtendedAdminByUsername :one
SELECT * FROM betonz."User" WHERE username = $1 AND role IN ('ADMIN'::betonz."Role", 'SUPERADMIN'::betonz."Role");

