package repository

import (
	"context"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
	"time"
)

type PostgresHorseRepository struct {
	db *gorm.DB
}

type PostgresUserRepository struct {
	db *gorm.DB
}

type PostgresPregnancyRepository struct {
	db *gorm.DB
}

type PostgresHealthRepository struct {
	db *gorm.DB
}

type PostgresBreedingRepository struct {
	db *gorm.DB
}

func NewHorseRepository(db *gorm.DB) HorseRepository {
	return &PostgresHorseRepository{db: db}
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func NewPregnancyRepository(db *gorm.DB) PregnancyRepository {
	return &PostgresPregnancyRepository{db: db}
}

func NewHealthRepository(db *gorm.DB) HealthRepository {
	return &PostgresHealthRepository{db: db}
}

func NewBreedingRepository(db *gorm.DB) BreedingRepository {
	return &PostgresBreedingRepository{db: db}
}

func (r *PostgresHorseRepository) Create(ctx context.Context, horse *models.Horse) error {
	return r.db.WithContext(ctx).Create(horse).Error
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *PostgresPregnancyRepository) AddPreFoaling(ctx context.Context, sign *models.PreFoalingSign) error {
	return r.db.WithContext(ctx).Create(sign).Error
}

func (r *PostgresPregnancyRepository) AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PostgresPregnancyRepository) AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error {
	return r.db.WithContext(ctx).Create(sign).Error
}

func (r *PostgresPregnancyRepository) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *PostgresPregnancyRepository) Create(ctx context.Context, pregnancy *models.Pregnancy) error {
	return r.db.WithContext(ctx).Create(pregnancy).Error
}

func (r *PostgresHealthRepository) CreateRecord(ctx context.Context, record *models.HealthRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *PostgresBreedingRepository) Create(ctx context.Context, cost *models.BreedingCost) error {
	return r.db.WithContext(ctx).Create(cost).Error
}

func (r *PostgresBreedingRepository) CreateRecord(ctx context.Context, record *models.BreedingRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *PostgresHorseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Horse{}, id).Error
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *PostgresHorseRepository) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
    var horse models.Horse
    if err := r.db.WithContext(ctx).First(&horse, id).Error; err != nil {
        return nil, err
    }
    return &horse, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
    var user models.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *PostgresHealthRepository) DeleteRecord(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.HealthRecord{}, id).Error
}

func (r *PostgresHealthRepository) GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
    var records []models.HealthRecord
    if err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&records).Error; err != nil {
        return nil, err
    }
    return records, nil
}

func (r *PostgresBreedingRepository) GetCosts(ctx context.Context, horseID uint) ([]models.BreedingCost, error) {
    var costs []models.BreedingCost
    if err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&costs).Error; err != nil {
        return nil, err
    }
    return costs, nil
}

func (r *PostgresBreedingRepository) GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error) {
    var records []models.BreedingRecord
    if err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&records).Error; err != nil {
        return nil, err
    }
    return records, nil
}

func (r *PostgresHorseRepository) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
    var horse models.Horse
    if err := r.db.WithContext(ctx).First(&horse, horseID).Error; err != nil {
        return nil, err
    }

    var mother *models.Horse
    var father *models.Horse
    if horse.DamID != nil {
        mother = &models.Horse{}
        if err := r.db.WithContext(ctx).First(mother, *horse.DamID).Error; err != nil {
            return nil, err
        }
    }
    if horse.SireID != nil {
        father = &models.Horse{}
        if err := r.db.WithContext(ctx).First(father, *horse.SireID).Error; err != nil {
            return nil, err
        }
    }

    var offspring []*models.Horse
    if err := r.db.WithContext(ctx).Where("dam_id = ? OR sire_id = ?", horseID, horseID).Find(&offspring).Error; err != nil {
        return nil, err
    }

    return &models.FamilyTree{
        Horse:     &horse,
        Mother:    mother,
        Father:    father,
        Offspring: offspring,
    }, nil
}

func (r *PostgresUserRepository) GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error) {
    var stats models.DashboardStats
    
    // Get total number of horses
    if err := r.db.WithContext(ctx).Model(&models.Horse{}).Where("owner_id = ?", userID).Count(&stats.TotalHorses).Error; err != nil {
        return nil, err
    }
    
    // Get number of pregnant mares
    if err := r.db.WithContext(ctx).Model(&models.Horse{}).
        Where("owner_id = ? AND is_pregnant = ?", userID, true).
        Count(&stats.PregnantMares).Error; err != nil {
        return nil, err
    }
    
    // Get total expenses
    var totalExpenses float64
    if err := r.db.WithContext(ctx).Model(&models.BreedingCost{}).
        Joins("JOIN horses ON horses.id = breeding_costs.horse_id").
        Where("horses.owner_id = ?", userID).
        Select("COALESCE(SUM(amount), 0)").
        Scan(&totalExpenses).Error; err != nil {
        return nil, err
    }
    stats.TotalExpenses = totalExpenses
    
    // Get upcoming foalings
    if err := r.db.WithContext(ctx).Model(&models.Horse{}).
        Where("owner_id = ? AND is_pregnant = ? AND conception_date IS NOT NULL", userID, true).
        Count(&stats.UpcomingFoalings).Error; err != nil {
        return nil, err
    }
    
    return &stats, nil
}

func (r *PostgresPregnancyRepository) DeletePreFoalingChecklistItem(ctx context.Context, itemID uint) error {
    return r.db.WithContext(ctx).Delete(&models.PreFoalingChecklistItem{}, itemID).Error
}

func (r *PostgresHealthRepository) UpdateRecord(ctx context.Context, record *models.HealthRecord) error {
    return r.db.WithContext(ctx).Save(record).Error
}

func (r *PostgresHorseRepository) GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error) {
    var offspring []models.Horse
    err := r.db.WithContext(ctx).
        Where("dam_id = ? OR sire_id = ?", horseID, horseID).
        Find(&offspring).Error
    return offspring, err
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
    return r.db.WithContext(ctx).
        Model(&models.User{}).
        Where("id = ?", userID).
        UpdateColumn("last_login", time.Now()).
        Error
}

func (r *PostgresPregnancyRepository) GetActive(ctx context.Context, userID string) ([]models.Pregnancy, error) {
    var pregnancies []models.Pregnancy
    err := r.db.WithContext(ctx).
        Joins("JOIN horses ON horses.id = pregnancies.horse_id").
        Where("horses.owner_id = ? AND pregnancies.status = ?", userID, models.PregnancyStatusActive).
        Find(&pregnancies).Error
    return pregnancies, err
}

func (r *PostgresHorseRepository) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
    var horses []models.Horse
    err := r.db.WithContext(ctx).
        Where("owner_id = ? AND is_pregnant = true", userID).
        Find(&horses).Error
    return horses, err
}

func (r *PostgresPregnancyRepository) GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
    var pregnancy models.Pregnancy
    err := r.db.WithContext(ctx).
        Where("horse_id = ?", horseID).
        First(&pregnancy).Error
    if err != nil {
        return nil, err
    }
    return &pregnancy, nil
}

func (r *PostgresHorseRepository) GetPregnant(ctx context.Context, userID string) ([]models.Horse, error) {
    var horses []models.Horse
    err := r.db.WithContext(ctx).
        Where("owner_id = ? AND is_pregnant = true", userID).
        Find(&horses).Error
    return horses, err
}

func (r *PostgresPregnancyRepository) GetByUserID(ctx context.Context, userID string) ([]models.Pregnancy, error) {
    var pregnancies []models.Pregnancy
    err := r.db.WithContext(ctx).
        Joins("JOIN horses ON horses.id = pregnancies.horse_id").
        Where("horses.owner_id = ?", userID).
        Find(&pregnancies).Error
    return pregnancies, err
}

func (r *PostgresPregnancyRepository) GetCurrentPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
    var pregnancy models.Pregnancy
    err := r.db.WithContext(ctx).
        Where("horse_id = ? AND status = ?", horseID, models.PregnancyStatusActive).
        First(&pregnancy).Error
    if err != nil {
        return nil, err
    }
    return &pregnancy, nil
}

func (r *PostgresHorseRepository) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
    var horses []models.Horse
    err := r.db.WithContext(ctx).
        Where("owner_id = ?", userID).
        Find(&horses).Error
    return horses, err
}

func (r *PostgresPregnancyRepository) GetEvents(ctx context.Context, pregnancyID uint) ([]models.PregnancyEvent, error) {
    var events []models.PregnancyEvent
    err := r.db.WithContext(ctx).
        Where("pregnancy_id = ?", pregnancyID).
        Order("created_at DESC").
        Find(&events).Error
    return events, err
}

func (r *PostgresHorseRepository) Update(ctx context.Context, horse *models.Horse) error {
    return r.db.WithContext(ctx).Save(horse).Error
}

func (r *PostgresPregnancyRepository) GetPreFoaling(ctx context.Context, pregnancyID uint) ([]models.PreFoalingSign, error) {
    var signs []models.PreFoalingSign
    err := r.db.WithContext(ctx).
        Where("pregnancy_id = ?", pregnancyID).
        Order("created_at ASC").
        Find(&signs).Error
    return signs, err
}

func (r *PostgresPregnancyRepository) GetPreFoalingChecklist(ctx context.Context, pregnancyID uint) ([]models.PreFoalingChecklistItem, error) {
    var items []models.PreFoalingChecklistItem
    err := r.db.WithContext(ctx).
        Where("pregnancy_id = ?", pregnancyID).
        Order("created_at ASC").
        Find(&items).Error
    return items, err
}

func (r *PostgresPregnancyRepository) GetPreFoalingChecklistItem(ctx context.Context, itemID uint) (*models.PreFoalingChecklistItem, error) {
    var item models.PreFoalingChecklistItem
    err := r.db.WithContext(ctx).
        Where("id = ?", itemID).
        First(&item).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

func (r *PostgresPregnancyRepository) GetPreFoalingSigns(ctx context.Context, pregnancyID uint) ([]models.PreFoalingSign, error) {
    var signs []models.PreFoalingSign
    err := r.db.WithContext(ctx).
        Where("pregnancy_id = ?", pregnancyID).
        Order("created_at ASC").
        Find(&signs).Error
    return signs, err
}

func (r *PostgresPregnancyRepository) GetPregnancy(ctx context.Context, pregnancyID uint) (*models.Pregnancy, error) {
    var pregnancy models.Pregnancy
    err := r.db.WithContext(ctx).
        Where("id = ?", pregnancyID).
        First(&pregnancy).Error
    if err != nil {
        return nil, err
    }
    return &pregnancy, nil
}

func (r *PostgresPregnancyRepository) Update(ctx context.Context, pregnancy *models.Pregnancy) error {
    return r.db.WithContext(ctx).Save(pregnancy).Error
}

func (r *PostgresPregnancyRepository) UpdatePreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
    return r.db.WithContext(ctx).Save(item).Error
}

func (r *PostgresPregnancyRepository) InitializePreFoalingChecklist(ctx context.Context, pregnancyID uint) error {
    defaultItems := []models.PreFoalingChecklistItem{
        {HorseID: pregnancyID, Description: "Prepare foaling kit", Priority: models.PriorityHigh, DueDate: time.Now(), Notes: "Gather necessary supplies", Completed: false, Season: models.SeasonAll, IsRequired: true, Category: "Preparation"},
        {HorseID: pregnancyID, Description: "Clean and disinfect foaling area", Priority: models.PriorityHigh, DueDate: time.Now(), Notes: "Ensure a clean environment", Completed: false, Season: models.SeasonAll, IsRequired: true, Category: "Preparation"},
        {HorseID: pregnancyID, Description: "Check camera system", Priority: models.PriorityHigh, DueDate: time.Now(), Notes: "Verify camera functionality", Completed: false, Season: models.SeasonAll, IsRequired: true, Category: "Preparation"},
        {HorseID: pregnancyID, Description: "Prepare emergency contacts", Priority: models.PriorityHigh, DueDate: time.Now(), Notes: "Have emergency numbers ready", Completed: false, Season: models.SeasonAll, IsRequired: true, Category: "Preparation"},
        {HorseID: pregnancyID, Description: "Stock up on supplies", Priority: models.PriorityHigh, DueDate: time.Now(), Notes: "Ensure adequate supplies", Completed: false, Season: models.SeasonAll, IsRequired: true, Category: "Preparation"},
        {
            Description: "Check mare's vital signs",
            Priority:    models.PriorityHigh,
            DueDate:     time.Now(),
            Notes:       "Monitor temperature, heart rate, and overall health",
            Completed:   false,
            Season:      models.SeasonAll,
            IsRequired:  true,
            Category:    "Health",
        },
        {
            Description: "Schedule veterinary checkup",
            Priority:    models.PriorityHigh,
            DueDate:     time.Now(),
            Notes:       "Regular checkup to monitor pregnancy progress",
            Completed:   false,
            Season:      models.SeasonAll,
            IsRequired:  true,
            Category:    "Health",
        },
        {
            Description: "Update vaccination records",
            Priority:    models.PriorityHigh,
            DueDate:     time.Now(),
            Notes:       "Ensure all necessary vaccinations are up to date",
            Completed:   false,
            Season:      models.SeasonAll,
            IsRequired:  true,
            Category:    "Health",
        },
        {
            Description: "Review nutrition plan",
            Priority:    models.PriorityMedium,
            DueDate:     time.Now(),
            Notes:       "Adjust feed and supplements as needed",
            Completed:   false,
            Season:      models.SeasonAll,
            IsRequired:  true,
            Category:    "Nutrition",
        },
        {
            Description: "Exercise schedule review",
            Priority:    models.PriorityMedium,
            DueDate:     time.Now(),
            Notes:       "Maintain appropriate exercise routine",
            Completed:   false,
            Season:      models.SeasonAll,
            IsRequired:  true,
            Category:    "Exercise",
        },
    }
    
    return r.db.WithContext(ctx).Create(&defaultItems).Error
}

func (r *PostgresPregnancyRepository) UpdatePregnancyStatus(ctx context.Context, horseID uint, isPregnant bool, conceptionDate *time.Time) error {
    pregnancy := &models.Pregnancy{
        HorseID:        horseID,
        Status:         "ACTIVE",
        ConceptionDate: conceptionDate,
    }
    
    if !isPregnant {
        pregnancy.Status = "INACTIVE"
        now := time.Now()
        pregnancy.EndDate = &now
    }
    
    return r.db.WithContext(ctx).Save(pregnancy).Error
}

func (r *PostgresBreedingRepository) UpdateRecord(ctx context.Context, record *models.BreedingRecord) error {
    return r.db.WithContext(ctx).Save(record).Error
}

func (r *PostgresBreedingRepository) DeleteRecord(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.BreedingRecord{}, id).Error
}