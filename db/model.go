package db

import (
	"context"
	"database/sql"
	"sensemon/model"
	"sensemon/sensor"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	queryTimeout = time.Duration(10 * time.Second)
)

func (dbc *Connection) InsertDhtData(data *sensor.DhtSensorData) error {
	sql := `insert into sensorreads (
       sr_date,
	   sr_device_id,
	   sr_farenheit,
	   sr_humidity
	) values (
	   :sr_date,
	   :sr_device_id,
	   round(:sr_farenheit,2),
	   round(:sr_humidity,2)
	)`
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tx, err := dbc.BeginTxx(ctx, nil)
	if err != nil {
		log.Errorf("Failed to begin transaction: %s", err)
		return err
	}

	_, err = tx.NamedExec(sql, data)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (dbc *Connection) Sensors() ([]*sensor.Sensor, error) {
	sensors := make([]*sensor.Sensor, 0)
	sql := `select * from sensor`
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := dbc.QueryxContext(ctx, sql)
	if err != nil {
		log.Errorf("Can't query: %s", err)
		return sensors, err
	}
	defer rows.Close()

	for rows.Next() {
		row := &sensor.Sensor{}
		if err := rows.StructScan(row); err != nil {
			return sensors, err
		}
		sensors = append(sensors, row)
	}
	return sensors, nil
}

func (dbc *Connection) AllDhtDataForSensor(deviceId string) ([]*sensor.DhtSensorData, error) {
	return dbc.AllDhtDataForSensorInterval(deviceId, 10)
}

func (dbc *Connection) AllDhtDataInterval(minuteInterval int) ([]*sensor.DhtSensorData, error) {
	allData := make([]*sensor.DhtSensorData, 0)
	q := `
	WITH q_params as (
	select
		$1::int as p_interval,
		2 as p_days_back
	),
	intervals AS (
		SELECT
			generate_series(
				(SELECT date_bin(interval '1 min' * p_interval, current_timestamp, current_date)-
				interval '1 day'* p_days_back),
				(SELECT MAX(sr_date) FROM sensemon.sensorreads),
				INTERVAL '1 minute' * p_interval
			) AS start_time
			from q_params
	)
	SELECT
		intervals.start_time as sr_date,
		coalesce(sr_device_id, '') as sr_device_id,
		--intervals.start_time + INTERVAL '1 minute' * p_interval AS end_time,
		coalesce(AVG(sr_farenheit),0) AS sr_farenheit,
		coalesce(AVG(sr_humidity),0) AS sr_humidity
		--,COUNT(*) AS reading_count
	FROM
		intervals
	left join sensemon.sensorreads
	on
		(sr_date >= intervals.start_time
		AND sr_date < intervals.start_time + INTERVAL '1 minute' * (select p_interval from q_params))
	GROUP BY
		sr_device_id, intervals.start_time
	ORDER BY
		intervals.start_time
	`
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := dbc.QueryxContext(ctx, q, minuteInterval)
	if err != nil {
		log.Errorf("Can't query: %s", err)
		return allData, err
	}
	defer rows.Close()

	for rows.Next() {
		row := &sensor.DhtSensorData{}
		if err := rows.StructScan(row); err != nil {
			return allData, err
		}
		allData = append(allData, row)
	}
	return allData, nil
}

func (dbc *Connection) AllDhtDataForSensorInterval(deviceId string, minuteInterval int) ([]*sensor.DhtSensorData, error) {
	// https://stackoverflow.com/questions/6195439/postgres-how-do-you-round-a-timestamp-up-or-down-to-the-nearest-minute
	allData := make([]*sensor.DhtSensorData, 0)
	q := `
	WITH q_params as (
	select
		$1::int as p_interval,
		$2 as p_device_id,
		7 as p_days_back
	),
	intervals AS (
		SELECT
			generate_series(
				(SELECT date_bin(interval '1 min' * p_interval, current_timestamp, current_date)-
				interval '1 day'* p_days_back),
				(SELECT MAX(sr_date) FROM sensemon.sensorreads where sr_device_id = p_device_id),
				INTERVAL '1 minute' * p_interval
			) AS start_time
			from q_params
	)
	SELECT
		intervals.start_time as sr_date,
		coalesce(sr_device_id, '') as sr_device_id,
		coalesce(AVG(sr_farenheit),0) AS sr_farenheit,
		coalesce(AVG(sr_humidity),0) AS sr_humidity
	FROM
		intervals
	inner join sensemon.sensorreads
	on
		    sr_date >= intervals.start_time
		AND sr_date < intervals.start_time + INTERVAL '1 minute' * (select p_interval from q_params)
		AND sr_device_id = (select p_device_id from q_params)
	GROUP BY
		sr_device_id, intervals.start_time
	ORDER BY
		intervals.start_time
	`
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := dbc.QueryxContext(ctx, q, minuteInterval, deviceId)
	if err != nil {
		log.Errorf("Can't query: %s", err)
		return allData, err
	}
	defer rows.Close()

	for rows.Next() {
		row := &sensor.DhtSensorData{}
		if err := rows.StructScan(row); err != nil {
			return allData, err
		}
		allData = append(allData, row)
	}
	return allData, nil
}

func (dbc *Connection) LatestDhtReadings() ([]*model.LatestDhtSensorData, error) {
	allData := make([]*model.LatestDhtSensorData, 0)
	q := `
	select sensor_name,
       (select sr_farenheit
          from sensorreads f
         where f.sr_device_id = latest.sensor_device_id
           and f.sr_date = latest.last_entry_date
           limit 1
         ) as fahrenheit,
         (select sr_humidity
          from sensorreads h
         where h.sr_device_id = latest.sensor_device_id
           and h.sr_date = latest.last_entry_date
           limit 1
         ) as humidity,
         latest.LAST_ENTRY_DATE as last_entry_date
  from ( select sensor_name, sensor_device_id, max(sr_date) last_entry_date
           from sensor, sensorreads
          where sr_device_id = sensor_device_id
          group by sensor_name, sensor_device_id
          order by sensor_name ) latest
	`
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := dbc.QueryxContext(ctx, q)
	if err != nil {
		log.Errorf("Can't query: %s", err)
		return allData, err
	}
	defer rows.Close()

	for rows.Next() {
		row := &model.LatestDhtSensorData{}
		if err := rows.StructScan(&row); err != nil {
			return allData, err
		}
		allData = append(allData, row)
	}

	return allData, nil
}

func (dbc *Connection) TableExists(table_name string) (bool, error) {
	log.Tracef("Checking if table \"%s\" exists", table_name)
	t := time.Now()
	stmt, err := dbc.Prepare(`
	select 'Y'
	  from all_tables
	 where table_name = upper(:1)
	   and owner = upper(user)
	`)

	if err != nil {
		return false, err
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Errorf("Can't close dataset: %s", err)
		}
	}()

	rows, err := stmt.Query(table_name)
	if err != nil {
		log.Errorf("Can't query: %s", err)
		return false, err
	}
	defer rows.Close()

	res := ""

	if rows.Next() {
		err := rows.Scan(&res)
		switch {
		case err == sql.ErrNoRows:
			log.Debugf("ErrNoRows - Table \"%s\" does not exists %s", table_name, time.Since(t))
			return false, nil
		case err != nil:
			return false, err
		}
	} else {
		log.Debugf("No next() record. Table \"%s\" does not exists %s", table_name, time.Since(t))
		return false, nil
	}
	log.Debugf("Table \"%s\" exists %s", table_name, time.Since(t))

	return true, nil
}
