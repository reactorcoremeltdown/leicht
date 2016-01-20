GOC := /usr/bin/go build
FETCHLIBS=/usr/bin/go get

MAKEMAN=/usr/bin/ronn --roff --pipe

BUILDDIR=$(CURDIR)/build
SRCDIR=src/leicht
MANSRCDIR=src/man
GOBINDIR=$(BUILDDIR)/bin
GOPATHDIR=$(BUILDDIR)/golibs

INSTALL=install
INSTALL_BIN=$(INSTALL) -m755
INSTALL_LIB=$(INSTALL) -m644
INSTALL_CONF=$(INSTALL) -m400

PREFIX?=$(DESTDIR)/usr
SYSTEMDCONFDIR?=$(DESTDIR)/lib/systemd/system
BINDIR?=$(PREFIX)/bin
LIBDIR?=$(PREFIX)/lib/leicht
CONFDIR?=$(DESTDIR)/etc/leicht
MANPAGEDIR?=$(DESTDIR)/usr/share/man

all: leicht manpages

leicht: Makefile src/leicht/main.go
	mkdir -p $(GOPATHDIR) && \
	mkdir -p $(GOBINDIR) && \
	export GOPATH=$(GOPATHDIR) && \
	export GOBIN=$(GOBINDIR) && \
	cd $(SRCDIR) && \
	$(FETCHLIBS) && \
	$(GOC)

manpages:
	export LANG=en_US.UTF-8 && \
	export LC_ALL=en_US.UTF-8 && \
	cd $(MANSRCDIR) && \
	$(MAKEMAN) leicht.1.md > leicht.1

clean:
	rm -fr build/
	rm -f src/leicht/leicht
	rm -f src/man/leicht.1

install:
	mkdir -p $(BINDIR)
	mkdir -p $(LIBDIR)
	mkdir -p $(CONFDIR)
	mkdir -p $(SYSTEMDCONFDIR)
	mkdir -p $(MANPAGEDIR)/man1
	$(INSTALL_BIN) src/leicht/leicht $(BINDIR)/
	$(INSTALL_LIB) src/lib/bash/leicht.sh $(LIBDIR)/
	$(INSTALL_LIB) src/lib/node.js/leicht.js $(LIBDIR)/
	$(INSTALL_LIB) src/man/leicht.1 $(MANPAGEDIR)/man1/
	$(INSTALL_CONF) src/leicht/default.json $(CONFDIR)/
	$(INSTALL_CONF) src/init/leicht.target $(SYSTEMDCONFDIR)/
