include $(GOROOT)/src/Make.$(GOARCH)

TARG=beanstalk
GOFILES=\
	beanstalk.go\

include $(GOROOT)/src/Make.pkg
