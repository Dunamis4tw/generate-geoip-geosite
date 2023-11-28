package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func Downloader(configs *Config) error {

	// Проверяем наличие папки
	if _, err := os.Stat(configs.Path); os.IsNotExist(err) {
		// Если папки нет, создаем её
		err := os.MkdirAll(configs.Path, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return err
		}
		fmt.Println("Folder successfully created:", configs.Path)
	}

	for _, source := range configs.Sources {

		var StartFilename string
		if source.IsExclude {
			StartFilename = "exclude"
		} else {
			StartFilename = "include"
		}

		if len(source.IpFilename) == 0 {
			source.IpFilename = configs.Path + StartFilename + "-ip-" + source.Category + ".lst"
		}
		if len(source.DomainFilename) == 0 {
			source.DomainFilename = configs.Path + StartFilename + "-domain-" + source.Category + ".lst"
		}
		if len(source.DownloadedFilename) == 0 {
			source.DownloadedFilename = configs.Path + uuid.New().String() + ".tmp"
		}

		log.Printf("INFO: Downloading the file '%s' to '%s'...", source.URL, source.DownloadedFilename)
		err := downloadFile(source.URL, source.DownloadedFilename)
		if err != nil {
			return fmt.Errorf("error downloading file: %v", err)
		}

		data, err := os.ReadFile(source.DownloadedFilename)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		parserFunc, ok := parsers[source.ContentType]
		if !ok {
			parserFunc = parseDefaultList
			return fmt.Errorf("invalid data handler type: %s", source.ContentType)
		}

		log.Printf("INFO: Parsing the file '%s'...", source.DownloadedFilename)
		ipAddresses, domains := parserFunc(string(data))

		if len(ipAddresses) != 0 {
			err = writeToFile(ipAddresses, source.IpFilename)
			if err != nil {
				return fmt.Errorf("error writing IP addresses to file: %v", err)
			}
			log.Printf("INFO: Parsed IP addresses are written in '%s'", source.IpFilename)
		}

		if len(domains) != 0 {
			err = writeToFile(domains, source.DomainFilename)
			if err != nil {
				return fmt.Errorf("error writing domains to file: %v", err)
			}
			log.Printf("INFO: Parsed domains are written in '%s'", source.DomainFilename)
		}

		err = os.Remove(source.DownloadedFilename)
		if err != nil {
			return fmt.Errorf("error removing file: %v", err)
		}
	}
	return nil
}

func downloadFile(url, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func parseJsonListDomains(jsonData string) ([]string, []string) {
	var domains []string
	err := json.Unmarshal([]byte(jsonData), &domains)
	if err != nil {
		log.Printf("WARNING: %v", err)
		return nil, nil
	}
	return nil, domains
}

func parseJsonListIPs(jsonData string) ([]string, []string) {
	var ips []string
	err := json.Unmarshal([]byte(jsonData), &ips)
	if err != nil {
		log.Printf("WARNING: %v", err)
		return nil, nil
	}
	return ips, nil
}

func parseJsonRublacklistDPI(jsonData string) ([]string, []string) {
	type Data struct {
		Domains     []string `json:"domains"`
		Name        string   `json:"name"`
		Restriction struct {
			Code string `json:"code"`
		} `json:"restriction"`
	}

	var data []Data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Printf("WARNING: %v", err)
		return nil, nil
	}

	var domains []string
	for _, item := range data {
		domains = append(domains, item.Domains...)
	}

	return nil, domains
}

func parseCsvDumpAntizapret(input string) ([]string, []string) {
	var ipAddresses []string
	var domains []string

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// Разделяем строку на части по символу ";"
		parts := strings.Split(line, ";")
		if len(parts) == 1 {
			continue
		}

		ips := strings.Split(parts[0], "|")
		// Извлекаем IP-адреса из первой части
		// ipMatches := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}|[0-9a-fA-F:]+`).FindAllString(parts[0], -1)
		ipAddresses = append(ipAddresses, ips...)

		// Если есть вторая часть, извлекаем домены из нее
		if len(parts) > 1 {
			domainMatches := strings.Split(parts[1], "|")
			// domainMatches := regexp.MustCompile(`\*?[^\s;|]+`).FindAllString(parts[1], -1)
			domains = append(domains, domainMatches...)
		}
	}

	// Убираем дубликаты
	ipAddresses = uniqueSlice(ipAddresses)
	domains = uniqueSlice(domains)

	return ipAddresses, domains
}

// func isDomain(s string) bool {
// 	// Регулярное выражение для проверки домена
// 	regex := `^(([a-zA-Z0-9А-яёЁ\*]|[a-zA-Z0-9А-яёЁ][a-zA-Z0-9А-яёЁ\-]*[a-zA-Z0-9А-яёЁ])\.)*([A-Za-z0-9А-яёЁ]|[A-Za-z0-9А-яёЁ][A-Za-z0-9А-яёЁ\-]*[A-Za-z0-9А-яёЁ])$`

// 	match, _ := regexp.MatchString(regex, s)
// 	return match
// }

func parseDefaultList(input string) ([]string, []string) {
	var ipAddresses []string
	var domains []string

	lines := strings.Split(input, "\n")
	// startTime := time.Now()
	// lastIndex := 0

	var rgxIPv4 = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(/(3[0-2]|2[0-9]|1[0-9]|[0-9]))?$`)
	var rgxIPv6 = regexp.MustCompile(`^s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(%.+)?s*(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?$`)
	var rgxDomain = regexp.MustCompile(`^(([a-zA-Z0-9А-яёЁ\*]|[a-zA-Z0-9А-яёЁ][a-zA-Z0-9А-яёЁ\-]*[a-zA-Z0-9А-яёЁ])\.)*([A-Za-z0-9А-яёЁ]|[A-Za-z0-9А-яёЁ][A-Za-z0-9А-яёЁ\-]*[A-Za-z0-9А-яёЁ])$`)

	for _, line := range lines {

		// Извлекаем IP-адреса из первой части
		if rgxIPv4.MatchString(line) {
			ipAddresses = append(ipAddresses, line)
			continue
		}

		// Извлекаем IP-адреса из первой части
		if rgxIPv6.MatchString(line) {
			ipAddresses = append(ipAddresses, line)
			continue
		}

		// Извлекаем IP-адреса из первой части
		if rgxDomain.MatchString(line) {
			domains = append(domains, line)
			continue
		}

		// // Извлекаем IP-адреса из первой части
		// ipMatches := rgxIPv4.FindAllString(line, -1)
		// if len(ipMatches) == 0 {
		// 	ipMatches = rgxIPv6.FindAllString(line, -1)
		// }
		// if len(ipMatches) != 0 {
		// 	ipAddresses = append(ipAddresses, ipMatches...)
		// } else {

		// 	domainMatches := rgxDomain.FindAllString(line, -1)
		// 	if len(domainMatches) != 0 {
		// 		domains = append(domains, domainMatches...)
		// 	}
		// }

		// if time.Since(startTime).Seconds() > 1 {
		// 	var speed = float64(i-lastIndex) / float64(time.Since(startTime).Seconds())
		// 	var prog = float64(i*100) / float64(len(lines))
		// 	fmt.Printf("\rParse speed: %.2f lines in second (%.2f%%)\n", speed, prog)
		// 	startTime = time.Now()
		// 	lastIndex = i
		// }
	}
	// fmt.Println()

	// Убираем дубликаты
	ipAddresses = uniqueSlice(ipAddresses)
	domains = uniqueSlice(domains)

	return ipAddresses, domains
}

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
