AC_PREREQ(2.60)
AC_INIT([libkcoidc], [1.0.0], [development@kopano.io])

AC_PROG_CC
if ! which "${CC}" >/dev/null ; then
    AC_MSG_ERROR([Please install a C compiler before trying to build libkcoidc])
fi
AC_PROG_CXX

AC_PATH_PROG([GO],[go])
if ! which "${GO}" >/dev/null ; then
    AC_MSG_ERROR([Please install Go (https://golang.org) before trying to build libkcoidc])
fi

AC_PATH_PROG([DEP],[dep])
if ! which "${DEP}" >/dev/null ; then
    AC_MSG_ERROR([Please install dep (https://golang.github.io/dep/) before trying to build libkcoidc])
fi

AC_PATH_PROG([GOFMT],[gofmt])
AC_PATH_PROG([GOLINT],[golint])
AC_PATH_PROG([GO2XUNIT],[go2xunit])

AC_ARG_WITH([python],
	[AS_HELP_STRING([--with-python=PATH], [path to Python binary])],
	[PYTHON=$withval],
	[AC_PATH_PROGS([PYTHON],[python3 python python2],[python])])

AC_OUTPUT(Makefile)
