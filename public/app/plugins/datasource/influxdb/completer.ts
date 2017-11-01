///<reference path="../../../headers/common.d.ts" />

export class InfluxdbCompleter {
  labelQueryCache: any;
  labelNameCache: any;
  labelValueCache: any;

  identifierRegexps = [/\[/, /[a-zA-Z0-9_:]/];

  constructor() {
    this.labelQueryCache = {};
    this.labelNameCache = {};
    this.labelValueCache = {};
  }

  getCompletions(editor, session, pos, prefix, callback) {
    return {
        caption: "name",
        value: "value",
        meta: 'metric',
    };
  }

//   getLabelNameAndValueForMetric(metricName) {
//     if (this.labelQueryCache[metricName]) {
//       return Promise.resolve(this.labelQueryCache[metricName]);
//     }
//     var op = '=~';
//     if (/[a-zA-Z_:][a-zA-Z0-9_:]*/.test(metricName)) {
//       op = '=';
//     }
//     var expr = '{__name__' + op + '"' + metricName + '"}';
//     return this.datasource.performInstantQuery({ expr: expr }, new Date().getTime() / 1000).then(response => {
//       this.labelQueryCache[metricName] = response.data.data.result;
//       return response.data.data.result;
//     });
//   }

//   transformToCompletions(words, meta) {
//     return words.map(name => {
//       return {
//         caption: name,
//         value: name,
//         meta: meta,
//         score: Number.MAX_VALUE
//       };
//     });
//   }

//   findMetricName(session, row, column) {
//     var metricName = '';

//     var tokens;
//     var nameLabelNameToken = this.findToken(session, row, column, 'entity.name.tag', '__name__', 'paren.lparen');
//     if (nameLabelNameToken) {
//       tokens = session.getTokens(nameLabelNameToken.row);
//       var nameLabelValueToken = tokens[nameLabelNameToken.index + 2];
//       if (nameLabelValueToken && nameLabelValueToken.type === 'string.quoted') {
//         metricName = nameLabelValueToken.value.slice(1, -1); // cut begin/end quotation
//       }
//     } else {
//       var metricNameToken = this.findToken(session, row, column, 'identifier', null, null);
//       if (metricNameToken) {
//         tokens = session.getTokens(metricNameToken.row);
//         if (tokens[metricNameToken.index + 1].type === 'paren.lparen') {
//           metricName = metricNameToken.value;
//         }
//       }
//     }

//     return metricName;
//   }

  findToken(session, row, column, target, value, guard) {
    return null;
  }
}
