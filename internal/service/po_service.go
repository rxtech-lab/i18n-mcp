package service

import "github.com/leonelquinteros/gotext"

type PoService struct {
	poFile *gotext.Po
}

func NewPoService(poFile gotext.Po) *PoService {
	return &PoService{poFile: &poFile}
}
