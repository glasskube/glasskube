package manifestvalues

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/adapter"
	"go.uber.org/multierr"
)

type Resolver struct {
	pkg    adapter.PackageClientAdapter
	client adapter.KubernetesClientAdapter
}

func NewResolver(pkg adapter.PackageClientAdapter, client adapter.KubernetesClientAdapter) *Resolver {
	return &Resolver{pkg: pkg, client: client}
}

func (r *Resolver) Resolve(ctx context.Context, values map[string]v1alpha1.ValueConfiguration) (
	map[string]string,
	error,
) {
	var errComposite error
	resolvedValues := make(map[string]string)
	for name, value := range values {
		if resolved, err := r.ResolveValue(ctx, value); err != nil {
			multierr.AppendInto(&errComposite, fmt.Errorf("cannot resolve value %v: %w", name, err))
		} else {
			resolvedValues[name] = resolved
		}
	}
	return resolvedValues, errComposite
}

func (r *Resolver) ResolveValue(ctx context.Context, value v1alpha1.ValueConfiguration) (string, error) {
	if value.Value != nil {
		return *value.Value, nil
	} else if value.ValueFrom != nil {
		if r, err := r.resolveReference(ctx, *value.ValueFrom); err != nil {
			return "", err
		} else {
			return r, nil
		}
	} else {
		return "", errors.New("cannot resolve empty value")
	}
}

func (r *Resolver) resolveReference(ctx context.Context, ref v1alpha1.ValueReference) (string, error) {
	if ref.ConfigMapRef != nil {
		return r.resolveConfigMapRef(ctx, *ref.ConfigMapRef)
	} else if ref.SecretRef != nil {
		return r.resolveSecretRef(ctx, *ref.SecretRef)
	} else if ref.PackageRef != nil {
		return r.resolvePackageRef(ctx, *ref.PackageRef)
	} else {
		return "", errors.New("cannot resolve empty reference")
	}
}

func (r *Resolver) resolveConfigMapRef(ctx context.Context, ref v1alpha1.ObjectKeyValueSource) (string, error) {
	if c, err := r.client.GetConfigMap(ctx, ref.Name, ref.Namespace); err != nil {
		return "", NewConfigMapRefError(ref, err)
	} else if v, ok := c.Data[ref.Key]; !ok {
		return "", NewConfigMapRefError(ref, NewKeyError(ref.Key))
	} else {
		return v, nil
	}
}

func (r *Resolver) resolveSecretRef(ctx context.Context, ref v1alpha1.ObjectKeyValueSource) (string, error) {
	if c, err := r.client.GetSecret(ctx, ref.Name, ref.Namespace); err != nil {
		return "", NewSecretRefError(ref, err)
	} else if v, ok := c.Data[ref.Key]; !ok {
		return "", NewSecretRefError(ref, NewKeyError(ref.Key))
	} else if decoded, err := base64.StdEncoding.DecodeString(string(v)); err != nil {
		return "", NewSecretRefError(ref, err)
	} else {
		return string(decoded), nil
	}
}

func (r *Resolver) resolvePackageRef(ctx context.Context, ref v1alpha1.PackageValueSource) (string, error) {
	if pkg, err := r.pkg.GetClusterPackage(ctx, ref.Name); err != nil {
		return "", NewPackageRefError(ref, err)
	} else if value, ok := pkg.Spec.Values[ref.Value]; !ok {
		return "", NewPackageRefError(ref, NewKeyError(ref.Value))
	} else if resolved, err := r.ResolveValue(ctx, value); err != nil {
		return "", NewPackageRefError(ref, err)
	} else {
		return resolved, nil
	}
}
