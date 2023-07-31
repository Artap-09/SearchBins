package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"searchbin/interfaces"
	"searchbin/logger"
	"searchbin/ranges"
	"searchbin/tree"
	"strconv"
)

func FindBinTree() {
	logger := logger.NewLogger()
	treeBinsRange := tree.NewTree()

	errorsFile := make([][]string, 0)
	errorsFile = append(errorsFile, []string{"EFFECTIVE_DATE", "MEMBER_ID", "LO_RANGE", "HI_RANGE", "BANK_NAME", "PRODUCT_CODE", "COUNTRY_CODE", "BILCURRENCY_CODE"})

	resultFile := make([][]string, 0)

	binsFile, err := os.Open("./bins.csv")
	if err != nil {
		logger.Fatal(err)
	}

	binsReader := csv.NewReader(binsFile)
	record, err := binsReader.Read()
	if err != nil {
		logger.Fatal(err)
	}

	for {
		record, err = binsReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Fatal(err)
		}

		low, err := strconv.ParseUint(record[4], 10, 64)
		if err != nil {
			logger.Error(err)
			continue
		}

		high, err := strconv.ParseUint(record[5], 10, 64)
		if err != nil {
			logger.Error(err)
			continue
		}

		rangeBin, err := ranges.NewRange(record[0], low, high)
		if err != nil {
			logger.Error(err)
			continue
		}
		rangeBin.RU = record[1]
		rangeBin.EN = record[2]

		treeBinsRange.Insert(rangeBin)

	}

	banksNSPKIDFile, err := os.OpenFile("./nspk_2.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	banksNSPKIDReader := csv.NewReader(banksNSPKIDFile)
	recordsAll, err := banksNSPKIDReader.ReadAll()
	if err != nil {
		logger.Fatal(err)
	}

	count := 0
	idsSave := make(map[string]struct{})

	for lineIdx, record := range recordsAll {
		var save bool

		if lineIdx == 0 {
			continue
		}

		bin := record[2]
		var rangeFind interfaces.Range

		for i := 0; i < 13; i++ {
			bin = bin[:19-i]
			if len(bin) != 19 {
				bin += "000000000000000000"
				bin = bin[:19]
			}

			binUint, _ := strconv.ParseUint(bin, 10, 64)
			if err != nil {
				logger.Error(err)
			}

			rangeFind = treeBinsRange.Find(binUint)
			if rangeFind != nil {
				break
			}
		}

		if rangeFind == nil {
			logger.Info("Record:", lineIdx+1, "not found")
		} else {
			count++

			fmt.Println()
			fmt.Println("----------------------------------------------------")
			fmt.Println("Match №", count)
			fmt.Println("Result:", rangeFind.Code(), rangeFind.GetLow(), "-", rangeFind.GetHigh())
			fmt.Println("Search:", record[2])
			fmt.Printf("%s, %s, %s, %s, %s\n", rangeFind.Code(), rangeFind.GetEN(), record[4], record[1], "")
			fmt.Println("----------------------------------------------------")
			save = true
			_, ok := idsSave[record[1]]
			if ok {
				continue
			}

			idsSave[record[1]] = struct{}{}
			result := make([]string, 0, 2)
			result = append(result, rangeFind.Code(), rangeFind.GetEN(), record[4], record[1], "")
			resultFile = append(resultFile, result)
		}

		if !save {
			errorsFile = append(errorsFile, record)
		}

	}

	banksIDResult, err := os.OpenFile("./banks_nspk_by_bin_id.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	banksIDResultWriter := csv.NewWriter(banksIDResult)
	err = banksIDResultWriter.WriteAll(resultFile)
	if err != nil {
		logger.Fatal(err)
	}

	err = os.Truncate(banksNSPKIDFile.Name(), 0)
	if err != nil {
		logger.Fatal(err)
	}

	banksNSPKIDWriter := csv.NewWriter(banksNSPKIDFile)

	err = banksNSPKIDWriter.WriteAll(errorsFile)
	if err != nil {
		logger.Fatal(err)
	}

}

func checkID() {
	idsSave := make(map[string]string)

	errorsFile := make([][]string, 0)
	errorsFile = append(errorsFile, []string{"EFFECTIVE_DATE", "MEMBER_ID", "LO_RANGE", "HI_RANGE", "BANK_NAME", "PRODUCT_CODE", "COUNTRY_CODE", "BILCURRENCY_CODE"})
	logger := logger.NewLogger()

	banksIDResult, err := os.OpenFile("./banks_nspk_by_bin_id.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	banksIDResultReader := csv.NewReader(banksIDResult)

	for {
		record, err := banksIDResultReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Fatal(err)
		}

		bank, ok := idsSave[record[3]]
		if ok {
			logger.Error(bank, record[3], record[2])
		} else {
			idsSave[record[3]] = record[2]
		}

	}

	banksNSPKIDFile, err := os.OpenFile("./nspk_2.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	banksNSPKIDReader := csv.NewReader(banksNSPKIDFile)
	for {
		record, err := banksNSPKIDReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Fatal(err)
		}

		if bank, ok := idsSave[record[1]]; ok {
			logger.Info("Уже есть:", record, "Банк:", bank)
		} else {
			errorsFile = append(errorsFile, record)
		}
	}

	err = os.Truncate(banksNSPKIDFile.Name(), 0)
	if err != nil {
		logger.Fatal(err)
	}

	banksNSPKIDWriter := csv.NewWriter(banksNSPKIDFile)

	err = banksNSPKIDWriter.WriteAll(errorsFile)
	if err != nil {
		logger.Fatal(err)
	}
}
