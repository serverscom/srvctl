deps:
	go mod tidy
	go mod vendor

test: deps
	go test ./...  -mod vendor -coverprofile cp.out

generate: deps
	go generate ./...
	mockgen --destination ./internal/mocks/collection.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/collection.go	
	mockgen --destination ./internal/mocks/hosts_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/hosts.go
	mockgen --destination ./internal/mocks/ssh_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/ssh_keys.go
	sed -i '' 's|github.com/serverscom/srvctl/vendor/github.com/serverscom/serverscom-go-client/pkg|github.com/serverscom/serverscom-go-client/pkg|g' \
	./internal/mocks/ssh_service.go \
	./internal/mocks/hosts_service.go \
	./internal/mocks/collection.go

