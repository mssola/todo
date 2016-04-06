test :: unit_test checks

# The vet, fmt and lint rules have been extracted from:
# 	github.com/docker/distribution

vet ::
	@echo "+ $@"
		@go list ./... | grep -v vendor | go vet

fmt ::
	@echo "+ $@"
		@test -z "$$(gofmt -s -l . | grep -v vendor | tee /dev/stderr)" || \
			echo "+ please format Go code with 'gofmt -s'"

lint ::
	@echo "+ $@"
		@test -z "$$(golint ./... | grep -v vendor | tee /dev/stderr)"

climate ::
	@echo "+ $@"
		@(./script/climate -o -a app && ./script/climate -o -a -t 80.0 lib)

unit_test ::
	@echo "+ $@"
		@go test -v ./...

checks :: vet fmt lint climate
