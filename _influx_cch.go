package logic

import (
	"encoding/json"
	"fmt"
	"mg/config"
	v2 "mg/models/v2"
	"mg/repo"
	"strings"
	"sync"
	"time"

	"advantech.com/ensaas/errors"

	"advantech.com/ensaas/logger"
	"github.com/influxdata/influxdb/client/v2"
)

//const Monthly_DATE_FORMAT = "2006-01"
const DATE_FORMAT = "2006-01"

var (
	usage            *Usage
	usageServiceOnce sync.Once
	partsRepo        *repo.PartsRepo
	numberOfRetry    int
)

func init() {
	usageServiceOnce.Do(func() {
		numberOfRetry = 5
		usage = new(Usage)
		usage.db = repo.GetInfluxUsageRepo()
		usage.RedisRepo = repo.GetRedisRepo()
		partsRepo = repo.GetPartsRepo()
	})
}

//Interval has start and end
type Interval struct {
	Start int64
	End   int64
}

type Usage struct {
	db        *repo.InfluxUsageRepo
	RedisRepo *repo.RedisRepo
}

// NewUsage is used for create usage
func NewUsage() *Usage {
	object := new(Usage)

	return object
}

func GetUsage() *Usage {
	return usage
}

func (p *Usage) getInterval(usageTime int64) (monthly Interval, daily Interval, hourly Interval) {
	//millisecond to second
	// dbTime := config.Location.Unix(usageTime/1000, 0)
	dbTime := time.Unix(usageTime/1000, 0)
	// dbTime.In(config.Location).Date()
	year, month, day := dbTime.In(config.Location).Date()
	hour := dbTime.In(config.Location).Hour()
	//time.FixedZone()
	monthly.Start = time.Date(year, month, 1, 0, 0, 0, 0, config.Location).UnixNano() / 1000000
	monthly.End = time.Date(year, month, 1, 0, 0, 0, 0, config.Location).AddDate(0, 1, 0).UnixNano() / 1000000
	daily.Start = time.Date(year, month, day, 0, 0, 0, 0, config.Location).UnixNano() / 1000000
	daily.End = time.Date(year, month, day, 0, 0, 0, 0, config.Location).AddDate(0, 0, 1).UnixNano() / 1000000
	hourly.Start = time.Date(year, month, day, hour, 0, 0, 0, config.Location).UnixNano() / 1000000
	hourly.End = hourly.Start + 3600000

	return monthly, daily, hourly
}

func (p *Usage) getMeasurementPrefix(usageTime int64) (measurement string) {
	dbTime := time.Unix(usageTime/1000, 0)
	year, month, _ := dbTime.In(config.Location).Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, config.Location)
	measurement = thisMonth.AddDate(0, 0, 0).Format(DATE_FORMAT)

	return measurement
}

func (p *Usage) getMeasurementName(pn string, usageTime int64) (monthlyName string, dailyName string, hourlyName string, rawName string) {
	prefix := p.getMeasurementPrefix(usageTime)
	monthlyName = pn
	dailyName = fmt.Sprintf("%s_%s-daily", prefix, pn)
	hourlyName = fmt.Sprintf("%s_%s-hourly", prefix, pn)
	rawName = fmt.Sprintf("%s_%s-raw", prefix, pn)
	return monthlyName, dailyName, hourlyName, rawName
}

func (p *Usage) getMeasurementNameByGranularity(pn string, time int64, granularity string) (name string) {
	monthlyName, dailyName, hourlyName, _ := p.getMeasurementName(pn, time)
	switch granularity {
	case "hourly":
		return hourlyName
	case "daily":
		return dailyName
	case "monthly":
		return monthlyName
	default:
		logger.Error("Granularity should not be " + granularity + ".")
		return ""
	}
}

func (p *Usage) filterIllegalUsage(usages []v2.Usages) (results []v2.Usages) {
	for i := range usages {
		if usages[i].ConsumerID == "" {
			logger.Error("Consumer Id is null.")
			passUsages := usages[0:i]
			if i != len(usages)-1 {
				//not last record
				results = p.filterIllegalUsage(usages[i+1 : len(usages)])
				passUsages = append(passUsages, results...)
			}
			return passUsages
		}
	}
	//all pass
	return usages
}

func (p *Usage) Insert(usage *v2.Usage) (err error) {
	if len(usage.Usages) == 0 {
		logger.Info("The length of usages is 0. part number: " + usage.PN)
		return nil
	}

	//verify usage
	usage.Usages = p.filterIllegalUsage(usage.Usages)

	monthInterval, dailyInterval, hourlyInterval := p.getInterval(usage.Time)
	monthlyName, dailyName, hourlyName, rawName := p.getMeasurementName(usage.PN, usage.Time)
	//insert raw data format 2020-01_pn-raw

	if err := p.insertUsage(rawName, usage); err != nil {
		return err
	}

	go p.RedisRepo.InsertRealTime(usage)
	logger.Info(usage.Time)

	usageInfo, err := partsRepo.FindPartFromPN(usage.PN)
	if err != nil {
		logger.Error(err)
		return err
	}

	if usageInfo == nil {
		logger.Warningf("Part Number %s can not be found.", usage.PN)
		return nil
	}
	//hourly
	//no hourly sql
	if usageInfo.AggregatedHourly == "" {
		return nil
	}
	hourlyUsage, err := p.queryAggregatedUsageWithInterval(usageInfo.AggregatedHourly, rawName, hourlyInterval)
	if err != nil {
		logger.Error(err)
		return err
	}
	hourlyUsage.PN = usage.PN
	config.Publish(hourlyUsage, "mg.usages."+usage.PN+".hourly")
	logger.Info(hourlyUsage.Time)
	if err := p.insertUsage(hourlyName, hourlyUsage); err != nil {
		return err
	}

	//daily
	if usageInfo.AggregatedDaily == "" {
		return nil
	}
	dailyUsage, err := p.queryAggregatedUsageWithInterval(usageInfo.AggregatedDaily, hourlyName, dailyInterval)
	if err != nil {
		logger.Error(err)
		return err
	}
	dailyUsage.PN = usage.PN
	config.Publish(dailyUsage, "mg.usages."+usage.PN+".daily")
	// logger.Info(dailyUsage.Time)
	if err := p.insertUsage(dailyName, dailyUsage); err != nil {
		return err
	}
	//monthly
	if usageInfo.AggregatedMonthly == "" {
		return nil
	}
	monthlyUsage, err := p.queryAggregatedUsageWithInterval(usageInfo.AggregatedMonthly, dailyName, monthInterval)
	if err != nil {
		logger.Error(err)
		return err
	}
	monthlyUsage.PN = usage.PN
	// logger.Info(monthlyUsage.Time)
	config.Publish(monthlyUsage, "mg.usages."+usage.PN+".monthly")
	p.insertUsage(monthlyName, monthlyUsage)

	//for cm
	cmInfo, _ := partsRepo.FindCMSql(usage.PN)
	if cmInfo == nil {
		logger.Warningf("Part Number %s can not be found.", usage.PN)
		return
	}
	sql := cmInfo.CMSQL
	//sql := "select mean(storage), sum(calls) from raw_db hourly_db daily_db monthly_db group by consumerId"
	sql = strings.Replace(sql, "raw_db", rawName, -1)
	sql = strings.Replace(sql, "hourly_db", hourlyName, -1)
	sql = strings.Replace(sql, "daily_db", dailyName, -1)
	sql = strings.Replace(sql, "monthly_db", monthlyName, -1)
	logger.Debug("sql: ", sql)

	cmUsage, err := p.queryAggregatedUsage(sql)
	if err != nil {
		logger.Error(err)
		return
	}
	cmUsage.PN = usage.PN
	if err != nil {
		logger.Error(err)
		return
	}

	if cmUsage.Usages != nil {
		config.Publish(cmUsage, "mg.usages."+usage.PN+".cm")
	}
	logger.Debug("end")

	return nil
}

func (p *Usage) insertUsage(measurement string, usage *v2.Usage) (err error) {
	dbTime := time.Unix(usage.Time/1000, 0)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "ms",
		Database:  p.db.Propertities.DBName,
	})

	for _, usageElement := range usage.Usages {
		//add tag, tag is consumer_id
		tag := make(map[string]string)
		tag["consumerId"] = usageElement.ConsumerID

		field := make(map[string]interface{})

		//add filed, filed means metric
		for _, measuredUsageElement := range usageElement.MeasuredUsage {
			field[measuredUsageElement.Measure] = measuredUsageElement.Quantity
		}

		point, err := client.NewPoint(
			measurement,
			tag,
			field, dbTime)
		if err != nil {
			logger.Debug(err)
		}
		bp.AddPoint(point)
	}

	for i := 0; i < numberOfRetry; i++ {
		if err := p.db.DB.Write(bp); err == nil {
			return nil
		}
		logger.Debug(err)
		time.Sleep(3 * time.Second)
	}

	logger.Debug(err)
	return errors.NewEnsaasError(errors.DB_Unavailable)

}

// This function only for migration, not normal operation.

func (p *Usage) Migration() (err error) {

	// hourlyStart := 1585666800000
	// hourlyEnd := 1587654000000

	// for ; hourlyStart <= hourlyEnd; hourlyStart += 3600000 {
	// 	hourlyUsage, err := p.queryAggregatedUsage(fmt.Sprintf("SELECT *  FROM \"2020-04_980GEDMA000-hourly\"  WHERE time >= %d000000 and time < %d000000 group by consumerId", hourlyStart, hourlyEnd))
	// 	hourlyUsage.Time = int64(hourlyStart)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	// 	// queryAggregatedUsage("SELECT *  FROM \"2020-04_980GEDMA000-hourly\"  WHERE time >= 1585670400000000000 and time < 1586926800000000000 group by consumerId") (usage *v2.Usage, err error) {
	// 	logger.Info(hourlyUsage.Time)

	// 	if err := p.insertUsage("2020-04_980GEDPA000-hourly", hourlyUsage); err != nil {
	// 		logger.Error(err)
	// 		// return err
	// 	}
	// }

	// hourlyStart = 1585666800000

	// for ; hourlyStart <= 1587654000000; hourlyStart += 86400000 {
	// 	daily, err := p.queryAggregatedUsage(fmt.Sprintf("SELECT *  FROM \"2020-04_980GEDMA000-daily\"  WHERE time >= %d000000 and time < %d000000 group by consumerId", hourlyStart, hourlyEnd))
	// 	daily.Time = int64(hourlyStart)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return err
	// 	}
	// 	// queryAggregatedUsage("SELECT *  FROM \"2020-04_980GEDMA000-daily\"  WHERE time >= 1585670400000000000 and time < 1586926800000000000 group by consumerId") (usage *v2.Usage, err error) {
	// 	logger.Info(daily.Time)

	// 	if err := p.insertUsage("2020-04_980GEDPA000-daily", daily); err != nil {
	// 		logger.Error(err)
	// 		// return err
	// 	}
	// }
	//month
	hourlyStart := 1582988400000
	hourlyEnd := 1583074800000
	month, err := p.queryAggregatedUsage(fmt.Sprintf("SELECT *  FROM \"980GEDMA000\"  WHERE time >= %d000000 and time < %d000000 group by consumerId", hourlyStart, hourlyEnd))
	month.Time = int64(hourlyStart)
	if err != nil {
		logger.Error(err)
		return err
	}
	// queryAggregatedUsage("SELECT *  FROM \"2020-04_980GEDMA000-daily\"  WHERE time >= 1585670400000000000 and time < 1586926800000000000 group by consumerId") (usage *v2.Usage, err error) {
	logger.Info(month.Time)

	if err := p.insertUsage("980GEDPA000", month); err != nil {
		logger.Error(err)
		// return err
	}

	return nil

}

func (p *Usage) Recalculate(pn string, start int64, end int64) (err error) {

	if end > start {
		return
	}

	usageInfo, err := partsRepo.FindPartFromPN(pn)
	if err != nil {
		logger.Error(err)
		return err
	}

	if usageInfo == nil {
		logger.Warningf("Part Number %s can not be found.", pn)
		return nil
	}

	for ; start < end; start += 3600 {
		monthInterval, dailyInterval, hourlyInterval := p.getInterval(start)
		monthlyName, dailyName, hourlyName, rawName := p.getMeasurementName(pn, start)
		//hourly
		//no hourly sql
		if usageInfo.AggregatedHourly == "" {
			return nil
		}
		hourlyUsage, err := p.queryAggregatedUsageWithInterval(usageInfo.AggregatedHourly, rawName, hourlyInterval)
		if err != nil {
			logger.Error(err)
			return err
		}
		hourlyUsage.PN = pn

		logger.Info(hourlyUsage.Time)
		if err := p.insertUsage(hourlyName, hourlyUsage); err != nil {
			return err
		}
		//daily
		if usageInfo.AggregatedDaily == "" {
			return nil
		}
		dailyUsage, err := p.queryAggregatedUsageWithInterval(usageInfo.AggregatedDaily, hourlyName, dailyInterval)
		if err != nil {
			logger.Error(err)
			return err
		}
		dailyUsage.PN = pn

		if err := p.insertUsage(dailyName, dailyUsage); err != nil {
			return err
		}
		//monthly
		if usageInfo.AggregatedMonthly == "" {
			return nil
		}
		monthlyUsage, err := p.queryAggregatedUsageWithInterval(usageInfo.AggregatedMonthly, dailyName, monthInterval)
		if err != nil {
			logger.Error(err)
			return err
		}
		monthlyUsage.PN = pn
		// logger.Info(monthlyUsage.Time)
		p.insertUsage(monthlyName, monthlyUsage)
	}

	return nil
}

func (p *Usage) Write(bp client.BatchPoints) error {
	return p.db.DB.Write(bp)
}

func (p *Usage) queryAggregatedUsage(cmd string) (usage *v2.Usage, err error) {
	usage = new(v2.Usage)
	// logger.Println(p.db.Propertities.DBName)
	logger.Debug("sql: " + cmd)
	q := client.Query{
		Command:  cmd,
		Database: p.db.Propertities.DBName,
	}
	for i := 0; i < numberOfRetry; i++ {
		if response, err := p.db.DB.Query(q); err == nil {
			if response.Error() != nil {
				logger.Error("Influx db query failed.")
				logger.Debug(response.Error())
				return nil, errors.NewEnsaasError(errors.DB_Unavailable)
			}
			//		response.Results
			for _, result := range response.Results {
				for _, series := range result.Series {
					usages := v2.Usages{}
					usages.ConsumerID = series.Tags["consumerId"]
					for i := 1; i < len(series.Columns); i++ {
						measuredUsage := v2.MeasuredUsage{}
						measuredUsage.Measure = series.Columns[i]
						if series.Values[0][i] != nil {
							measuredUsage.Quantity, _ = series.Values[0][i].(json.Number).Float64()
						}

						usages.MeasuredUsage = append(usages.MeasuredUsage, measuredUsage)
					}
					usage.Usages = append(usage.Usages, usages)
				}
			}
			return usage, nil
		}
	}
	logger.Error("Influx db query failed.")
	logger.Debug(err)
	return nil, errors.NewEnsaasError(errors.DB_Unavailable)
}

func (p *Usage) queryAggregatedUsageWithInterval(cmd string, measurement string, interval Interval) (usage *v2.Usage, err error) {
	sql := fmt.Sprintf("%s FROM \"%s\" WHERE time >= %d000000 and time < %d000000 group by consumerId", cmd, measurement, interval.Start, interval.End)
	logger.Debug("Interval sql: " + sql)
	usage, err = p.queryAggregatedUsage(sql)
	if err != nil {
		logger.Error("Query interval usage failed.")
		logger.Debug(err)
		return nil, errors.NewEnsaasError(errors.DB_Unavailable)
	}
	usage.Time = interval.Start
	return usage, nil
}

func (p *Usage) queryUsage(cmd string) (usages *v2.QueryUsages, err error) {
	usages = new(v2.QueryUsages)
	logger.Debug("sql: " + cmd)
	q := client.Query{
		Command:  cmd,
		Database: p.db.Propertities.DBName,
	}

	if response, err := p.db.DB.Query(q); err == nil {
		if response.Error() != nil {
			logger.Error("Influx db query failed.")
			logger.Debug(response.Error())
			return nil, errors.NewEnsaasError(errors.DB_Unavailable)
		}

		//		response.Results
		for _, result := range response.Results {
			//series is group of result (like group by consumerId, each conusmerId in one serie)
			for _, series := range result.Series {
				usage := v2.QueryUsage{}
				consumerId := series.Tags["consumerId"]
				usage.ConsumerId = consumerId
				usage.Data = []v2.Data{}

				//values[row][col]
				if len(series.Columns) > 0 {
					for row := 0; row < len(series.Values); row++ {
						dataUsages := v2.Data{}

						// dataUsages.Date = series.Values[row][0].(string)
						epoch, err := time.Parse(time.RFC3339, series.Values[row][0].(string))
						if err == nil {
							//no error
							dataUsages.DateEpoch = epoch.Unix() * 1000
							dataUsages.Date = epoch.In(config.GetLocation()).String()
						}
						dataUsages.Usages = []v2.DataUsage{}

						//columns
						for i := 1; i < len(series.Values[row]); i++ {
							dataUsage := v2.DataUsage{}
							dataUsage.MetricName = series.Columns[i]
							if series.Values[row][i] != nil {
								dataUsage.Quantity, _ = series.Values[row][i].(json.Number).Float64()
							}

							dataUsages.Usages = append(dataUsages.Usages, dataUsage)
						}
						usage.Data = append(usage.Data, dataUsages)
					}
				}
				usages.Usages = append(usages.Usages, usage)
			}
		}
		return usages, nil
	} else {
		logger.Error("Influx db query failed.")
		logger.Debug(err)
		return nil, errors.NewEnsaasError(errors.DB_Unavailable)
	}
}

func (p *Usage) QueryUsageWithInterval(pn string, startTime int64, endTime int64, granularity string, consumerId string) (usage *v2.QueryUsages, err error) {
	measurement := p.getMeasurementNameByGranularity(pn, startTime, granularity)
	sql := fmt.Sprintf("select * FROM \"%s\" WHERE time >= %d000000 and time <= %d000000", measurement, startTime, endTime)
	//has consumerId
	if len(consumerId) > 0 {
		sql = sql + " and \"consumerId\" = '" + strings.ReplaceAll(consumerId, ",", "' or \"consumerId\" = '") + "'"
	}
	sql = sql + " group by consumerId"

	logger.Debug("Interval sql: " + sql)

	usage, err = p.queryUsage(sql)
	if err != nil {
		logger.Error("Query interval usage failed.")
		logger.Debug(err)
		return nil, errors.NewEnsaasError(errors.DB_Unavailable)
	}

	if usage.Usages == nil {
		usage.Usages = []v2.QueryUsage{}
	}

	return usage, nil
}

func (p *InfluxUsageRepo) QueryMonthUsage(time int64, pn string, consumerId string) (usages map[string]float64, err error) {
	sql := fmt.Sprintf("select * FROM \"%s\"  WHERE time = %d000000 and consumerId = '%s'", pn, time, consumerId)
	logger.Debug("sql: " + sql)
	results, err := p.QueryUsage(sql)
	usages = map[string]float64{}
	for _, result := range results {
		for _, series := range result.Series {
			for i := 1; i < len(series.Columns); i++ {
				if series.Values[0][i] != nil {
					if value, ok := series.Values[0][i].(json.Number); ok {
						usages[series.Columns[i]], _ = value.Float64()
					}
				}
			}
		}
	}
	return usages, nil
}
