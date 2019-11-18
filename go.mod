module github.com/Benbentwo/bb

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

go 1.12.4

require (
	github.com/AlecAivazis/survey/v2 v2.0.4
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667 // indirect
	github.com/Pallinder/go-randomdata v1.2.0 // indirect
	github.com/fatih/color v1.7.0
	github.com/heptio/sonobuoy v0.16.1 // indirect
	github.com/jenkins-x/jx v0.0.0-20191002101425-246bdbf20015
	github.com/petergtz/pegomock v2.6.0+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/rickar/props v0.0.0-20170718221555-0b06aeb2f037
	github.com/shurcooL/githubv4 v0.0.0-20190718010115-4ba037080260 // indirect
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	go.opencensus.io v0.21.0 // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20190718183219-b59d8169aab5 // indirect
	k8s.io/apiextensions-apiserver v0.0.0-20190718185103-d1ef975d28ce // indirect
	k8s.io/apimachinery v0.0.0-20190703205208-4cfb76a8bf76
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible // indirect
	sigs.k8s.io/yaml v1.1.0
)
