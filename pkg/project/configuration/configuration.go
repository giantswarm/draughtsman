package configuration

var (
	commonProjectList = []string{
		"shared-app-collection",
		"draughtsman",
	}

	providerProjectList = map[string][]string{
		"aws": {
			"aws-app-collection",
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
