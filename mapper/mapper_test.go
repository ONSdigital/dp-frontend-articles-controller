package mapper

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	ctx := context.Background()

	Convey("Blank maps correctly", t, func() {
		cfg := config.Config{
			BindAddr:                   "1234",
			GracefulShutdownTimeout:    0,
			HealthCheckInterval:        0,
			HealthCheckCriticalTimeout: 0,
		}

		bulletin := Bulletin{
			Name: "test",
		}

		model := Blank(ctx, bulletin, cfg)
		So(model.Name, ShouldEqual, bulletin.Name)
	})
}
