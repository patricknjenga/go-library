package library

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	Models []interface{}
	Name   string
}

func (p Postgres) New(r Redis) (Postgres, error) {
	var err error
	p.DB, err = gorm.Open(postgres.Open(fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=disable TimeZone=%s", p.Name, r.GetSecret("DB_HOST"), r.GetSecret("DB_USER"), r.GetSecret("DB_PASS"), r.GetSecret("DB_PORT"), "Africa/Nairobi")), &gorm.Config{})
	if err != nil {
		return p, err
	}
	err = Postgres.AutoMigrate(p, p.Models...)
	if err != nil {
		return p, err
	}
	return p, nil
}
