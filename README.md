[![Build Status](https://travis-ci.org/christiangalsterer/httpbeat.svg?branch=master)](https://travis-ci.org/christiangalsterer/httpbeat)
[![codecov.io](http://codecov.io/github/christiangalsterer/httpbeat/coverage.svg?branch=master)](http://codecov.io/github/christiangalsterer/httpbeat?branch=master)

# Httpbeat

Httpbeat is the [Beat](https://www.elastic.co/products/beats) used to call HTTP endpoints.
Multiple endpoints can be configured which are polled in a regular interval and the result is shipped to the configured output channel.

Httpbeat is inspired by the Logstash [http_poller](https://www.elastic.co/guide/en/logstash/current/plugins-inputs-http_poller.html) input filter but doesn't require that the endpoint is reachable by Logstash as Httpbeat pushes the data to Logstash or Elasticsearch.
This is often necessary in security restricted network setups, where Logstash is not able to reach all servers. Instead the server to be monitored itself has Httpbeat installed and can send the data or a collector server has Httpbeat installed which is deployed in the secured network environment and can reach all servers to be monitored.

Examples are:
* [Apache Stats](https://httpd.apache.org/docs/2.4/mod/mod_status.html)
* [Jolokia](https://jolokia.org)
* [Spring Boot Actuator](http://docs.spring.io/spring-boot/docs/current/reference/htmlsingle/#production-ready)

## Exported Document Types

There is exactly one document type exported:

- `type: httpbeat` http request and response information

### httpbeat type

<pre>
{
  "_index": "httpbeat-2015.10.02",
  "_type": "container",
  "_id": "AVAow1NYKDyuAT4RG9KO",
  "_score": null,
  "_source": {
    "container": {
      "command": "/docker-entrypoint.sh kibana",
      "created": "2015-08-10T13:33:10Z",
      "id": "7e91fbb0c7885f55ef8bf9402bbe4b366f88224c8baf31d36265061aa5ba2735",
      "image": "kibana",
      "labels": {},
      "names": [
        "/kibana"
      ],
      "ports": [
        {
          "ip": "0.0.0.0",
          "privatePort": 5601,
          "publicPort": 5601,
          "type": "tcp"
        }
      ],
      "sizeRootFs": 0,
      "sizeRw": 0,
      "status": "Up 2 minutes"
    },
    "containerID": "7e91fbb0c7885f55ef8bf9402bbe4b366f88224c8baf31d36265061aa5ba2735",
    "containerNames": [
      "/kibana"
    ],
    "count": 1,
    "shipper": "0b42b9dded44",
    "timestamp": "2015-10-02T13:35:00.338Z",
    "type": "http"
  },
  "fields": {
    "timestamp": [
      1443792900338
    ]
  },
  "sort": [
    1443792900338
  ]
}
</pre>


## Elasticsearch Template

To apply Httpbeat template:

    curl -XPUT 'http://localhost:9200/_template/httpbeat' -d@etc/httpbeat.template.json
    
## Run Httpbeat

To launch Httpbeat, run the following command:
 
```
/etc/init.d/httpbeat -c /etc/httpbeat/httpbeat.yml
```
