package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type HorseRepository struct {
	db *gorm.DB
}

// SecurityAuditLog represents unauthorized access attempts
type SecurityAuditLog struct {
	Timestamp    time.Time
	UserID       string
	ResourceType string
	ResourceID   uint
	Action       string
	IPAddress    string
}

func NewHorseRepository(db *gorm.DB) *HorseRepository {
	return &HorseRepository{db: db}
}

// logSecurityViolation records unauthorized access attempts
func (r *HorseRepository) logSecurityViolation(ctx context.Context, userID string, resourceID uint, action string) {
	// In a real-world scenario, this would be sent to a secure logging system
	log.Printf("SECURITY VIOLATION: User %s attempted unauthorized %s on horse ID %d", 
		userID, action, resourceID)
	
	// Optionally, you could store this in a database table
	auditLog := SecurityAuditLog{
		Timestamp:    time.Now(),
		UserID:       userID,
		ResourceType: "horse",
		ResourceID:   resourceID,
		Action:       action,
	}
	r.db.WithContext(ctx).Create(&auditLog)
}

// validateUserAccess checks if the user has access to the resource
func (r *HorseRepository) validateUserAccess(ctx context.Context, userID string, horseID uint) error {
	var horse models.Horse
	if err := r.db.WithContext(ctx).First(&horse, horseID).Error; err != nil {
		return err
	}

	if horse.UserID != userID {
		r.logSecurityViolation(ctx, userID, horseID, "access")
		return fmt.Errorf("unauthorized access: user %s cannot access horse %d", userID, horseID)
	}
	return nil
}

func (r *HorseRepository) Create(ctx context.Context, horse *models.Horse) error {
	// Ensure the horse is created with the correct user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return fmt.Errorf("no user ID found in context")
	}
	horse.UserID = userID
	return r.db.WithContext(ctx).Create(horse).Error
}

func (r *HorseRepository) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, fmt.Errorf("no user ID found in context")
	}

	if err := r.validateUserAccess(ctx, userID, id); err != nil {
		return nil, err
	}

	var horse models.Horse
	err := r.db.WithContext(ctx).First(&horse, id).Error
	if err != nil {
		return nil, err
	}
	return &horse, nil
}

func (r *HorseRepository) Update(ctx context.Context, horse *models.Horse) error {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return fmt.Errorf("no user ID found in context")
	}

	if err := r.validateUserAccess(ctx, userID, horse.ID); err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(horse).Error
}

func (r *HorseRepository) Delete(ctx context.Context, id uint) error {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return fmt.Errorf("no user ID found in context")
	}

	if err := r.validateUserAccess(ctx, userID, id); err != nil {
		return err
	}

	return r.db.WithContext(ctx).Delete(&models.Horse{}, id).Error
}

func (r *HorseRepository) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&horses).Error
	return horses, err
}

func (r *HorseRepository) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_pregnant = true", userID).Find(&horses).Error
	return horses, err
}

func (r *HorseRepository) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, fmt.Errorf("no user ID found in context")
	}

	if err := r.validateUserAccess(ctx, userID, horseID); err != nil {
		return nil, err
	}

	var horse models.Horse
	if err := r.db.WithContext(ctx).First(&horse, horseID).Error; err != nil {
		return nil, err
	}

	tree := &models.FamilyTree{
		Horse: &horse,
	}

	// Get mother if exists
	if horse.MotherId != nil {
		var mother models.Horse
		if err := r.db.WithContext(ctx).First(&mother, *horse.MotherId).Error; err == nil {
			tree.Mother = &mother
		}
	}

	// Get father if exists
	if horse.FatherId != nil {
		var father models.Horse
		if err := r.db.WithContext(ctx).First(&father, *horse.FatherId).Error; err == nil {
			tree.Father = &father
		}
	}

	// Get offspring
	var offspring []*models.Horse
	if err := r.db.WithContext(ctx).Where("mother_id = ? OR father_id = ?", horseID, horseID).Find(&offspring).Error; err != nil {
		return nil, err
	}
	tree.Offspring = offspring

	return tree, nil
}

func (r *HorseRepository) GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, fmt.Errorf("no user ID found in context")
	}

	if err := r.validateUserAccess(ctx, userID, horseID); err != nil {
		return nil, err
	}

	var offspring []models.Horse
	err := r.db.WithContext(ctx).Where("mother_id = ? OR father_id = ?", horseID, horseID).Find(&offspring).Error
	return offspring, err
}