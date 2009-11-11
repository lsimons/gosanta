#!/bin/bash
# Copyright 2009 The Go Authors. All rights reserved.
# Copyright 2009 Leo Simons. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if [[ "x$GOROOT" == "x" ]]; then
        echo '$GOROOT is not set' 1>&2
        exit 1
fi

case "$GOARCH" in
amd64 | 386 | arm)
        ;;
*)
        echo '$GOARCH is set to <'$GOARCH'>, must be amd64, 386, or arm' 1>&2
        exit 1
esac

case "$GOOS" in
darwin | linux | nacl)
        ;;
*)
        echo '$GOOS is set to <'$GOOS'>, must be darwin, linux, or nacl' 1>&2
        exit 1
esac

rm -f santa.6 santa || exit 1
6g -o santa.6 santa.go || exit 1
6l -o santa santa.6 || exit 1
exit $?

