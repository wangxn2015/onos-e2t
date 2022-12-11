module github.com/wangxn2015/onos-e2t

go 1.16

require (
	github.com/atomix/atomix-go-client v0.6.2
	github.com/atomix/atomix-go-framework v0.10.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.2.0
	github.com/onosproject/onos-api/go v0.9.29
	github.com/onosproject/onos-test v0.6.4
	github.com/prometheus/common v0.32.1
	github.com/stretchr/testify v1.7.2
	github.com/wangxn2015/helmit v1.6.21-0.20221211101934-72a55699c433
	github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_kpm_v2_go v1.8.91-0.20221211112622-747f3a81e726 // indirect
	github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc_pre_go v0.8.11-0.20221211112622-747f3a81e726 // indirect
	//github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_kpm_v2_go v0.8.6
	//github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc_pre_go v0.8.6
	github.com/wangxn2015/onos-lib-go v0.8.16-0.20221211095953-88ba820eb386
	github.com/wangxn2015/onos-ric-sdk-go v0.8.9-0.20221211103601-f383199a0bf2 // indirect
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.24.2
	k8s.io/apimachinery v0.24.2
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
)

//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_kpm_v2_go => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_kpm_v2_go
//
//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_mho_go => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_mho_go
//
//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_rc
//
//replace github.com/wangxn2015/onos-e2-sm/servicemodels/e2sm_rc_pre_go => /home/baicells/go_project/modified-onos-module/onos-e2-sm/servicemodels/e2sm_rc_pre_go

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
