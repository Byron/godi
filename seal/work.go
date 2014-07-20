package seal

import (
	"sync"

	"github.com/Byron/godi/api"
)

func (s *SealCommand) Gather(files <-chan godi.FileInfo, results chan<- godi.Result, wg *sync.WaitGroup, done <-chan bool) {
	makeResult := func(f *godi.FileInfo) (godi.Result, *godi.BasicResult) {
		res := godi.BasicResult{
			Finfo: f,
			Prio:  godi.Progress,
		}
		return &res, &res
	}

	godi.Gather(files, results, wg, done, makeResult, &s.pCtrl)
}
