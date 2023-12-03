package main

import "log"

func main() {

	configs, err := loadConfigsFromJSON("configCustom.json")
	if err == nil {
		log.Print("==== DOWNLOADING ====")
		Downloader(configs)
	}

	log.Print("==== READING FILE LISTS ====")
	// Получаем всю информацию из файла
	fileDataArray := processFiles(configs.Path)

	log.Print("==== GENERATING GEOSITE & GEOIP ====")
	generate(fileDataArray, *configs)

}
