package db

import (
	"context"
	"database/sql"
	"sensemon/sensor"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	queryTimeout = time.Duration(10 * time.Second)
)

type TabInfo struct {
	Owner     string `db:"Owner" json:"owner"`
	TableName string `db:"TableName" json:"tableName"`
	RowCount  int64  `db:"RowCount" json:"rowCount"`
}

func (dbc *Connection) InsertDhtData(data *sensor.DhtSensorData) error {
	sql := `insert into sensemon.sensorreads (
       sr_date,
	   sr_device_id, 
	   sr_farenheit,
	   sr_humidity
	) values (
	   :SR_DATE,
	   :SR_DEVICE_ID,
	   round(:SR_FARENHEIT,2),
	   round(:SR_HUMIDITY,2) 
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

func (dbc *Connection) AllDhtDataForSensor(deviceId string) ([]*sensor.DhtSensorData, error) {
	return dbc.AllDhtDataForSensorInterval(deviceId, 5)
}

func (dbc *Connection) AllDhtDataForSensorInterval(deviceId string, minuteInterval int) ([]*sensor.DhtSensorData, error) {
	allData := make([]*sensor.DhtSensorData, 0)
	q := `
		with rws as (
		  select trunc ( sr_date ) dy,
				 trunc ( sr_date, 'mi' ) mins,
				 :1 / 1440 time_interval,
				 sr_device_id,
				 sr_farenheit,
				 sr_humidity
		  from   sensemon.sensorreads
         where sr_device_id = :2
		   and sr_date >= sysdate-2
		), intervals as (
		  select dy + (
				   floor ( ( mins - dy ) / time_interval ) * time_interval
				 ) start_datetime,
				 sr_device_id,
				 sr_farenheit,
				 sr_humidity
		  from   rws
		)
		  select start_datetime as sr_date,
				 sr_device_id,
				 round(avg(sr_farenheit),2) as sr_farenheit,
				 round(avg(sr_humidity),2) as sr_humidity
			from intervals
		  group  by start_datetime, sr_device_id
		  order  by start_datetime
	`
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := dbc.QueryxContext(ctx, q, minuteInterval, deviceId)
	defer rows.Close()
	if err != nil {
		log.Errorf("Can't query: %s", err)
		return allData, err
	}

	for rows.Next() {
		row := &sensor.DhtSensorData{}
		if err := rows.StructScan(row); err != nil {
			return allData, err
		}
		allData = append(allData, row)
	}
	return allData, nil
}

func (dbc *Connection) AllTables() ([]TabInfo, error) {
	allData := make([]TabInfo, 0)
	q := `
	select
		   owner as "Owner",
		   table_name as "TableName",
		   num_rows as "RowCount"
      from all_tables
	 where owner = USER
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
		row := TabInfo{}
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
