package main

import (
	"flag"
	"fmt"
)

// CmdLineOptions содержит значения параметров командной строки
type CmdLineOptions struct {
	ConfigPath string
	// WorkingDirectory string
	ShowHelp bool
}

// ParseCommandLine парсит параметры командной строки и возвращает структуру CmdLineOptions
func ParseCommandLine() *CmdLineOptions {
	var options CmdLineOptions

	flag.StringVar(&options.ConfigPath, "c", "", "set configuration file path")
	flag.StringVar(&options.ConfigPath, "config", "", "set configuration file path (shorthand)")

	// flag.StringVar(&options.WorkingDirectory, "d", "", "set working directory")
	// flag.StringVar(&options.WorkingDirectory, "directory", "", "set working directory (shorthand)")

	flag.BoolVar(&options.ShowHelp, "h", false, "help")
	flag.BoolVar(&options.ShowHelp, "help", false, "help (shorthand)")

	flag.Parse()

	return &options
}

// PrintHelp выводит информацию о параметрах
func PrintHelp() {
	// Выводим информацию о параметрах
	fmt.Println("Usage:")
	fmt.Println("  generate-geoip-geosite [flags]")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -c, --config string             set configuration file path")
	// fmt.Println("  -d, --directory string          set working directory")
	fmt.Println("  -h, --help                      help")
}

// ValidateOptions проверяет валидность параметров
func ValidateOptions(options *CmdLineOptions) error {
	// Здесь вы можете добавить дополнительные проверки, если необходимо

	// Пример: если требуется указать конфигурационный файл
	if options.ConfigPath == "" {
		return fmt.Errorf("configuration file path is required")
	}

	return nil
}
