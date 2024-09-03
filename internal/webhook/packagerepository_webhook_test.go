/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"context"

	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newPackageRepositoryValidatingWebhook(objects ...client.Object) *PackageRepositoryValidatingWebhook {
	fakeClient := fake.NewClientBuilder().
		WithObjects(objects...).
		Build()
	return &PackageRepositoryValidatingWebhook{
		Client: fakeClient,
	}
}

var (
	glasskubev1Repo = v1alpha1.PackageRepository{
		ObjectMeta: v1.ObjectMeta{Name: "glasskube"},
		Spec:       v1alpha1.PackageRepositorySpec{Url: "https://packages.dl.glasskube.dev/packages"},
	}
	glasskubev2Repo = v1alpha1.PackageRepository{
		ObjectMeta: v1.ObjectMeta{Name: glasskubev1Repo.Name},
		Spec:       v1alpha1.PackageRepositorySpec{Url: "https://packages.dl.glasskube.new/packages"},
	}
	hinterseerv1Package = v1alpha1.Package{
		ObjectMeta: v1.ObjectMeta{Name: "hinterseer"},
		Spec:       v1alpha1.PackageSpec{PackageInfo: v1alpha1.PackageInfoTemplate{Name: "hinterseer", RepositoryName: glasskubev1Repo.Name}},
	}
	hinterseerv1PackageInfo = v1alpha1.PackageInfo{
		ObjectMeta: v1.ObjectMeta{Name: "glasskube--hinterseer--v0.8.1--5"},
		Spec:       v1alpha1.PackageInfoSpec{Name: hinterseerv1Package.Name, RepositoryName: hinterseerv1Package.Spec.PackageInfo.RepositoryName, Version: "v0.8.1+5"},
	}
	metalkubev1Repo = v1alpha1.PackageRepository{
		ObjectMeta: v1.ObjectMeta{Name: "metalkube"},
		Spec:       v1alpha1.PackageRepositorySpec{Url: "https://packages.dl.metalkube.dev/packages"},
	}
	metalkubev2Repo = v1alpha1.PackageRepository{
		ObjectMeta: v1.ObjectMeta{Name: metalkubev1Repo.Name},
		Spec:       v1alpha1.PackageRepositorySpec{Url: "https://packages.dl.metalkube.new/packages"},
	}
)

var _ = Describe("PackageRepositoryValidatingWebhook", Ordered, func() {
	BeforeAll(func() {
		err := v1alpha1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())
	})
	Context("ValidateDelete", func() {
		When("no packages installed", func() {
			It("should not return an error", func(ctx context.Context) {
				webhook := newPackageRepositoryValidatingWebhook(&metalkubev1Repo)
				_, err := webhook.ValidateDelete(ctx, &metalkubev1Repo)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("packages installed", func() {
			It("should return an error", func(ctx context.Context) {
				webhook := newPackageRepositoryValidatingWebhook(&glasskubev1Repo, &hinterseerv1Package, &hinterseerv1PackageInfo)
				_, err := webhook.ValidateDelete(ctx, &glasskubev1Repo)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ErrPackagesInstalled))
			})
		})
	})
	Context("ValidateUpdate", func() {
		When("url update and no packages installed", func() {
			It("should not return an error", func(ctx context.Context) {
				webhook := newPackageRepositoryValidatingWebhook(&metalkubev1Repo)
				_, err := webhook.ValidateUpdate(ctx, &metalkubev1Repo, &metalkubev2Repo)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("url update and packages installed", func() {
			It("should return an error", func(ctx context.Context) {
				webhook := newPackageRepositoryValidatingWebhook(&glasskubev1Repo, &hinterseerv1Package, &hinterseerv1PackageInfo)
				_, err := webhook.ValidateUpdate(ctx, &glasskubev1Repo, &glasskubev2Repo)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ErrPackagesInstalled))
			})
		})
		When("url not changed", func() {
			It("should not return an error", func(ctx context.Context) {
				webhook := newPackageRepositoryValidatingWebhook(&glasskubev1Repo, &hinterseerv1Package, &hinterseerv1PackageInfo)
				_, err := webhook.ValidateUpdate(ctx, &glasskubev1Repo, &glasskubev1Repo)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
