site: _site/index.html

index.html: beanstalk.go
	_bin/gen beanstalk > $@.part
	mv $@.part $@

beanstalk.go: force
	git cat-file -p master:$@ > $@.part
	mv $@.part $@

_site/index.html: index.html
	rm -rf _site
	jekyll

server: index.html
	rm -rf _site
	jekyll --server

.PHONY: force site server
