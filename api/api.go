package api

import (
    . "github.com/pcchecker/model"
	"github.com/pcchecker/mlabConnector"
)

func GetAllDocs() ([]PcItem, error){
	pcItems, err := mlabConnector.GetMLab()
	if err != nil {
		return nil, err
	}
	return pcItems, nil
}

