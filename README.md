[![Build Status](https://travis-ci.org/christiangalsterer/httpbeat.svg?branch=master)](https://travis-ci.org/christiangalsterer/httpbeat)
[![codecov.io](http://codecov.io/github/christiangalsterer/httpbeat/coverage.svg?branch=master)](http://codecov.io/github/christiangalsterer/httpbeat?branch=master)

# Overview

Httpbeat is the [Beat](https://www.elastic.co/products/beats) used to call HTTP endpoints.
Multiple endpoints can be configured which are polled in a regular interval and the result is shipped to the configured output channel.

Httpbeat is inspired by the Logstash [http_poller](https://www.elastic.co/guide/en/logstash/current/plugins-inputs-http_poller.html) input filter but doesn't require that the endpoint is reachable by Logstash as Httpbeat pushes the data to Logstash or Elasticsearch.
This is often necessary in security restricted network setups, where Logstash is not able to reach all servers. Instead the server to be monitored itself has Httpbeat installed and can send the data or a collector server has Httpbeat installed which is deployed in the secured network environment and can reach all servers to be monitored.

Example use cases are:
* Monitor [Apache Stats](https://httpd.apache.org/docs/2.4/mod/mod_status.html)
* Monitor Java application with [Jolokia](https://jolokia.org)
* Monitor [Spring Boot Actuators](http://docs.spring.io/spring-boot/docs/current/reference/htmlsingle/#production-ready)

# Releases

2.0.0-alpha.5 (2016-08-XX) Work in Progress

Feature release containing the following changes:
* Update to beats v5.0.0-alpha5

1.2.0 (2016-07-19)

Feature release containing the following changes:
* Update to Go 1.6
* Update to libbeat 1.2.3
* Use [Glide](https://github.com/Masterminds/glide) for dependency management

1.1.0 (2016-04-06)

Feature release containing the following changes:
* [Provide output directly as JSON](https://github.com/christiangalsterer/httpbeat/issues/2)

1.0.1 (2016-02-17)

Bugfix release containing the following changes:
* Fix: [Infinite loop when using logstash output](https://github.com/christiangalsterer/httpbeat/issues/4)
* Fix: [Hanging during shutdown](https://github.com/christiangalsterer/httpbeat/issues/5)

1.0.0 (2015-12-29)
* Initial release

# Configuration

## Configuration Options

See [here](docs/configuration.asciidoc) for more information.

## Exported Document Types

There is exactly one document type exported:

- `type: httpbeat` http request and response information

## Exported Fields

See [here](docs/fields.asciidoc) for a detailed description of all exported fields.

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

# Build, Test, Run

```
# Build
GOPATH=<your go path> make httpbeat

# Test
GOPATH=<your go path> make test

# Run
./httpbeat -c /etc/httpbeat/httpbeat.yml
```
# Contribution
All sorts of contributions are welcome. Please create a pull request and/or issue.
