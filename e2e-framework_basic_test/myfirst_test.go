package main

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestMyFirst(t *testing.T) {

	f1 := features.New("Crossplane").
		Assess("Successfully installed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

			client, err := cfg.NewClient()
			if err != nil {
				t.Fatalf("Unable to create client: %v/n", err)
			}
			var pods v1.PodList
			err = client.Resources("kube-system").List(context.TODO(), &pods)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(pods)

			return ctx
		}).Feature()

	testenv.Test(t, f1)

}
