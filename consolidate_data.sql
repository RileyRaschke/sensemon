

create table sensor_reads_temp as
        with rws as (
          select trunc ( sr_date ) dy,
                 trunc ( sr_date, 'mi' ) mins,
                 5 / 1440 time_interval,
                 sr_device_id,
                 sr_farenheit,
                 sr_humidity
          from   sensorreads
         where sr_date < trunc(sysdate)
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
;

delete from sensorreads where sr_date < trunc(sysdate);

insert into sensorreads select * from sensor_reads_temp;

commit;

drop table sensor_reads_temp;

