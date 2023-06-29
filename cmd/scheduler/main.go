package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"k8s.io/component-base/logs"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"github.com/jingchu000/kube-scheduler-framework-sample/pkg/plugins/networktraffic"
	"github.com/jingchu000/kube-scheduler-framework-sample/pkg/plugins/sample"
)

/*
Reference:
https://github.com/kubernetes-sigs/scheduler-plugins
https://mp.weixin.qq.com/s/pRnXeRGw-5YpEDAjk5HF6g
https://mp.weixin.qq.com/s/OwSK-taa6zIwKsfyeudPwg
Welcome to communicate With My Chat: Base64 Code{bGItQWJpbmc=}
*/
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	logs.InitLogs()
	defer logs.FlushLogs()

	cmd := app.NewSchedulerCommand(
		app.WithPlugin(sample.Name, sample.New),
		app.WithPlugin(networktraffic.Name, networktraffic.New),
	)

	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
