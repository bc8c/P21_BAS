package store_test

import (
	"context"
	"testing"

	"../models"
	"../store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClientStore(t *testing.T) {
	Convey("Test client store", t, func() {
		clientStore := store.NewClientStore()

		err := clientStore.Set("1", &models.Client{ID: "1", Secret: "2"})
		So(err, ShouldBeNil)

		cli, err := clientStore.GetByID(context.Background(), "1")
		So(err, ShouldBeNil)
		So(cli.GetID(), ShouldEqual, "1")
	})
}
