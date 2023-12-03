package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/sagernet/sing-box/common/geosite"
)

// Rule структура для представления правил в JSON
type Rule struct {
	Domain       []string `json:"domain"`
	DomainSuffix []string `json:"domain_suffix"`
	// DomainKeyword []string `json:"domain_keyword"`
	// DomainRegex   []string `json:"domain_regex"`
	// SourceIPCIDR  []string `json:"source_ip_cidr"`
	IPCIDR []string `json:"ip_cidr"`
}

// RuleSet структура для представления всего JSON файла
type RuleSet struct {
	Version int    `json:"version"`
	Rules   []Rule `json:"rules"`
}

// SaveRuleSetToFile сохраняет набор правил в файл
func SaveRuleSetToFile(ruleSet RuleSet, filename string) error {
	jsonData, err := json.MarshalIndent(ruleSet, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ReadRuleSetFromFile читает набор правил из файла
func ReadRuleSetFromFile(filename string) (RuleSet, error) {
	var ruleSet RuleSet

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return ruleSet, err
	}

	err = json.Unmarshal(fileData, &ruleSet)
	if err != nil {
		return ruleSet, err
	}

	return ruleSet, nil
}

func generate(fileDataArray []FileData, config Config) error {

	// Подготавливаем Writter для записи данных в бинарный формат баз данных MaxMind DB (MMDB).
	mmdb, err := mmdbwriter.New(mmdbwriter.Options{
		// Задаём тип БД (Просто строка, которая видимо нужна СингБоксу)
		DatabaseType: "sing-geoip",
		// Указываем языки (категории в случае с СингБоксом)
		Languages: extractCategories(fileDataArray),
	})
	if err != nil {
		log.Fatalf("ERROR: cannot create new mmdb: %v", err)
	}

	// Пытаемся создать файл geosite.db
	outSites, err := os.Create(config.GeositeFilename)
	if err != nil {
		log.Fatalf("ERROR: cannot create geosite file: %v", err)
	}
	defer outSites.Close()

	// Перебираем файлы
	for _, fileData := range fileDataArray {

		// Перебираем только Include (Пропускаем Exclude)
		if !fileData.IsInclude {
			continue
		}

		// Подготавливаем списки для rule-set
		RuleSetDomain := []string{}
		RuleSetDomainSuffix := []string{}
		// RuleSetDomainKeyword:= []string{}
		// RuleSetDomainRegex:=   []string{}
		// RuleSetSourceIPCIDR:=  []string{}
		RuleSetIPCIDR := []string{}

		// Если файл с IP-адресами
		if fileData.IsIP {

			log.Printf("INFO: Adding IP addresses from the '%s' file...", fileData.Path)
			startTime := time.Now()
			lastIndex := 0

			// Находим исключающий файл с IP-адерсами этой же категории
			ExcludeFileData := findFileData(fileDataArray, false, true, false, fileData.Category)
			ExcludeFileDataRegex := findFileData(fileDataArray, false, true, true, fileData.Category)

			// Добавляем IP из файла include в итоговый массив
			for i, ipStr := range fileData.Content {

				var ipNet *net.IPNet
				// Если в IP адресе есть символ /, то парсим его как адрес записанный методом бесклассовой адресации (CIDR)
				if strings.Contains(ipStr, "/") {
					// Пробуем спарсить, если не получается - пропускаем
					_, ipNet, err = net.ParseCIDR(ipStr)
					if err != nil {
						log.Println(err)
						continue
					}
				} else {
					// Пробуем парсить как обычный IP адрес, если не получается - пропускаем
					addr, err := netip.ParseAddr(ipStr)
					if err != nil {
						log.Println(err)
						continue
					}

					// Если IP адрес нашёлся в списках исключения, то пропускаем его
					if ExcludeFileData != nil && containsString(ipStr, *ExcludeFileData) {
						continue
					}
					// Если IP адрес нашёлся в списках исключения с регулярками, то пропускаем его
					if ExcludeFileDataRegex != nil && containsString(ipStr, *ExcludeFileDataRegex) {
						continue
					}

					// Конвертируем IP адрес в строку байтов (4 или 16 в зависимости от версии IP)
					ipNet = &net.IPNet{
						IP: addr.AsSlice(),
					}
					// Добавляем маску /32 или /128 в зависимости от версии IP
					if addr.Is4() {
						ipNet.Mask = net.CIDRMask(32, 32)
					} else if addr.Is6() {
						ipNet.Mask = net.CIDRMask(128, 128)
					}
				}

				RuleSetIPCIDR = append(RuleSetIPCIDR, ipNet.String())

				// Вставляем полученный IP адрес в указанную категорию в MMDB GeoIP
				if err := mmdb.Insert(ipNet, mmdbtype.String(fileData.Category)); err != nil {
					log.Printf("WARNING: Cannot insert '%s' into mmdb: %v", ipNet, err)
					continue
				}

				// Выводит в консоль информацию о скорости добавления
				if time.Since(startTime).Seconds() > 1 {
					var speed = float64(i-lastIndex) / float64(time.Since(startTime).Seconds())
					var prog = float64(i*100) / float64(len(fileData.Content))
					fmt.Printf("\rGeneration speed: %.2f lines per second (%.2f%% complete)", speed, prog)
					startTime = time.Now()
					lastIndex = i
				}
			}
			if lastIndex != 0 {
				fmt.Println()
			}

			log.Printf("INFO: IP addresses from file '%s' added!", fileData.Path)
		} else { // Если файл с доменами

			log.Printf("INFO: Adding domains from the '%s' file...", fileData.Path)
			// Создаём массив айтемов (доменов) geosite (объект из библиотеки сингбокса)
			var domains []geosite.Item
			// Находим исключающий файл с доменами этой же категории
			ExcludeFileData := findFileData(fileDataArray, false, false, false, fileData.Category)
			ExcludeFileDataRegex := findFileData(fileDataArray, false, false, true, fileData.Category)

			startTime := time.Now()
			lastIndex := 0

			// Добавляем домены из файла include в итоговый массив
			for i, domain := range fileData.Content {

				// Если Домен нашёлся в списках исключения, то пропускаем его
				if ExcludeFileData != nil && containsString(domain, *ExcludeFileData) {
					continue
				}
				// Если Домен нашёлся в списках исключения с регулярками, то пропускаем его
				if ExcludeFileDataRegex != nil && containsString(domain, *ExcludeFileDataRegex) {
					continue
				}

				// Если домен начинается с символа "*" (Например *.domain.com)
				if strings.HasPrefix(domain, "*") {
					// То добавляем строку, убрав * (Получится .domain.com) и задав тип, означающий что эта запись - суффикс (окончание) домена
					// Другими словами, эта запись позволит проксировать все поддомены указанного домена
					domains = append(domains, geosite.Item{
						Type:  geosite.RuleTypeDomainSuffix,
						Value: strings.Replace(domain, "*", "", 1),
					})
					RuleSetDomainSuffix = append(RuleSetDomainSuffix, strings.Replace(domain, "*", "", 1))
					// А также добавляем сам домен без "*" и "." задав тип, означающий что эта запись - домен
					domains = append(domains, geosite.Item{
						Type:  geosite.RuleTypeDomain,
						Value: strings.Replace(domain, "*.", "", 1),
					})
					RuleSetDomain = append(RuleSetDomain, strings.Replace(domain, "*.", "", 1))
				} else {
					// в случае, если нет символа "*", то просто добавляем домен задав тип, означающий что эта запись - домен
					domains = append(domains, geosite.Item{
						Type:  geosite.RuleTypeDomain,
						Value: domain,
					})
					RuleSetDomain = append(RuleSetDomain, domain)
				}

				// Выводит в консоль информацию о скорости добавления
				if time.Since(startTime).Seconds() > 1 {
					var speed = float64(i-lastIndex) / float64(time.Since(startTime).Seconds())
					var prog = float64(i*100) / float64(len(fileData.Content))
					fmt.Printf("\rGeneration speed: %.2f lines per second (%.2f%% complete)", speed, prog)
					startTime = time.Now()
					lastIndex = i
				}

			}
			if lastIndex != 0 {
				fmt.Println()
			}
			log.Printf("INFO: Domains from file '%s' added!", fileData.Path)

			// Сохраняем в файл GeoSite.db полученные домены с указанной категорией
			if err := geosite.Write(outSites, map[string][]geosite.Item{
				fileData.Category: domains,
			}); err != nil {
				log.Println("ERROR: cannot write into geosite file:", err)
			}
		}

		// Создаем rule-set и заполняем его получившимися списками
		ruleSet := RuleSet{
			Version: 1,
			Rules: []Rule{
				{
					Domain:       RuleSetDomain,
					DomainSuffix: RuleSetDomainSuffix,
					// DomainKeyword: []string{},
					// DomainRegex:   []string{},
					// SourceIPCIDR:  []string{},
					IPCIDR: RuleSetIPCIDR,
				},
			},
		}

		// Сохраняем rule-set в файл
		if len(RuleSetIPCIDR) != 0 {
			if err := SaveRuleSetToFile(ruleSet, config.Path+"ip-"+fileData.Category+".json"); err != nil {
				log.Println("ERROR: Error while saving rule-set:", err)
			}
		} else if len(RuleSetDomain) != 0 || len(RuleSetDomainSuffix) != 0 {
			if err := SaveRuleSetToFile(ruleSet, config.Path+"domain-"+fileData.Category+".json"); err != nil {
				log.Println("ERROR: Error while saving rule-set:", err)
			}
		}

	}

	// Пытаемся создать файл geoip.db
	outIPs, err := os.Create(config.GeoipFilename)
	if err != nil {
		log.Println("ERROR: cannot create geoip file:", err)
	}
	defer outIPs.Close()

	// Сохраняем в файл GeoIP.db полученные IP-адреса
	if _, err := mmdb.WriteTo(outIPs); err != nil {
		log.Println("ERROR: cannot write into geoip file:", err)
	}

	return nil
}

// findFileData вызвращает те FileData, у которых параметры равны isInclude, isIP, isRegexp, category
func findFileData(files []FileData, isInclude, isIP, isRegexp bool, category string) *FileData {
	for _, fileData := range files {
		if fileData.IsInclude == isInclude &&
			fileData.IsIP == isIP &&
			fileData.IsRegexp == isRegexp &&
			fileData.Category == category {
			return &fileData
		}
	}
	return nil
}

// containsString проверяет inputStr, есть ли она в файле fileData (Нужно для проверки на исключение)
func containsString(inputStr string, fileData FileData) bool {
	if fileData.IsRegexp {
		// Если файл с регуляркам - проверяем совпадение на регулярку
		for _, regex := range fileData.Regex {
			if regex.MatchString(inputStr) {
				// log.Printf("DEBUG: domain/ip '%s' was excluded by the regular expression '%s' from the file '%s'", inputStr, regex, fileData.Path)
				return true
			}
		}
	} else {
		// Иначе просто ищем совпадение
		for _, content := range fileData.Content {
			if content == inputStr {
				return true
			}
		}
	}
	return false
}

// extractCategories выводит список массив с категориями, прочитанным из папки
func extractCategories(fileDataArray []FileData) []string {
	// Используем map для хранения уникальных значений Category
	categoryMap := make(map[string]bool)

	// Перебираем массив и добавляем каждую категорию в map
	for _, fileData := range fileDataArray {
		categoryMap[fileData.Category] = true
	}

	// Формируем массив уникальных категорий из map
	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}

	return categories
}
