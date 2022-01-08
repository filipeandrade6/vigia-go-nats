package cmd

import (
	"fmt"
	"os"

	// "path/filepath"

	"github.com/filipeandrade6/vigia-go/internal/gravacao/config"
	"github.com/filipeandrade6/vigia-go/internal/gravacao/core"
	"github.com/filipeandrade6/vigia-go/internal/sys/logger"

	"github.com/spf13/cobra"
)

var cfg config.Configuration

var rootCmd = &cobra.Command{
	Use:   "gravacao",
	Short: "Servico de gravacao",
	Run: func(cmd *cobra.Command, args []string) {
		log, err := logger.New("GRAVACAO")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer log.Sync()

		if err := core.Run(log, cfg); err != nil {
			log.Errorw("startup", "ERROR", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	var err error
	cfg, err = config.ParseConfig("cli")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.Flags().StringVar(&cfg.Database.Host, "db-host", cfg.Database.Host, "host do banco de dados")
	rootCmd.Flags().StringVar(&cfg.Database.User, "db-user", cfg.Database.User, "usuario do banco de dados")
	rootCmd.Flags().StringVar(&cfg.Database.Password, "db-password", cfg.Database.Password, "senha do banco de dados")
	rootCmd.Flags().StringVar(&cfg.Database.Name, "db-name", cfg.Database.Name, "nome do banco de dados")
	rootCmd.Flags().IntVar(&cfg.Database.MaxIDLEConns, "db-maxidleconns", cfg.Database.MaxIDLEConns, "numero maximo de conexoes ociosas")
	rootCmd.Flags().IntVar(&cfg.Database.MaxOpenConns, "db-maxopenconns", cfg.Database.MaxOpenConns, "numero maximo de conexoes abertas")
	rootCmd.Flags().BoolVar(&cfg.Database.DisableTLS, "db-sslmode", cfg.Database.DisableTLS, "desabilitar conexao com SSL ao banco de dados")
}
