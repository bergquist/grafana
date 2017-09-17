import {describe, beforeEach, it, sinon, expect, angularMocks} from 'test/lib/common';
import moment from 'moment';
import helpers from 'test/specs/helpers';
import {PrometheusDatasource} from '../datasource';
import {ResponseTransformer} from '../result_transformer';

describe('ResultTransformer', function() {
  var response;
  var transformer;
  var start = 1234567788;
  var end = 12334567776;

  beforeEach(() => {
    transformer = new ResponseTransformer(null);
  });

  describe.only('can transform to table', function() {
    response = [ {
        metric: {"__name__": "metric", job: "testjob", le: "10"},
        values: [ [1443454528.000, "10"], [1443454528.000, "10"] ],
      }, {
        metric: {"__name__": "metric", job: "testjob", le: "20"},
        values: [ [1443454528.000, "20"], [1443454528.000, "10"] ],
      }, {
        metric: {"__name__": "metric", job: "testjob", le: "30"},
        values: [ [1443454528.000, "25"], [1443454528.000, "10"] ],
      }
    ];
    var want = [ {
        target: "10",
        datapoints: [ [10, 1443454528000], [10, 1443454528000] ],
      }, {
        target: "20",
        datapoints: [ [10, 1443454528000], [0, 1443454528000] ],
      }, {
        target: "30",
        datapoints: [ [5, 1443454528000], [0, 1443454528000] ],
      }
    ];

    it('label_values(resource) should generate label search query', function() {
      var got = transformer.transFormToHistogramOverTime(response);

      expect(got).to.eql(want);
    });
  });
});
