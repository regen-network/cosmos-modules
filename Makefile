# modules to build in CI
SUBDIRS = incubator/orm  incubator/group

test: test-unit

test-unit:
	for dir in $(SUBDIRS); do \
        $(MAKE) -C "$$dir" test; \
    done

.PHONY: test test-unit