+++
title = "Developing Backend Plugins v2"
keywords = ["grafana", "plugins", "backend", "plugin", "backend-plugins", "documentation"]
type = "docs"
[menu.docs]
name = "Developing Backend Plugins"
parent = "developing"
weight = 5
status = "draft"
+++

<!--
 This will replace backend-plugins-guide.md
 But Im writing this as a new doc to make the PR cleaner.

 Relevant links
 * Kyles tips for writing BE https://docs.google.com/document/d/1dg5kq9DqKZNTptse6wk_EobQCtVDEpjBb-mId0coUio/edit#heading=h.smzt5e3i1cq
* https://pkg.go.dev/github.com/grafana/grafana-plugin-sdk-go/data?tab=doc
* 
-->

# Backend Plugins

> In Grafana 7.0 we released a new version of backend plugins. Docs for previous versions of backend plugins can be found in [Grafana 6.7 docs](https://grafana.com/docs/grafana/v6.7/plugins/developing/backend-plugins-guide/).

Grafana added support for plugins in Grafana 3.0 and this enabled the Grafana community to create panel plugins and data source plugins. It was wildly successful and has made Grafana much more useful as you can integrate it with anything and do any type of custom visualization that you want. However, these plugin hooks are on the frontend only and we also want to provide hooks into the Grafana backend to allow the community to extend and improve Grafana in new ways.

Grafana have been supporting backend plugins for a while but we weren't really happy with them. In Grafana 7.0 we released a new version of backend plugins that we are much more comfortable with supporting for a longer period of time. The three main differences in the new version of backend plugins are:
We now offer an Grafana plugin SDK written in Go that removes all the wiring and plumbing when building Grafana backend plugins. 
We introduced a new endpoint called Resources that plugins can implement. Resource requests are untypes requests that are passed through Grafanas backend to the backend plugins without knowing anything about the structure of the data. Grafana adds extra contextual information to the requests without looking at the actual body. This allows plugin developers to send custom requests to their backend service without Grafana knowing about its shape. 
[Dataframes]() is the new data structure in Grafana for passing data between datasources and and panels. This data structure can hold both time series data and tabular data in different shapes. 

### Grafana's Backend Plugin System

The backend plugin feature is implemented with the [HashiCorp plugin system](https://github.com/hashicorp/go-plugin) which is a Go plugin system over RPC. Grafana server launches each plugin as a subprocess and communicates with it over RPC. This approach has a number of benefits:

- Plugins can't crash your grafana process: a panic in a plugin doesn't panic the server.
- Plugins are easy to develop: just write a Go application and `go build` 
- Plugins can be relatively secure: The plugin only has access to the interfaces and args given to it, not to the entire memory space of the process.

<!-- link to tutorial -->

### Datasource backend 

<!-- what does datasource backend plugin mean? -->

### Resource requests
<!-- What's the purpose of Resource requests + Examples use-cases? -->

Example 1
Example 2
Example 3


### Grafana plugin SDK
 Link: https://pkg.go.dev/github.com/grafana/grafana-plugin-sdk-go/backend?tab=doc

### Backend plugins with support for Alerting

<!-- 
  Whats required for a plugin to work with Grafanas alerting?
    * "alerting": true
    * return data that can be turned into timeseries
-->




