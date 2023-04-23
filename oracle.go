package library

import (
	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

type Oracle struct {
	*gorm.DB `gorm:"-"`
	Cluster  string
	Host     string
	Password string
	Port     int
	Service  string
	User     string
}

func (o Oracle) New() (Oracle, error) {
	var err error
	o.DB, err = gorm.Open(oracle.Open(oracle.BuildUrl(o.Host, o.Port, o.Service, o.User, o.Password, nil)), &gorm.Config{})
	if err != nil {
		return o, err
	}
	return o, err
}

func GetActiveStandbyOracleDbs(dbs []Oracle) ([]Oracle, []Oracle, error) {
	var active, standby []Oracle
	for _, v := range dbs {
		v, err := v.New()
		if err != nil {
			return active, standby, err
		}
		var status string
		err = v.DB.Raw("select controlfile_type from v$database").Pluck("controlfile_type", &status).Error
		if err != nil {
			return active, standby, err
		}
		if status == "CURRENT" {
			active = append(active, v)
		} else if status == "STANDBY" {
			standby = append(standby, v)
		}
	}
	return active, standby, nil
}
