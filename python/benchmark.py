#!/bin/env python

from __future__ import print_function

from multiprocessing import cpu_count
import multiprocessing.dummy as multiprocessing
import sys
import time

import pykcoidc


def bench_validateTokenS(id, count, token_s):
    print("> Info : thread %s started ..." % id)

    success = 0
    failed = 0
    for c in range(count):
        try:
            sub = validate_and_get_subject(token_s)
            success += 1
        except pykcoidc.Error as e:
            print("> Error: validation failed: 0x%x" % e.args[0])
            failed += 1
        except Exception as e:
            print("> Error: unknown exception: %s" % e)
            failed += 1

    print("> Info : thread %s done:%d failed:%d" % (id, success, failed))


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

    concurentThreadsSupported = cpu_count()
    count = 100000
    pool = multiprocessing.Pool(concurentThreadsSupported)

    # Wait until oidc validation becomes ready.
    try:
        pykcoidc.wait_until_ready(10)
    except pykcoidc.Error as e:
        print("> Error: failed to get ready in time: 0x%x" % e.args[0])
        return -1

    print("> Info : using %d threads with %d runs per thread" % (concurentThreadsSupported, count))
    begin_time = time.time()
    for i in range(concurentThreadsSupported):
        pool.apply_async(bench_validateTokenS, [i+1, count, token_s])

    pool.close()
    pool.join()
    end_time = time.time()
    duration = end_time - begin_time
    rate = (count * concurentThreadsSupported) / duration
    print("> Time : %fs" % duration)
    print("> Rate : %f op/s" % rate)

    try:
        pykcoidc.uninitialize()
    except pykcoidc.Error as e:
        print("> Error: failed to uninitialize 0x%x" % e.args[0])

    return 0


def validate_and_get_subject(token_s):
    return pykcoidc.validate_token_s(token_s)


if __name__ == "__main__":
    status = main(sys.argv[1:])
    sys.exit(status)
