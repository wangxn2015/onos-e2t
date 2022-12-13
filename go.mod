module github.com/onosproject/onos-e2t

go 1.16

require (
	cloud.google.com/go v0.93.3 // indirect
	github.com/atomix/atomix-go-client v0.6.2
	github.com/atomix/atomix-go-framework v0.10.1
	github.com/bits-and-blooms/bloom/v3 v3.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.6.3
	github.com/fatih/color v1.9.0 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/uuid v1.2.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onosproject/helmit v0.6.19
	github.com/onosproject/onos-api/go v0.9.3
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go v0.8.6
	github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go v0.8.6
	github.com/onosproject/onos-lib-go v0.8.13
	//github.com/wangxn2015/onos-lib-go v0.8.13
	//github.com/wangxn2015/onos-ric-sdk-go v0.8.7
	github.com/onosproject/onos-test v0.6.4
	github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e // indirect /////
	github.com/prometheus/common v0.26.0
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/wangxn2015/onos-lib-go v0.8.16-0.20221213045740-e38a2ad92701
	github.com/wangxn2015/onos-ric-sdk-go v0.8.9-0.20221212153731-3320d18f773e
	google.golang.org/genproto v0.0.0-20210828152312-66f60bf46e71 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/ini.v1 v1.63.2 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/utils v0.0.0-20210707171843-4b05e18ac7d9
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
