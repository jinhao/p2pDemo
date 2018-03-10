PREFIX=./
DESTDIR=
GOFLAGS=
BINDIR=${PREFIX}/bin

BLDDIR = build
EXT=
ifeq (${GOOS},windows)
	    EXT=.exe
	endif

APPS = server client
all: $(APPS)

$(BLDDIR)/server:        $(wildcard apps/server/*.go       proto/*.go       )
$(BLDDIR)/client:        $(wildcard apps/client/*.go       proto/*.go       )

$(BLDDIR)/%:
		@mkdir -p $(dir $@)
		go build ${GOFLAGS} -o $@ ./apps/$*

$(APPS): %: $(BLDDIR)/%

clean:
		rm -fr $(BLDDIR)

.PHONY: install clean all
	.PHONY: $(APPS)
