index.md: beanstalk.go
	(\
	  printf -- '---\n---\n';\
	  godoc -html beanstalk;\
	) > $@.part
	mv $@.part $@

beanstalk.go: force
	git cat-file -p master:$@ > $@.part
	mv $@.part $@

.PHONY: force
