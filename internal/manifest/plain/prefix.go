package plain

import (
	"encoding/json"
	"maps"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ScopeChecker interface {
	IsObjectNamespaced(runtime.Object) (bool, error)
}

type prefixer struct {
	ScopeChecker
}

func newPrefixer(client ScopeChecker) *prefixer {
	return &prefixer{ScopeChecker: client}
}

type nameMappingKey struct {
	Kind, Name string
}

type nameMapping map[nameMappingKey]string

func (p *prefixer) prefixAndUpdateReferences(
	pkg ctrlpkg.Package, manifest *v1alpha1.PackageManifest, objects []client.Object) error {
	if pkg.IsNamespaceScoped() {
		if nameMapping, err := p.prefix(objects, pkg.GetName()); err != nil {
			return err
		} else {
			maps.Copy(nameMapping, mappingKeyForTransitiveResources(manifest.TransitiveResources, pkg.GetName()))
			if err = updateReferences(objects, nameMapping); err != nil {
				return err
			} else {
				return nil
			}
		}
	} else {
		return nil
	}
}

func (p *prefixer) prefix(objs []client.Object, prefix string) (nameMapping, error) {
	nameMapping := make(nameMapping)
	for _, obj := range objs {
		if isNamespaced, err := p.IsObjectNamespaced(obj); err != nil {
			return nil, err
		} else if isNamespaced {
			newName := prefixedObjectName(prefix, obj.GetName())
			nameMapping[mappingKeyForObj(obj)] = newName
			obj.SetName(newName)
		}
	}
	return nameMapping, nil
}

func updateReferences(objs []client.Object, nameMapping nameMapping) error {
	for i, obj := range objs {
		switch obj.GetObjectKind().GroupVersionKind().Kind {
		case "Deployment":
			if deployment, err := toRealObj(obj, appsv1.Deployment{}); err != nil {
				return err
			} else {
				updateObjectReferencesPodTemplate(&deployment.Spec.Template, nameMapping)
				objs[i] = deployment
			}
		case "StatefulSet":
			if statefulSet, err := toRealObj(obj, appsv1.StatefulSet{}); err != nil {
				return err
			} else {
				updateObjectReferencesPodTemplate(&statefulSet.Spec.Template, nameMapping)
				updateObjectReferenceTarget("Service", &statefulSet.Spec.ServiceName, nameMapping)
				objs[i] = statefulSet
			}
		case "CronJob":
			if cronJob, err := toRealObj(obj, batchv1.CronJob{}); err != nil {
				return err
			} else {
				updateObjectReferencesPodTemplate(&cronJob.Spec.JobTemplate.Spec.Template, nameMapping)
				objs[i] = cronJob
			}
		case "Job":
			if job, err := toRealObj(obj, batchv1.Job{}); err != nil {
				return err
			} else {
				updateObjectReferencesPodTemplate(&job.Spec.Template, nameMapping)
				objs[i] = job
			}
		case "Ingress":
			if ingress, err := toRealObj(obj, netv1.Ingress{}); err != nil {
				return err
			} else {
				for _, rule := range ingress.Spec.Rules {
					if rule.HTTP == nil {
						continue
					}
					for _, path := range rule.HTTP.Paths {
						updateObjectReferenceTarget("Service", &path.Backend.Service.Name, nameMapping)
					}
				}
				objs[i] = ingress
			}
		case "RoleBinding":
			if roleBinding, err := toRealObj(obj, v1.RoleBinding{}); err != nil {
				return err
			} else {
				for _, subj := range roleBinding.Subjects {
					updateObjectReferenceTarget(subj.Kind, &subj.Name, nameMapping)
				}
				updateObjectReferenceTarget(roleBinding.RoleRef.Kind, &roleBinding.RoleRef.Name, nameMapping)
				objs[i] = roleBinding
			}
		}
	}
	return nil
}

func toRealObj[T any](obj client.Object, target T) (*T, error) {
	if data, err := json.Marshal(obj); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(data, &target); err != nil {
			return nil, err
		} else {
			return &target, nil
		}
	}
}

func updateObjectReferencesPodTemplate(tpl *corev1.PodTemplateSpec, nameMapping nameMapping) {
	for _, container := range tpl.Spec.Containers {
		updateObjectReferencesContainer(&container, nameMapping)
	}
	for _, container := range tpl.Spec.InitContainers {
		updateObjectReferencesContainer(&container, nameMapping)
	}
	for _, volume := range tpl.Spec.Volumes {
		updateObjectReferencesPodVolume(volume, nameMapping)
	}
}

func updateObjectReferencesContainer(container *corev1.Container, nameMapping nameMapping) {
	for _, it := range container.EnvFrom {
		if it.SecretRef != nil {
			updateObjectReferencesLocalReference("Secret", &it.SecretRef.LocalObjectReference, nameMapping)
		}
		if it.ConfigMapRef != nil {
			updateObjectReferencesLocalReference("ConfigMap", &it.SecretRef.LocalObjectReference, nameMapping)
		}
	}
	for _, it := range container.Env {
		if it.ValueFrom != nil {
			if it.ValueFrom.SecretKeyRef != nil {
				updateObjectReferencesLocalReference("Secret", &it.ValueFrom.SecretKeyRef.LocalObjectReference, nameMapping)
			}
			if it.ValueFrom.ConfigMapKeyRef != nil {
				updateObjectReferencesLocalReference("ConfigMap", &it.ValueFrom.SecretKeyRef.LocalObjectReference, nameMapping)
			}
		}
	}
}

func updateObjectReferencesPodVolume(volume corev1.Volume, nameMapping nameMapping) {
	if volume.Secret != nil {
		updateObjectReferenceTarget("Secret", &volume.Secret.SecretName, nameMapping)
	}
	if volume.ConfigMap != nil {
		updateObjectReferenceTarget("ConfigMap", &volume.ConfigMap.Name, nameMapping)
	}
	if volume.PersistentVolumeClaim != nil {
		updateObjectReferenceTarget("PersistentVolumeClaim", &volume.PersistentVolumeClaim.ClaimName, nameMapping)
	}
}

func updateObjectReferencesLocalReference(
	kind string,
	ref *corev1.LocalObjectReference,
	nameMapping nameMapping,
) {
	updateObjectReferenceTarget(kind, &ref.Name, nameMapping)
}

func updateObjectReferenceTarget(kind string, target *string, nameMapping nameMapping) {
	if name, ok := nameMapping[mappingKey(kind, *target)]; ok {
		*target = name
	}
}

func mappingKeyForTransitiveResources(ress []v1alpha1.TransitiveResource, prefix string) nameMapping {
	mapping := make(nameMapping, len(ress))
	for _, res := range ress {
		mapping[mappingKey(res.Kind, res.Name)] = prefixedObjectName(prefix, res.Name)
	}
	return mapping
}

func mappingKeyForObj(obj client.Object) nameMappingKey {
	return mappingKey(obj.GetObjectKind().GroupVersionKind().Kind, obj.GetName())
}

func mappingKey(kind, name string) nameMappingKey {
	return nameMappingKey{Kind: kind, Name: name}
}

func prefixedObjectName(prefix, name string) string {
	return strings.Join([]string{prefix, name}, "-")
}
