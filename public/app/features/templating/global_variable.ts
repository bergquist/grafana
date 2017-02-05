///<reference path="../../headers/common.d.ts" />

import _ from 'lodash';
import kbn from 'app/core/utils/kbn';
import {Variable, assignModelProperties, variableTypes} from './variable';
import {VariableSrv} from './variable_srv';

var mockGlobalVariable = {
  "allValue": null,
  "current": {
    "text": "desktop",
    "value": "desktop"
  },
  "datasource": "graphite",
  "hide": 0,
  "includeAll": false,
  "label": null,
  "multi": true,
  "name": "globalDevice",
  "options": [
    //{
    //  "selected": true,
    //  "text": "desktop",
    //  "value": "desktop"
    //},
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
};

export class GlobalVariable implements Variable {
  query: string;
  options: any;
  includeAll: boolean;
  multi: boolean;
  current: any;

  defaults = {
    type: 'global',
    name: '',
    hide: 0,
  };

  globalTemplate: any;

  /** @ngInject **/
  constructor(private model, private timeSrv, private templateSrv, private variableSrv, private backendSrv) {
    //should fetch values from /templates global and init
    console.log('creating global variable', model);
    console.log('using mock data', mockGlobalVariable);

    //this.backendSrv.get("/templates/100")
    //  then()
    this.globalTemplate = variableSrv.createVariableFromModel(mockGlobalVariable);
    this.globalTemplate.type = 'global';
    return this.globalTemplate;
    //assignModelProperties(this, model, this.defaults);
  }

  setValue(option) {
    return this.variableSrv.setOptionAsCurrent(this, option);
  }

  getSaveModel() {
    assignModelProperties(this.model, this, this.defaults);
    return this.model;
  }

  updateOptions() {
    // extract options in comma separated string
    this.options = _.map(this.query.split(/[,]+/), function(text) {
      return { text: text.trim(), value: text.trim() };
    });

    if (this.includeAll) {
      this.addAllOption();
    }

    return this.variableSrv.validateVariableSelectionState(this);
  }

  //TODO REMOVE OR THROW
  addAllOption() {
    this.options.unshift({text: 'All', value: "$__all"});
  }

  dependsOn(variable) {
    return false; //TODO
  }

  setValueFromUrl(urlValue) {
    return this.variableSrv.setOptionFromUrl(this, urlValue);
  }

  getValueForUrl() {
    if (this.current.text === 'All') {
      return 'All';
    }
    return this.current.value;
  }
}

variableTypes['global'] = {
  name: 'Global',
  ctor: GlobalVariable,
  description: 'Global template variable' ,
  supportsMulti: true,
};
