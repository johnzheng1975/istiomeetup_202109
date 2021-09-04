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
 
