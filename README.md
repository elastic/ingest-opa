# HOW TO RUN

On this repo:
```
cd ~/workspace/javaopa/ingest-opa
cd src/main/java/org/elasticsearch/plugin/ingest/opa
go build -o awesome.so -buildmode=c-shared awesome.go
mkdir -p /Users/eyal/Library/Frameworks/./awesome.so.framework/./
cp ~/workspace/javaopa/ingest-opa/src/main/java/org/elasticsearch/plugin/ingest/opa/awesome.so /Users/eyal/Library/Frameworks/./awesome.so.framework/./
# cp ~/workspace/javaopa/ingest-opa/src/main/java/org/elasticsearch/plugin/ingest/opa/awesome.h /Users/eyal/Library/Frameworks/./awesome.so.framework/./

cd ~/workspace/javaopa/ingest-opa
gradle clean check
```

On your Kibana repo:
```
cd ~/Projects/kibana
yarn es snapshot --license trial --version 7.16.3
```

After the relevant version of elasticsearch has beed downloaded and started stop it.

(Note: this is the required version of elasticsearch and kibana because of `gradle.properties`: `elasticsearchVersion = 7.16.3` which fail when changing to 8.1.0)

```
cd .es/7.16.3
bin/elasticsearch-plugin install file:///Users/eyal/workspace/javaopa/ingest-opa/build/distributions/ingest-opa-0.0.1-SNAPSHOT.zip
# approve with y when asked
# if already installed run:
bin/elasticsearch-plugin remove ingest-opa
# note to self: edit config files?
```

create here the file `custom-elasticsearch.policy`:
```
grant {
  permission java.io.FilePermission "/usr/lib/lib./awesome.so.dylib", "read";
  permission java.io.FilePermission "/usr/lib/lib./awesome.so.jnilib", "read";
  permission java.io.FilePermission "/Users/eyal/Library/Frameworks/./awesome.so.framework/./awesome.so", "read";
  permission java.io.FilePermission "/Library/Frameworks/./awesome.so.framework/./awesome.so", "read";
  permission java.io.FilePermission "/System/Library/Frameworks/./awesome.so.framework/./awesome.so", "read";
  permission java.lang.reflect.ReflectPermission "OpaProcessor.org.elasticsearch.plugin.ingest.opa";
  permission java.lang.RuntimePermission "accessDeclaredMembers";
};
```

now start elasticsearch using it:
```
export ES_JAVA_OPTS=-Djava.security.policy=file:/Users/eyal/Projects/kibana/.es/7.16.3/custom-elasticsearch.policy
bin/elasticsearch
```

Now you should start kibana with version 7.16.3 to connect to this elasticsearch.

I created a new profile with elastic-stack `eyal2`, reconfigured the file `profiles/eyal2/stack/kibana.config.default.yml` to contain `elasticsearch.hosts: [ "http://host.docker.internal:9200" ]` and ran:
```
elastic-package stack up -p eyal2 -d --version=7.16.3-SNAPSHOT -s kibana
```

Now on the developer console of kibana run these by order:

```
GET _nodes/ingest?filter_path=nodes.*.ingest.processors

PUT _ingest/pipeline/opa-pipeline
{
  "description": "A pipeline to do whatever",
  "processors": [
    {
      "opa" : {
        "field" : "resource",
        "target_field" : "finding"
      }
    }
  ]
}

PUT /my-index/my-type/1?pipeline=opa-pipeline
{
  "resource" : "Some content"
}

GET /my-index/my-type/1
```


### Relevant Links:

#### Running Go from Java:

https://medium.com/learning-the-go-programming-language/calling-go-functions-from-other-languages-4c7d8bcc69bf

https://github.com/java-native-access/jna#:~:text=Elasticsearch%3A%20Large%2Dscale%20distributed%20search%20and%20analytics%20engine - JNA Already used by Elasticsearch to run native code

https://en.wikipedia.org/wiki/Java_Native_Access

https://github.com/java-native-access/jna/blob/master/www/GettingStarted.md

https://stackoverflow.com/questions/49986729/how-can-i-call-a-go-function-from-java-using-the-java-native-interface

https://stackoverflow.com/questions/49547293/jna-to-go-dll-how-do-i-get-string-returned-from-go-func

https://github.com/jbuberel/buildmodeshared

#### Elasticsearch Plugins

https://www.elastic.co/guide/en/elasticsearch/plugins/master/index.html

https://www.elastic.co/guide/en/elasticsearch/plugins/master/ingest.html - Ingest specifically

https://www.elastic.co/guide/en/elasticsearch/plugins/master/integrations.html - Example plugins

https://www.elastic.co/guide/en/elasticsearch/plugins/master/installation.html - Installing

https://www.elastic.co/guide/en/elasticsearch/plugins/master/plugin-management-custom-url.html - Installing local

https://www.elastic.co/guide/en/elasticsearch/plugins/master/manage-plugins-using-configuration-file.html - install with config (for future use)

https://github.com/spinscale/cookiecutter-elasticsearch-ingest-processor - Injest plugin generator (I used it to generate the plugin)


#### Elastic Plugins using JNA (Running native code) / Java security related (This was something to tackle...)

https://github.com/EvidentSolutions/elasticsearch-analysis-voikko/blob/master/README.md - GOLD! an example

https://www.elastic.co/guide/en/elasticsearch/plugins/current/plugin-authors.html#plugin-authors-jsm - java security permissions

https://stackoverflow.com/a/35402714

https://discuss.elastic.co/t/es-5-0-security-policy-for-plugin-using-jna-to-load-libs/67278 

https://discuss.elastic.co/u/ninesalt - so many unanswered questions

https://www.elastic.co/guide/en/elasticsearch/reference/current/executable-jna-tmpdir.html

https://github.com/opendistro-for-elasticsearch/k-NN/blob/c2ac595ce5e8999878de05251e14dae6f59de9fb/jni/src/com_amazon_opendistroforelasticsearch_knn_index_v2011_KNNIndex.cpp - CPP example (old?)

https://discuss.elastic.co/t/questions-about-non-java-c-elasticsearch-native-clients/50719

https://github.com/elastic/elasticsearch/issues/62143

#### OPA API

https://www.openpolicyagent.org/docs/latest/integration/#integrating-with-the-go-api

https://pkg.go.dev/github.com/open-policy-agent/opa/rego#example-Rego.Eval-Input


-----------------------


# Elasticsearch opa Ingest Processor

Explain the use case of this processor in a TLDR fashion.

## Usage


```
PUT _ingest/pipeline/opa-pipeline
{
  "description": "A pipeline to do whatever",
  "processors": [
    {
      "opa" : {
        "field" : "my_field"
      }
    }
  ]
}

PUT /my-index/my-type/1?pipeline=opa-pipeline
{
  "my_field" : "Some content"
}

GET /my-index/my-type/1
{
  "my_field" : "Some content"
  "potentially_enriched_field": "potentially_enriched_value"
}
```

## Configuration

| Parameter | Use |
| --- | --- |
| some.setting   | Configure x |
| other.setting  | Configure y |

## Setup

In order to install this plugin, you need to create a zip distribution first by running

```bash
gradle clean check
```

This will produce a zip file in `build/distributions`.

After building the zip file, you can install it like this

```bash
bin/elasticsearch-plugin install file:///path/to/ingest-opa/build/distribution/ingest-opa-0.0.1-SNAPSHOT.zip
```

## Bugs & TODO

* There are always bugs
* and todos...

