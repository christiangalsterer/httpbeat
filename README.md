[![Build Status](https://travis-ci.org/christiangalsterer/httpbeat.svg?branch=master)](https://travis-ci.org/christiangalsterer/httpbeat)
[![codecov.io](http://codecov.io/github/christiangalsterer/httpbeat/coverage.svg?branch=master)](http://codecov.io/github/christiangalsterer/httpbeat?branch=master)

# Httpbeat

Httpbeat is the [Beat](https://www.elastic.co/products/beats) used to call HTTP endpoints.
Multiple endpoints can be configured which are polled in a regular interval and the result is shipped to the configured output channel.

Httpbeat is inspired by the Logstash [http_poller](https://www.elastic.co/guide/en/logstash/current/plugins-inputs-http_poller.html) input filter but doesn't require that the endpoint is reachable by Logstash as Httpbeat pushes the data to Logstash or Elasticsearch.
This is often necessary in security restricted network setups, where Logstash is not able to reach all servers. Instead the server to be monitored itself has Httpbeat installed and can send the data or a collector server has Httpbeat installed which is deployed in the secured network environment and can reach all servers to be monitored.

Example use cases are:
* Monitor [Apache Stats](https://httpd.apache.org/docs/2.4/mod/mod_status.html)
* Monitor Java application with [Jolokia](https://jolokia.org)
* Monitor [Spring Boot Actuators](http://docs.spring.io/spring-boot/docs/current/reference/htmlsingle/#production-ready)

## Exported Document Types

There is exactly one document type exported:

- `type: httpbeat` http request and response information

### httpbeat type

<pre>
{
  "_index": "httpbeat-2015.12.05",
  "_type": "httpbeat",
  "_source": {
    "@timestamp": "2015-12-05T11:16:13.070Z",
    "beat": {
      "hostname": "mbp.box",
      "name": "mbp.box"
    },
    "count": 1,
    "fields": {
      "host": "test"
    },
    "request": {
      "url": "http://httpbin.org/headers",
      "method": "get",
      "headers": {
        "Accept": "application/json",
        "Foo": "bar"
      }
    },
    "response": {
      "statusCode": 200,
      "headers": {
        "Access-Control-Allow-Credentials": "true",
        "Access-Control-Allow-Origin": "*",
        "Connection": "keep-alive",
        "Content-Length": "220",
        "Content-Type": "application/json",
        "Date": "Sat, 05 Dec 2015 11:16:13 GMT",
        "Server": "nginx"
      },
      "body": "{\n  \"headers\": {\n    \"Accept\": \"application/json\", \n    \"Accept-Encoding\": \"gzip\", \n    \"Authorization\": \"Basic Zm9vOmJhcg==\", \n    \"Foo\": \"bar\", \n    \"Host\": \"httpbin.org\", \n    \"User-Agent\": \"Go-http-client/1.1\"\n  }\n}\n"
    },
    "type": "httpbeat"
  },
  "fields": {
    "timestamp": [
      1449314173
    ]
  },
  "sort": [
    1449314173
  ]
}
</pre>


## Elasticsearch Template

To apply the Httpbeat template:

    curl -XPUT 'http://localhost:9200/_template/httpbeat' -d@etc/httpbeat.template.json

## Build, Test, Run

```
# Build
GOPATH=<your go path> make httpbeat

# Test
GOPATH=<your go path> make test

# Run
./httpbeat -c /etc/httpbeat/httpbeat.yml
```
## Contribution
All sorts of contributions are welcome. Please create a pull request and/or issue.
