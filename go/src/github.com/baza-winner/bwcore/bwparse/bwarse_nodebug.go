// +build !debug

package bwparse

type Start struct {
	ps     *PosInfo
	suffix string
}

func (p *P) Stop(start *Start) {
	delete(p.starts, start.ps.pos)
}

func (p *proxy) Stop(start *Start) {
	delete(p.starts, start.ps.pos)
}
