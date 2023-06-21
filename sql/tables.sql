
--drop table sensemon.sensorreads;

create table sensemon.sensorreads (
   sr_date date,
   sr_device_id varchar2(30 char),
   sr_farenheit number,
   sr_humidity number
);

grant select,insert on sensemon.sensorreads to app_sensemon;

