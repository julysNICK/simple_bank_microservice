-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
--method normal
-- UPDATE users SET
--   username = $2,
--   hashed_password = $3,
--   full_name = $4,
--   email = $5
-- WHERE username = $1  
-- RETURNING *;
-- method when not all fields are updated
-- UPDATE users 
--   Set 
--   hashed_password = CASE
--     WHEN @set_hashed_password::boolean = TRUE THEN @hashed_password
--     ELSE hashed_password 
--   END,
--   full_name = CASE
--     WHEN @set_full_name = TRUE THEN @full_name
--     ELSE full_name
--   END,
--   email = CASE
--     WHEN  @set_email = TRUE THEN @email
--     ELSE email 
--   END
-- WHERE username = @username
-- RETURNING *;
-- method when not all fields are updated method 3
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email),
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
  
WHERE
  username = sqlc.arg(username)
RETURNING *;