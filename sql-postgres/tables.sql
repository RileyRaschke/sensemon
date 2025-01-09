
drop table sensemon.sensorreads;

-- Create sensorreads table
CREATE TABLE sensemon.sensorreads (
    sr_date timestamp with time zone,
    sr_device_id VARCHAR(30),
    sr_farenheit NUMERIC,
    sr_humidity NUMERIC
);

-- Create sensor table
CREATE TABLE sensemon.sensor (
    sensor_device_id VARCHAR(30),
    sensor_type_id NUMERIC,
    sensor_address VARCHAR(128),
    sensor_name VARCHAR(64)
);

-- Create sensortype table
CREATE TABLE sensemon.sensortype (
    type_id NUMERIC,
    type_name VARCHAR(32)
);

GRANT select,insert,update,delete ON ALL TABLES IN SCHEMA sensemon to sensemon;

-- Grant privileges to app_sensemon

GRANT SELECT, INSERT ON sensemon.sensorreads TO app_sensemon;
GRANT SELECT ON sensemon.sensor TO app_sensemon;
GRANT SELECT ON sensemon.sensortype TO app_sensemon;

