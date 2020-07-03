package annotations

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
)

// Service is responsible for storing annotations in the database
type Service struct {
	log      log.Logger
	Cfg      *setting.Cfg       `inject:""`
	SQLStore *sqlstore.SqlStore `inject:""`
}

func init() {
	registry.RegisterService(&Service{})
}

// Init initilizes Service
func (as *Service) Init() error {
	as.log = log.New("annotation")
	return nil
}

func (as *Service) CleanupAnnotations(ctx context.Context) error {
	return nil
}

// Update the item so that EpochEnd >= Epoch
func validateTimeRange(item *Item) error {
	if item.EpochEnd == 0 {
		if item.Epoch == 0 {
			return errors.New("Missing Time Range")
		}
		item.EpochEnd = item.Epoch
	}
	if item.Epoch == 0 {
		item.Epoch = item.EpochEnd
	}
	if item.EpochEnd < item.Epoch {
		tmp := item.Epoch
		item.Epoch = item.EpochEnd
		item.EpochEnd = tmp
	}
	return nil
}

// Save persit an annotation in the database
func (as *Service) Save(item *Item) error {
	return as.SQLStore.WithTransactionalDbSession(context.Background(), func(sess *sqlstore.DBSession) error {
		tags := models.ParseTagPairs(item.Tags)
		item.Tags = models.JoinTagPairs(tags)
		item.Created = time.Now().UnixNano() / int64(time.Millisecond)
		item.Updated = item.Created
		if item.Epoch == 0 {
			item.Epoch = item.Created
		}
		if err := validateTimeRange(item); err != nil {
			return err
		}

		if _, err := sess.Table("annotation").Insert(item); err != nil {
			return err
		}

		if item.Tags != nil {
			tags, err := sqlstore.EnsureTagsExist(sess, tags)
			if err != nil {
				return err
			}
			for _, tag := range tags {
				if _, err := sess.Exec("INSERT INTO annotation_tag (annotation_id, tag_id) VALUES(?,?)", item.Id, tag.Id); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Update updates an annotation in the database
func (as *Service) Update(item *Item) error {
	return as.SQLStore.WithTransactionalDbSession(context.Background(), func(sess *sqlstore.DBSession) error {
		var (
			isExist bool
			err     error
		)
		existing := new(Item)

		isExist, err = sess.Table("annotation").Where("id=? AND org_id=?", item.Id, item.OrgId).Get(existing)

		if err != nil {
			return err
		}
		if !isExist {
			return errors.New("Annotation not found")
		}

		existing.Updated = time.Now().UnixNano() / int64(time.Millisecond)
		existing.Text = item.Text

		if item.Epoch != 0 {
			existing.Epoch = item.Epoch
		}
		if item.EpochEnd != 0 {
			existing.EpochEnd = item.EpochEnd
		}

		if err := validateTimeRange(existing); err != nil {
			return err
		}

		if item.Tags != nil {
			tags, err := sqlstore.EnsureTagsExist(sess, models.ParseTagPairs(item.Tags))
			if err != nil {
				return err
			}
			if _, err := sess.Exec("DELETE FROM annotation_tag WHERE annotation_id = ?", existing.Id); err != nil {
				return err
			}
			for _, tag := range tags {
				if _, err := sess.Exec("INSERT INTO annotation_tag (annotation_id, tag_id) VALUES(?,?)", existing.Id, tag.Id); err != nil {
					return err
				}
			}
		}

		existing.Tags = item.Tags

		_, err = sess.Table("annotation").ID(existing.Id).Cols("epoch", "text", "epoch_end", "updated", "tags").Update(existing)
		return err
	})
}

// Find searches for annotation in the database
func (as *Service) Find(query *ItemQuery) ([]*ItemDTO, error) {
	var sql bytes.Buffer
	params := make([]interface{}, 0)
	dialect := as.SQLStore.Dialect

	sql.WriteString(`
		SELECT
			annotation.id,
			annotation.epoch as time,
			annotation.epoch_end as time_end,
			annotation.dashboard_id,
			annotation.panel_id,
			annotation.new_state,
			annotation.prev_state,
			annotation.alert_id,
			annotation.text,
			annotation.tags,
			annotation.data,
			annotation.created,
			annotation.updated,
			usr.email,
			usr.login,
			alert.name as alert_name
		FROM annotation
		LEFT OUTER JOIN ` + dialect.Quote("user") + ` as usr on usr.id = annotation.user_id
		LEFT OUTER JOIN alert on alert.id = annotation.alert_id
		INNER JOIN (
			SELECT a.id from annotation a
		`)

	sql.WriteString(`WHERE a.org_id = ?`)
	params = append(params, query.OrgId)

	if query.AnnotationId != 0 {
		// fmt.Print("annotation query")
		sql.WriteString(` AND a.id = ?`)
		params = append(params, query.AnnotationId)
	}

	if query.AlertId != 0 {
		sql.WriteString(` AND a.alert_id = ?`)
		params = append(params, query.AlertId)
	}

	if query.DashboardId != 0 {
		sql.WriteString(` AND a.dashboard_id = ?`)
		params = append(params, query.DashboardId)
	}

	if query.PanelId != 0 {
		sql.WriteString(` AND a.panel_id = ?`)
		params = append(params, query.PanelId)
	}

	if query.UserId != 0 {
		sql.WriteString(` AND a.user_id = ?`)
		params = append(params, query.UserId)
	}

	if query.From > 0 && query.To > 0 {
		sql.WriteString(` AND a.epoch <= ? AND a.epoch_end >= ?`)
		params = append(params, query.To, query.From)
	}

	if query.Type == "alert" {
		sql.WriteString(` AND a.alert_id > 0`)
	} else if query.Type == "annotation" {
		sql.WriteString(` AND a.alert_id = 0`)
	}

	if len(query.Tags) > 0 {
		keyValueFilters := []string{}

		tags := models.ParseTagPairs(query.Tags)
		for _, tag := range tags {
			if tag.Value == "" {
				keyValueFilters = append(keyValueFilters, "(tag."+dialect.Quote("key")+" = ?)")
				params = append(params, tag.Key)
			} else {
				keyValueFilters = append(keyValueFilters, "(tag."+dialect.Quote("key")+" = ? AND tag."+dialect.Quote("value")+" = ?)")
				params = append(params, tag.Key, tag.Value)
			}
		}

		if len(tags) > 0 {
			tagsSubQuery := fmt.Sprintf(`
        SELECT SUM(1) FROM annotation_tag at
          INNER JOIN tag on tag.id = at.tag_id
          WHERE at.annotation_id = a.id
            AND (
              %s
            )
      `, strings.Join(keyValueFilters, " OR "))

			if query.MatchAny {
				sql.WriteString(fmt.Sprintf(" AND (%s) > 0 ", tagsSubQuery))
			} else {
				sql.WriteString(fmt.Sprintf(" AND (%s) = %d ", tagsSubQuery, len(tags)))
			}

		}
	}

	if query.Limit == 0 {
		query.Limit = 100
	}

	// order of ORDER BY arguments match the order of a sql index for performance
	sql.WriteString(" ORDER BY a.org_id, a.epoch_end DESC, a.epoch DESC" + dialect.Limit(query.Limit) + " ) dt on dt.id = annotation.id")

	items := make([]*ItemDTO, 0)
	err := as.SQLStore.WithTransactionalDbSession(context.Background(), func(sess *sqlstore.DBSession) error {
		return sess.SQL(sql.String(), params...).Find(&items)
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (as *Service) Delete(params *DeleteParams) error {
	return as.SQLStore.WithTransactionalDbSession(context.Background(), func(sess *sqlstore.DBSession) error {
		var (
			sql         string
			annoTagSql  string
			queryParams []interface{}
		)

		as.log.Info("delete", "orgId", params.OrgId)
		if params.Id != 0 {
			annoTagSql = "DELETE FROM annotation_tag WHERE annotation_id IN (SELECT id FROM annotation WHERE id = ? AND org_id = ?)"
			sql = "DELETE FROM annotation WHERE id = ? AND org_id = ?"
			queryParams = []interface{}{params.Id, params.OrgId}
		} else {
			annoTagSql = "DELETE FROM annotation_tag WHERE annotation_id IN (SELECT id FROM annotation WHERE dashboard_id = ? AND panel_id = ? AND org_id = ?)"
			sql = "DELETE FROM annotation WHERE dashboard_id = ? AND panel_id = ? AND org_id = ?"
			queryParams = []interface{}{params.DashboardId, params.PanelId, params.OrgId}
		}

		sqlOrArgs := append([]interface{}{annoTagSql}, queryParams...)

		if _, err := sess.Exec(sqlOrArgs...); err != nil {
			return err
		}

		sqlOrArgs = append([]interface{}{sql}, queryParams...)

		if _, err := sess.Exec(sqlOrArgs...); err != nil {
			return err
		}

		return nil
	})
}

type Repository interface {
	Save(item *Item) error
	Update(item *Item) error
	Find(query *ItemQuery) ([]*ItemDTO, error)
	Delete(params *DeleteParams) error
}

type ItemQuery struct {
	OrgId        int64    `json:"orgId"`
	From         int64    `json:"from"`
	To           int64    `json:"to"`
	UserId       int64    `json:"userId"`
	AlertId      int64    `json:"alertId"`
	DashboardId  int64    `json:"dashboardId"`
	PanelId      int64    `json:"panelId"`
	AnnotationId int64    `json:"annotationId"`
	Tags         []string `json:"tags"`
	Type         string   `json:"type"`
	MatchAny     bool     `json:"matchAny"`

	Limit int64 `json:"limit"`
}

type PostParams struct {
	DashboardId int64  `json:"dashboardId"`
	PanelId     int64  `json:"panelId"`
	Epoch       int64  `json:"epoch"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	Icon        string `json:"icon"`
}

type DeleteParams struct {
	OrgId       int64
	Id          int64
	AlertId     int64
	DashboardId int64
	PanelId     int64
}

type Item struct {
	Id          int64            `json:"id"`
	OrgId       int64            `json:"orgId"`
	UserId      int64            `json:"userId"`
	DashboardId int64            `json:"dashboardId"`
	PanelId     int64            `json:"panelId"`
	Text        string           `json:"text"`
	AlertId     int64            `json:"alertId"`
	PrevState   string           `json:"prevState"`
	NewState    string           `json:"newState"`
	Epoch       int64            `json:"epoch"`
	EpochEnd    int64            `json:"epochEnd"`
	Created     int64            `json:"created"`
	Updated     int64            `json:"updated"`
	Tags        []string         `json:"tags"`
	Data        *simplejson.Json `json:"data"`

	// needed until we remove it from db
	Type  string
	Title string
}

type ItemDTO struct {
	Id          int64            `json:"id"`
	AlertId     int64            `json:"alertId"`
	AlertName   string           `json:"alertName"`
	DashboardId int64            `json:"dashboardId"`
	PanelId     int64            `json:"panelId"`
	UserId      int64            `json:"userId"`
	NewState    string           `json:"newState"`
	PrevState   string           `json:"prevState"`
	Created     int64            `json:"created"`
	Updated     int64            `json:"updated"`
	Time        int64            `json:"time"`
	TimeEnd     int64            `json:"timeEnd"`
	Text        string           `json:"text"`
	Tags        []string         `json:"tags"`
	Login       string           `json:"login"`
	Email       string           `json:"email"`
	AvatarUrl   string           `json:"avatarUrl"`
	Data        *simplejson.Json `json:"data"`
}
