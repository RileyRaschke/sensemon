
drop table sensor_reads_temp;

create table sensor_reads_temp as
WITH q_params as (
   select 30::int as p_interval,
           EXTRACT(epoch from 
               (SELECT current_date-date_bin(interval '1 day', min(sr_date), current_date) from sensemon.sensorreads)
           )::int/(3600*24) as p_days_back
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
),
tabed_data as (
    SELECT
        intervals.start_time as sr_date,
        coalesce(sr_device_id, '') as sr_device_id,
        coalesce(round(AVG(sr_farenheit),1),0) AS sr_farenheit,
        coalesce(round(AVG(sr_humidity),1),0) AS sr_humidity
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
		)
select * 
  from tabed_data 
 where sr_device_id is not null 
   and sr_device_id != ''
order by sr_date
;

select * from sensor_reads_temp;

truncate table sensemon.sensorreads;

insert into sensemon.sensorreads select * from sensor_reads_temp;

commit;

drop table sensor_reads_temp;

