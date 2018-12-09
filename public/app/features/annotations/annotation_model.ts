import _ from 'lodash';
import { assignModelProperties } from 'app/core/utils/model_utils';

export class AnnotationModel {
  name: string;
  datasource: any;
  iconColor: string;
  enable: boolean;
  showIn: number;
  hide: boolean;

  defaults: any = {
    name: '',
    datasource: null,
    iconColor: 'rgba(255, 96, 96, 1)',
    enable: true,
    showIn: 0,
    hide: false,
  };

  constructor(private model) {
    assignModelProperties(this, model, this.defaults);
  }

  getSaveModel() {
    assignModelProperties(this.model, this, this.defaults);
    return this.model;
  }
}
