package project

var (
	description = "draughtsman is an in-cluster agent that handles Helm based deployments."
	gitSHA      = "n/a"
	name        = "draughtsman"
	source      = "https://github.com/giantswarm/draughtsman"
	version     = "n/a"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
