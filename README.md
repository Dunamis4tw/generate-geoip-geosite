# Geoip and geosite generator for sing-box

It generates geoip and geosite from lists of domains and IP addresses

## Features

- Supports various input list formats (csv, json, list and others)
- Supports exclusion lists to exclude unnecessary domains and IP addresses
- Supports combining lists into categories to configure routes in sing-box
- Supports downloading current lists from public sources

<!-- ## Getting start -->

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

1. **Path** (string, mandatory)
   - *Description*: The base path where the application will store its files.
   - *Default Value*: "./"

2. **GeositeFilename** (string, mandatory)
   - *Description*: Output Geosite file name.
   - *Default Value*: "./geosite.db"

3. **GeoipFilename** (string, mandatory)
   - *Description*: Output GeoIP file name.
   - *Default Value*: "./geoip.db"

4. **Sources** (array of Source, mandatory)
   - *Description*: An array of sources containing information about the data to be fetched and processed.

### Source Fields

Each source within the "Sources" array is defined by the following fields:

1. **URL** (string, mandatory)
   - *Description*: The URL from which the data will be fetched.
   - *Example*: "https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv"

2. **Category** (string, mandatory)
   - *Description*: The name of the category associated with the source. Used in Sing-Box routes.
   - *Example*: "antizapret"

3. **ContentType** (string, optional)
   - *Description*: The type of content in the source file. Must match one of the predefined content types.
   - *Options*:
     - "DefaultList": Ordinary list with either IP addresses or domains.
     - "CsvDumpAntizapret": CSV file from Antizapret containing domains and IP addresses.
     - "JsonRublacklistDPI": JSON file from Rublacklist with domains blocked via DPI.
     - "JsonListDomains": JSON file with a list of domains.
     - "JsonListIPs": JSON file with a list of IP addresses.
   - *Example*: "CsvDumpAntizapret"
   - *Default Value*: "DefaultList"

4. **IsExclude** (bool, optional)
   - *Description*: If set to true, the data from this source will be included; otherwise, it will be excluded.
   - *Default Value*: false

5. **DownloadedFilename** (string, optional)
   - *Description*: The filename used to save temporary downloaded data.
   - *Default Value*: "{random_UUID}.tmp"

6. **IpFilename** (string, optional)
   - *Description*: The filename for storing IP-related data.
   - *Default Value*: "{include/exclude}-ip-{category_name}.lst"

7. **DomainFilename** (string, optional)
   - *Description*: The filename for storing domain-related data.
   - *Default Value*: "{include/exclude}-domain-{category_name}.lst"

### Parser Functions

The application uses specific parser functions for different content types. These functions are defined in the `parsers` map.

**Note**: Ensure that the configuration file adheres to the specified structure and content type values for successful application execution.

## Config examples
