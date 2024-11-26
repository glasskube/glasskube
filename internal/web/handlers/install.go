package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/clicontext"
	"github.com/glasskube/glasskube/internal/clientutils"
	"github.com/glasskube/glasskube/internal/cliutils"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/namespaces"
	repoerror "github.com/glasskube/glasskube/internal/repo/error"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/internal/web/responder"
	"github.com/glasskube/glasskube/pkg/client"
	"github.com/glasskube/glasskube/pkg/install"
	"github.com/glasskube/glasskube/pkg/manifest"
	"go.uber.org/multierr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PostPackageDetail(w http.ResponseWriter, r *http.Request) {
	p := getPackageContext(r).request
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)
	namespace := r.FormValue("namespace")
	name := r.FormValue("name")
	autoUpdate := strings.ToLower(r.FormValue("autoUpdate")) == "on"
	dryRun, _ := strconv.ParseBool(r.FormValue("dryRun"))

	var err error
	pkg := &v1alpha1.Package{}
	var mf *v1alpha1.PackageManifest
	if err := pkgClient.Packages(p.namespace).Get(ctx, p.name, pkg); err != nil && !errors.IsNotFound(err) {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch package %v/%v: %w", p.namespace, p.name, err)))
		return
	} else if err != nil {
		pkg = nil
	} else {
		// because disabled form elements are not submitted, we need to fall back on repo and version if advanced options disabled
		if p.repositoryName == "" {
			p.repositoryName = pkg.Spec.PackageInfo.RepositoryName
		}
		if p.version == "" {
			p.version = pkg.Spec.PackageInfo.Version
		}
	}

	mf, err = resolveManifest(ctx, pkg, p.repositoryName, p.manifestName, p.version)
	if repoerror.IsPartial(err) {
		fmt.Fprintf(os.Stderr, "problem fetching manifest and repo, but installation can continue: %v", err)
	} else if err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to get manifest and repo of %v: %w", p.manifestName, err)))
		return
	}

	if values, err := extractValues(r, mf); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to parse values: %w", err)))
		return
	} else if pkg == nil {
		opts := v1.CreateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		k8sClient := clicontext.KubernetesClientFromContext(ctx)
		if exists, err := namespaces.Exists(ctx, k8sClient, namespace); err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to check namespace: %w", err)))
			return
		} else if !exists {
			ns := corev1.Namespace{
				ObjectMeta: v1.ObjectMeta{
					Name: namespace,
				},
			}
			if _, err := k8sClient.CoreV1().Namespaces().Create(ctx, &ns, opts); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to create namespace: %w", err)))
				return
			}
		}
		pkg = client.PackageBuilder(p.manifestName).
			WithVersion(p.version).
			WithRepositoryName(p.repositoryName).
			WithAutoUpdates(autoUpdate).
			WithValues(values).
			WithNamespace(namespace).
			WithName(name).
			BuildPackage()
		err := install.NewInstaller(pkgClient).Install(ctx, pkg, opts)
		if err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to install %v: %w", p.manifestName, err)))
		} else if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				responder.SendYamlModal(w, yamlOutput, nil)
			}
		} else {
			responder.Redirect(w, "/packages")
			w.WriteHeader(http.StatusAccepted)
		}
	} else {
		pkg.Spec.PackageInfo.Version = p.version
		pkg.Spec.PackageInfo.RepositoryName = p.repositoryName
		pkg.Spec.Values = values
		pkg.SetAutoUpdatesEnabled(autoUpdate)
		opts := v1.UpdateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		if err := pkgClient.Packages(pkg.GetNamespace()).Update(ctx, pkg, opts); err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to configure %v: %w", p.manifestName, err)))
			return
		}
		valueResolver := cliutils.ValueResolver(ctx)
		_, resolveErr := valueResolver.Resolve(ctx, values)
		if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				responder.SendYamlModal(w, yamlOutput, resolveErr)
			}
		} else if resolveErr != nil {
			responder.SendToast(w,
				toast.WithErr(fmt.Errorf("some values could not be resolved: %w", resolveErr)),
				toast.WithSeverity(toast.Warning),
				toast.WithStatusCode(http.StatusAccepted))
		} else {
			responder.SendToast(w, toast.WithMessage("Configuration updated successfully"))
		}
	}
}

func PostClusterPackageDetail(w http.ResponseWriter, r *http.Request) {
	p := getPackageContext(r).request
	ctx := r.Context()
	autoUpdate := strings.ToLower(r.FormValue("autoUpdate")) == "on"
	dryRun, _ := strconv.ParseBool(r.FormValue("dryRun"))

	var err error
	pkg := &v1alpha1.ClusterPackage{}
	var mf *v1alpha1.PackageManifest
	pkgClient := clicontext.PackageClientFromContext(ctx)
	if err = pkgClient.ClusterPackages().Get(ctx, p.manifestName, pkg); err != nil && !errors.IsNotFound(err) {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch clusterpackage %v: %w", p.manifestName, err)))
		return
	} else if err != nil {
		pkg = nil
	} else {
		// because disabled form elements are not submitted, we need to fall back on repo and version if advanced options disabled
		if p.repositoryName == "" {
			p.repositoryName = pkg.Spec.PackageInfo.RepositoryName
		}
		if p.version == "" {
			p.version = pkg.Spec.PackageInfo.Version
		}
	}

	mf, err = resolveManifest(ctx, pkg, p.repositoryName, p.manifestName, p.version)
	if repoerror.IsPartial(err) {
		fmt.Fprintf(os.Stderr, "problem fetching manifest and repo, but installation can continue: %v", err)
	} else if err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to get manifest and repo of %v: %w", p.manifestName, err)))
		return
	}

	if values, err := extractValues(r, mf); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to parse values: %w", err)))
		return
	} else if pkg == nil {
		pkg = client.PackageBuilder(p.manifestName).
			WithVersion(p.version).
			WithRepositoryName(p.repositoryName).
			WithAutoUpdates(autoUpdate).
			WithValues(values).
			BuildClusterPackage()
		opts := v1.CreateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		err := install.NewInstaller(pkgClient).Install(ctx, pkg, opts)
		if err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to install %v: %w", p.manifestName, err)))
			return
		} else if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				responder.SendYamlModal(w, yamlOutput, nil)
			}
		}
	} else {
		pkg.Spec.PackageInfo.Version = p.version
		pkg.Spec.PackageInfo.RepositoryName = p.repositoryName
		pkg.Spec.Values = values
		pkg.SetAutoUpdatesEnabled(autoUpdate)
		opts := v1.UpdateOptions{}
		if dryRun {
			opts.DryRun = []string{v1.DryRunAll}
		}
		if err := pkgClient.ClusterPackages().Update(ctx, pkg, opts); err != nil {
			responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to configure %v: %w", p.manifestName, err)))
			return
		}
		valueResolver := cliutils.ValueResolver(ctx)
		_, resolveErr := valueResolver.Resolve(ctx, values)
		if dryRun {
			if yamlOutput, err := clientutils.Format(clientutils.OutputFormatYAML, false, pkg); err != nil {
				responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to render yaml: %w", err)))
			} else {
				responder.SendYamlModal(w, yamlOutput, resolveErr)
			}
		} else if resolveErr != nil {
			responder.SendToast(w,
				toast.WithErr(fmt.Errorf("some values could not be resolved: %w", resolveErr)),
				toast.WithSeverity(toast.Warning),
				toast.WithStatusCode(http.StatusAccepted))
		} else {
			responder.SendToast(w, toast.WithMessage("Configuration updated successfully"))
		}
	}
}

func resolveManifest(ctx context.Context, pkg ctrlpkg.Package, repositoryName string, manifestName string, selectedVersion string) (
	*v1alpha1.PackageManifest, error) {

	var mf v1alpha1.PackageManifest
	var repoErr error
	if pkg.IsNil() ||
		(pkg.GetSpec().PackageInfo.RepositoryName != repositoryName || pkg.GetSpec().PackageInfo.Version != selectedVersion) {
		repoClientset := clicontext.RepoClientsetFromContext(ctx)
		repoClient := repoClientset.ForRepoWithName(repositoryName)
		if err := repoClient.FetchPackageManifest(manifestName, selectedVersion, &mf); err != nil {
			return nil, multierr.Append(err, repoErr)
		}
	} else {
		if installedMf, err := manifest.GetInstalledManifestForPackage(ctx, pkg); err != nil {
			return nil, multierr.Append(err, repoErr)
		} else {
			mf = *installedMf
		}
	}
	return &mf, repoErr
}
