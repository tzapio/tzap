release:
	rm -rf release
	rm -rf bin
	rm -rf public
	mkdir -p release
	cp -r ../release ./
	cp -r ../public ./
	cp ../README.md README.md
	npm publish
.PHONY: release