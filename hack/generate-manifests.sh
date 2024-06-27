#! /bin/sh
CONFIG_DIR="config"
OUT_DIR="build/manifest"
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"
for VARIANT in "aio" "slim" "autoupdater"; do
  OUT_FILE="$OUT_DIR/manifest-$VARIANT.yaml ..."
  echo "Generating $OUT_FILE"
  kubectl kustomize "$CONFIG_DIR/$VARIANT" > "$OUT_FILE"
done
