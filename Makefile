# Setup basic variables
GOROOT=$(shell go env | grep GOROOT | awk -F= '{print substr($$2,2,length($$2)-2)}')
#GOROOT=/usr/local/go
export GO=$(GOROOT)/bin/go

APP=app
CORE=server
SERVICES=provisioning

services=$(SERVICES)
core=$(CORE)

.PHONY: clean $(CORE) $(SERVICES)

# The axlrate app
$(APP): $(CORE) $(SERVICES)
	$(MAKE) -C $@

# Clean
clean:
	# Clean the app
	$(MAKE) -C app $@
	# Clean the core package sub-directories
	for subdir in $(CORE); do \
		$(MAKE) -C core/$$subdir $@; \
	done
	# Clean the service package sub-directories
	for subdir in $(SERVICES); do \
		$(MAKE) -C service/$$subdir $@; \
	done

# Build the core server package
$(CORE): 
	$(MAKE) -C core/$@

# Build the services
services: $(SERVICES)

$(SERVICES):
	$(MAKE) -C service/$@