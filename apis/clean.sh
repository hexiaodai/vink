set -eu

PATTERNS=".validate.go _deepcopy.gen.go .gen.json gr.gen.go .pb.go _json.gen.go .pb.gw.go .swagger.json .deepcopy.go"

for p in $PATTERNS; do
    rm -f ./**/**/v3alpha1/*"${p}"
    rm -f ./**/**/*"${p}"
    rm -f ./**/*"${p}"
done

find ./sdks/ts | grep -v  package.json | awk "NR != 1" | xargs rm -rf
