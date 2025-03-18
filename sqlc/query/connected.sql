-- name: GetCommentsByPage :many
SELECT Comments.CommentId, Comments.UserId, Users.UserName, Comments.CreatedAt,  Comments.CommentData 
FROM Comments
JOIN Users ON Comments.UserId = Users.UserId
WHERE Comments.PageId = $1 
AND Comments.ParentId = 0
ORDER BY Comments.CreatedAt ASC
LIMIT $2;


-- name: GetUserIDByCommentID :one
SELECT UserId FROM Comments WHERE CommentId = $1;