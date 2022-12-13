package main

import (
	"fmt"
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
)

var (
	testenv         env.Environment
	kindClusterName string
)

type Readier func() (bool, error)

func TestMain(m *testing.M) {

	cfg := envconf.New()
	testenv = env.NewWithConfig(cfg)
	kindClusterName = envconf.RandomName("test-", 16)

	fmt.Println("Setting up environment...")

	testenv.Setup(
		envfuncs.CreateKindCluster(kindClusterName),
		envfuncs.CreateNamespace("tester"),
	)

	testenv.Finish(
		envfuncs.DeleteNamespace("tester"),
		envfuncs.DestroyKindCluster(kindClusterName),
	)

	os.Exit(testenv.Run(m))
}
