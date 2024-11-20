#!/bin/bash

set -euo pipefail

(
  cd /tmp
  curl -LO https://github.com/jepsen-io/maelstrom/releases/download/v0.2.3/maelstrom.tar.bz2
  tar -xf maelstrom.tar.bz
)

cp /tmp/maelstrom/maelstrom .
