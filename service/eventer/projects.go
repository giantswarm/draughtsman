package eventer

var (
	commonProjectList = []string{
		"cluster-operator",
		"draughtsman",
	}
	awsProjectList = []string{
		"aws-app-collection",
		"aws-operator",
	}
	kvmProjectList = []string{
		"kvm-app-collection",
		"kvm-operator",
	}
	azureProjectList = []string{
		"azure-app-collection",
	}
	perInstallationProjectLists = map[string][]string{
		"gaia": []string{"app-operator", "conformance-app-collection"},
		"gauss": []string{
			"app-operator",
		},
		"geckon": []string{"app-operator"},
		"ghost":  []string{"app-operator"},
		"ginger": []string{"app-operator"},
		"giraffe": []string{
			"app-operator",
		},
		"godsmack": []string{"app-operator"},
		"gorgoth":  []string{"app-operator"},
	}
)
