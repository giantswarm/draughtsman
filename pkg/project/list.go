package project

var (
	commonProjectList = []string{
		"cluster-operator",
		"draughtsman",
	}

	providerProjectList = map[string][]string{
		"aws": {
			"aws-app-collection",
			"aws-operator",
		},
		"kvm": {
			"kvm-app-collection",
			"kvm-operator",
		},
		"azure": {
			"azure-app-collection",
			"azure-operator",
		},
	}

	perInstallationProjectLists = map[string][]string{
		"gaia": {"app-operator"},
		"gauss": {
			"app-operator",
		},
		"geckon": {"app-operator"},
		"ghost":  {"app-operator"},
		"ginger": {"app-operator"},
		"giraffe": {
			"app-operator",
		},
		"godsmack": {"app-operator"},
		"gorgoth":  {"app-operator"},
	}
)

func GetProjectList(provider, installation string) []string {
	list := append(commonProjectList, providerProjectList[provider]...)
	return append(list, perInstallationProjectLists[installation]...)
}
