package log

import (
	"context"
	"testing"
)

func TestInfo(t *testing.T) {
	ctx := context.Background()
	Infof(ctx, "123 :%s", "666")
	//logrus.Info(123)
}
