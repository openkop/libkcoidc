from distutils.core import setup, Extension

setup(name='pykcoidc', version='1.0',
      ext_modules=[
        Extension('_pykcoidc',
                  ['pykcoidc.c'],
                  include_dirs=['../.libs/include/kcoidc'],
                  library_dirs=['../.libs'],
                  libraries=['kcoidc'])
      ],
      packages=['pykcoidc'],
      package_dir={'pykcoidc': 'src'},
      )
