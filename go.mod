module github.com/Benbentwo/gh

go 1.12

require (
	github.com/Benbentwo/go-utils v0.0.0-20191121035730-023054be6742
	github.com/fatih/color v1.7.0
	github.com/jenkins-x/jx v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.8.1
	github.com/rickar/props v0.0.0-20170718221555-0b06aeb2f037
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	sigs.k8s.io/yaml v1.1.0
)

exclude knative.dev/pkg v0.0.0-20191002055904-849fcc967b59

exclude knative.dev/pkg v0.0.0-20191001225505-346b0abf16cd

// exclude knative.dev/serving

replace k8s.io/api => k8s.io/api v0.0.0-20181128191700-6db15a15d2d3

replace k8s.io/metrics => k8s.io/metrics v0.0.0-20181128195641-3954d62a524d

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190122181752-bebe27e40fb7

replace k8s.io/client-go => k8s.io/client-go v2.0.0-alpha.0.0.20190115164855-701b91367003+incompatible

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181128195303-1f84094d7e8e

replace github.com/banzaicloud/bank-vaults => github.com/banzaicloud/bank-vaults v0.0.0-20190508130850-5673d28c46bd

replace github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v21.1.0+incompatible

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v10.15.5+incompatible

replace github.com/jenkins-x/jx => github.com/jenkins-x/jx v0.0.0-20191002101425-246bdbf20015
