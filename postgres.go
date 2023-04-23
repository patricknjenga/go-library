package library

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	Host     string
	Models   []interface{}
	Name     string
	Password string
	Port     string
	Timezone string
	User     string
}

func (p Postgres) New(r Redis) (Postgres, error) {
	var err error
	p.DB, err = gorm.Open(postgres.Open(fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=disable TimeZone=%s", p.Name, p.Host, p.User, r.GetSecret(p.Password), p.Port, p.Timezone)), &gorm.Config{})
	if err != nil {
		return p, err
	}
	err = Postgres.AutoMigrate(p, p.Models...)
	if err != nil {
		return p, err
	}
	return p, nil
}
