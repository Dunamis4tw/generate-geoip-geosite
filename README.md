# GeoIP, Geosite and Rule-Set generator for Sing-Box

Generates GeoIP, Geosite and Rule-Set files (used by Sing-Box to configure routes) from lists of IP addresses and domains.

<!--
# Генератор Geoip и Geosite для Sing-Box

Генерирует файлы GeoIP и Geosite (используются Sing-Box'ом для настройки маршрутов) из списков IP-адресов и доменов.
-->

## Program Features

- **List Grouping into Categories:** Enables the grouping of lists into categories for more refined configuration of routes in Sing-Box.

- **Downloading Up-to-date Lists:** Supports the capability to download up-to-date lists from publicly available sources.

- **Various Downloadable List Formats Support:** The program can handle lists in formats such as csv, json, list, and others.

- **Exclusion of Unnecessary Domains and IP Addresses:** It provides the option to use exclusion lists to eliminate redundant domains and IP addresses.

- **Rule-Set Generation:** The program is capable of generating new Rule-Sets in both .json and .srs formats, replacing GeoSite and GeoIP in Sing-Box starting from version v1.8.0-alpha.

<!-- 
## Возможности программы

- **Группировка списков в категории:** Позволяет объединять списки в категории для более тонкой настройки маршрутов в Sing-Box.

- **Загрузка актуальных списков:** Поддерживает возможность загружать актуальные списки из общедоступных источников.

- **Поддержка различных форматов загружаемых списков:** Программа способна обрабатывать списки в форматах csv, json, list и других.

- **Исключение ненужных доменов и IP-адресов:** Предоставляет возможность использовать списки исключений для исключения избыточных доменов и IP-адресов.

- Генерация Rule-Set: программа способна генерировать новые Rule-Set в форматах .json и .srs (пришли на замену Geosite и GeoIP в Sing-Box начиная с версии v1.8.0-альфа).
-->

## Usage

To use this program, you can run it from the command line with the following options:

```bash
generate-geoip-geosite -c /path/to/config.json
```

## Build

Given that the Go Language compiler (version 1.11 or greater is required) is installed, you can build it with:

```bash
go get github.com/Dunamis4tw/generate-geoip-geosite
cd $GOPATH/src/github.com/Dunamis4tw/generate-geoip-geosite
go build .
```

## Configuration File Description

The configuration file (`config.json`) for the application is designed to manage various sources of data for building and updating GeoIP and Geosite files. Below is a detailed description of each field in the JSON configuration file.

### Config Fields

1. **path** (string, mandatory)
   - *Description*: The base path where the application will store its files.
   - *Default Value*: "./"

2. **geositeFilename** (string, optional)
   - *Description*: Output Geosite file name.
   - *Default Value*: "./geosite.db"

3. **geoipFilename** (string, optional)
   - *Description*: Output GeoIP file name.
   - *Default Value*: "./geoip.db"

4. **sources** (array of Source, optional)
   - *Description*: An array of sources containing information about the data to be fetched and processed. If not specified, existing files at "path" will be processed.

### Source Fields

Each source within the "Sources" array is defined by the following fields:

1. **url** (string, mandatory)
   - *Description*: The URL from which the data will be fetched.
   - *Example*: "<https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv>"

2. **category** (string, mandatory)
   - *Description*: The name of the category associated with the source. Used in Sing-Box routes.
   - *Example*: "antizapret"

3. **contentType** (string, optional)
   - *Description*: The type of content in the source file. Must match one of the predefined content types.
   - *Options*:
     - "DefaultList": Ordinary list with either IP addresses or domains.
     - "CsvDumpAntizapret": CSV file from Antizapret containing domains and IP addresses.
     - "JsonRublacklistDPI": JSON file from Rublacklist with domains blocked via DPI.
     - "JsonListDomains": JSON file with a list of domains.
     - "JsonListIPs": JSON file with a list of IP addresses.
   - *Example*: "CsvDumpAntizapret"
   - *Default Value*: "DefaultList"

4. **isExclude** (bool, optional)
   - *Description*: If set to true, the data from this source will be included; otherwise, it will be excluded.
   - *Default Value*: false

5. **downloadedFilename** (string, optional)
   - *Description*: The filename used to save temporary downloaded data.
   - *Default Value*: "{random_UUID}.tmp"

6. **ipFilename** (string, optional)
   - *Description*: The filename for storing IP-related data.
   - *Default Value*: "{include/exclude}-ip-{category_name}.lst"

7. **domainFilename** (string, optional)
   - *Description*: The filename for storing domain-related data.
   - *Default Value*: "{include/exclude}-domain-{category_name}.lst"

## Config Examples

The project contains example configuration files:

- **`configAntifilter.json`**
  - Downloads lists of IP addresses and domains provided by [Antifilter](https://antifilter.download/), then categorizes each list.

- **`configRublacklist.json`**
  - Downloads lists of IP addresses and domains provided by [Roskomsvoboda](https://reestr.rublacklist.net/ru/article/api/), then categorizes each list.

- **`configAntizapret.json`**
  - Downloads lists of IP addresses and domains provided by [zapret-info/z-i](https://github.com/zapret-info/z-i). Unnecessary domains are then excluded using regular expressions from the file `antizapret\exclude-domain-antizapret.rgx` (Slightly modified [exclude-regexp-dist.awk](https://bitbucket.org/anticensority/antizapret-pac-generator-light/src/master/config/exclude-regexp-dist.awk)). The result is a list of IP addresses and domains roughly corresponding to the Antizapret lists.

- **`configCustom.json`**
  - A configuration file with custom lists of IP addresses and domains. You only need to specify the path to the directory where your lists are stored in files named "{include/exclude}-{ip/domain}-{category_name}.{lst/rgx}".

<!--
## Примеры конфигурации

Проект содержит примеры файлов конфигурации:

- **`configAntifilter.json`**
  - Скачивает списки IP-адресов и Доменов, предоставленные [Antifilter](https://antifilter.download/), разбивает каждый из списков на категории.

- **`configRublacklist.json`**
  - Скачивает списки IP-адресов и Доменов, предоставленные [Roskomsvoboda](https://reestr.rublacklist.net/ru/article/api/), разбивает каждый из списков на категории.

- **`configAntizapret.json`**
  - Скачивает списки IP-адресов и Доменов, предоставленные [zapret-info/z-i](https://github.com/zapret-info/z-i). Затем из них исключаются ненужные домены регулярными выражениями из файла `antizapret\exclude-domain-antizapret.rgx` (Немного изменённый [exclude-regexp-dist.awk](https://bitbucket.org/anticensority/antizapret-pac-generator-light/src/master/config/exclude-regexp-dist.awk)). В итоге получается список IP-адресов и Доменов, примерно соотвествующий спискам Antizapret.

- **`configCustom.json`**
  - Конфиг файл со своими списками IP-адресов и Доменов. Вы лишь указываете путь до директории, где хранятся ваши списки в файлах с названием "{include/exclude}-{ip/domain}-{category_name}.{lst/rgx}".
-->