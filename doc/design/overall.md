# overall design

Initial requirements:
- Polling for GitHub deployment events
- Download charts from a registry
- Handling installation configuration from a configmap to charts
- Install charts to the cluster

Design decisions:
- Use Helm charts for application definition, using a standard format is better than defining our own
- Use Tiller for application templating and installation, this should save us replicating effort on templating and installation
- Use GitHub deployment events for launching deployments, this is a decision that should fit in with our overall workflow
- Open source by default. Don't force a specific configuration format, that is, organisations should be able to specify their own configuration schema

Open questions:
- What does the chart registry look like? GitHub repository or CNR registry? How do charts get there?
- Are we okay with not abstracting Helm charts? We can abstract Tiller and GitHub deployment events, but abstracting charts feels wrong
- What does the bootstrap process look like? Do we install `draughtsman` as part of Kubernetes provisioning, and then use it to update itself?
- Do we need to define some GitHub deployment event format, to interface with a deployment bot?

