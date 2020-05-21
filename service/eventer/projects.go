package eventer

var (
	commonProjectList = []string{
		"api",
		"cluster-operator",
		"credentiald",
		"draughtsman",
		"happa",
		"passage",
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
		"azure-operator",
	}
	perInstallationProjectLists = map[string][]string{
		"argali":  []string{"route53-manager"},
		"axolotl": []string{"route53-manager"},
		"gaia":    []string{"app-operator"},
		"gauss":   []string{"gauss, release-bot"},
		"geckon":  []string{"app-operator"},
		"ghost":   []string{"app-operator"},
		"ginger":  []string{"app-operator"},
		"giraffe": []string{
			"app-operator",
			"route53-manager",
		},
		"godsmack": []string{"app-operator"},
		"gorgoth":  []string{"app-operator"},
	}
)
