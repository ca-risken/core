#! /bin/bash
mockery_path=$(which mockery)
if [ "${mockery_path}" = "" ]; then
  echo "mockery is not found. mockery must be installed."
  exit 1
fi

generation_target_dir=$1
if [ $# != 1 ]; then
  echo "only one args, target go module dir, is needed."
  exit 1
fi

pushd "${generation_target_dir}" && mockery --all && popd
