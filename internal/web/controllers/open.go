package controllers

/*
func PostOpenClusterPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)
	p := getPackageContext(r).request

	var pkg v1alpha1.ClusterPackage
	if err := pkgClient.ClusterPackages().Get(ctx, p.manifestName, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch clusterpackage %v: %w", p.manifestName, err)))
		return
	}
	handleOpen(ctx, w, &pkg)
}

func PostOpenPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pkgClient := clicontext.PackageClientFromContext(ctx)
	p := getPackageContext(r).request

	var pkg v1alpha1.Package
	if err := pkgClient.Packages(p.namespace).Get(ctx, p.name, &pkg); err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to fetch package %v/%v: %w", p.namespace, p.name, err)))
		return
	}
	handleOpen(ctx, w, &pkg)
}

func handleOpen(ctx context.Context, w http.ResponseWriter, pkg ctrlpkg.Package) {
	pkgClient := clicontext.PackageClientFromContext(ctx)
	fwName := cache.NewObjectName(pkg.GetNamespace(), pkg.GetName()).String()
	s.forwardersMutex.Lock()
	defer s.forwardersMutex.Unlock()
	if result, ok := s.forwarders[fwName]; ok {
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		return
	}

	result, err := open.NewOpener().Open(ctx, pkg, "", s.Host, 0)
	if err != nil {
		responder.SendToast(w, toast.WithErr(fmt.Errorf("failed to open %v: %w", pkg.GetName(), err)))
	} else {
		s.forwarders[fwName] = result
		result.WaitReady()
		_ = cliutils.OpenInBrowser(result.Url)
		w.WriteHeader(http.StatusAccepted)

		go func() {
			ctx = context.WithoutCancel(ctx)
		resultLoop:
			for {
				select {
				case <-s.stopCh:
					break resultLoop
				case fwErr := <-result.Completion:
					// note: this does not happen in "realtime" (e.g. when the forwarded-to-pod is deleted), but only
					// the next time a connection on that port is requested, e.g. when the user reloads the forwarded
					// page or clicks open again – only then we will end up here.
					if fwErr != nil {
						fmt.Fprintf(os.Stderr, "forwarder %v completed with error: %v\n", fwName, fwErr)
					} else {
						fmt.Fprintf(os.Stderr, "forwarder %v completed without error\n", fwName)
					}
					// try to re-open if the package is still installed
					if pkg.IsNamespaceScoped() {
						var p v1alpha1.Package
						if err := pkgClient.Packages(pkg.GetNamespace()).Get(ctx, pkg.GetName(), &p); err != nil {
							if !apierrors.IsNotFound(err) {
								fmt.Fprintf(os.Stderr, "can not reopen %v: %v\n", fwName, err)
							}
							break resultLoop
						} else {
							pkg = &p
						}
					} else {
						var cp v1alpha1.ClusterPackage
						if err := pkgClient.ClusterPackages().Get(ctx, pkg.GetName(), &cp); err != nil {
							if !apierrors.IsNotFound(err) {
								fmt.Fprintf(os.Stderr, "can not reopen %v: %v\n", fwName, err)
							}
							break resultLoop
						} else {
							pkg = &cp
						}
					}
					result, err = open.NewOpener().Open(ctx, pkg, "", s.Host, 0)
					s.forwardersMutex.Lock()
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to reopen forwarder %v: %v\n", fwName, err)
						delete(s.forwarders, fwName)
						s.forwardersMutex.Unlock()
						break resultLoop
					} else {
						fmt.Fprintf(os.Stderr, "reopened forwarder %v\n", fwName)
						s.forwarders[fwName] = result
						s.forwardersMutex.Unlock()
					}
				}
			}
		}()
	}
}
*/
