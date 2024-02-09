package open

import "github.com/glasskube/glasskube/pkg/future"

type OpenResult struct {
	opener     *opener
	Completion future.Future
	Url        string
}

func (r *OpenResult) Stop() {
	r.opener.stop()
}

func (r *OpenResult) WaitReady() {
	for _, c := range r.opener.readyCh {
		<-c
	}
}
