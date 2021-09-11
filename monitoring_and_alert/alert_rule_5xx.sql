-- # 500 count > 10%, on the same time, api invoke count > 10
WITH table_all AS (
select time, service_name, api_path, api_method, sum(measure_value::double) as allcount
FROM devops.horizon_monitor
where measure_name = 'api_count'
group by time, service_name, api_path, api_method 
), table_5xx AS(
select time, service_name, api_path, api_method, response_code, sum(measure_value::double) as errorCount
FROM devops.horizon_monitor
where measure_name = 'api_count'
and response_code in ('500','501','502','503','504')
group by time, service_name, api_path, api_method, response_code
)
select a.time, a.service_name, a.api_path, a.api_method, b.response_code, errorCount, allcount,  errorCount/allcount as failrate
from table_all a, table_5xx b
where a.time=b.time and a.service_name = b.service_name and a.api_path = b.api_path and a.api_method = b.api_method
and errorCount/allcount > 0.1
and allcount > 10
