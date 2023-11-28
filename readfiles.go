package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileData структура для хранения информации о файле
type FileData struct {
	Path      string           // полный путь к файлу
	IsInclude bool             // true, если файл "include" и false, если файл "exclude"
	IsIP      bool             // true, если файл c IP-адресами и false, если файл с доменами
	IsRegexp  bool             // true, если файл с регулярными выражениями
	Category  string           // категория файла
	Content   []string         // содержимое файла
	Regex     []*regexp.Regexp // содержимое файла
	// ExcludeData []string // содержимое файла exclude с регулярными выражениями
}

func processFiles(folderPath string) []FileData {
	// Получаем список файлов в папке
	files, err := getFilesInFolder(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем массив структур FileData
	fileDataArray := make([]FileData, 0)

	// Обрабатываем каждый файл и заполняем массив структур
	for _, file := range files {
		log.Printf("INFO: Reading the file '%s'...\n", file)
		fileData, err := getFileInfo(file)
		if err != nil {
			log.Printf("WARNING: the file '%s' does not match the format: %v\n", file, err)
			continue
		}

		fileDataArray = append(fileDataArray, *fileData)
	}

	return fileDataArray
}

// getFilesInFolder возвращает список .lst и .rgx файлов в заданной папке
func getFilesInFolder(folderPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// getFileInfo по названию файла определяет параметры файла и читает его, возвращает структуру с данными и содержимым
func getFileInfo(filePath string) (*FileData, error) {
	// Получаем имя файла без пути к нему и убираем расширение
	fileName := filepath.Base(filePath)
	fileExtension := filepath.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExtension)

	// Проверяем расширение на .lst и .rgx
	if fileExtension != ".lst" && fileExtension != ".rgx" {
		return nil, errors.New("'" + fileExtension + "' is invalid extension, expected '.lst' or '.rgx'")
	}

	// Разделяем имя файла получая 3 значения
	parts := strings.Split(fileNameWithoutExt, "-")
	if len(parts) < 3 {
		return nil, errors.New("expected at least 3 values in the file name: include/exclude, ip/domain, category_name")
	}

	// Определяем файл типа include или exclude
	include, err := checkIncludeExclude(parts[0])
	if err != nil {
		return nil, err
	}

	// Определяем файл с IP-адресами или доменами
	ip, err := checkIpDomain(parts[1])
	if err != nil {
		return nil, err
	}
	// Считываем категорию
	category := parts[2]

	if fileExtension == ".rgx" {
		// Если файл с регулярками, получаем массив скомпилированных регулярных выражений
		content, err := readRegexFile(filePath)
		if err != nil {
			return nil, err
		}
		return &FileData{
			Path:      filePath,
			IsInclude: include,
			IsIP:      ip,
			IsRegexp:  true,
			Category:  category,
			Regex:     content,
		}, nil
	} else {
		// Если обычный, получаем массив строк
		content, err := readFile(filePath)
		if err != nil {
			return nil, err
		}
		return &FileData{
			Path:      filePath,
			IsInclude: include,
			IsIP:      ip,
			IsRegexp:  false,
			Category:  category,
			Content:   content,
		}, nil
	}
}

// checkIncludeExclude проверяет входную строку на include и exclude
func checkIncludeExclude(input string) (bool, error) {
	lowerInput := strings.ToLower(input)

	switch lowerInput {
	case "include":
		return true, nil
	case "exclude":
		return false, nil
	default:
		return false, errors.New("invalid value, expected 'include' or 'exclude'")
	}
}

// checkIpDomain проверяет входную строку на ip и domain
func checkIpDomain(input string) (bool, error) {
	lowerInput := strings.ToLower(input)

	switch lowerInput {
	case "ip":
		return true, nil
	case "domain":
		return false, nil
	default:
		return false, errors.New("invalid value, expected 'ip' or 'domain'")
	}
}

// readFile читает строки из файла
func readFile(filePath string) ([]string, error) {
	// Читаем файл filePath
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Переменная со строками для результата
	var content []string

	// Запускаем чтение файла построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Пропускаем комментарии
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Пропускаем пустые строки
		if len(line) == 0 {
			continue
		}
		// Добавляем в результирующую переменную
		content = append(content, line)

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return content, nil
}

// readFile читает строки из файла
func readRegexFile(filePath string) ([]*regexp.Regexp, error) {
	// Читаем файл filePath
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Переменная со строками для результата
	var regex []*regexp.Regexp

	// Запускаем чтение файла построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Пропускаем комментарии
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Пропускаем пустые строки
		if len(line) == 0 {
			continue
		}

		// Компилируем полученную регулярку
		rx, err := regexp.Compile(line)
		if err != nil {
			log.Println(err)
			continue
		}

		// Если удачно, добавляем регулярку в исключающий массив
		regex = append(regex, rx)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return regex, nil
}
