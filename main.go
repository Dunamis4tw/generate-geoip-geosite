package main

import (
	"log"
	"os"
)

var (
	logWarn  = log.New(os.Stdout, "WARN  ", log.LstdFlags)
	logInfo  = log.New(os.Stdout, "INFO  ", log.LstdFlags)
	logError = log.New(os.Stderr, "ERROR ", log.LstdFlags|log.Lshortfile)
)

func main() {
	// Используем функцию из cmdlineparser для парсинга параметров
	options := ParseCommandLine()

	// Если указан флаг для вывода справки, выводим справку и завершаем программу
	if options.ShowHelp {
		PrintHelp()
		os.Exit(0)
	}

	// Валидируем параметры
	err := ValidateOptions(options)
	if err != nil {
		PrintHelp()
		logError.Fatal(err)
	}

	// options.ConfigPath = "./configCustom.json" // For DEBUG

	logInfo.Print("==== READING CONFIG FILE ====")
	configs, err := loadConfigsFromJSON(options.ConfigPath)
	if err != nil {
		logError.Fatal(err)
	}

	logInfo.Print("==== DOWNLOADING ====")
	if err := Downloader(configs); err != nil {
		logError.Fatal(err)
	}

	logInfo.Print("==== READING FILE LISTS ====")
	fileDataArray, err := processFiles(configs.Path)
	if err != nil {
		logError.Fatal(err)
	}

	logInfo.Print("==== GENERATING GEOSITE & GEOIP ====")
	err = generate(fileDataArray, *configs)
	if err != nil {
		logError.Fatal(err)
	}
}
