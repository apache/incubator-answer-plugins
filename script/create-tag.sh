# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

#!/bin/bash

project_root="."

for dir in $(find "$project_root" -type d); do
    if [[ "$dir" == *node_modules* ]] || [[ "$dir" == *.git* ]] || [[ "$dir" == *.vscode* ]]; then
        continue
    fi

    # 当 info.yaml 不存在时，跳过
    if [ ! -f "$dir/info.yaml" ]; then
        continue
    fi

    if git diff HEAD^^ HEAD $dir/info.yaml | grep -q "^\+version:"; then
        project_name=$(echo $dir | awk -F'/' '{print $NF}')
        version=$(grep "version:" $dir/info.yaml | awk '{print $2}')
        commit_msg=$(git log -1 --pretty=%B)
        echo "create tag for $project_name/v$version"
        git tag -a $project_name/v$version -m "$commit_msg"
    fi    
done