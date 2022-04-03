TEST_KW?=

debug-build:
	go build -o ./bin/terraform_provider_pfptmeta_linux_amd64 -gcflags="all=-N -l"
	./bin/terraform_provider_pfptmeta_linux_amd64 --debug

mod-tidy:
	go mod tidy


unittest:
	go test ./... $(if $(VERBOSE),"-v") -timeout 120m $(if $(TEST_KW),-run $(TEST_KW))


acc_tests:
	TF_ACC=1 go test ./... $(if $(VERBOSE),"-v") -run "TestAcc*" -timeout 120m $(if $(TEST_KW),-run $(TEST_KW))


generate:
	go generate -v -x


tests: verify_clean acc_tests unittest

# generate is necessary here because it generates the documentation from the code and formats the .go and .tf files
# we verify git is clean after that to make sure the documentation, .tf and .go files were updated
verify_clean: mod-tidy generate
	! git status -s | grep "??" || (echo "Uncommitted files found" && exit 1)
	git diff --stat --exit-code || (echo "Uncommitted files found" && exit 1)


release:
	gpg --batch --import $(GPG_SECRET_PATH) && goreleaser release --rm-dist


LAST_VERSION:=$(shell git describe --tags --abbrev=0)
VERSION_PARTS:=$(subst ., ,$(LAST_VERSION))
MAJOR:=0
MINOR:=1
PATCH_AND_LABEL:=$(word 3, $(VERSION_PARTS))
PATCH_AND_LABEL_PARTS:=$(subst -, ,$(PATCH_AND_LABEL))
PATCH:=$(word 1, $(PATCH_AND_LABEL_PARTS))
NEXT_PATCH:=$(shell echo $$(($(PATCH)+1)))

tag_version:
	git tag "$(MAJOR).$(MINOR).$(NEXT_PATCH)$(if $(LABEL),-$(LABEL))"
