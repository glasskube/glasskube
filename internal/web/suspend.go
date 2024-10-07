package web

import (
	"fmt"
	"net/http"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	"github.com/glasskube/glasskube/internal/web/components/toast"
	"github.com/glasskube/glasskube/pkg/suspend"
	"github.com/gorilla/mux"
)

func (s *server) handleSuspend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if pkg, err := s.getPackageFromRequest(r); err != nil {
		s.sendToast(w, toast.WithErr(err))
	} else if result, err := suspend.Suspend(r.Context(), pkg); err != nil {
		s.sendToast(w, toast.WithErr(err))
	} else {
		s.sendToast(w, toast.WithMessage(fmt.Sprintf("suspend result: %v", result)))
	}
}

func (s *server) handleResume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if pkg, err := s.getPackageFromRequest(r); err != nil {
		s.sendToast(w, toast.WithErr(err))
	} else if result, err := suspend.Resume(r.Context(), pkg); err != nil {
		s.sendToast(w, toast.WithErr(err))
	} else {
		s.sendToast(w, toast.WithMessage(fmt.Sprintf("resume result: %v", result)))
	}
}

func (s *server) getPackageFromRequest(r *http.Request) (ctrlpkg.Package, error) {
	var pkg ctrlpkg.Package
	if name := mux.Vars(r)["pkgName"]; name != "" {
		var cp v1alpha1.ClusterPackage
		if err := s.pkgClient.ClusterPackages().Get(r.Context(), name, &cp); err != nil {
			return nil, err
		} else {
			pkg = &cp
		}
	} else {
		name = mux.Vars(r)["name"]
		namespace := mux.Vars(r)["namespace"]
		var p v1alpha1.Package
		if err := s.pkgClient.Packages(namespace).Get(r.Context(), name, &p); err != nil {
			return nil, err
		} else {
			pkg = &p
		}
	}
	return pkg, nil
}
