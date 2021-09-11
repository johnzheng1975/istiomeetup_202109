-- # Alert in case one API P90 response time more than 1 second in last minute
select *
from devops.horizon_monitor 
where response_code in ('200', '201') 
and p90_duration > 1000
