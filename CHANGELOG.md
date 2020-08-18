# CHANGELOG

## Unreleased



## v0.9.2 (2020-08-18)

- Add changelog support
- Build with all hardening flags turned on


## v0.9.1 (2020-08-17)

- Allow to set additional EXTLDFLAGS
- Eliminate build warnings in Python module
- Fix invalid include path when building the Python module
- Ensure to create pkg-config folder on install


## v0.9.0 (2020-08-05)

- Fix compile warning in example
- Fix unit test Jenkins reporting
- Move Jenkins pipleline to warnings-ng plugin
- Bump C API version and include real version in pc
- Add support to set logger to Go and C API
- Fix linter
- Replace static global C constants with #define
- Add support for pkg-config
- Do not strip C library so by default, leave that to packaging
- Add version information and expose in Go and C API
- Use SPDX license identifier header
- Fix license text to match the template exactly


## v0.8.2 (2020-07-14)

- Update license ranger
- Build with Go 1.14.4
- Remove obsolete files
- Remove Go 1.11 and 1.12 from Travis
- Build with Go 1.13.3


## v0.8.1 (2019-09-18)

- Update external dependencies
- Force to use Go modules for Go < 1.13
- Fixup travis


## v0.8.0 (2019-09-18)

- Update to Go 1.13 and use Go modules


## v0.7.2 (2019-08-26)

- Add comment about using Debian Stretch to build
- Update Go Dep to 0.5.4
- Compile on stretch to ensure compatibility


## v0.7.1 (2019-08-23)

- Ensure provider.Definition is never nil
- Remove Go 1.8, since its no longer supported by libkcoidc


## v0.7.0 (2019-04-25)

- Bump dep to 0.5.1
- Add TLS session resumption support and ensure http2
- Add HTTP client settings with proxy support


## v0.6.0 (2019-04-23)

- Use debug log string prefix for logger of c library
- Silence default logger of C library
- Remove obsolete file
- Fixup TravisCI
- Fixup TravisCI
- Use dep with TravisCI
- Update Go dependencies
- Update docs/configure to match
- Use oidc-go
- Ignore .vscode
- Migrate from Glide to Dep
- Log using logger instead of printing directly
- Add Go 1.12, remove 1.x
- Cleanup
- Add import comment
- Make linter happy
- Add unit test to show all kcoidc errors


## v0.5.0 (2019-02-04)

- Add support to retrieve guest claim info to Go API


## v0.4.4 (2019-01-24)

- Bump copyright year to 2019
- Add Go 1.11
- Add Go 1.10 as minimal requirement
- Lint after built to ensure that dependencies have loaded
- Print Go version in Jenkins


## v0.4.3 (2018-09-21)

- Ensure correct salt length of RSA-PSS signing methods


## v0.4.2 (2018-09-17)

- Fix segfault when API is used without initialization


## v0.4.1 (2018-09-06)

- Ensure identity claims value


## v0.4.0 (2018-09-06)

- Build on Jenkins with Go 1.10
- Add validation with required scopes check
- Validate response content-type header
- Add -f to mv so rebuilds work without interactive prompt
- Use backend identifier claim for ID in validation
- Create symlink in .lib so our own C stuff compiles


## v0.3.0 (2018-06-04)

- Add soname to shared c-lib


## v0.2.1 (2018-05-24)

- Fixup building examples with Go 1.10


## v0.2.0 (2018-03-12)

- Set permissions in dist tarball and include symlink too
- Install library properly so linking works
- Really use Go 1.10
- Add Go 1.10
- Update to Go 1.9
- Add 3rd party license information


## v0.1.0 (2018-02-09)

- Use autoconf
- Add make install/uninstall targets
- Add vanity import and Travis CI


## v0.0.1 (2018-02-06)

- Fix up python module memory and some warnings
- Update benchmark results
- Update Python validate example with claims
- Add support get claims and to fetch userinfo
- Fix typo 'untill' -> 'until'
- Add Jenkinsfile
- Include Python and Go in README
- Auto detect cpu count
- Add python make target
- Implement Go importable module
- Implement Python wrapper module
- Fix time measurement in examples
- Do not run linter in default target
- Properly end start function after success
- Add simple benchmark cpp example
- Fix token validation check
- Add some more validations
- Update README
- Add real validation
- Reorganize project
- Auto generate more stuff in Makefile
- Update for .go files
- Initial commit

