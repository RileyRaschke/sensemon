set linesize 150
spool pg_conversion_prod.sql
SELECT
    'INSERT INTO sensemon.sensorreads (sr_date, sr_device_id, sr_farenheit, sr_humidity) VALUES ('
    || '''' || TO_CHAR(sr_date, 'YYYY-MM-DD HH24:MI:SS') || '''::timestamp, '
    ||CASE WHEN sr_device_id IS NOT NULL THEN
        '''' || REPLACE(sr_device_id, '''', '''''') || ''''
       ELSE 'NULL'
       END || ', '
    || sr_farenheit || ', '
    || sr_humidity || ');'
FROM
    sensemon.sensorreads;
spool off
