index.html: beanstalk.go
	godoc -html beanstalk > $@.part
	mv $@.part $@

beanstalk.go: force
	git cat-file -p master:$@ > $@.part
	mv $@.part $@

.PHONY: force
