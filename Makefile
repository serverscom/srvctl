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
	mockgen --destination ./internal/mocks/ssl_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/ssl_certificates.go
	mockgen --destination ./internal/mocks/load_balancers_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/load_balancers.go
	mockgen --destination ./internal/mocks/load_balancer_clusters_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/load_balancer_clusters.go
	mockgen --destination ./internal/mocks/racks_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/racks.go
	mockgen --destination ./internal/mocks/invoices_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/invoices.go
	mockgen --destination ./internal/mocks/account_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/accounts.go
	mockgen --destination ./internal/mocks/locations_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/locations.go
	mockgen --destination ./internal/mocks/kubernetes_clusters_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/kubernetes_clusters.go
	mockgen --destination ./internal/mocks/l2_segment_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/l2_segments.go
	mockgen --destination ./internal/mocks/network_pool_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/network_pools.go
	mockgen --destination ./internal/mocks/cloud_instances_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/cloud_computing_instances.go
	mockgen --destination ./internal/mocks/cloud_computing_regions_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/cloud_computing_regions.go
	mockgen --destination ./internal/mocks/cloud_block_storage_volumes_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/cloud_block_storage_volumes.go
	mockgen --destination ./internal/mocks/cloud_block_storage_backups_service.go --package=mocks --source ./vendor/github.com/serverscom/serverscom-go-client/pkg/cloud_block_storage_backups.go
	sed -i '' 's|github.com/serverscom/srvctl/vendor/github.com/serverscom/serverscom-go-client/pkg|github.com/serverscom/serverscom-go-client/pkg|g' \
	./internal/mocks/ssh_service.go \
	./internal/mocks/hosts_service.go \
	./internal/mocks/ssl_service.go \
	./internal/mocks/load_balancers_service.go \
	./internal/mocks/load_balancer_clusters_service.go \
	./internal/mocks/racks_service.go \
	./internal/mocks/invoices_service.go \
	./internal/mocks/account_service.go \
	./internal/mocks/collection.go \
	./internal/mocks/locations_service.go \
	./internal/mocks/kubernetes_clusters_service.go \
	./internal/mocks/l2_segment_service.go \
	./internal/mocks/network_pool_service.go \
	./internal/mocks/cloud_instances_service.go \
	./internal/mocks/cloud_computing_regions_service.go \
	./internal/mocks/cloud_block_storage_volumes_service.go \
	./internal/mocks/cloud_block_storage_backups_service.go

docs:
  go run cmd/gendoc/main.go
