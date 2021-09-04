package main

import (
	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	jaegerlog "github.com/opentracing/opentracing-go/log"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"io"
	"log"
	"net/http"
	"os"
)

var servicename string
var jhost string
var jport string
var nextserviceurl string
var err error
var tracer opentracing.Tracer
var closer io.Closer
var cfg *jaegercfg.Configuration

func init() {
	log.Println("Initialize ... ...")
	//os.Setenv("JAEGER_SERVICE_NAME", "fromCode")
	//os.Setenv("JAEGER_AGENT_HOST", "localhost")
	//os.Setenv("JAEGER_AGENT_PORT", "6831")
	//os.Setenv("nextserviceurl", "http://localhost:8081/servicetest/v1/jaegertest")

	if os.Getenv("JAEGER_SERVICE_NAME") != "" {
		servicename = os.Getenv("JAEGER_SERVICE_NAME")
	} else {
		servicename = "default"
	}
	if os.Getenv("JAEGER_AGENT_HOST") != "" {
		jhost = os.Getenv("JAEGER_AGENT_HOST")
	} else {
		jhost = "localhost"
	}
	if os.Getenv("JAEGER_AGENT_PORT") != "" {
		jport = os.Getenv("JAEGER_AGENT_PORT")
	} else {
		jport = "6831"
	}

	if os.Getenv("nextserviceurl") != "" {
		nextserviceurl = os.Getenv("nextserviceurl")
	} else {
		nextserviceurl = ""
	}

	tracer, closer = InitJaeger(servicename, jhost, jport)
	//defer closer.Close()
}

// Init returns an instance of Jaeger Tracer
func InitJaeger(service string, jhost string, jport string) (opentracing.Tracer, io.Closer) {
	jLocalAgentHostPort := jhost + ":" + jport

	cfg := &jaegercfg.Configuration{
		ServiceName: service,

		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jLocalAgentHostPort,
		},
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator(zipkin.BaggagePrefix("x-request-id"))
	opts := []jaegercfg.Option{
		jaegercfg.ZipkinSharedRPCSpan(true),
		jaegercfg.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
	}

	tracer, closer, err := cfg.NewTracer(opts...)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func jaegertest(w http.ResponseWriter, r *http.Request) {
	for key, value := range r.Header {
		fmt.Println("Key:", key, "Value:", value)
	}
	x_request_id := r.Header.Get("x-request-id")
	if(x_request_id == "") {x_request_id="no_x-request-id"}
	log.Println(x_request_id + "  Started jaegertest API now.")

	span_context, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	spanMain := tracer.StartSpan("spainMain-jaegertest", opentracing.ChildOf(span_context))
	spanMain.LogFields(
		jaegerlog.String("mainfield1", "value1"),
	)
	spanMain.LogKV("maink1", "v1")
	spanMain.SetTag("maintag", "value")
	invokeAnotherFunction(&spanMain)

	spanInvokeAPI := tracer.StartSpan("spanInvokeOtherAPI", opentracing.ChildOf(spanMain.Context()))
	spanInvokeAPI.LogFields(
		jaegerlog.String("invokeOtherAPI", "invokeOtherAPIValue"),
	)
	if nextserviceurl != "" {
		client := &http.Client{
			//CheckRedirect: redirectPolicyFunc,
		}
		req, err := http.NewRequest("GET", nextserviceurl, nil)

		tracer.Inject(spanInvokeAPI.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		resp, err := client.Do(req)

		if err != nil {
			// handle error
		}

		//defer resp.Body.Close()
		log.Println(resp.Body)

	}

	spanInvokeAPI.Finish()

	spanMain.Finish()

	log.Println(x_request_id + "  Finished jaegertest API.")
	fmt.Fprintln(w, "hello!")
}

func invokeAnotherFunction(sp *opentracing.Span) {
	span := tracer.StartSpan("spanInAnotherFunc", opentracing.ChildOf(opentracing.Span(*sp).Context()))
	span.SetTag("tagInAnotherFunc", "funnyTagValue")

	span.LogFields(
		jaegerlog.String("key1InOtherFunc", "cool"),
		jaegerlog.String("key2InOtherFunc", "try"),
	)

	span.LogKV("logKV-key", "logKV-value")

	span.Finish()

}

func main() {
	http.HandleFunc("/servicetest/v1/jaegertest", jaegertest)
	http.ListenAndServe(":8081", nil)
}
