#!/bin/bash
# Head-to-head scorer: replicates guess-it-dockerized/server.js scoring exactly.
# per hit: round(10000000 / (1 + width) / (N-1)); miss = 0.
SRC="/mnt/c/Users/Dell/Documents/VS Code/Zone01 Study/1st Semester/guess-it-2"
WORK="$HOME/gtest"
rm -rf "$WORK"; mkdir -p "$WORK/ai" "$WORK/student"
cp "$SRC/.resources/guess-it-dockerized/ai/"* "$WORK/ai/" 2>/dev/null
cp "$SRC/student/guess-it-2" "$WORK/student/"
cp -r "$SRC/.resources/guess-it-dockerized/data_sets" "$WORK/"
chmod +x "$WORK/ai/"* "$WORK/student/guess-it-2"
cd "$WORK"

score() { # $1 program, $2 data file -> prints integer score
  "$1" < "$2" > /tmp/o.txt 2>/dev/null
  awk '
    FNR==NR { if (NF>0) { n++; d[n]=$1+0 } next }
    { if (NF>=2) { m++; lo[m]=$1+0; hi[m]=$2+0 } }
    END {
      for (j=1; j<=n-1; j++) {
        v=d[j+1]
        if (j in lo && v>=lo[j] && v<=hi[j]) {
          w=hi[j]-lo[j]; s+=int(1e7/(1+w)/(n-1)+0.5)
        }
      }
      printf "%d", s
    }
  ' "$2" /tmp/o.txt
}

OPP="big-range linear-regr correlation-coef mse nic average median huge-range"
for ds in 4 5; do
  echo ""
  echo "########## DATA $ds  (files 1-5) ##########"
  unset ST; declare -A ST
  srow=$(printf '%-18s' "student")
  ssum=0
  for f in 1 2 3 4 5; do
    ST[$f]=$(score student/guess-it-2 "data_sets/$ds/$f.txt")
    srow="$srow $(printf '%9d ' ${ST[$f]})"
    ssum=$((ssum+ST[$f]))
  done
  echo "$srow  mean=$((ssum/5))"
  echo "-------------------------------------------------------------------------"
  for o in $OPP; do
    row=$(printf '%-18s' "$o")
    wins=0; osum=0
    for f in 1 2 3 4 5; do
      os=$(score "ai/$o" "data_sets/$ds/$f.txt")
      osum=$((osum+os))
      if [ "${ST[$f]}" -gt "$os" ]; then m="W"; wins=$((wins+1)); else m="L"; fi
      row="$row $(printf '%9d%s' $os $m)"
    done
    echo "$row  mean=$((osum/5)) studentWins=$wins/5"
  done
done
