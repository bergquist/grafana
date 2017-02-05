package api

import (
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/middleware"
)

func GetTemplate(c *middleware.Context) Response {
	id := c.ParamsInt(":id")
	if id == 0 {
		return ApiError(400, "missing template id", nil)
	}

	json := `
  {
    "allValue": null,
    "current": {
      "text": "desktop",
      "value": "desktop"
    },
    "datasource": "graphite",
    "hide": 0,
    "includeAll": false,
    "label": null,
    "multi": false,
    "name": "global template",
    "options": [
      {
        "selected": true,
        "text": "desktop",
        "value": "desktop"
      },
      {
        "selected": false,
        "text": "mobile",
        "value": "mobile"
      },
      {
        "selected": false,
        "text": "tablet",
        "value": "tablet"
      }
    ],
    "query": "statsd.fakesite.counters.session_start.*",
    "refresh": 0,
    "regex": "",
    "sort": 0,
    "tagValuesQuery": "",
    "tags": [],
    "tagsQuery": "",
    "type": "query",
    "useTags": false
  }`

	jsonObj, _ := simplejson.NewJson([]byte(json))
	return Json(200, jsonObj)
}
