-- name: CreateContact :one
INSERT INTO contact (first_name, last_name, age, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetContactById :one
SELECT * FROM contact
WHERE Id = $1;

-- name: GetContactByPhone :many
SELECT * FROM contact
WHERE Phone = $1;

-- name: UpdateContact :one
UPDATE contact
	set first_name = $2,
	last_name = $3,
	age = $4,
	phone = $5
WHERE id = $1
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contact
WHERE Id = $1;

-- name: GetAllContacts :many
SELECT * FROM contact;
