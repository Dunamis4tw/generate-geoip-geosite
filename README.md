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

## How To Use

```bash
generate-geoip-geosite -i /path/to/input -o /path/to/output
```

To use, you need to specify the path to the directory containing your lists of domains and IP addresses (`"-i ./path/to/input/directory"`) and the path to the directory where the generated files will be saved (`"-o ./path/to/output/directory"`).

In the directory specified by `"-i ./path/to/input/directory"`, files of the following format should be present: `"{include/exclude}-{ip/domain}-{category_name}.{lst/rgx}"`. During the generation process, these files will be processed as follows:

- `{include/exclude}`. IP addresses and domains in a file with "include" in the name will be included in the final file during generation. IP addresses and domains in a file with "exclude" in the name will be excluded from the final file during generation. In other words, all matches of IP addresses and domains in "include" files and "exclude" files of the same category will not be included in the final file.
- `{ip/domain}` in the name indicates that the file contains IP addresses or domains, respectively.
- `{category_name}` is any category name. A category allows combining multiple domains or IP addresses into one list, which appears in the final GeoIP and Geosite files. In the case of Rule-set, each category will create one Rule-set. The same category name can be given for IP addresses and domains, resulting in two different categories (for IP and for domains).
- `{lst/rgx}` is the file extension, indicating the format of the entries in the file: a regular string or a regular expression. Currently, it makes sense to use it for excluding domains.

<!-- 
## Как использовать

Для использования, вам нужно указать путь к каталогу, в котором хранятся ваши списки доменов и ip адресов (`-i ./path/to/input/directory`), и путь к каталогу, куда будут сохраняться генерируемые файлы (`-o ./path/to/output/directory`).
В директории по пути `./path/to/input/directory` должны лежать файлы вида `{include/exclude}-{ip/domain}-{category_name}.{lst/rgx}`, которые во время генерации будут обрабатываться следующим образом:
- IP-адреса и домены в файле с "include" в названии будут включены в итоговый файл во время генерации.
- IP-адреса и домены в файле с "exclude" в названии будут исключены из итогового файла во время генерации. То есть, все совпадения ip-адресов и доменов в "include" файлах и "exclude" файлах одной категории не будут включены в итоговый файл.
- {ip/domain} в названии означает, что файл содержит ip адреса или домены, соотвественно.
- {category_name} - любое имя категории. Категория позволяет объединить несколько доменов или ip-адресов в один список, который фигурирует в итоговых GeoIP и Geosite. В случае с Rule-set, каждая категория создаст один Rule-set. Одно и то же имя категории может быть дано для ip-адресов и для доменов, в итоге всё равно будет две разных категории (для ip и для доменов).
- {lst/rgx} - расширение файла, означает в каком виде представлены записи в файле: обычная строка или регулярное выражение. На данный момент есть смысл использовать для исключения доменов.
-->

To always have up-to-date GeoIP, Geosite, and Rule-set files, functionality for downloading files and subsequent parsing into `{include/exclude}-{ip/domain}-{category_name}.{lst/rgx}` format files has been added to the program. The list of files to download is specified in the source file. To specify the path to the source file, add the `-s ./path/to/source.json` flag. Below are ready-made source files for downloading up-to-date lists of IP addresses and domains useful for users in Russia.

<!--
Для того, чтобы всегда иметь актуальные GeoIP, Geosite и Rule-set'ы, в программу был добавлен функционал скачивания файлов и их последующего парсинга в файлы формата `{include/exclude}-{ip/domain}-{category_name}.{lst/rgx}`. Список файлов, необходимых для скачивания указываются в source-файле. Чтобы указать путь к source-файлу, добавьте флаг `-s ./path/to/source.json`. Ниже представлены уже готовые source-файлы, для скачивания актуальных списков ip-адресов и доменов, полезных для пользователей из РФ.
-->

### Source File Examples

The project contains example source files:

- **`sourceAdAway.json`**
  - Downloads lists of domains provided by [AdAway](https://4pda.to/forum/index.php?showtopic=275091&view=findpost&p=89665467), then categorizes each list.
  - Use: `generate-geoip-geosite -s sourceAdAway.json -i ./adaway -o ./adaway`

- **`sourceAntifilter.json`**
  - Downloads lists of IP addresses and domains provided by [Antifilter](https://antifilter.download/), then categorizes each list.
  - Use: `generate-geoip-geosite -s sourceAntifilter.json -i ./antifilter -o ./antifilter`

- **`sourceRublacklist.json`**
  - Downloads lists of IP addresses and domains provided by [Roskomsvoboda](https://reestr.rublacklist.net/ru/article/api/), then categorizes each list.
  - Use: `generate-geoip-geosite -s sourceRublacklist.json -i ./rublacklist -o ./rublacklist`

- **`sourceAntizapret.json`**
  - Downloads lists of IP addresses and domains provided by [zapret-info/z-i](https://github.com/zapret-info/z-i). Unnecessary domains are then excluded using regular expressions from the file `antizapret\exclude-domain-antizapret.rgx` (Slightly modified [exclude-regexp-dist.awk](https://bitbucket.org/anticensority/antizapret-pac-generator-light/src/master/config/exclude-regexp-dist.awk)). The result is a list of IP addresses and domains roughly corresponding to the Antizapret lists.
  - Use: `generate-geoip-geosite -s sourceAntizapret.json -i ./antizapret -o ./antizapret`

- **`sourceTorrents.json`**
  - Downloads lists domains provided by [github.com/SM443](https://github.com/SM443/Pi-hole-Torrent-Blocklist), then categorizes each list.
  - Use: `generate-geoip-geosite -s sourceTorrents.json -i ./torrents -o ./torrents`

<!--
## Примеры конфигурации

Проект содержит примеры файлов конфигурации:

- **`sourceAdAway.json`**
  - Скачивает списки Доменов, предоставленные [AdAway](https://4pda.to/forum/index.php?showtopic=275091&view=findpost&p=89665467), разбивает каждый из списков на категории.

- **`configAntifilter.json`**
  - Скачивает списки IP-адресов и Доменов, предоставленные [Antifilter](https://antifilter.download/), разбивает каждый из списков на категории.

- **`configRublacklist.json`**
  - Скачивает списки IP-адресов и Доменов, предоставленные [Roskomsvoboda](https://reestr.rublacklist.net/ru/article/api/), разбивает каждый из списков на категории.

- **`configAntizapret.json`**
  - Скачивает списки IP-адресов и Доменов, предоставленные [zapret-info/z-i](https://github.com/zapret-info/z-i). Затем из них исключаются ненужные домены регулярными выражениями из файла `antizapret\exclude-domain-antizapret.rgx` (Немного изменённый [exclude-regexp-dist.awk](https://bitbucket.org/anticensority/antizapret-pac-generator-light/src/master/config/exclude-regexp-dist.awk)). В итоге получается список IP-адресов и Доменов, примерно соотвествующий спискам Antizapret.
-->

## Source File Description

The source file (`source.json`) is intended to provide the program with URLs containing files that include domains and IP addresses, along with parsing information for each file. The JSON file is structured as an array with a "Source" object, defined by the following fields:

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

## Build

Given that the Go Language compiler (version 1.11 or greater is required) is installed, you can build it with:

```bash
go get github.com/Dunamis4tw/generate-geoip-geosite
cd $GOPATH/src/github.com/Dunamis4tw/generate-geoip-geosite
go build .
```

## Flags

- **-i, --inputDir string:** Set the path to the input directory for listing files (`{include/exclude}-{ip/domain}-{category_name}.{lst/rgx}`).
- **-o, --outputDir string:** Set the path to the output directory for GeoIP, Geosite, Rule-set files (`.db`, `rule-set.json`, `rule-set.srs`).
- **-s, --sources string:** Set the path to the `sources.json` file containing an array of URLs for download.
- **--gen-geoip:** Generate GeoIP file.
- **--gen-geosite:** Generate Geosite file.
- **--gen-rule-set-json:** Generate Rule-Set JSON files.
- **--gen-rule-set-srs:** Generate Rule-Set SRS files.
- **-h, --help:** Help.

*Note: If none of the four flags (`--gen-geoip`, `--gen-geosite`, `--gen-rule-set-json`, `--gen-rule-set-srs`) are specified, all four types of final files will be generated. If at least one flag is specified, only the files corresponding to the specified flags will be generated.*

<!--
*Примечаение: Если ни один из четырёх флагов (`--gen-geoip`, `--gen-geosite`, `--gen-rule-set-json`, `--gen-rule-set-srs`) не указан, будут генерироваться все четыре типа финальных файла. Если указан хотя бы один флаг, то будут генерироваться только те файлы, которые были заданы соответствующими флагами.*
-->
