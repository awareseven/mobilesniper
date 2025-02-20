SDCX_BIN_PATH = assets/sdcX-Bin/

# Check if patchelf is installed (evaluated at parse time)
ifeq ($(shell which patchelf),)
    $(error "patchelf is not installed. Please install patchelf before running this Makefile.")
endif

.PHONY: default
default:
	@find $(SDCX_BIN_PATH) -type f -executable | while IFS= read -r bin; do \
	    echo "Patching $$bin"; \
	    patchelf --set-rpath '$$ORIGIN' "$$bin"; \
	done