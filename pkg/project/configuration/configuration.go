package configuration

var (
	commonProjectList = []string{
		"draughtsman",
	}

	providerProjectList = map[string][]string{
		"aws": {
			"aws-app-collection",
			"aws-operator",
			"cluster-operator",
		},
		"kvm": {
			"cluster-operator",
			"kvm-app-collection",
			"kvm-operator",
		},
		"azure": {
			"azure-app-collection",
		},
	}
)

func GetProjectList(provider, installation string) []string {
	return append(commonProjectList, providerProjectList[provider]...)
}
