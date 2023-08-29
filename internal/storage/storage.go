package storage

import (
	"auth-service/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type Postgres struct {
	Db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, cfDb *config.ConfigDatabase) (*Postgres, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfDb.User, cfDb.Password, cfDb.Host, cfDb.Port, cfDb.Name)
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connStr)
		if err != nil {
			panic(err)
		}

		pgInstance = &Postgres{db}
	})

	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.Db.Close()
}

//func Initialize(cfDb *config.ConfigDatabase) *gorm.DB {
//	databaseUrl := fmt.Sprintf("pohost=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfDb.Host, cfDb.Port, cfDb.Name, cfDb.Password, cfDb.Name)
//	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)
//	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfDb.Host, cfDb.Port, cfDb.Name, cfDb.Password, cfDb.Name)
//	db, err := gorm.Open(postgres.New(postgres.Config{
//		DSN:                  dsn,
//		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
//	}), &gorm.Config{})
//
//	if err != nil {
//		log.Println(err)
//		panic(err)
//	}
//	fmt.Println("Database connected")
//
//	return db
//}
