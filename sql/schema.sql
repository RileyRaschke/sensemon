
drop user sensemon;
drop user app_sensemon;

create user sensemon;
create user app_sensemon identified by &credential;

grant unlimited tablespace to sensemon;
grant create session to app_sensemon;
grant connect,resource to app_sensemon;

