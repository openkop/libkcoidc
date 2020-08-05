AC_PREREQ(2.60)
AC_INIT([libkcoidc], m4_esyscmd([scripts/git-version-gen .version]), [development@kopano.io])

AC_CONFIG_MACRO_DIR([m4])

AC_CONFIG_FILES([libkcoidc.pc])

AC_PROG_CC
if ! which "${CC}" >/dev/null ; then
    AC_MSG_ERROR([Please install a C compiler before trying to build libkcoidc])
fi
AC_PROG_CXX

AC_PATH_PROG([GO],[go])
if ! which "${GO}" >/dev/null ; then
    AC_MSG_ERROR([Please install Go (https://golang.org) before trying to build libkcoidc])
fi

GO_VERSION_MIN=1.13
GO_VERSION=$(${GO} version | sed 's/^go version go//' | sed 's/ .*//')
AX_COMPARE_VERSION([$GO_VERSION], [ge], [$GO_VERSION_MIN],
    AC_MSG_NOTICE([Go ${GO_VERSION} found]),
    AC_MSG_ERROR([Please install Go ${GO_VERSION_MIN} or newer before trying to build libkcoidc (found Go ${GO_VERSION}).]))

AC_PATH_PROG([GOFMT],[gofmt])
AC_PATH_PROG([GOLINT],[golangci-lint])
AC_PATH_PROG([GO2XUNIT],[go2xunit])

AC_PATH_PROGS([PYTHON],[python3 python python2],[python])

AC_OUTPUT(Makefile)
