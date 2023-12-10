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

	// Добавляем параметры в основной конфиг файл
	var config = Config{
		InputDir:   options.InputDir,
		OutputDir:  options.OutputDir,
		SourceFile: options.SourceFile,
		Generate:   options.Generate,
		Sources:    []Source{},
	}

	// Если указан файл с источникам, то
	if len(config.SourceFile) != 0 {
		// Читаем источники
		logInfo.Print("==== READING SOURCE FILE ====")
		Sources, err := loadSourcesFromJSON(config.SourceFile)
		if err != nil {
			logError.Fatal(err)
		}
		config.Sources = Sources
		// Скачиваем их
		logInfo.Print("==== DOWNLOADING ====")
		if err := Downloader(&config); err != nil {
			logError.Fatal(err)
		}
	}

	// Читаем скачанные файлы + те, которые уже были
	logInfo.Print("==== READING FILE LISTS ====")
	fileDataArray, err := processFiles(config.InputDir)
	if err != nil {
		logError.Fatal(err)
	}

	// Генерируем итоговые файлы (GeoIP, Geosite, Rule-Set)
	logInfo.Print("==== GENERATING GEOSITE & GEOIP & RULE-SET ====")
	err = generate(fileDataArray, config)
	if err != nil {
		logError.Fatal(err)
	}
}
