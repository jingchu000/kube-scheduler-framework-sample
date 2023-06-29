package networktraffic

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
)

const (
	nodeMeasureQueryTemplate = "sum_over_time(node_network_receive_bytes_total{kubernetes_node=\"%s\",device=\"%s\"}[%s])"
)

type PrometheusHandler struct {
	networkInterface string
	timeRange        time.Duration
	address          string
	api              v1.API
}

func NewPrometheus(address, networkInterface string, timeRange time.Duration) *PrometheusHandler {
	client, err := api.NewClient(api.Config{Address: address})
	if err != nil {
		klog.Fatal("[NewWorkTraffic] Error creating prometheuss client :%s", err.Error())
	}
	return &PrometheusHandler{
		networkInterface: networkInterface,
		timeRange:        timeRange,
		address:          address,
		api:              v1.NewAPI(client),
	}
}

func (p *PrometheusHandler) GetNodeBandwidthMeasure(node string) (*model.Sample, error) {
	query := getNodeBandwidthQuery(node, p.networkInterface, p.timeRange)
	res, err := p.query(query)
	if err != nil {
		return nil, fmt.Errorf("[New] Error querying prometheus:%w", err)
	}

	nodeMeasure := res.(model.Vector)
	if len(nodeMeasure) != 1 {
		return nil, fmt.Errorf("[New] Invalid response, expected 1 value , got :%d", len(nodeMeasure))
	}
	return nodeMeasure[0], nil
}

func getNodeBandwidthQuery(node, networkInterface string, timeRage time.Duration) string {
	return fmt.Sprintf(nodeMeasureQueryTemplate, node, networkInterface, timeRage)
}

func (p *PrometheusHandler) query(query string) (model.Value, error) {
	results, warnings, err := p.api.Query(context.Background(), query, time.Now())

	if len(warnings) > 0 {
		klog.Warningf("[New]  Warnings:%v \n", warnings)
	}
	return results, err
}
