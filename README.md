[![Build Status](https://travis-ci.org/Ingensi/dockerbeat.svg?branch=develop)](https://travis-ci.org/Ingensi/dockerbeat)
[![codecov.io](http://codecov.io/github/Ingensi/dockerbeat/coverage.svg?branch=develop)](http://codecov.io/github/Ingensi/dockerbeat?branch=develop)

# Httpbeat

Httpbeat is the [Beat](https://www.elastic.co/products/beats) used to call HTTP endpoints.

## Exported Document Types

There is exactly one document type exported:

- `type: http` for container attributes

### http type

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
