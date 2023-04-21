package main

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
	Port     int
	Timezone string
	User     string
}

func (p Postgres) New() (Postgres, error) {
	var err error
	p.DB, err = gorm.Open(postgres.Open(fmt.Sprintf("dbname=aaa_reconcilliation host=%s user=%s password=%s port=%d sslmode=disable TimeZone=%s", p.Host, p.User, p.Password, p.Port, p.Timezone)), &gorm.Config{})
	if err != nil {
		return p, err
	}
	err = Postgres.AutoMigrate(p, p.Models...)
	if err != nil {
		return p, err
	}
	return p, nil
}
