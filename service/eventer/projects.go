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
		"gauss":   []string{"release-bot"},
		"giraffe": []string{
			"app-operator",
			"route53-manager",
		},
	}
)
