package notifiers

import (
	"strings"

	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/models"

	"github.com/grafana/grafana/pkg/services/alerting"
)

func init() {
	alerting.RegisterNotifier(&alerting.NotifierPlugin{
		Type:        "prometheus_alertmanager",
		Name:        "Prometheus alertmanager",
		Description: "Send notifications to Prometheus alert manager",
		Factory:     NewPrometheusAlertmanagerNotifier,
		OptionsTemplate: `
      <h3 class="page-heading">Prometheus alertmanager servers</h3>
      <div class="gf-form">
         <textarea rows="7" class="gf-form-input width-25" required ng-model="ctrl.model.settings.servers"></textarea>
      </div>
      <div class="gf-form">
      <span>You can enter multiple servers addresses using a ";" separator</span>
    `,
	})
}

func NewPrometheusAlertmanagerNotifier(model *models.AlertNotification) (alerting.Notifier, error) {
	serverString := model.Settings.Get("servers").MustString()

	if serverString == "" {
		return nil, alerting.ValidationError{Reason: "Could not find servers in settings"}
	}

	// split addresses with a few different ways
	servers := strings.FieldsFunc(serverString, func(r rune) bool {
		switch r {
		case ',', ';', '\n':
			return true
		}
		return false
	})

	return &prometheusAlertmanagerNotifier{
		NotifierBase: NewNotifierBase(model.Id, model.IsDefault, model.Name, model.Type, model.Settings),
		log:          log.New("alerting.notifier.alertmanager"),
		servers:      servers,
	}, nil
}

type prometheusAlertmanagerNotifier struct {
	NotifierBase
	log     log.Logger
	servers []string
}

func (this *prometheusAlertmanagerNotifier) Notify(evalContext *alerting.EvalContext) error {
	return nil
}
