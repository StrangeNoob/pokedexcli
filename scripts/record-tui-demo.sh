#!/usr/bin/env bash
# Record a TUI demo with asciinema and render it to assets/tui-demo.gif with agg.
# Drives the full-screen Bubble Tea UI with tmux send-keys.
# Requires: asciinema, agg, tmux, go.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SESS=pokedextui
CAST="$(mktemp -t pokedex-tui-XXXX).cast"
BIN="$(mktemp -t pokedexcli-XXXX)"
GIF="$ROOT/assets/tui-demo.gif"

for tool in go asciinema agg tmux; do
  command -v "$tool" >/dev/null || { echo "missing required tool: $tool" >&2; exit 1; }
done

go build -o "$BIN" "$ROOT"
tmux kill-session -t "$SESS" 2>/dev/null || true

# A wide pty so the battle screen (sprite + stats per side) fits.
tmux new-session -d -s "$SESS" -x 120 -y 38 "asciinema record '$CAST' --overwrite -c '$BIN tui'"
sleep 3 # let asciinema and the TUI start

tk() { tmux send-keys -t "$SESS" "$1"; sleep "$2"; } # tk <key> <seconds>

# Pokédex: open and scroll (sprites + stats load per selection).
tk Enter 3
tk Down 2.5
tk Down 2.5
tk Escape 1.2

# Bag.
tk Down 0.4
tk Down 0.4
tk Down 0.4
tk Enter 2.5
tk Escape 1.2

# Battle: pick two Pokémon, watch the animated fight.
tk Up 0.5
tk Enter 2.5    # -> choose your Pokémon
tk Down 1.8
tk Enter 2      # first picked -> choose opponent
tk Down 1.8
tk Enter 7      # start battle, animate HP bars
tk Space 2.5    # fast-forward to the result
tk Escape 1.2

# Quit from the menu.
tk q 2

for _ in $(seq 1 25); do
  tmux has-session -t "$SESS" 2>/dev/null || break
  sleep 1
done
tmux kill-session -t "$SESS" 2>/dev/null || true

mkdir -p "$ROOT/assets"
agg --speed 1.5 --font-size 14 --idle-time-limit 2 "$CAST" "$GIF"
rm -f "$CAST" "$BIN"
echo "Wrote $GIF ($(wc -c < "$GIF") bytes)"
