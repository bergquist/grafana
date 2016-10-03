package elasticsearch

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestElastic(t *testing.T) {
	Convey("Elastic", t, func() {

		Convey("get indices", func() {
			ds := &ElasticDatasource{
				Index:    "[metrics-]YYYY-MM-DD",
				Interval: "Daily",
			}
			start := time.Date(2016, 10, 7, 0, 0, 0, 0, time.Local)
			end := time.Date(2016, 10, 10, 0, 0, 0, 0, time.Local)

			indices := ds.GetIndices(start, end)

			So(indices[0], ShouldEqual, "metrics-2016.10.7")
			//So(indices[1], ShouldEqual, "metrics-2016.10.8")
			//So(indices[2], ShouldEqual, "metrics-2016.10.9")
		})
	})
}
