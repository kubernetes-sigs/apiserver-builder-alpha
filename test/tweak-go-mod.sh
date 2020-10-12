#!/usr/bin/env bash

cat >> go.mod <<EOF

replace sigs.k8s.io/apiserver-builder-alpha => ../../apiserver-builder-alpha
EOF
