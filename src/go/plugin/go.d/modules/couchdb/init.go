// SPDX-License-Identifier: GPL-3.0-or-later

package couchdb

import (
	"errors"
	"net/http"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/web"
)

func (cdb *CouchDB) validateConfig() error {
	if cdb.URL == "" {
		return errors.New("URL not set")
	}
	if cdb.Node == "" {
		return errors.New("'node' not set")
	}
	if _, err := web.NewHTTPRequest(cdb.RequestConfig); err != nil {
		return err
	}
	return nil
}

func (cdb *CouchDB) initHTTPClient() (*http.Client, error) {
	return web.NewHTTPClient(cdb.ClientConfig)
}

func (cdb *CouchDB) initCharts() (*Charts, error) {
	charts := module.Charts{}

	if err := charts.Add(*dbActivityCharts.Copy()...); err != nil {
		return nil, err
	}
	if err := charts.Add(*httpTrafficBreakdownCharts.Copy()...); err != nil {
		return nil, err
	}
	if err := charts.Add(*serverOperationsCharts.Copy()...); err != nil {
		return nil, err
	}
	if len(cdb.databases) != 0 {
		dbCharts := dbSpecificCharts.Copy()

		if err := charts.Add(*dbCharts...); err != nil {
			return nil, err
		}

		for _, chart := range *dbCharts {
			for _, db := range cdb.databases {
				if err := chart.AddDim(&module.Dim{ID: "db_" + db + "_" + chart.ID, Name: db}); err != nil {
					return nil, err
				}
			}
		}

	}
	if err := charts.Add(*erlangStatisticsCharts.Copy()...); err != nil {
		return nil, err
	}

	if len(charts) == 0 {
		return nil, errors.New("zero charts")
	}
	return &charts, nil
}
