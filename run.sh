#!/bin/bash

set -euo pipefail

export PATH="$(pwd)/maelstrom:${PATH}"

if [[ "$#" -eq 0 ]]; then
  maelstrom serve
  exit 0
fi

readonly command="${1}"

# install first
(
  cd "maelstrom-${command}"
  go install .
)

case "${command}" in
echo)
  maelstrom test -w echo --bin ~/go/bin/maelstrom-echo --node-count 1 --time-limit 10
  ;;
unique-ids)
  maelstrom test -w unique-ids \
    --bin ~/go/bin/maelstrom-unique-ids \
    --time-limit 30 \
    --rate 1000 \
    --node-count 3 \
    --availability total \
    --nemesis partition
  ;;
broadcast)
  maelstrom test -w broadcast \
    --bin ~/go/bin/maelstrom-broadcast \
    --node-count 5 \
    --time-limit 20 \
    --rate 10
  ;;
counter)
  maelstrom test -w g-counter \
    --bin ~/go/bin/maelstrom-counter \
    --node-count 3 \
    --rate 100 \
    --time-limit 20 \
    --nemesis partition
  ;;
*)
  echo "Unknown command: ${command}"
  exit 1
  ;;
esac
