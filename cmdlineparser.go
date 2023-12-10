package main

import (
	"flag"
	"fmt"
	"os"
)

// CmdLineOptions содержит значения параметров командной строки
type CmdLineOptions struct {
	// ConfigPath string
	InputDir   string
	OutputDir  string
	SourceFile string
	Generate   GenerateOptions
	ShowHelp   bool
}

// GenerateOptions содержит параметры для генерации
type GenerateOptions struct {
	GeoIP       bool
	Geosite     bool
	RuleSetJSON bool
	RuleSetSRS  bool
}

// ParseCommandLine парсит параметры командной строки и возвращает структуру CmdLineOptions
func ParseCommandLine() *CmdLineOptions {
	var options CmdLineOptions

	// flag.StringVar(&options.ConfigPath, "c", "", "set configuration file path")
	// flag.StringVar(&options.ConfigPath, "config", "", "set configuration file path (shorthand)")

	flag.StringVar(&options.SourceFile, "s", "", "set sources.json file path containing an array of URLs for download")
	flag.StringVar(&options.SourceFile, "sources", "", "set sources.json file path containing an array of URLs for download (shorthand)")

	flag.StringVar(&options.InputDir, "i", "", "set input directory path for listing files ({include/exclude}-{ip/domain}-{category_name}.{lst/rgx})")
	flag.StringVar(&options.InputDir, "inputDir", "", "set input directory path for listing files ({include/exclude}-{ip/domain}-{category_name}.{lst/rgx}) (shorthand)")

	flag.StringVar(&options.OutputDir, "o", "", "set output directory path for Geosite, GeoIP, Rule-set files (.db, rule-set.json, rule-set.srs)")
	flag.StringVar(&options.OutputDir, "outputDir", "", "set output directory path for Geosite, GeoIP, Rule-set files (.db, rule-set.json, rule-set.srs) (shorthand)")

	flag.BoolVar(&options.Generate.GeoIP, "gen-geoip", false, "generate GeoIP files")
	flag.BoolVar(&options.Generate.Geosite, "gen-geosite", false, "generate Geosite files")
	flag.BoolVar(&options.Generate.RuleSetJSON, "gen-rule-set-json", false, "generate Rule-set JSON file")
	flag.BoolVar(&options.Generate.RuleSetSRS, "gen-rule-set-srs", false, "generate Rule-set SRS file")

	flag.BoolVar(&options.ShowHelp, "h", false, "help")
	flag.BoolVar(&options.ShowHelp, "help", false, "help (shorthand)")

	flag.Parse()

	// Если указан флаг для вывода справки, выводим справку и завершаем программу
	if options.ShowHelp {
		PrintHelp()
		os.Exit(0)
	}

	// Валидируем параметры
	err := ValidateOptions(&options)
	if err != nil {
		PrintHelp()
		logError.Fatal(err)
	}

	return &options
}

// PrintHelp выводит информацию о параметрах
func PrintHelp() {
	// Выводим информацию о параметрах
	fmt.Println("Generate-geoip-geosite is a tool that generates GeoIP, Geosite, and Rule-Set files (used by Sing-Box to configure routes) from lists of IP addresses and domains.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  generate-geoip-geosite [flags]")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -i, --inputDir string           set input directory path for listing files ({include/exclude}-{ip/domain}-{category_name}.{lst/rgx})")
	fmt.Println("  -o, --outputDir string          set output directory path for Geosite, GeoIP, Rule-set files (.db, rule-set.json, rule-set.srs)")
	fmt.Println("  -s, --sources string            set sources.json file path containing an array of URLs for download")
	fmt.Println("      --gen-geoip                 generate GeoIP file")
	fmt.Println("      --gen-geosite               generate Geosite file")
	fmt.Println("      --gen-rule-set-json         generate Rule-Set JSON files")
	fmt.Println("      --gen-rule-set-srs          generate Rule-Set SRS files")
	fmt.Println("  -h, --help                      help")
}

// ValidateOptions проверяет валидность параметров
func ValidateOptions(options *CmdLineOptions) error {

	// Если не указан InputDir, выдаём ошибку
	if options.InputDir == "" {
		return fmt.Errorf("input directory path is required")
	}

	// Если не указан OutputDir, выдаём ошибку
	if options.OutputDir == "" {
		return fmt.Errorf("output directory path is required")
	}

	// Если не выбран ни один из параметров, выбираем их все
	if !options.Generate.GeoIP && !options.Generate.Geosite && !options.Generate.RuleSetJSON && !options.Generate.RuleSetSRS {
		options.Generate = GenerateOptions{true, true, true, true}
	}

	return nil
}
