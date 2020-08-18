package configuration

var (
	commonProjectList = []string{
		"cluster-operator",
		"draughtsman",
	}

	conformanceProjectList = []string{
		"conformance-app-collection",
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
	if installation == "gaia" {
		projects := []string{}
		projects = append(projects, commonProjectList...)
		projects = append(projects, providerProjectList[provider]...)
		projects = append(projects, conformanceProjectList...)
		return projects
	}

	return append(commonProjectList, providerProjectList[provider]...)
}
