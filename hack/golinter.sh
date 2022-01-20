#! /bin/bash
linter_path=$(which golangci-lint)
if [ "${linter_path}" = "" ]; then
  echo "linter is not found. linter must be installed."
  exit 1
fi

linter_target_dir=$1
if [ $# != 1 ]; then
  echo "only one args, target go module dir, is needed."
  exit 1
fi

pushd "${linter_target_dir}" && golangci-lint run --timeout 5m && popd
