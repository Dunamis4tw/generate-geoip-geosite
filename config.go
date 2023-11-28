package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Path            string   `json:"path" default:"./"`
	GeositeFilename string   `json:"geositeFilename" default:"./geosite.db"`
	GeoipFilename   string   `json:"geoipFilename" default:"./geoip.db"`
	Sources         []Source `json:"sources"`
}

// Source структура с информацией о источнике списка
type Source struct {
	URL                string      `json:"url"`
	Category           string      `json:"category"`
	ContentType        ContentType `json:"contentType"`
	IsExclude          bool        `json:"isExclude"` // false = exclude, true = include. Default: false
	DownloadedFilename string      `json:"downloadedFilename"`
	IpFilename         string      `json:"ipFilename"`
	DomainFilename     string      `json:"domainFilename"`
	// ListType      ListType    `json:"listType"`
	// CheckValidity bool        `json:"checkValidity"` // Если true - необходима проверка ip и доменов на валидность, если false вставлять как есть. Default: false
}

// ContentType перечисление для определения типа обработчика данных
type ContentType string

// // ContentType перечисление для определения типа обработчика данных
// type ListType string

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

// const (
// 	Ip     ListType = "ip"
// 	Domain ListType = "domain"
// 	Mixed  ListType = "mixed"
// )

// func loadConfigsFromJSON(jsonFile string) ([]Config, error) {
// 	data, err := os.ReadFile(jsonFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при чтении JSON-файла: %v", err)
// 	}

// 	var configs []Config
// 	err = json.Unmarshal(data, &configs)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при разборе JSON: %v", err)
// 	}
// 	// fmt.Println(configs[0].ListType)
// 	return configs, nil
// }

// loadConfigsFromJSON читает файл jsonFile и возвращает заполненную содержимым файла структуру Config
func loadConfigsFromJSON(jsonFile string) (*Config, error) {
	// Читаем файл
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	// Переменная для хренения конфига
	var config Config

	// Парсим содержимое файла
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Добавляем слеш только если его нет в конце пути
	if !strings.HasSuffix(config.Path, "/") {
		config.Path += "/"
	}

	return &config, nil
}
