#!/bin/bash
set -e

su - clickhouse -c "tar --skip-old-files -xvzf  /bundle.tar.gz -C /"

for f in /var/lib/clickhouse/user_defined/*.sql; do clickhouse client --queries-file $f; done