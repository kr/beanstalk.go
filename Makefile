index.md: beanstalk.go
	_bin/gen beanstalk > $@.part
	mv $@.part $@

beanstalk.go: force
	git cat-file -p master:$@ > $@.part
	mv $@.part $@

.PHONY: force
