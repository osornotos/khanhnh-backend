package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"os"
	"path"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		host := appConfig.Postgres.Host
		ctx := context.Background()
		dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", appConfig.Postgres.User, appConfig.Postgres.Password,
			host, appConfig.Postgres.Database)
		log.Info().Msgf("postgres: host=%s db=%s", appConfig.Postgres.Host, appConfig.Postgres.Database)
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		db := bun.NewDB(sqldb, pgdialect.New())
		defer db.Close()
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Msgf("fail to get current dir %v", err)
		}
		migrationDir := path.Join(pwd, "migrations")
		log.Info().Msgf("migrate from %s", migrationDir)
		mg := migrate.NewMigrations()
		err = mg.Discover(os.DirFS(migrationDir))
		if err != nil {
			log.Fatal().Msgf("fail to discover migrations %v", err)
		}
		migrator := migrate.NewMigrator(db, mg)
		err = migrator.Init(ctx)
		if err != nil {
			log.Fatal().Msgf("fail to init migrator %v", err)
		}
		group, err := migrator.Migrate(ctx)
		if err != nil {
			panic(err)
		}
		if group.ID == 0 {
			log.Info().Msgf("there are no new migrations to run")
			return
		}
		log.Info().Msgf("migrated to %v", group)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
