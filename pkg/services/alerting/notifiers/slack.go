package notifiers

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"

	"net/http"

	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
)

func init() {
	alerting.RegisterNotifier(&alerting.NotifierPlugin{
		Type:        "slack",
		Name:        "Slack",
		Description: "Sends notifications using Grafana server configured STMP settings",
		Factory:     NewSlackNotifier,
		OptionsTemplate: `
      <h3 class="page-heading">Slack settings</h3>
      <div class="gf-form max-width-30">
        <span class="gf-form-label width-6">Url</span>
        <input type="text" required class="gf-form-input max-width-30" ng-model="ctrl.model.settings.url" placeholder="Slack incoming webhook url"></input>
      </div>
      <div class="gf-form max-width-30">
        <span class="gf-form-label width-6">Recipient</span>
        <input type="text"
          class="gf-form-input max-width-30"
          ng-model="ctrl.model.settings.recipient"
          data-placement="right">
        </input>
        <info-popover mode="right-absolute">
          Override default channel or user, use #channel-name or @username
        </info-popover>
      </div>
      <div class="gf-form max-width-30">
        <span class="gf-form-label width-6">Mention</span>
        <input type="text"
          class="gf-form-input max-width-30"
          ng-model="ctrl.model.settings.mention"
          data-placement="right">
        </input>
        <info-popover mode="right-absolute">
          Mention a user or a group using @ when notifying in a channel
        </info-popover>
      </div>
    `,
	})

}

func NewSlackNotifier(model *m.AlertNotification) (alerting.Notifier, error) {
	url := model.Settings.Get("url").MustString()
	if url == "" {
		return nil, alerting.ValidationError{Reason: "Could not find url property in settings"}
	}

	recipient := model.Settings.Get("recipient").MustString()
	mention := model.Settings.Get("mention").MustString()

	return &SlackNotifier{
		NotifierBase: NewNotifierBase(model.Id, model.IsDefault, model.Name, model.Type, model.Settings),
		Url:          url,
		Recipient:    recipient,
		Mention:      mention,
		log:          log.New("alerting.notifier.slack"),
	}, nil
}

type SlackNotifier struct {
	NotifierBase
	Url       string
	Recipient string
	Mention   string
	log       log.Logger
}

func (this *SlackNotifier) Notify(evalContext *alerting.EvalContext) error {
	this.log.Info("Executing slack notification", "ruleId", evalContext.Rule.Id, "notification", this.Name)
	// metrics.M_Alerting_Notification_Sent_Slack.Inc(1)

	// ruleUrl, err := evalContext.GetRuleUrl()
	// if err != nil {
	// 	this.log.Error("Failed get rule link", "error", err)
	// 	return err
	// }

	// fields := make([]map[string]interface{}, 0)
	// fieldLimitCount := 4
	// for index, evt := range evalContext.EvalMatches {
	// 	fields = append(fields, map[string]interface{}{
	// 		"title": evt.Metric,
	// 		"value": evt.Value,
	// 		"short": true,
	// 	})
	// 	if index > fieldLimitCount {
	// 		break
	// 	}
	// }

	// if evalContext.Error != nil {
	// 	fields = append(fields, map[string]interface{}{
	// 		"title": "Error message",
	// 		"value": evalContext.Error.Error(),
	// 		"short": false,
	// 	})
	// }

	// message := this.Mention
	// if evalContext.Rule.State != m.AlertStateOK { //dont add message when going back to alert state ok.
	// 	message += " " + evalContext.Rule.Message
	// }

	// body := map[string]interface{}{
	// 	"attachments": []map[string]interface{}{
	// 		{
	// 			"color":       evalContext.GetStateModel().Color,
	// 			"title":       evalContext.GetNotificationTitle(),
	// 			"title_link":  ruleUrl,
	// 			"text":        message,
	// 			"fields":      fields,
	// 			"image_url":   evalContext.ImagePublicUrl,
	// 			"footer":      "Grafana v" + setting.BuildVersion,
	// 			"footer_icon": "http://grafana.org/assets/img/fav32.png",
	// 			"ts":          time.Now().Unix(),
	// 		},
	// 	},
	// 	"parse": "full", // to linkify urls, users and channels in alert message.
	// }

	// //recipient override
	// if this.Recipient != "" {
	// 	body["channel"] = this.Recipient
	// }

	// data, _ := json.Marshal(&body)
	// cmd := &m.SendWebhookSync{Url: this.Url, Body: string(data)}

	// if err := bus.DispatchCtx(evalContext.Ctx, cmd); err != nil {
	// 	this.log.Error("Failed to send slack notification", "error", err, "webhook", this.Name)
	// 	return err
	// }

	//if evalContext.ImagePublicUrl == "" && this.NotifierBase.NeedsImage() {
	if this.NotifierBase.NeedsImage() {
		/*
						    jsonPayload := simplejson.New()

			          jsonPayload.Set("filename", "image.png")
			          jsonPayload.Set("filetype", "png")
			          jsonPayload.Set("channels", "#grafana")
			          jsonPayload.Set("title", evalContext.GetNotificationTitle())
			          jsonPayload.Set("token", "T02S4RCS0/B06AGLK5H/L5xITbLWG2eVTRw4jsDP7AD9")
		*/
		//jsonPayload.Set("initial_comment", "asdf")

		//content, _ := ioutil.ReadFile(evalContext.ImageOnDiskPath)
		//jsonPayload.Set("content", content)

		//body, _ := jsonPayload.MarshalJSON()

		this.log.Info("Trying to upload image to slack!", "image", evalContext.ImageOnDiskPath)
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		w.WriteField("filename", "image.png")
		w.WriteField("filetype", "png")
		w.WriteField("channels", "#grafana")
		w.WriteField("title", evalContext.GetNotificationTitle())
		w.WriteField("token", "T02S4RCS0/B06AGLK5H/L5xITbLWG2eVTRw4jsDP7AD9")

		// Add your image file
		f, err := os.Open(evalContext.ImageOnDiskPath)
		if err != nil {
			this.log.Error("failed to send image", "error", err)
			return err
		}
		defer f.Close()
		fw, err := w.CreateFormFile("file", evalContext.ImageOnDiskPath)
		if err != nil {
			return err
		}
		if _, err = io.Copy(fw, f); err != nil {
			return err
		}
		// Don't forget to close the multipart writer.
		// If you don't close it, your request will be missing the terminating boundary.
		//res, err := http.DefaultClient.Post("https://slack.com/api/files.upload", "application/json", bytes.NewReader(body))
		req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", &b)
		if err != nil {
			return err
		}
		// Don't forget to set the content type, this will contain the boundary.
		req.Header.Set("Content-Type", w.FormDataContentType())
		w.Close()

		//res, err := http.DefaultClient.PostForm("https://slack.com/api/files.upload", &b)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			this.log.Error("failed to send image", "error", err)
		} else {
			this.log.Info("Success?!", "response", res)
		}

	}

	return nil
}
