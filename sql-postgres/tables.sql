
-- Create sensorreads table
CREATE TABLE sensemon.sensorreads (
    sr_date DATE,
    sr_device_id VARCHAR(30),
    sr_farenheit NUMERIC, -- Replace "number" with NUMERIC
    sr_humidity NUMERIC   -- Replace "number" with NUMERIC
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

-- Grant privileges to app_sensemon
GRANT SELECT, INSERT ON sensemon.sensorreads TO app_sensemon;
GRANT SELECT ON sensemon.sensor TO app_sensemon;
GRANT SELECT ON sensemon.sensortype TO app_sensemon;

-- Create synonyms (PostgreSQL uses views for this)
CREATE VIEW app_sensemon.sensorreads AS SELECT * FROM sensemon.sensorreads;
CREATE VIEW app_sensemon.sensor AS SELECT * FROM sensemon.sensor;
CREATE VIEW app_sensemon.sensortype AS SELECT * FROM sensemon.sensortype;

