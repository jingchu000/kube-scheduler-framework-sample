package networktraffic

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"

	"github.com/jingchu000/kube-scheduler-framework-sample/pkg/apis/config"
)

const Name = "NetworkTraffic"

var _ = framework.ScorePlugin(&NetworkTraffic{})

// NetworkTraffic is a score plugin that favors nodes based on their
// network traffic amount. Nodes with less traffic are favored.
// Implements framework.ScorePlugin
type NetworkTraffic struct {
	handle     framework.Handle
	prometheus *PrometheusHandler
}

func (n *NetworkTraffic) Name() string {
	return Name
}

func New(obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	args := config.NetworkTrafficArgs{}
	if obj != nil {
		klog.Infof("[Net] obj != nil ,got %v", obj)
		err := frameworkruntime.DecodeInto(obj, &args)
		if err != nil {
			klog.Infof("[Net]  frameworkruntime.DecodeInto(obj, &args),err %v", err)
			return nil, fmt.Errorf("[Net] want args to be of type Net ,got %v", err)
		}
	}

	klog.Infof("[Net] args received.NetworkInterface:%s;  TimeRangeInMinutes: %d,Address:%s", args.NetworkInterface, args.TimeRangeInMinutes, args.Address)
	return &NetworkTraffic{
		handle:     h,
		prometheus: NewPrometheus(args.Address, args.NetworkInterface, time.Minute*time.Duration(args.TimeRangeInMinutes)),
		//prometheus: nil,
	}, nil
}

func (n *NetworkTraffic) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeBandwidth, err := n.prometheus.GetNodeBandwidthMeasure(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error getting node bandwidth measure:%s ,", err))
	}
	klog.Info("[Net] node'%s' bandwidth:%s", nodeName, nodeBandwidth.Value)
	return int64(nodeBandwidth.Value), nil
}
func (n *NetworkTraffic) ScoreExtensions() framework.ScoreExtensions {
	return n
}

func (n *NetworkTraffic) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var higherScore int64

	for _, node := range scores {
		klog.Infof("[Net]  Nodes final NodeName:%v,Node.score:%v", node.Name, node.Score)
		if higherScore < node.Score {
			higherScore = node.Score
		}
	}

	for i, node := range scores {
		scors := framework.MaxNodeScore - (node.Score * framework.MaxNodeScore / higherScore)
		scores[i].Score = scors
	}
	klog.Infof("[Net]  Nodes final score:%v", scores)
	return nil
}
