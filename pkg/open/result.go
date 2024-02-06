package open

import "github.com/glasskube/glasskube/pkg/future"

type openResult struct {
	opener     *opener
	Completion future.Future
	Url        string
}

func (r *openResult) Stop() {
	r.opener.stop()
}

func (r *openResult) WaitReady() {
	for _, c := range r.opener.readyCh {
		<-c
	}
}
