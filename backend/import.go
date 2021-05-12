package main

import (
	"encoding/csv"
	"github.com/Jleagle/gym-tracker/log"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/Jleagle/gym-tracker/influx"
	"go.uber.org/zap"
)

// DROP MEASUREMENT GymTracker
// SELECT max("people") AS "mean_people" FROM "GymTracker"."alltime"."gyms" WHERE time > now() - 1000w GROUP BY time(1m) FILL(none)

func importFromChronograf() {

	csvfile, err := os.Open("import.csv")
	if err != nil {
		log.Instance.Error("reading file", zap.Error(err))
		return
	}

	r := csv.NewReader(csvfile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Instance.Error("reading csv", zap.Error(err))
		}

		if len(record) == 2 {

			if record[0] == "time" || record[1] == "" {
				continue
			}

			t, err := time.Parse(time.RFC3339, record[0])
			if err != nil {
				log.Instance.Error("reading csv", zap.Error(err), zap.String("time", record[0]))
				continue
			}

			i, err := strconv.Atoi(record[1])
			if err != nil {
				log.Instance.Error("reading csv", zap.Error(err))
				continue
			}

			p := calculatePercent(i, 120)

			_, err = influx.Write("fareham", i, 120, p, t)
			if err != nil {
				log.Instance.Error("reading csv", zap.Error(err))
				continue
			}
		}
	}

	log.Instance.Info("Done")
}
