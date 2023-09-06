.PHONY: build
build:
	@./scripts/2_build_zip_upload.sh

deploy: cd scripts
	@ ./3_terraform_apply_auto_version.sh

.PHONY: test
test:
	@rm -rf test/main;
	@cp lambda-go/main test/main
	@sam local invoke MyLambdaFunction -t test/template.yml -e test/event.json;