package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const VPI_CSV_URL = "https://data.statistik.gv.at/data/OGD_vpi15_VPI_2015_1.csv"

var (
	logger            = log.New(log.Writer(), "", log.LstdFlags)
	errCacheNotExists = fmt.Errorf("cache does not exist")
	errParseRecord    = fmt.Errorf("could not parse record")
)

type VPIRecord struct {
	Year  int
	Month int
	VPI   float64
}

func parseRecords(records [][]string) ([]VPIRecord, error) {
	var parsedRecords []VPIRecord
	for _, record := range records {
		if parsedRecord, err := parseVPIRecord(record); err != nil {
			if err == errParseRecord {
				continue
			} else {
				return nil, err
			}
		} else {
			parsedRecords = append(parsedRecords, parsedRecord)
		}
	}
	return parsedRecords, nil
}

func getRawCSV() ([][]string, error) {
	if records, err := readFromCache(); err != nil {
		if err != errCacheNotExists {
			return nil, err
		} else if records, err := downloadCSV(VPI_CSV_URL); err != nil {
			return nil, err
		} else if err := saveToCache(records); err != nil {
			return nil, err
		} else {
			return records, nil
		}
	} else {
		return records, nil
	}
}

func cachePath() string {
	return fmt.Sprintf("%s/vpi_%d%d.csv", os.TempDir(), time.Now().Year(), time.Now().Month())
}

func readFromCache() ([][]string, error) {
	path := cachePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errCacheNotExists
	} else if file, err := os.Open(path); err != nil {
		return nil, err
	} else {
		logger.Println("Reading from cache")
		defer file.Close()
		return readCSV(file)
	}
}

func saveToCache(records [][]string) error {
	path := cachePath()
	if file, err := os.Create(path); err != nil {
		return err
	} else {
		defer file.Close()
		writer := csv.NewWriter(file)
		writer.Comma = ';'
		writer.WriteAll(records)
		writer.Flush()
		return nil
	}
}

func downloadCSV(url string) ([][]string, error) {
	logger.Println("Downloading CSV file")
	if resp, err := http.Get(url); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		return readCSV(resp.Body)
	}
}

func readCSV(input io.Reader) ([][]string, error) {
	reader := csv.NewReader(input)
	reader.Comma = ';'
	if records, err := reader.ReadAll(); err != nil {
		return nil, err
	} else {
		return records, nil
	}
}

func (r VPIRecord) String() string {
	return fmt.Sprintf("%d-%02d: %f", r.Year, r.Month, r.VPI)
}

func parseVPIRecord(record []string) (VPIRecord, error) {
	periodPattern := regexp.MustCompile(`^VPIZR-(\d{4})(\d{2})$`)
	matches := periodPattern.FindAllStringSubmatch(record[0], -1)
	if len(matches) == 0 {
		return VPIRecord{}, errParseRecord
	} else {
		year, err := strconv.Atoi(matches[0][1])
		if err != nil {
			return VPIRecord{}, err
		}
		month, err := strconv.Atoi(matches[0][2])
		if err != nil {
			return VPIRecord{}, err
		}
		vpi, err := strconv.ParseFloat(strings.ReplaceAll(record[2], ",", "."), 64)
		if err != nil {
			return VPIRecord{}, err
		}

		return VPIRecord{
			Year:  year,
			Month: month,
			VPI:   vpi,
		}, nil
	}
}

func filterForVPI0(records [][]string) [][]string {
	var filtered [][]string
	for _, record := range records {
		if record[1] == "VPI-0" {
			filtered = append(filtered, record)
		}
	}
	return filtered
}
