package database

import (
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	// Configure logger
	newLogger := logger.Default.LogMode(logger.Info)

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Connection pooling and performance settings
		PrepareStmt: true, // Prepare statement for better performance
		Logger:      newLogger,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)           // Maximum number of connections in idle pool
	sqlDB.SetMaxOpenConns(100)           // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // Maximum lifetime of a connection

	return &PostgresDB{DB: db}, nil
}

// Horse operations
func (p *PostgresDB) GetHorse(id int64) (models.Horse, error) {
	var horse models.Horse
	err := p.DB.First(&horse, id).Error
	if err != nil {
		return horse, fmt.Errorf("failed to get horse: %w", err)
	}
	return horse, nil
}

func (p *PostgresDB) GetAllHorses() ([]models.Horse, error) {
	var horses []models.Horse
	err := p.DB.Find(&horses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get horses: %w", err)
	}
	return horses, nil
}

func (p *PostgresDB) GetHorsesByUser(userID string) ([]models.Horse, error) {
	var horses []models.Horse
	err := p.DB.Where("user_id = ?", userID).Find(&horses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user's horses: %w", err)
	}
	return horses, nil
}

// Health record methods
func (p *PostgresDB) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	var records []models.HealthRecord
	err := p.DB.Where("horse_id = ?", horseID).Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get health records: %w", err)
	}
	return records, nil
}

func (p *PostgresDB) AddHealthRecord(record *models.HealthRecord) error {
	err := p.DB.Create(record).Error
	if err != nil {
		return fmt.Errorf("failed to add health record: %w", err)
	}
	return nil
}

// Pregnancy methods
func (p *PostgresDB) GetPregnancyEvents(horseID int64) ([]models.PregnancyEvent, error) {
	var events []models.PregnancyEvent
	err := p.DB.Where("horse_id = ?", horseID).Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancy events: %w", err)
	}
	return events, nil
}

func (p *PostgresDB) AddPregnancyEvent(event *models.PregnancyEvent) error {
	err := p.DB.Create(event).Error
	if err != nil {
		return fmt.Errorf("failed to add pregnancy event: %w", err)
	}
	return nil
}

// Pre-foaling methods
func (p *PostgresDB) GetPreFoalingSigns(horseID int64) ([]models.PreFoalingSign, error) {
	var signs []models.PreFoalingSign
	err := p.DB.Where("horse_id = ?", horseID).Find(&signs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pre-foaling signs: %w", err)
	}
	return signs, nil
}

func (p *PostgresDB) AddPreFoalingSign(sign *models.PreFoalingSign) error {
	err := p.DB.Create(sign).Error
	if err != nil {
		return fmt.Errorf("failed to add pre-foaling sign: %w", err)
	}
	return nil
}

// Pre-foaling checklist methods
func (p *PostgresDB) GetPreFoalingChecklist(horseID int64) ([]models.PreFoalingChecklistItem, error) {
	var items []models.PreFoalingChecklistItem
	err := p.DB.Where("horse_id = ?", horseID).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pre-foaling checklist: %w", err)
	}
	return items, nil
}

func (p *PostgresDB) AddPreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error {
	err := p.DB.Create(item).Error
	if err != nil {
		return fmt.Errorf("failed to add pre-foaling checklist item: %w", err)
	}
	return nil
}

func (p *PostgresDB) UpdatePreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error {
	err := p.DB.Save(item).Error
	if err != nil {
		return fmt.Errorf("failed to update pre-foaling checklist item: %w", err)
	}
	return nil
}

func (p *PostgresDB) DeletePreFoalingChecklistItem(id int64) error {
	err := p.DB.Delete(&models.PreFoalingChecklistItem{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete pre-foaling checklist item: %w", err)
	}
	return nil
}

// Horse operations (missing methods)
func (p *PostgresDB) AddHorse(horse *models.Horse) error {
	err := p.DB.Create(horse).Error
	if err != nil {
		return fmt.Errorf("failed to add horse: %w", err)
	}
	return nil
}

func (p *PostgresDB) UpdateHorse(horse *models.Horse) error {
	err := p.DB.Save(horse).Error
	if err != nil {
		return fmt.Errorf("failed to update horse: %w", err)
	}
	return nil
}

func (p *PostgresDB) DeleteHorse(id int64) error {
	err := p.DB.Delete(&models.Horse{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete horse: %w", err)
	}
	return nil
}

// Health record operations (missing method)
func (p *PostgresDB) UpdateHealthRecord(record *models.HealthRecord) error {
	err := p.DB.Save(record).Error
	if err != nil {
		return fmt.Errorf("failed to update health record: %w", err)
	}
	return nil
}

// Pregnancy operations (missing methods)
func (p *PostgresDB) GetPregnancies(userID string) ([]models.Pregnancy, error) {
	var pregnancies []models.Pregnancy
	err := p.DB.Where("user_id = ?", userID).Find(&pregnancies).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancies: %w", err)
	}
	return pregnancies, nil
}

func (p *PostgresDB) GetPregnancy(id int64) (models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	err := p.DB.First(&pregnancy, id).Error
	if err != nil {
		return pregnancy, fmt.Errorf("failed to get pregnancy: %w", err)
	}
	return pregnancy, nil
}

func (p *PostgresDB) AddPregnancy(pregnancy *models.Pregnancy) error {
	err := p.DB.Create(pregnancy).Error
	if err != nil {
		return fmt.Errorf("failed to add pregnancy: %w", err)
	}
	return nil
}

func (p *PostgresDB) UpdatePregnancy(pregnancy *models.Pregnancy) error {
	err := p.DB.Save(pregnancy).Error
	if err != nil {
		return fmt.Errorf("failed to update pregnancy: %w", err)
	}
	return nil
}

// Expense operations
func (p *PostgresDB) GetExpenses(userID string, from, to time.Time) ([]models.Expense, error) {
	var expenses []models.Expense
	err := p.DB.Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to).Find(&expenses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}
	return expenses, nil
}

func (p *PostgresDB) GetHorseExpenses(horseID int64, from, to time.Time) ([]models.Expense, error) {
	var expenses []models.Expense
	err := p.DB.Where("horse_id = ? AND date BETWEEN ? AND ?", horseID, from, to).Find(&expenses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get horse expenses: %w", err)
	}
	return expenses, nil
}

func (p *PostgresDB) AddExpense(expense *models.Expense) error {
	err := p.DB.Create(expense).Error
	if err != nil {
		return fmt.Errorf("failed to add expense: %w", err)
	}
	return nil
}

func (p *PostgresDB) UpdateExpense(expense *models.Expense) error {
	err := p.DB.Save(expense).Error
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}
	return nil
}

func (p *PostgresDB) DeleteExpense(id int64) error {
	err := p.DB.Delete(&models.Expense{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}
	return nil
}

// Recurring expense operations
func (p *PostgresDB) GetRecurringExpenses(userID string) ([]models.RecurringExpense, error) {
	var expenses []models.RecurringExpense
	err := p.DB.Where("user_id = ?", userID).Find(&expenses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get recurring expenses: %w", err)
	}
	return expenses, nil
}

func (p *PostgresDB) AddRecurringExpense(expense *models.RecurringExpense) error {
	err := p.DB.Create(expense).Error
	if err != nil {
		return fmt.Errorf("failed to add recurring expense: %w", err)
	}
	return nil
}

func (p *PostgresDB) UpdateRecurringExpense(expense *models.RecurringExpense) error {
	err := p.DB.Save(expense).Error
	if err != nil {
		return fmt.Errorf("failed to update recurring expense: %w", err)
	}
	return nil
}

func (p *PostgresDB) DeleteRecurringExpense(id int64) error {
	err := p.DB.Delete(&models.RecurringExpense{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete recurring expense: %w", err)
	}
	return nil
}

// Summary methods
func (p *PostgresDB) GetExpenseSummary(userID string, from, to time.Time) (map[string]float64, error) {
	var expenses []models.Expense
	err := p.DB.Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to).Find(&expenses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get expense summary: %w", err)
	}

	summary := make(map[string]float64)
	for _, exp := range expenses {
		summary[string(exp.Type)] += exp.Amount
	}
	return summary, nil
}

// Add these methods
func (p *PostgresDB) GetBreedingCosts(horseID uint) ([]models.BreedingCost, error) {
	var costs []models.BreedingCost
	err := p.DB.Where("horse_id = ?", horseID).Find(&costs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get breeding costs: %w", err)
	}
	return costs, nil
}

func (p *PostgresDB) AddBreedingCost(cost *models.BreedingCost) error {
	err := p.DB.Create(cost).Error
	if err != nil {
		return fmt.Errorf("failed to add breeding cost: %w", err)
	}
	return nil
}

func (p *PostgresDB) GetOffspring(horseID int64) ([]models.Horse, error) {
	var offspring []models.Horse
	if err := p.DB.Where("mother_id = ? OR father_id = ?", horseID, horseID).Find(&offspring).Error; err != nil {
		return nil, fmt.Errorf("failed to get offspring: %w", err)
	}
	return offspring, nil
}

func (p *PostgresDB) GetDashboardStats(userID string) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	
	// Get total horses
	var totalHorses int64
	if err := p.DB.Model(&models.Horse{}).Where("user_id = ?", userID).Count(&totalHorses).Error; err != nil {
		return nil, fmt.Errorf("failed to count horses: %v", err)
	}
	stats.TotalHorses = int(totalHorses)

	// Get pregnant mares
	var pregnantMares int64
	if err := p.DB.Model(&models.Horse{}).Where("user_id = ? AND is_pregnant = true", userID).Count(&pregnantMares).Error; err != nil {
		return nil, fmt.Errorf("failed to count pregnant mares: %w", err)
	}
	stats.PregnantMares = int(pregnantMares)

	// Get total expenses
	var totalExpenses float64
	if err := p.DB.Model(&models.Expense{}).Where("user_id = ?", userID).Select("COALESCE(SUM(amount), 0)").Scan(&totalExpenses).Error; err != nil {
		return nil, fmt.Errorf("failed to sum expenses: %w", err)
	}
	stats.TotalExpenses = totalExpenses

	// Get upcoming foalings
	var upcomingFoalings int64
	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)
	if err := p.DB.Model(&models.Pregnancy{}).Where("user_id = ? AND status = ? AND end_date BETWEEN ? AND ?", userID, models.PregnancyStatusActive, time.Now(), thirtyDaysFromNow).Count(&upcomingFoalings).Error; err != nil {
		return nil, fmt.Errorf("failed to count upcoming foalings: %w", err)
	}
	stats.UpcomingFoalings = int(upcomingFoalings)

	return stats, nil
}

type FamilyTree struct {
	Horse     models.Horse      `json:"horse"`
	Mother    *models.Horse     `json:"mother,omitempty"`
	Father    *models.Horse     `json:"father,omitempty"`
	Offspring []models.Horse    `json:"offspring,omitempty"`
	Siblings  []models.Horse    `json:"siblings,omitempty"`
}

func (p *PostgresDB) GetFamilyTree(horseID int64) (*FamilyTree, error) {
	var horse models.Horse
	if err := p.First(&horse, horseID).Error; err != nil {
		return nil, fmt.Errorf("failed to get horse: %w", err)
	}

	tree := &FamilyTree{Horse: horse}

	// Get mother if exists
	if horse.MotherId != nil {
		var mother models.Horse
		if err := p.First(&mother, *horse.MotherId).Error; err == nil {
			tree.Mother = &mother
		}
	}

	// Get father if exists
	if horse.FatherId != nil {
		var father models.Horse
		if err := p.First(&father, *horse.FatherId).Error; err == nil {
			tree.Father = &father
		}
	}

	// Get offspring
	if err := p.Where("mother_id = ? OR father_id = ?", horseID, horseID).Find(&tree.Offspring).Error; err != nil {
		return nil, fmt.Errorf("failed to get offspring: %w", err)
	}

	// Get siblings
	if horse.MotherId != nil || horse.FatherId != nil {
		siblingQuery := p.Where("id != ?", horseID)
		if horse.MotherId != nil {
			siblingQuery = siblingQuery.Where("mother_id = ?", *horse.MotherId)
		}
		if horse.FatherId != nil {
			siblingQuery = siblingQuery.Or("father_id = ?", *horse.FatherId)
		}
		if err := siblingQuery.Find(&tree.Siblings).Error; err != nil {
			return nil, fmt.Errorf("failed to get siblings: %w", err)
		}
	}

	return tree, nil
}



func (p *PostgresDB) GetBreedingRecords(horseID int64) ([]models.BreedingRecord, error) {
	var records []models.BreedingRecord
	err := p.DB.Where("mare_id = ? OR stallion_id = ?", horseID, horseID).Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get breeding records: %w", err)
	}
	return records, nil
}

func (p *PostgresDB) AddBreedingRecord(record *models.BreedingRecord) error {
	err := p.DB.Create(record).Error
	if err != nil {
		return fmt.Errorf("failed to add breeding record: %w", err)
	}
	return nil
}

// Add these methods to PostgresDB
func (p *PostgresDB) AutoMigrate(models ...interface{}) error {
	return p.DB.AutoMigrate(models...)
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.DB
}

// Add these methods to delegate to the underlying GORM DB
func (p *PostgresDB) Create(value interface{}) error {
	err := p.DB.Create(value).Error
	if err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}
	return nil
}

func (p *PostgresDB) First(dest interface{}, conds ...interface{}) error {
	err := p.DB.First(dest, conds...).Error
	if err != nil {
		return fmt.Errorf("failed to get first record: %w", err)
	}
	return nil
}

func (p *PostgresDB) Find(dest interface{}, conds ...interface{}) error {
	return p.DB.Find(dest, conds...).Error
}

func (p *PostgresDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return p.DB.Where(query, args...)
}

func (p *PostgresDB) Save(value interface{}) error {
	return p.DB.Save(value).Error
}

func (p *PostgresDB) Delete(value interface{}, conds ...interface{}) error {
	return p.DB.Delete(value, conds...).Error
}

func (p *PostgresDB) Model(value interface{}) *gorm.DB {
	return p.DB.Model(value)
}