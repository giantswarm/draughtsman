package configuration

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
)

func GetProjectList(provider, installation string) []string {
	return append(commonProjectList, providerProjectList[provider]...)
}
