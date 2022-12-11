module github.com/wangxn2015/onos-e2t

go 1.16

require (
	github.com/atomix/atomix-go-client v0.6.2
	github.com/atomix/atomix-go-framework v0.10.1
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/ericchiang/oidc v0.0.0-20160908143337-11f62933e071 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-logr/zapr v1.2.0 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/google/cel-go v0.10.1 // indirect
	github.com/google/uuid v1.2.0
	github.com/joncalhoun/pipe v0.0.0-20170510025636-72505674a733 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/onosproject/onos-api/go v0.9.29
	github.com/onosproject/onos-test v0.6.4
	github.com/prometheus/client_golang v1.11.1 // indirect
	github.com/prometheus/common v0.26.0
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/spf13/cobra v1.4.0 // indirect
	github.com/stretchr/testify v1.7.2
	github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_kpm_v2_go v1.8.91-0.20221211112622-747f3a81e726 // indirect
	github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc_pre_go v0.8.11-0.20221211112622-747f3a81e726 // indirect
	//github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_kpm_v2_go v0.8.6
	//github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc_pre_go v0.8.6
	github.com/wangxn2015/onos-lib-go v0.8.16-0.20221211095953-88ba820eb386
	github.com/wangxn2015/onos-ric-sdk-go v0.8.9-0.20221211103601-f383199a0bf2 // indirect
	go.etcd.io/etcd/client/v3 v3.5.1 // indirect
	google.golang.org/genproto v0.0.0-20220107163113-42d7afdf6368 // indirect
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.7.2 // indirect
	k8s.io/api v0.24.2
	k8s.io/apiextensions-apiserver v0.23.0-alpha.0 // indirect
	k8s.io/apimachinery v0.24.2
	k8s.io/client-go v0.24.2 // indirect
	k8s.io/code-generator v0.24.2 // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.0.30 // indirect
)

//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_kpm_v2_go => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_kpm_v2_go
//
//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_mho_go => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_mho_go
//
//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_rc
//
//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc_pre_go => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_rc_pre_go

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
