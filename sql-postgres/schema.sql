
DROP USER sensemon;
CREATE USER sensemon;
CREATE SCHEMA IF NOT EXISTS sensemon;
GRANT ALL PRIVILEGES ON SCHEMA sensemon TO sensemon;

\password sensemon

DROP USER app_sensemon;
CREATE USER app_sensemon;
CREATE SCHEMA IF NOT EXISTS app_sensemon;
GRANT ALL PRIVILEGES ON SCHEMA app_sensemon TO app_sensemon;
\password app_sensemon

