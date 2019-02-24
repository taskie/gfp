package gfp

import (
	"bufio"
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taskie/gfp"
	"github.com/taskie/osplus"
)

type Config struct {
	FromType, ToType, FromPreset, FromFormat, ToFormat, ToPreset, LogLevel string
	Error                                                                  bool
}

var configFile string
var config Config
var (
	verbose, debug, version bool
)

const CommandName = "gfp"

func init() {
	Command.PersistentFlags().StringVarP(&configFile, "config", "c", "", `config file (default "`+CommandName+`.yml")`)
	Command.Flags().StringP("from-type", "f", "", "from type")
	Command.Flags().StringP("to-type", "t", "", "to type")
	Command.Flags().StringP("from-format", "F", "", "from format string")
	Command.Flags().StringP("to-format", "T", "", "to format string")
	Command.Flags().StringP("from-preset", "p", "", "to type")
	Command.Flags().StringP("to-preset", "P", "", "to type")
	Command.Flags().BoolP("error", "e", false, "exit if error")
	Command.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	Command.Flags().BoolVar(&debug, "debug", false, "debug output")
	Command.Flags().BoolVarP(&version, "version", "V", false, "show Version")

	for _, s := range []string{"from-type", "to-type", "from-format", "to-format", "from-preset", "to-preset", "error"} {
		envKey := strcase.ToSnake(s)
		structKey := strcase.ToCamel(s)
		viper.BindPFlag(envKey, Command.Flags().Lookup(s))
		viper.RegisterAlias(structKey, envKey)
	}

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(CommandName)
		conf, err := osplus.GetXdgConfigHome()
		if err != nil {
			log.Info(err)
		} else {
			viper.AddConfigPath(filepath.Join(conf, CommandName))
		}
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix(CommandName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Warn(err)
	}
}

func Main() {
	Command.Execute()
}

var Command = &cobra.Command{
	Use:  CommandName + ` [INPUT] [OUTPUT]`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		err := run(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func run(cmd *cobra.Command, args []string) error {
	if version {
		fmt.Println(gfp.Version)
		return nil
	}
	if config.LogLevel != "" {
		lv, err := log.ParseLevel(config.LogLevel)
		if err != nil {
			log.Warn(err)
		} else {
			log.SetLevel(lv)
		}
	}
	if debug {
		if viper.ConfigFileUsed() != "" {
			log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		}
		log.Debug(pp.Sprint(config))
	}

	input := ""
	output := ""
	switch len(args) {
	case 0:
		break
	case 1:
		input = args[0]
	case 2:
		input = args[0]
		output = args[1]
	default:
		return fmt.Errorf("invalid arguments: %v", args[2:])
	}

	opener := osplus.NewOpener()
	r, err := opener.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, commit, err := opener.CreateTempFileWithDestination(output, "", CommandName+"-")
	if err != nil {
		return err
	}
	defer w.Close()

	sc, err := gfp.NewScanner(&gfp.ScannerConfig{
		Type:   config.FromType,
		Format: config.FromFormat,
		Preset: config.FromPreset,
	})
	if err != nil {
		return err
	}

	pr, err := gfp.NewPrinter(&gfp.PrinterConfig{
		Type:   config.ToType,
		Format: config.ToFormat,
		Preset: config.ToPreset,
	})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		v, err := sc.Scan(scanner.Text())
		if err != nil {
			if config.Error {
				return err
			}
			log.Error(err)
			continue
		}
		s, err := pr.Print(v)
		if err != nil {
			if config.Error {
				return err
			}
			log.Error(err)
			continue
		}
		fmt.Println(s)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	commit(true)
	return nil
}
