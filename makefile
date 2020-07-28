build: ## Build the binary file
	@GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s" test-harness-framework.go shared-data-struct.go
	@zip test-harness-framework.zip test-harness-framework
	@GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s" worker-handler.go shared-data-struct.go *Performancer.go
	@zip worker-handler.zip worker-handler

clear: ## Clear binary file
	@rm test-harness-framework worker-handler test-harness-framework.zip worker-handler.zip

run: ## Run to generate reports
	@go run auto-run.go shared-data-struct.go -bucket-name $$BUCKET_NAME -time-to-wait $$TIME_TO_WAIT

auto_provision: ## Provision infrustructures
	@terraform init
	@terraform apply -var="bucket_name=$$BUCKET_NAME" -auto-approve

auto_destroy: ## Destroy infrustructures
	@terraform destroy -var="bucket_name=$$BUCKET_NAME" -auto-approve