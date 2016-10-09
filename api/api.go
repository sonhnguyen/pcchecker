package api

import (
	"github.com/sonhnguyen/pcchecker/mlabConnector"
	. "github.com/sonhnguyen/pcchecker/model"
)

func GetAllDocs() ([]PcItem, error) {
	pcItems, err := mlabConnector.GetMLab()
	if err != nil {
		return nil, err
	}
	return pcItems, nil
}
