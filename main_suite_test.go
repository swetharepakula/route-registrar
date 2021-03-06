package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/cloudfoundry-incubator/route-registrar/Godeps/_workspace/src/github.com/onsi/ginkgo"
	gconfig "github.com/cloudfoundry-incubator/route-registrar/Godeps/_workspace/src/github.com/onsi/ginkgo/config"
	. "github.com/cloudfoundry-incubator/route-registrar/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry-incubator/route-registrar/Godeps/_workspace/src/github.com/onsi/gomega/gexec"
	"github.com/cloudfoundry-incubator/route-registrar/config"

	"testing"
)

const (
	routeRegistrarPackage = "github.com/cloudfoundry-incubator/route-registrar/"
)

var (
	routeRegistrarBinPath string
	pidFile               string
	configFile            string
	rootConfig            config.ConfigSchema
	natsPort              int

	tempDir string
)

func TestRouteRegistrar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = BeforeSuite(func() {
	var err error
	routeRegistrarBinPath, err = gexec.Build(routeRegistrarPackage, "-race")
	Expect(err).ShouldNot(HaveOccurred())

	tempDir, err = ioutil.TempDir(os.TempDir(), "route-registrar")
	Expect(err).ShouldNot(HaveOccurred())

	pidFile = filepath.Join(tempDir, "route-registrar.pid")

	natsPort = 40000 + gconfig.GinkgoConfig.ParallelNode

	configFile = filepath.Join(tempDir, "registrar_settings.yml")
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(tempDir)
	Expect(err).ShouldNot(HaveOccurred())

	gexec.CleanupBuildArtifacts()
})
