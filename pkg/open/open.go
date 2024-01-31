package open

import (
	"context"

	"github.com/glasskube/glasskube/pkg/client"
)

func Portforward(client *client.PackageV1Alpha1Client, ctx context.Context, packageName string) error {
	// portforward.New(httpstream.Dialer())
	return nil
}
