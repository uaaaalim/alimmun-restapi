package repositories

import (
	"context"
	"time"
	"uuid"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (repo *SessionRepository) Create(
	ctx context.Context,
	session *models.Session,
) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User
		err := tx.Select("id").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&user, "id = ?", session.UserID).Error
		if err != nil {
			return err
		}

		if err := repo.deleteOldSessions(tx, session.UserID, 4); err != nil {
			return err
		}

		return tx.Create(session).Error
	})
}

func (repo *SessionRepository) deleteOldSessions(
	tx *gorm.DB,
	userID uuid.UUID,
	keep int,
) error {
	oldSessions := tx.Model(&models.Session{}).
		Select("id").
		Where("user_id = ?", userID).
		Order("created_at DESC, id DESC").
		Offset(keep)

	return tx.Where("id IN (?)", oldSessions).
		Delete(&models.Session{}).Error
}

func (repo *SessionRepository) GetByHashedToken(ctx context.Context, hashedToken string) (*models.Session, error) {
	var session models.Session

	err := repo.db.WithContext(ctx).
		Where("token_hash = ?", hashedToken).
		First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (repo *SessionRepository) Revoke(
	ctx context.Context,
	sessionID uuid.UUID,
) error {
	return repo.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Update("revoked_at", time.Now()).Error
}
