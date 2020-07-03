package annotations

import (
	"testing"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/stretchr/testify/require"
)

func TestAnnotations(t *testing.T) {
	fakeSQLStore := sqlstore.InitTestDB(t)
	repo := &Service{
		SQLStore: fakeSQLStore,
		log:      log.New("test logger"),
	}

	t.Cleanup(func() {
		session := repo.SQLStore.NewSession()
		defer session.Close()
		_, err := session.Exec("DELETE FROM annotation WHERE 1=1")
		require.Nil(t, err)
		_, err = session.Exec("DELETE FROM annotation_tag WHERE 1=1")
		require.Nil(t, err)
	})

	annotation := &Item{
		OrgId:       1,
		UserId:      1,
		DashboardId: 1,
		Text:        "hello",
		Type:        "alert",
		Epoch:       10,
		Tags:        []string{"outage", "error", "type:outage", "server:server-1"},
	}
	err := repo.Save(annotation)

	require.Nil(t, err)
	require.Greater(t, annotation.Id, int64(0))
	require.Equal(t, annotation.Epoch, annotation.EpochEnd)

	annotation2 := &Item{
		OrgId:       1,
		UserId:      1,
		DashboardId: 2,
		Text:        "hello",
		Type:        "alert",
		Epoch:       21, // Should swap epoch & epochEnd
		EpochEnd:    20,
		Tags:        []string{"outage", "error", "type:outage", "server:server-1"},
	}
	err = repo.Save(annotation2)
	require.Nil(t, err)
	require.Greater(t, annotation2.Id, int64(0))
	require.Equal(t, annotation2.Epoch, int64(20))
	require.Equal(t, annotation2.EpochEnd, int64(21))

	globalAnnotation1 := &Item{
		OrgId:  1,
		UserId: 1,
		Text:   "deploy",
		Type:   "",
		Epoch:  15,
		Tags:   []string{"deploy"},
	}
	err = repo.Save(globalAnnotation1)
	require.Nil(t, err)
	require.Greater(t, globalAnnotation1.Id, int64(0))
	//So(err, ShouldBeNil)
	//So(globalAnnotation1.Id, ShouldBeGreaterThan, 0)

	globalAnnotation2 := &Item{
		OrgId:  1,
		UserId: 1,
		Text:   "rollback",
		Type:   "",
		Epoch:  17,
		Tags:   []string{"rollback"},
	}
	err = repo.Save(globalAnnotation2)
	require.Nil(t, err)
	require.Greater(t, globalAnnotation2.Id, int64(0))
	//So(err, ShouldBeNil)
	//So(globalAnnotation2.Id, ShouldBeGreaterThan, 0)

	t.Run("Can query for annotation by dashboard id", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        0,
			To:          15,
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 1)
		//So(err, ShouldBeNil)
		//So(items, ShouldHaveLength, 1)

		//Convey("Can read tags", func() {
		expectedTags := []string{"outage", "error", "type:outage", "server:server-1"}
		require.Equal(t, items[0].Tags, expectedTags, "Can read the tags on the annotation")
		//So(items[0].Tags, ShouldResemble, []string{"outage", "error", "type:outage", "server:server-1"})
		//})

		//Convey("Has created and updated values", func() {
		require.Greater(t, items[0].Created, int64(0))
		require.Greater(t, items[0].Updated, int64(0))
		require.Equal(t, items[0].Updated, items[0].Created)

		//So(items[0].Created, ShouldBeGreaterThan, 0)
		//So(items[0].Updated, ShouldBeGreaterThan, 0)
		//So(items[0].Updated, ShouldEqual, items[0].Created)
		//})
	})

	t.Run("Can query for annotation by id", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:        1,
			AnnotationId: annotation2.Id,
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 1)
		require.Equal(t, items[0].Id, annotation2.Id)
		//So(err, ShouldBeNil)
		//So(items, ShouldHaveLength, 1)
		//So(items[0].Id, ShouldEqual, annotation2.Id)
	})

	t.Run("Should not find any when item is outside time range", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        12,
			To:          15,
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 0)
		//So(err, ShouldBeNil)
		//So(items, ShouldHaveLength, 0)
	})

	t.Run("Should not find one when tag filter does not match", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        1,
			To:          15,
			Tags:        []string{"asd"},
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 0)
		//So(err, ShouldBeNil)
		//So(items, ShouldHaveLength, 0)
	})

	t.Run("Should not find one when type filter does not match", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        1,
			To:          15,
			Type:        "alert",
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 0)
		//So(err, ShouldBeNil)
		//So(items, ShouldHaveLength, 0)
	})

	t.Run("Should find one when all tag filters does match", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        1,
			To:          15, //this will exclude the second test annotation
			Tags:        []string{"outage", "error"},
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 1)
		//So(err, ShouldBeNil)
		//So(items, ShouldHaveLength, 1)
	})

	t.Run("Should find two annotations using partial match", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:    1,
			From:     1,
			To:       25,
			MatchAny: true,
			Tags:     []string{"rollback", "deploy"},
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 2)
		// So(err, ShouldBeNil)
		// So(items, ShouldHaveLength, 2)
	})

	t.Run("Should find one when all key value tag filters does match", func(t *testing.T) {
		items, err := repo.Find(&ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        1,
			To:          15,
			Tags:        []string{"type:outage", "server:server-1"},
		})

		require.Nil(t, err)
		require.Equal(t, len(items), 1)
		// So(err, ShouldBeNil)
		// So(items, ShouldHaveLength, 1)
	})

	t.Run("Can update annotation and remove all tags", func(t *testing.T) {
		query := &ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        0,
			To:          15,
		}
		items, err := repo.Find(query)

		require.Nil(t, err)
		//So(err, ShouldBeNil)

		annotationId := items[0].Id

		err = repo.Update(&Item{
			Id:    annotationId,
			OrgId: 1,
			Text:  "something new",
			Tags:  []string{},
		})

		require.Nil(t, err)
		//So(err, ShouldBeNil)

		items, err = repo.Find(query)

		require.Nil(t, err)
		//So(err, ShouldBeNil)

		//Convey("Can read tags", func() {
		require.Equal(t, items[0].Id, annotationId)
		require.Equal(t, len(items[0].Tags), 0)
		require.Equal(t, items[0].Text, "something new")
		//So(items[0].Id, ShouldEqual, annotationId)
		//So(len(items[0].Tags), ShouldEqual, 0)
		//So(items[0].Text, ShouldEqual, "something new")
		//})
	})

	t.Run("Can update annotation with new tags", func(t *testing.T) {
		query := &ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        0,
			To:          15,
		}
		items, err := repo.Find(query)

		require.Nil(t, err)
		//So(err, ShouldBeNil)

		annotationId := items[0].Id

		err = repo.Update(&Item{
			Id:    annotationId,
			OrgId: 1,
			Text:  "something new",
			Tags:  []string{"newtag1", "newtag2"},
		})

		require.Nil(t, err)
		//So(err, ShouldBeNil)

		items, err = repo.Find(query)

		require.Nil(t, err)
		//So(err, ShouldBeNil)

		//Convey("Can read tags", func(t *testing.T) {
		require.Equal(t, items[0].Id, annotationId)
		require.Equal(t, items[0].Tags, []string{"newtag1", "newtag2"})
		require.Equal(t, items[0].Text, "something new")
		// So(items[0].Id, ShouldEqual, annotationId)
		// So(items[0].Tags, ShouldResemble, []string{"newtag1", "newtag2"})
		// So(items[0].Text, ShouldEqual, "something new")
		//})

		//Convey("Updated time has increased", func() {
		require.Greater(t, items[0].Updated, items[0].Created)
		//So(items[0].Updated, ShouldBeGreaterThan, items[0].Created)
		//})
	})

	t.Run("Can delete annotation", func(t *testing.T) {
		query := &ItemQuery{
			OrgId:       1,
			DashboardId: 1,
			From:        0,
			To:          15,
		}
		items, err := repo.Find(query)
		require.Nil(t, err)
		//So(err, ShouldBeNil)

		annotationId := items[0].Id

		err = repo.Delete(&DeleteParams{Id: annotationId, OrgId: 1})
		require.Nil(t, err)
		//So(err, ShouldBeNil)

		items, err = repo.Find(query)
		require.Nil(t, err)
		//So(err, ShouldBeNil)

		//Convey("Should be deleted", func() {
		require.Equal(t, len(items), 0)
		//So(len(items), ShouldEqual, 0)
		//})
	})
}
