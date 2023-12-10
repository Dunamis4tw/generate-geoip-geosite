package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	SourceFile string          // Json-файл со списком ссылок на списки
	Sources    []Source        // Содержимое файла SourceFile
	InputDir   string          // Директория, откуда будут браться списки для генерации (сюда же будут качаться файлы)
	OutputDir  string          // Директория, куда будут складываться сгенерированный файлы
	Generate   GenerateOptions // Массив с выбранными генерируемыми файлами
}

// Source структура с информацией о источнике списка
type Source struct {
	URL         string      `json:"url"`         // Ссылка на скачиваемый файл
	Category    string      `json:"category"`    // Название категории
	ContentType ContentType `json:"contentType"` // Как парсить файл
	IsExclude   bool        `json:"isExclude"`   // Список с исключением или включением. false = exclude, true = include. Default: false
	// DownloadedFilename string      `json:"downloadedFilename"` // Имя временного файла для скачивания
	// IpFilename         string      `json:"ipFilename"`         // Имя распарсенного файла с IP-адресами
	// DomainFilename     string      `json:"domainFilename"`     // Имя распарсенного файла с Доменами
}

// ContentType перечисление для определения типа обработчика данных
type ContentType string

// ParserFunc функция для обработки данных
type ParserFunc func(input string) ([]string, []string)

var parsers = map[ContentType]ParserFunc{
	DefaultList:        parseDefaultList,
	CsvDumpAntizapret:  parseCsvDumpAntizapret,
	JsonRublacklistDPI: parseJsonRublacklistDPI,
	JsonListDomains:    parseJsonListDomains,
	JsonListIPs:        parseJsonListIPs,
}

const (
	DefaultList        ContentType = "DefaultList"        // Обычный список, где в каждой строке указан либо IP-адрес, либо Домен
	CsvDumpAntizapret  ContentType = "CsvDumpAntizapret"  // CSV файл от Антизапрета с доменами и IP-адресами в виде CSV-файла
	JsonRublacklistDPI ContentType = "JsonRublacklistDPI" // JSON файл от Роскомсвободы (Rublacklist) со списком доменов, заблокированных по DPI
	JsonListDomains    ContentType = "JsonListDomains"    // JSON файл со списком доменов (Например: ["dom1.com","dom2.com","dom3.com"])
	JsonListIPs        ContentType = "JsonListIPs"        // JSON файл со списком IP-адресов (Например: ["1.1.1.1","2.2.2.2","3.3.3.3"])
)

func loadSourcesFromJSON(jsonFile string) ([]Source, error) {
	// Читаем файл
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading sources file: %v", err)
	}
	logInfo.Printf("sources file '%s' successfully read", jsonFile)

	// Переменная для хранения массива Source
	var sources []Source

	// Парсим содержимое файла
	err = json.Unmarshal(data, &sources)
	if err != nil {
		return nil, fmt.Errorf("sources file '%s' deserialization error: %v", jsonFile, err)
	}
	logInfo.Printf("sources file '%s' successfully deserialized", jsonFile)

	// Возвращаем массив Source
	return sources, nil
}
