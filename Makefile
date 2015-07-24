FIND ?= find

fmt:
	        $(FIND) ./src -iname \*.go -print0 | xargs -0 -I {} dirname {} | uniq | xargs go fmt
