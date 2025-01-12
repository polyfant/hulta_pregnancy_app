package repository

import "github.com/polyfant/hulta_pregnancy_app/internal/database"

type PostgresHorseRepository struct {
    db *database.PostgresDB
}

type PostgresUserRepository struct {
    db *database.PostgresDB
}

type PostgresPregnancyRepository struct {
    db *database.PostgresDB
}

type PostgresHealthRepository struct {
    db *database.PostgresDB
}

type PostgresBreedingRepository struct {
    db *database.PostgresDB
}

func NewHorseRepository(db *database.PostgresDB) HorseRepository {
    return &PostgresHorseRepository{db: db}
}

func NewUserRepository(db *database.PostgresDB) UserRepository {
    return &PostgresUserRepository{db: db}
}

func NewPregnancyRepository(db *database.PostgresDB) PregnancyRepository {
    return &PostgresPregnancyRepository{db: db}
}

func NewHealthRepository(db *database.PostgresDB) HealthRepository {
    return &PostgresHealthRepository{db: db}
}

func NewBreedingRepository(db *database.PostgresDB) BreedingRepository {
    return &PostgresBreedingRepository{db: db}
} 