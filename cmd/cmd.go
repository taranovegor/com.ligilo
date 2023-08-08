package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/taranovegor/com.ligilo/internal/config"
	"github.com/taranovegor/com.ligilo/internal/container"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"gorm.io/gorm"
	"log"
)

func Init(scope string) container.ServiceContainer {
	fmt.Println(fmt.Sprintf("[cmd/%s] %s! Version: %s", scope, config.AppName, config.Version))

	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	sc, err := container.Init()
	if err != nil {
		panic(err)
	}

	orm := sc.Get(container.Orm).(*gorm.DB)
	err = orm.AutoMigrate(
		&domain.Link{},
	)
	if err != nil {
		panic(err)
	}

	return sc
}
