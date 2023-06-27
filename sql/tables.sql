
--drop table sensemon.sensorreads;
--drop table sensemon.sensor;
--drop table sensemon.sensortype;

create table sensemon.sensorreads (
   SR_DATE date,
   SR_DEVICE_ID varchar2(30 char),
   SR_FARENHEIT number,
   SR_HUMIDITY number
);

create table sensemon.sensor (
   SENSOR_DEVICE_ID varchar2(30 char),
   SENSOR_TYPE_ID number,
   SENSOR_ADDRESS varchar2(128 char),
   SENSOR_NAME varchar2(64 char)
);

create table sensemon.sensortype (
   TYPE_ID number,
   TYPE_NAME varchar2(32 char)
);

grant select,insert on sensemon.sensorreads to app_sensemon;
grant select on sensemon.sensor to app_sensemon;
grant select on sensemon.sensortype to app_sensemon;

create synonym app_sensemon.sensorreads for sensemon.sensorreads;
create synonym app_sensemon.sensor for sensemon.sensor;
create synonym app_sensemon.sensortype for sensemon.sensortype;

