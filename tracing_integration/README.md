# sample service
This is a sample, about how to integrate  istio/ app /logs/ client together.

## For docker-compose debug.
- Run `docker-compose up`
- Put `http://localhost:16686/` in browser, to access jaeger
- Put `http://localhost:8081/servicetest/v1/jaegertest` in browser, to access jaegertest start by docker-compose
- Change code /cmd/main.go, 
  - Port from "8081" to "8080"
  - Uncomment below
  ```
  	//Uncomment below, for docker-compose debug
  	//os.Setenv("JAEGER_SERVICE_NAME", "fromCode")
  	//os.Setenv("JAEGER_AGENT_HOST", "localhost")
  	//os.Setenv("JAEGER_AGENT_PORT", "6831")
  	//os.Setenv("nextserviceurl", "http://localhost:8081/servicetest/v1/jaegertest")
  ```
- Start `/cmd/main.go`, Put `http://localhost:8080/servicetest/v1/jaegertest` in browser, to access jaegertest start by Goland.
- Check the log in docker-compose console, it had been invoked also
- Check the jaeger, all span had been integrated, include `service in docker` and `service in Goland`.

## For k8s test
- Run `kubectl create namespace jaegertest`
- Run `kubectl label namespace jaegertest istio-injection=enabled`
- Run `kubectl create -n jaegertest -f k8s_deploy/` to create two deploy
- Run `kubectl -n jaegertest create svc clusterip foo --tcp=8081:8081 `
- Run `kubectl -n jaegertest create svc clusterip bar --tcp=8081:8081 `
- Change virtual service, let foo service expose to outside acces.
```
 - match:
    - uri:
        prefix: /servicetest
    route:
    - destination:
        host: jaegertest2.hpbp-stg-core.svc.cluster.local
        port:
          number: 8081

```

## Test Result
### Jaeger Connect together - ingress-gateway, istio-proxy, applications
![result-of-jaeger](https://github.com/johnzheng1975/istiomeetup_202109/blob/main/tracing_integration/testresult/result-k8s.png?raw=true)

### Logs connect together - client, ingress-gateway, istio-proxy, applications
![result-of-logs](https://github.com/johnzheng1975/istiomeetup_202109/blob/main/tracing_integration/testresult/log-k8s-with-x-request-id.png?raw=true)

### Query trace with x-request-id
![result-of-logs](https://github.com/johnzheng1975/istiomeetup_202109/blob/main/tracing_integration/testresult/result-query-trace-based-on-x-request-id.png?raw=true)
