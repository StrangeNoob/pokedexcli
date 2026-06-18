#!/usr/bin/env bash
# Record a REPL demo with asciinema and render it to assets/demo.gif with agg.
# Requires: asciinema, agg, tmux, go.
#
#   asciinema -> https://github.com/asciinema/asciinema
#   agg       -> https://github.com/asciinema/agg
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SESS=pokedexdemo
CAST="$(mktemp -t pokedex-demo-XXXX).cast"
BIN="$(mktemp -t pokedexcli-XXXX)"
GIF="$ROOT/assets/demo.gif"

for tool in go asciinema agg tmux; do
  command -v "$tool" >/dev/null || { echo "missing required tool: $tool" >&2; exit 1; }
done

go build -o "$BIN" "$ROOT"
tmux kill-session -t "$SESS" 2>/dev/null || true

# Record the REPL inside a tmux pty so asciinema has a terminal to attach to.
tmux new-session -d -s "$SESS" -x 100 -y 30 "asciinema record '$CAST' --overwrite -c '$BIN'"
sleep 3 # let asciinema and the REPL start

send() { # send <text> <seconds-to-wait>
  tmux send-keys -t "$SESS" "$1"
  sleep 0.5
  tmux send-keys -t "$SESS" Enter
  sleep "$2"
}

send "help" 2.5
send "map" 3
send "explore pastoria-city-area" 3
send "catch tentacool ultraball" 4
send "inspect tentacool" 2.5
send "bag" 2.5
send "pokedex" 3
send "exit" 2

# Wait for the recording session to finish.
for _ in $(seq 1 25); do
  tmux has-session -t "$SESS" 2>/dev/null || break
  sleep 1
done
tmux kill-session -t "$SESS" 2>/dev/null || true

mkdir -p "$ROOT/assets"
agg --speed 1.4 --font-size 20 "$CAST" "$GIF"
rm -f "$CAST" "$BIN"
echo "Wrote $GIF"
