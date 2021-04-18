package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/Jleagle/puregym-tracker/influx"
	"go.uber.org/zap"
)

func importFromChronograf() {

	csvfile, err := os.Open("import.csv")
	if err != nil {
		logger.Error("reading file", zap.Error(err))
		return
	}

	r := csv.NewReader(csvfile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("reading csv", zap.Error(err))
		}

		if len(record) == 2 {

			if record[1] == "" {
				continue
			}

			t, err := time.Parse(time.RFC3339, record[0])
			if err != nil {
				logger.Error("reading csv", zap.Error(err))
				continue
			}

			i, err := strconv.Atoi(record[1])
			if err != nil {
				logger.Error("reading csv", zap.Error(err))
				continue
			}

			p := calculatePercent(i, 120)

			_, err = influx.Write("fareham", i, 120, p, t)
			if err != nil {
				logger.Error("reading csv", zap.Error(err))
				continue
			}
		}
	}
}
