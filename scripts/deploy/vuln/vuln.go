package vuln

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ulikunitz/xz"
)

var (
	currentScriptDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	targetFolder        = filepath.Join(currentScriptDir, "..", "opensca")
	databasePath        = filepath.Join(targetFolder, "local_cve_database.json")
)

type CVEItems struct {
	CVEItems []interface{} `json:"cve_items"`
}

func fetchCveData(year int, state int) {
	url := fmt.Sprintf("https://github.com/fkie-cad/nvd-json-data-feeds/releases/latest/download/CVE-%d.json.xz", year)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("请求失败，状态码: %d\n", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(fmt.Sprintf("CVE-%d.json-%d.xz", year, state))
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	xzFile, err := os.Open(fmt.Sprintf("CVE-%d.json-%d.xz", year, state))
	if err != nil {
		log.Println(err)
		return
	}
	defer xzFile.Close()

	lzmaReader, err := xz.NewReader(xzFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonFile, err := os.Create(fmt.Sprintf("CVE-%d-%d.json", year, state))
	if err != nil {
		log.Println(err)
		return
	}
	defer jsonFile.Close()
	io.Copy(jsonFile, lzmaReader)
	log.Printf("%d年的 CVE 数据已下载并解压到当前文件夹。\n", year)
}

func mergeData(filePath, databasePath string) {
	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		file, err := os.Create(databasePath)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		json.NewEncoder(file).Encode(CVEItems{CVEItems: []interface{}{}})
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var newData CVEItems
	json.NewDecoder(file).Decode(&newData)

	dbFile, err := os.OpenFile(databasePath, os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer dbFile.Close()

	var database CVEItems
	json.NewDecoder(dbFile).Decode(&database)

	database.CVEItems = append(database.CVEItems, newData.CVEItems...)

	dbFile.Seek(0, 0)
	json.NewEncoder(dbFile).Encode(database)

	os.Remove(filePath)
	log.Printf("%s 中的数据已合并到 %s 中。\n", filePath, databasePath)
}

func checkDataCount(databasePath string) int {
	dbFile, err := os.Open(databasePath)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer dbFile.Close()

	var database CVEItems
	json.NewDecoder(dbFile).Decode(&database)

	numCveItems := len(database.CVEItems)
	log.Printf("当前数据库中有 %d 条数据。\n", numCveItems)
	return numCveItems
}

var last_update_time time.Time

func InitDatabase() {
	last_update_time = time.Now()
	os.MkdirAll(targetFolder, os.ModePerm)

	currentYear := time.Now().Year()

	fetchCveData(currentYear, 0)
	mergeData(fmt.Sprintf("CVE-%d.json", currentYear), databasePath)

	numCveItems := checkDataCount(databasePath)
	for numCveItems < 10000 {
		currentYear--
		fetchCveData(currentYear, 0)
		mergeData(fmt.Sprintf("CVE-%d.json", currentYear), databasePath)
		numCveItems = checkDataCount(databasePath)
	}
}
