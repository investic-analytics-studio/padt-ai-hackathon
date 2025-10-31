package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type NotificationRepo struct {
	db *sqlx.DB
}

func NewCryptoNotificationRepo(db *sqlx.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) GetNotificationGroupList(ctx context.Context, uid string) ([]model.NotificationGroupList, error) {
	type flatGroupRow struct {
		GroupID    sql.NullString `db:"group_id"`
		GroupName  sql.NullString `db:"group_name"`
		AuthorID   sql.NullString `db:"id"`
		AuthorName sql.NullString `db:"name"`
		CreatedAt  time.Time      `db:"created_at"`
		UpdatedAt  time.Time      `db:"updated_at"`
	}

	query := `
		SELECT 
			grp.group_id, 
			grp.group_name,
			grp.created_at,
			grp.updated_at,
			auth.id,
			auth.authors_username AS name
		FROM crypto_notification_group grp
		LEFT JOIN crypto_notification_authors_list auth 
			ON auth.group_id = grp.group_id
		WHERE grp.crypto_user_id = $1
		ORDER BY grp.group_name;
	`

	var rows []flatGroupRow
	if err := r.db.SelectContext(ctx, &rows, query, uid); err != nil {
		return nil, fmt.Errorf("failed to get notification group: %w", err)
	}

	groupMap := make(map[string]*model.NotificationGroupList)

	for _, row := range rows {
		if !row.GroupID.Valid {
			continue // skip malformed rows
		}
		groupID := row.GroupID.String

		group, exists := groupMap[groupID]
		if !exists {
			group = &model.NotificationGroupList{
				GroupID:     groupID,
				GroupName:   row.GroupName.String,
				AuthorNames: []model.Author{},
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
			}
			groupMap[groupID] = group
		}

		if row.AuthorID.Valid {
			group.AuthorNames = append(group.AuthorNames, model.Author{
				ID:   row.AuthorID.String,
				Name: row.AuthorName.String,
			})
		}
	}

	result := make([]model.NotificationGroupList, 0, len(groupMap))
	for _, group := range groupMap {
		result = append(result, *group)
	}

	return result, nil
}

func (r *NotificationRepo) UpdateGroupName(ctx context.Context, groupID string, newGroupName string) error {
	query := `
		UPDATE crypto_notification_group
		SET group_name = $1
		WHERE group_id = $2;
	`
	_, err := r.db.ExecContext(ctx, query, newGroupName, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group name: %w", err)
	}
	return nil
}
func (r *NotificationRepo) AddGroup(ctx context.Context, uid string, groupName string) error {
	query := `
		INSERT INTO crypto_notification_group (group_name, crypto_user_id)
		VALUES ($1, $2);
	`
	_, err := r.db.ExecContext(ctx, query, groupName, uid)
	if err != nil {
		return fmt.Errorf("failed to add group: %w", err)
	}
	return nil
}
func (r *NotificationRepo) CountGroup(ctx context.Context, uid string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM crypto_notification_group
		WHERE crypto_user_id = $1;
	`
	var count int
	err := r.db.GetContext(ctx, &count, query, uid)
	if err != nil {
		return 0, fmt.Errorf("failed to count group: %w", err)
	}
	return count, nil
}
func (r *NotificationRepo) AddAuthor(ctx context.Context, uid string, groupID string, authorName string) error {
	const checkQuery = `
		SELECT COUNT(*) 
		FROM crypto_notification_group
		WHERE crypto_notification_group.crypto_user_id = $1 AND crypto_notification_group.group_id = $2;
	`
	var count int
	if err := r.db.GetContext(ctx, &count, checkQuery, uid, groupID); err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("group not found or does not belong to the user")
	}
	query := `
		INSERT INTO crypto_notification_authors_list (group_id, authors_username)
		VALUES ($1, $2);
	`
	_, err := r.db.ExecContext(ctx, query, groupID, authorName)
	if err != nil {
		return fmt.Errorf("failed to add author: %w", err)
	}
	return nil
}
func (r *NotificationRepo) RemoveAuthor(ctx context.Context, uid, groupAuthorID string) error {
	const checkQuery = `
		SELECT COUNT(*) 
		FROM crypto_notification_authors_list a
		INNER JOIN crypto_notification_group g ON a.group_id = g.group_id
		WHERE g.crypto_user_id = $1 AND a.id = $2;
	`

	var count int
	if err := r.db.GetContext(ctx, &count, checkQuery, uid, groupAuthorID); err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("author not found or does not belong to the user")
	}

	const deleteQuery = `DELETE FROM crypto_notification_authors_list WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, deleteQuery, groupAuthorID); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	return nil
}
func (r *NotificationRepo) UpdateTelegram(ctx context.Context, uid string, chatID string, userID string) error {

	query := `
	UPDATE crypto_user
	SET telegram_chat_id = $1,
	    telegram_user_id = $2
	WHERE uuid = $3
	`
	if _, err := r.db.ExecContext(ctx, query, chatID, userID, uid); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}
func (r *NotificationRepo) GetTelegram(ctx context.Context, uid string) (model.Telegram, error) {
	query := `
		SELECT 
			COALESCE(telegram_chat_id, '') AS telegram_chat_id,
			COALESCE(telegram_user_id, '') AS telegram_user_id
		FROM crypto_user
		WHERE uuid = $1
	`
	var telegram model.Telegram
	if err := r.db.GetContext(ctx, &telegram, query, uid); err != nil {
		fmt.Print(err.Error())
		return model.Telegram{}, fmt.Errorf("failed to get telegram: %w", err)
	}
	return telegram, nil
}
func (r *NotificationRepo) DeleteGroupAuthor(ctx context.Context, uid string, groupID string) error {
	const checkQuery = `
		SELECT COUNT(*) 
		FROM crypto_notification_group
		WHERE crypto_notification_group.crypto_user_id = $1 AND crypto_notification_group.group_id = $2;
	`

	var count int
	if err := r.db.GetContext(ctx, &count, checkQuery, uid, groupID); err != nil {

		return fmt.Errorf("query failed: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("group not found or does not belong to the user")
	}

	query := `
	DELETE FROM crypto_notification_authors_list
	WHERE group_id = $1
	`
	if _, err := r.db.ExecContext(ctx, query, groupID); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}
func (r *NotificationRepo) DisconnectNotification(ctx context.Context, uid string) error {
	query := `
	UPDATE crypto_user
	SET telegram_chat_id = '',
	    telegram_user_id = ''
	WHERE uuid = $1
	`
	if _, err := r.db.ExecContext(ctx, query, uid); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("update failed: %w", err)
	}

	query = `
	DELETE FROM crypto_notification_authors_list
	WHERE group_id IN ( 
		SELECT group_id FROM crypto_notification_group WHERE crypto_user_id = $1
	)
	`
	if _, err := r.db.ExecContext(ctx, query, uid); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("delete failed: %w", err)
	}

	query = `
	DELETE FROM crypto_notification_group
	WHERE crypto_user_id = $1
	`
	if _, err := r.db.ExecContext(ctx, query, uid); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}
