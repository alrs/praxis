.PHONY: release

release:
	mkdir -p /tmp/release/$(VERSION)
	jq '.Mappings.Constants.Values.Version|= "$(VERSION)"' formation/rack.json > /tmp/release/$(VERSION)/rack.json
	aws s3 cp /tmp/release/$(VERSION)/rack.json s3://praxis-releases/release/$(VERSION)/formation/rack.json --acl public-read
	aws s3 cp formation/network.json s3://praxis-releases/release/$(VERSION)/formation/network.json --acl public-read
	aws s3 cp formation/resource/*.json s3://praxis-releases/release/$(VERSION)/formation/resource/ --acl public-read
	aws s3 cp /tmp/release/$(VERSION)/rack.json s3://praxis-releases/release/latest/formation/rack.json --acl public-read
	aws s3 cp formation/network.json s3://praxis-releases/release/latest/formation/network.json --acl public-read
	aws s3 cp formation/resource/*.json s3://praxis-releases/release/latest/formation/resource/ --acl public-read
	
	make -C lambda/lifecycle release VERSION=$(VERSION)
