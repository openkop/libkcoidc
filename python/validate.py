#!/bin/env python

from __future__ import print_function

import sys
import time

import pykcoidc


def main(args):
    iss_s = len(args) > 0 and args[0] or ""
    token_s = len(args) > 1 and args[1] or ""

    # Allow insecure operations.
    try:
        pykcoidc.insecure_skip_verify(1)
    except pykcoidc.Error as e:
        print("> Error: insecure_skip_verify failed: 0x%x" % e.args[0])
        return -1
    # Initialize with issuer identifier.
    try:
        pykcoidc.initialize(iss_s)
    except pykcoidc.Error as e:
        print("> Error: initialize failed: 0x%x" % e.args[0])
        return -1
    # Wait until oidc validation becomes ready.
    try:
        pykcoidc.wait_until_ready(10)
    except pykcoidc.Error as e:
        print("> Error: failed to get ready in time: 0x%x" % e.args[0])
        return -1

    sub = None
    err = None
    begin = time.time()
    # Validate token passed from commandline.
    try:
        sub = validate_and_get_subject(token_s)
    except pykcoidc.Error as e:
        err = e
    end = time.time()
    time_spent = end - begin

    print("> Token subject : %s -> %s" % (sub, err is None and "valid" or "invalid"))
    print("> Time spent    : %fs" % time_spent)

    res = err and err.args[0] or 0
    print("> Result code   : 0x%x" % res)

    try:
        pykcoidc.uninitialize()
    except pykcoidc.Error as e:
        print("> Error: failed to uninitialize 0x%x" % e.args[0])

    return res != 0 and -1 or 0


def validate_and_get_subject(token_s):
    return pykcoidc.validate_token_s(token_s)


if __name__ == "__main__":
    status = main(sys.argv[1:])
    sys.exit(status)
