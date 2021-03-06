/*
Copyright 2021 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package observers

import (
	"fmt"
	"time"

	flaggerv1 "github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/fluxcd/flagger/pkg/metrics/providers"
)

var osmQueries = map[string]string{
	"request-success-rate": `
    sum(
        rate(
            osm_request_total{
				destination_namespace="{{ namespace }}",
				destination_kind="Deployment",
				destination_name="{{ target }}",
				response_code!~"5.*"
            }[{{ interval }}]
        )
    )
    /
    sum(
        rate(
            osm_request_total{
				destination_namespace="{{ namespace }}",
				destination_kind="Deployment",
				destination_name="{{ target }}"
            }[{{ interval }}]
        )
    )
	* 100`,
	"request-duration": `
	histogram_quantile(
		0.99,
		sum(
		  rate(
			osm_request_duration_ms_bucket{
				destination_namespace="{{ namespace }}",
				destination_kind="Deployment",
				destination_name="{{ target }}"
			}[{{ interval }}]
		  )
		) by (le)
	)`,
}

type OsmObserver struct {
	client providers.Interface
}

func (ob *OsmObserver) GetRequestSuccessRate(model flaggerv1.MetricTemplateModel) (float64, error) {
	query, err := RenderQuery(osmQueries["request-success-rate"], model)
	if err != nil {
		return 0, fmt.Errorf("rendering query failed: %w", err)
	}

	value, err := ob.client.RunQuery(query)
	if err != nil {
		return 0, fmt.Errorf("running query failed: %w", err)
	}

	return value, nil
}

func (ob *OsmObserver) GetRequestDuration(model flaggerv1.MetricTemplateModel) (time.Duration, error) {
	query, err := RenderQuery(osmQueries["request-duration"], model)
	if err != nil {
		return 0, fmt.Errorf("rendering query failed: %w", err)
	}

	value, err := ob.client.RunQuery(query)
	if err != nil {
		return 0, fmt.Errorf("running query failed: %w", err)
	}

	ms := time.Duration(int64(value)) * time.Millisecond
	return ms, nil
}
