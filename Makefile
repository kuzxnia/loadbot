

build-chart:
	helm package ./helm/workload
	mv workload-*.tgz ./orchiestrator/resourcemanager/workload-chart.tgz


build-dev:
	# only your platform
	goreleaser build --snapshot --clean --single-target

build-all:
	# only your platform
	goreleaser build --snapshot --clean

.PHONY: build-dev build-all build-chart
