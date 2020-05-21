package eventer

var (
	commonProjectList = []string{
		"api",
		"cluster-operator",
		"draughtsman",
		"happa",
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
		"gauss": []string{"release-bot"},
	}
)
