///<reference path="../../../headers/common.d.ts" />

export class InfluxConfigCtrl {
  static templateUrl = 'partials/config.html';
  current: any;

  /** @ngInject **/
  constructor() {
    this.current.jsonData = this.current.jsonData || {};
    this.current.jsonData.influxDBVersion = this.current.jsonData.influxDBVersion || 1;
  }

  influxDBVersions = [
    {name: "0.9.x - 0.11.1", value: 1},
    {name: "0.12.0", value: 2}
  ];
}
