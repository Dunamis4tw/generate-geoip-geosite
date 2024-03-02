package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func Downloader(configs *Config) error {
	// Проверяем наличие директории InputDir
	if _, err := os.Stat(configs.InputDir); os.IsNotExist(err) {
		// Если её нет, создаем
		err := os.MkdirAll(configs.InputDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory '%s': %v", configs.InputDir, err)
		}
		logInfo.Printf("the directory '%s' was missing, but it was created:", configs.InputDir)
	}

	// Перебираем источники
	for _, source := range configs.Sources {

		// Определяем тип источника (для названия файла)
		StartFilename := "include"
		if source.IsExclude {
			StartFilename = "exclude"
		}
		// Собираем имена файлов
		var IpFilename = configs.InputDir + StartFilename + "-ip-" + source.Category + ".lst"
		var DomainFilename = configs.InputDir + StartFilename + "-domain-" + source.Category + ".lst"

		// Cкачиваем файл
		logInfo.Printf("downloading the file '%s'...", source.URL)
		data, err := downloadURL(source.URL)
		if err != nil {
			return fmt.Errorf("error downloading file: %v", err)
		}

		// Парсим скачанный файл в зависимости от указанного source.ContentType
		logInfo.Printf("parsing the file...")
		parserFunc, ok := parsers[source.ContentType]
		if !ok {
			parserFunc = parseDefaultList
			return fmt.Errorf("invalid data handler type: %s", source.ContentType)
		}
		ipAddresses, domains := parserFunc(string(data))

		// Если были распарсены IP-адреса, то сохраняем их в файл
		if len(ipAddresses) != 0 {
			err = writeToFile(ipAddresses, IpFilename)
			if err != nil {
				return fmt.Errorf("error writing IP addresses to file: %v", err)
			}
			logInfo.Printf("parsed IP addresses are written in '%s'", IpFilename)
		}

		// Если были распарсены Домены, то сохраняем их в файл
		if len(domains) != 0 {
			err = writeToFile(domains, DomainFilename)
			if err != nil {
				return fmt.Errorf("error writing domains to file: %v", err)
			}
			logInfo.Printf("parsed domains are written in '%s'", DomainFilename)
		}

	}
	return nil
}

func downloadURL(url string) ([]byte, error) {
	// Получаем ответ от get запроса на указанный url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Если ответ не 200, выдаём ошибку
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP response error: %s", resp.Status)
	}

	// Читаем ответ в переменную data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseJsonListDomains(jsonData string) ([]string, []string) {
	var domains []string
	err := json.Unmarshal([]byte(jsonData), &domains)
	if err != nil {
		logWarn.Print(err)
		return nil, nil
	}
	return nil, domains
}

func parseJsonListIPs(jsonData string) ([]string, []string) {
	var ips []string
	err := json.Unmarshal([]byte(jsonData), &ips)
	if err != nil {
		logWarn.Print(err)
		return nil, nil
	}
	return ips, nil
}

func parseJsonRublacklistDPI(jsonData string) ([]string, []string) {
	// Создаём стркутуру
	type Data struct {
		Domains     []string `json:"domains"`
		Name        string   `json:"name"`
		Restriction struct {
			Code string `json:"code"`
		} `json:"restriction"`
	}

	// Парсим json в структуру
	var data []Data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		logWarn.Print(err)
		return nil, nil
	}

	// Добавлям домены в общий список
	var domains []string
	for _, item := range data {
		domains = append(domains, item.Domains...)
	}

	return nil, domains
}

func parseCsvDumpAntizapret(input string) ([]string, []string) {
	var ipAddresses []string
	var domains []string

	// Декодируем входную строку из Windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	decodedInput, _ := decoder.String(input)

	lines := strings.Split(decodedInput, "\n")
	for _, line := range lines {
		// Разделяем строку на столбцы по символу ";"
		columns := strings.Split(line, ";")

		// Пропускаем первую строку (в ней один столбец)
		if len(columns) == 1 {
			continue
		}

		// Извлекаем IP-адреса из первого столбца
		ips := strings.Split(columns[0], "|")
		ipAddresses = append(ipAddresses, ips...)

		// Если есть второй столбца, извлекаем домены из нее
		if len(columns) > 1 {
			domainMatches := strings.Split(columns[1], "|")
			domains = append(domains, domainMatches...)
		}
	}

	// Убираем дубликаты
	ipAddresses = uniqueSlice(ipAddresses)
	domains = uniqueSlice(domains)

	return ipAddresses, domains
}

var (
	rgxIPv4   = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(/(3[0-2]|2[0-9]|1[0-9]|[0-9]))?$`)
	rgxIPv6   = regexp.MustCompile(`^s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(%.+)?s*(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?$`)
	rgxDomain = regexp.MustCompile(`^(([a-zA-Z0-9А-яёЁ\*]|[a-zA-Z0-9А-яёЁ][a-zA-Z0-9А-яёЁ\-]*[a-zA-Z0-9А-яёЁ])\.)*([A-Za-z0-9А-яёЁ]|[A-Za-z0-9А-яёЁ][A-Za-z0-9А-яёЁ\-]*[A-Za-z0-9А-яёЁ])$`)
)

func parseDefaultList(input string) ([]string, []string) {
	var ipAddresses []string
	var domains []string

	lines := strings.Split(input, "\n")

	for _, line := range lines {

		// Извлекаем домен
		if rgxDomain.MatchString(line) {
			domains = append(domains, line)
			continue
		}

		// Извлекаем IPv4-адрес
		if rgxIPv4.MatchString(line) {
			ipAddresses = append(ipAddresses, line)
			continue
		}

		// Извлекаем IPv6-адрес
		if rgxIPv6.MatchString(line) {
			ipAddresses = append(ipAddresses, line)
			continue
		}

		// Если строка не была комментарием или пустой строкой, выводим предупреждение, что не удалось распарсить строку
		// (так как такие строки отсекаются регулярками, нет смысла делать проверку перед регулярками)
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			logWarn.Printf("Failed to parse '%s' as an IPv4, IPv6, or domain address", line)
		}
	}

	// Убираем дубликаты
	ipAddresses = uniqueSlice(ipAddresses)
	domains = uniqueSlice(domains)

	return ipAddresses, domains
}

func parseHostsFile(input string) ([]string, []string) {
	var ips []string
	var domains []string

	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()

		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}

		// Разбиваем строку на слова
		words := strings.Fields(line)

		// Если в строке два слова
		if len(words) == 2 {
			var ip = words[0]
			var domain = words[1]

			// Добавляем IP, есом он не адрес вида 127.0.0.1, 0.0.0.0, ::1 и т.п.
			if !isLoopbackIP(ip) {
				ips = append(ips, ip)
			}

			// Добавялем домен, если он не localhost
			if domain != "localhost" {
				domains = append(domains, domain)
			}
		}
	}
	return ips, domains
}

// Проверяем, является ли переданный IP адрес "зацикленным" (loopback)
func isLoopbackIP(ip string) bool {
	loopbackPatterns := []string{"127.", "0.0.0.0", "::1"}
	for _, pattern := range loopbackPatterns {
		if strings.HasPrefix(ip, pattern) {
			return true
		}
	}
	return false
}

// uniqueSlice удаляет дубликаты
func uniqueSlice(slice []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueSlice := make([]string, 0)

	for _, item := range slice {
		if _, found := uniqueMap[item]; !found {
			uniqueMap[item] = true
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice
}

func writeToFile(data []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, item := range data {
		_, err := fmt.Fprintln(writer, item)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}
