test :: unit_test checks

# The vet, fmt and lint rules have been extracted from:
# 	github.com/docker/distribution

vet ::
	@echo "+ $@"
		@go vet ./...

fmt ::
	@echo "+ $@"
		@test -z "$$(gofmt -s -l . | grep -v Godeps/_workspace/src/ | tee /dev/stderr)" || \
			echo "+ please format Go code with 'gofmt -s'"

lint ::
	@echo "+ $@"
		@test -z "$$(golint ./... | grep -v Godeps/_workspace/src/ | tee /dev/stderr)"

climate ::
	@echo "+ $@"
		@(./script/climate -o -a app && ./script/climate -o -a -t 80.0 lib)

unit_test ::
	@echo "+ $@"
		@godep go test -v ./...

checks :: vet fmt lint climate
