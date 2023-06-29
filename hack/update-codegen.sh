#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

# gen config:v1beta1
bash "${CODEGEN_PKG}"/generate-internal-groups.sh \
  "deepcopy,conversion,defaulter" \
  github.com/jingchu000/scheduler-framework-sample/pkg/generated \
  github.com/jingchu000/scheduler-framework-sample/pkg/apis \
  github.com/jingchu000/scheduler-framework-sample/pkg/apis \
  "config:v1beta1" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate/boilerplate.generatego.txt

#$ vendor/k8s.io/code-generator/generate-groups.sh
#   "deepcopy,client,informer,lister" \ // 指定生成的内容
 #  sample-controller/pkg/generated sample-controller/pkg/apis \ // 指定客户端，lister和informer代码框架的包名
 #  samplecontroller:v1alpha1 \ // API组的基础包名
 #  --output-base "${GOPATH}/src" \  用于代码生成器查找包输出的基本路径；
 #  --go-header-file "${GOPATH}/src/sample-controller/hack/boilerplate.go.txt" // 自定义代码用到的版权信息头