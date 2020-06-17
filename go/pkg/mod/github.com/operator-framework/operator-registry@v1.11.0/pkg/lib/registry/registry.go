package registry

import (
	"context"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/operator-framework/operator-registry/pkg/containertools"
	"github.com/operator-framework/operator-registry/pkg/image/execregistry"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/operator-framework/operator-registry/pkg/image"
	"github.com/operator-framework/operator-registry/pkg/image/containerdregistry"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/operator-framework/operator-registry/pkg/sqlite"
)

type RegistryUpdater struct {
	Logger   *logrus.Entry
}

type AddToRegistryRequest struct {
	Permissive    bool
	SkipTLS       bool
	CaFile        string
	InputDatabase string
	Bundles       []string
	Mode          registry.Mode
	ContainerTool containertools.ContainerTool
}

func (r RegistryUpdater) AddToRegistry(request AddToRegistryRequest) error {
	var errs []error

	db, err := sql.Open("sqlite3", request.InputDatabase)
	if err != nil {
		return err
	}
	defer db.Close()

	dbLoader, err := sqlite.NewSQLLiteLoader(db)
	if err != nil {
		return err
	}
	if err := dbLoader.Migrate(context.TODO()); err != nil {
		return err
	}

	graphLoader, err := sqlite.NewSQLGraphLoaderFromDB(db)
	if err != nil {
		return err
	}
	dbQuerier := sqlite.NewSQLLiteQuerierFromDb(db)

	// add custom ca certs to resolver
	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if len(request.CaFile) > 0 {
		certs, err := ioutil.ReadFile(request.CaFile)
		if err != nil {
			return fmt.Errorf("failed to append %q to RootCAs: %v", certs, err)
		}
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			return fmt.Errorf("unable to add certs specified in %s", request.CaFile)
		}
	}

	var reg image.Registry
	var rerr error
	switch request.ContainerTool {
	case containertools.NoneTool:
		reg, rerr = containerdregistry.NewRegistry(containerdregistry.SkipTLS(request.SkipTLS), containerdregistry.WithRootCAs(rootCAs))
	case containertools.PodmanTool:
	case containertools.DockerTool:
		reg, rerr = execregistry.NewRegistry(request.ContainerTool, r.Logger)
	}
	if rerr != nil {
		return rerr
	}
	defer func() {
		if err := reg.Destroy(); err != nil {
			r.Logger.WithError(err).Warn("error destroying local cache")
		}
	}()

	// TODO(njhale): Parallelize this once bundle add is commutative
	for _, ref := range request.Bundles {
		if err := populate(context.TODO(), dbLoader, graphLoader, dbQuerier, reg, image.SimpleReference(ref), request.Mode); err != nil {
			err = fmt.Errorf("error loading bundle from image: %s", err)
			if !request.Permissive {
				r.Logger.WithError(err).Error("permissive mode disabled")
				errs = append(errs, err)
			} else {
				r.Logger.WithError(err).Warn("permissive mode enabled")
			}
		}
	}

	return utilerrors.NewAggregate(errs) // nil if no errors
}

func populate(ctx context.Context, loader registry.Load, graphLoader registry.GraphLoader, querier registry.Query, reg image.Registry, ref image.Reference, mode registry.Mode) error {
	workingDir, err := ioutil.TempDir("./", "bundle_tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workingDir)

	if err = reg.Pull(ctx, ref); err != nil {
		return err
	}

	if err = reg.Unpack(ctx, ref, workingDir); err != nil {
		return err
	}

	populator := registry.NewDirectoryPopulator(loader, graphLoader, querier, ref, workingDir)

	return populator.Populate(mode)
}

type DeleteFromRegistryRequest struct {
	Permissive    bool
	InputDatabase string
	Packages      []string
}

func (r RegistryUpdater) DeleteFromRegistry(request DeleteFromRegistryRequest) error {
	db, err := sql.Open("sqlite3", request.InputDatabase)
	if err != nil {
		return err
	}
	defer db.Close()

	dbLoader, err := sqlite.NewSQLLiteLoader(db)
	if err != nil {
		return err
	}
	if err := dbLoader.Migrate(context.TODO()); err != nil {
		return err
	}

	for _, pkg := range request.Packages {
		remover := sqlite.NewSQLRemoverForPackages(dbLoader, pkg)
		if err := remover.Remove(); err != nil {
			err = fmt.Errorf("error deleting packages from database: %s", err)
			if !request.Permissive {
				logrus.WithError(err).Fatal("permissive mode disabled")
				return err
			}
			logrus.WithError(err).Warn("permissive mode enabled")
		}
	}

	return nil
}
