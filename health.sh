#!/bin/sh

set -e

[ $(/usr/bin/find /tmp/imalive -mmin -1 -type f -print | wc -l) -gt "0" ]

