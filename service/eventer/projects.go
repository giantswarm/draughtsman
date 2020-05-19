package eventer

var (
	commonProjectList = []string{
		"api",
		"cluster-operator",
		"credentiald",
		"draughtsman",
		"happa",
		"passage",
		"vault-exporter",
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
		"anubis":  []string{"pv-cleaner-operator"},
		"argali":  []string{"route53-manager"},
		"axolotl": []string{"route53-manager"},
		"centaur": []string{"pv-cleaner-operator"},
		"gauss":   []string{"release-bot"},
		"giraffe": []string{"route53-manager"},
	}
)
