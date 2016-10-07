[![Build Status](https://travis-ci.org/christiangalsterer/httpbeat.svg?branch=master)](https://travis-ci.org/christiangalsterer/httpbeat)
[![codecov.io](http://codecov.io/github/christiangalsterer/httpbeat/coverage.svg?branch=master)](http://codecov.io/github/christiangalsterer/httpbeat?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/christiangalsterer/httpbeat)](https://goreportcard.com/report/github.com/christiangalsterer/httpbeat)

# Httpbeat

Welcome to Httpbeat.

Httpbeat is a [Beat](https://www.elastic.co/products/beats) used to call HTTP endpoints.
Multiple endpoints can be configured which are polled in a regular interval and the result is shipped to the configured output channel.

Httpbeat is inspired by the Logstash [http_poller](https://www.elastic.co/guide/en/logstash/current/plugins-inputs-http_poller.html) input filter but doesn't require that the endpoint is reachable by Logstash as Httpbeat pushes the data to Logstash or Elasticsearch.
This is often necessary in security restricted network setups, where Logstash is not able to reach all servers. Instead the server to be monitored itself has Httpbeat installed and can send the data or a collector server has Httpbeat installed which is deployed in the secured network environment and can reach all servers to be monitored.

Example use cases are:
* Monitor [Apache Stats](https://httpd.apache.org/docs/2.4/mod/mod_status.html)
* Monitor Java application with [Jolokia](https://jolokia.org)
* Monitor [Spring Boot Actuators](http://docs.spring.io/spring-boot/docs/current/reference/htmlsingle/#production-ready)

Ensure that this folder is at the following location:
`${GOPATH}/github.com/christiangalsterer`

## Getting Started with Httpbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7.1
* [Glide](https://github.com/Masterminds/glide) >= 0.11.0

### Build

To build the binary for httpbeat run the command below. This will generate a binary
in the same directory with the name httpbeat.

```
make
```


### Run

To run httpbeat with debugging output enabled, run:

```
./httpbeat -c httpbeat.yml -e -d "*"
```


### Test

To test httpbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/httpbeat.template.json and etc/httpbeat.asciidoc

```
make update
```


### Cleanup

To clean httpbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone httpbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/christiangalsterer
cd ${GOPATH}/github.com/christiangalsterer
git clone https://github.com/christiangalsterer/httpbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.

# Releases

2.0.0-beta.1 (2016-10-XX) Work in Progress

Feature release containing the following changes:
* Update to beats v5.0.0-beta1

Please note that this release contains the following breaking changes introduced by beats 5.0.X, see also [Beats Changelog](https://github.com/elastic/beats/blob/v5.0.0-beta1/CHANGELOG.asciidoc)
* SSL Configuration
    * rename tls configurations section to ssl
    * rename certificate_key configuration to key.
    * replace tls.insecure with ssl.verification_mode setting.
    * replace tls.min/max_version with ssl.supported_protocols setting requiring full protocol name

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
```
curl -XPUT 'http://localhost:9200/_template/httpbeat' -d@etc/httpbeat.template.json
```
    
# Contribution
All sorts of contributions are welcome. Please create a pull request and/or issue.
