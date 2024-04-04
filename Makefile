.PHONY: build
build:
	@./scripts/2_build_zip_upload.sh

deploy: cd scripts
	@ ./3_terraform_apply_auto_version.sh
