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
			"kvm-app-collection",
		},
		"azure": {
			"azure-app-collection",
		},
		"vmware": {
			"vmware-app-collection",
		},
	}
)

func GetProjectList(provider, installation string) []string {
	return append(commonProjectList, providerProjectList[provider]...)
}
