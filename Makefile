site: _site/index.html

index.md: beanstalk.go
	_bin/gen beanstalk > $@.part
	mv $@.part $@

beanstalk.go: force
	git cat-file -p master:$@ > $@.part
	mv $@.part $@

_site/index.html: index.md
	rm -rf _site
	jekyll

server: index.md
	rm -rf _site
	jekyll --server

.PHONY: force site server
