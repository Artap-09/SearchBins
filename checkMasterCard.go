package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"searchbin/logger"
	"searchbin/ranges"
	"searchbin/tree"
	"strconv"
)

func checkMaster() {
	logger := logger.NewLogger()
	treeMasterCardRange := tree.NewTree()

	masterCardFile, err := os.OpenFile("./latest.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	masterCardFileReader := csv.NewReader(masterCardFile)
	record, err := masterCardFileReader.Read()
	if err != nil {
		logger.Fatal(err)
	}

	for {
		record, err = masterCardFileReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Fatal(err)
		}

		lowString := record[2]
		if len(lowString) != 19 {
			lowString += "000000000000000000"
			lowString = lowString[:19]
		}

		low, err := strconv.ParseUint(lowString, 10, 64)
		if err != nil {
			logger.Error(err)
			continue
		}

		highString := record[3]
		if len(highString) != 19 {
			highString += "99999999999999999999"
			highString = highString[:19]
		}

		high, err := strconv.ParseUint(highString, 10, 64)
		if err != nil {
			logger.Error(err)
			continue
		}

		rangeBin, err := ranges.NewRange(record[0], low, high)
		if err != nil {
			logger.Error(err)
			continue
		}

		treeMasterCardRange.Insert(rangeBin)

	}

	binsFile, err := os.OpenFile("./bins.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	binsFileReader := csv.NewReader(binsFile)
	record, err = binsFileReader.Read()
	if err != nil {
		logger.Fatal(err)
	}

	resultFileCSV, err := os.OpenFile("./result_check_mc.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	rw := csv.NewWriter(resultFileCSV)
	err = rw.Write([]string{"MC_name", "Bin_code", "MC_low", "MC_high", "Bin_low", "Bin_high"})
	if err != nil {
		logger.Fatal(err)
	}
	for {
		record, err = binsFileReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Fatal(err)
		}

		bin, err := strconv.ParseUint(record[4], 10, 64)
		if err != nil {
			logger.Error(err)
			continue
		}

		findRange := treeMasterCardRange.Find(bin)
		if findRange == nil {
			logger.Info("BIN:", bin, "not found")
		} else {
			err = rw.Write([]string{findRange.Code(), record[0], fmt.Sprint(findRange.GetLow()), fmt.Sprint(findRange.GetHigh()), record[4], record[5]})
			if err != nil {
				logger.Fatal(err)
			}
		}
	}
}
