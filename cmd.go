package chassis

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var DBConfig PostgresConfig

// serveCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the http server",
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Debug("Enumerating flags:")
		if cmd.HasFlags() {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Debug(flag.Name, " ", flag.Value)
			})
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// This gets overwritten to initialize the http service
		log.Info("Starting the http server in the chassis")
	},
}

// migrateCmd represents the migrate command
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("migrate called")
		db, err := sql.Open("postgres", cmd.Flags().Lookup("postgres").Value.String())
		if err != nil {
			panic(err)
		}
		defer db.Close()
		level, _ := strconv.Atoi(cmd.Flags().Lookup("target").Value.String())
		err = ApplyMigrationsUp(db, "./migrations/", uint(level))
		if err != nil {
			panic(err)
		}
	},
}

// rootCmd represents the base command when called without any subcommands and should be replaced with your own
var rootCmd = &cobra.Command{
	Use:   "chassis",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("debug").Value.String() == "true" {
			log.SetLevel(log.DebugLevel)
		}
		log.Debug("Setting up the database connection")
		log.Debug("viper thinks our uri should be :" + viper.GetString("postgres"))
		log.Debug("pflag thinks our uri should be :" + cmd.Flag("postgres").Value.String())
		DBConfig = ParsePostgresConfig(viper.GetString("postgres"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootCmd.Use = filepath.Base(path)
	rootCmd.AddCommand(ServeCmd)
	rootCmd.AddCommand(MigrateCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Activate environment variable parsing
	viper.BindEnv("postgres", "POSTGRES_URI")

	// This is an examnple of a global flag.  We don't use a config file.  We prefer ENV variables and flags.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chassis.yaml)")
	rootCmd.PersistentFlags().StringP("postgres", "P", "postgresql://localhost:5432/mydb?user=other&password=secret", "Postgres connection string")
	viper.BindPFlag("postgres", rootCmd.PersistentFlags().Lookup("postgres"))
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Enable debug logging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	ServeCmd.Flags().Int16P("port", "p", 8080, "Port to listen on")
	ServeCmd.Flags().StringP("host", "H", "0.0.0.0", "Host to listen on")
	ServeCmd.Flags().StringP("cert", "c", "", "Path to the cert file")
	ServeCmd.Flags().StringP("key", "k", "", "Path to the key file")
	ServeCmd.Flags().StringP("ca", "a", "", "Path to the ca file")
	MigrateCmd.Flags().StringP("direction", "d", "up", "Direction to migrate")
	MigrateCmd.Flags().IntP("target", "t", 1, "Target migration to migrate to")
}
