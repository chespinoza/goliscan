#!/bin/bash

PROJECT_NAME=${PROJECT_NAME:-"Project Name"}

S3_ALIAS=${S3_ALIAS:-"release"}
S3_BUCKET=${S3_BUCKET:-}
S3_HOST=${S3_HOST:-}
S3_ACCESS_KEY=${S3_ACCESS_KEY:-}
S3_SECRET_KEY=${S3_SECRET_KEY:-}
S3_API=${S3_API:-"S3v4"}

GIT_TREE_URL=${GIT_TREE_URL:-"unknown"}
SOURCE_URL="${GIT_TREE_URL}/${CI_BUILD_REF_NAME}"

CREATED_AT=$(date +%Y-%m-%dT%H:%M:%S%:z)

__check_variable() {
    if [[ "${!1}" == "" ]]; then
        echo "You need to set a value for '${1}' variable!"
        return 1
    fi
}

__check_variables() {
    EXIT=0
    __check_variable PROJECT_NAME || EXIT=1
    __check_variable S3_ALIAS || EXIT=1
    __check_variable S3_BUCKET || EXIT=1
    __check_variable S3_HOST || EXIT=1
    __check_variable S3_ACCESS_KEY || EXIT=1
    __check_variable S3_SECRET_KEY || EXIT=1
    __check_variable S3_API || EXIT=1


    if [[ "${EXIT}" -gt 0 ]]; then
        echo $EXIT
        exit 1
    fi
}

__usage_info() {
    echo "Usage ${0} (development|unstable|stable)"
    exit 1
}

__configure_s3_client() {
    go get -u github.com/minio/mc
    mc config host add  ${S3_ALIAS} ${S3_HOST} "${S3_ACCESS_KEY}" "${S3_SECRET_KEY}" ${S3_API}
}

__prepare_index() {
    index_file=index.html
    title="${PROJECT_NAME} :: Release for ${RELEASE} (${CI_BUILD_REF_NAME})"
    sources_url=${GIT_TREE_URL}/${CI_BUILD_REF_NAME}

    path=$(pwd)
    cd out

    echo -e "\033[1mPreparing index file\033[0m"

    cat > ${index_file} <<EOF
<html>
    <head>
        <meta charset="utf-8/">
        <title>${title}</title>
        <style type="text/css">
            body {font-family: arial, serif; font-size: 14px; margin: 40px; padding: 0;}
            h1 {border-bottom: 1px solid #aaa; padding: 10px;}
            p {font-size: 0.85em; margin: 5px 10px;}
            p span {display: inline-block; font-weight: bold; width: 75px;}
            p a {color: #000; font-weight: bold; text-decoration: none;}
            p a:hover {text-decoration: underline;}
            ul {background: #eee; border: 1px solid #aaa; border-radius: 3px; box-shadow: 0 0 5px #aaa inset; list-style-type: none; margin: 10px; padding: 20px;}
            li {margin: 5px 0; padding: 2px;}
            li:hover {background: #dedede;}
            li a {color: #000; text-decoration: none;}
            .file_name {display: inline-block; font-weight: bold; width: calc(100% - 100px)}
            .file_size {display: inline-block; width: 90px;}
        </style>
    </head>
    <body>
        <h1>${title}</h1>
        <ul>
EOF

    files=$(find * -type f ! -name "index.html" | sort -g)
    for file in ${files}; do
        size=$(du -k ${file} | awk '{printf "%0.2f", $1/1024}')
        echo -e "found \033[32m${file}\033[36m [${size} MiB]\033[0m"
        echo "            <li><a href=\"./${file}\"><span class=\"file_name\">${file}</span> <span class=\"file_size\">${size} MiB</span></a></li>" >> ${index_file}
    done

    cat >> ${index_file} <<EOF
        </ul>
        <p><span>Sources:</span> <a href="${sources_url}">${sources_url}</a></p>
        <p><span>Created at:</span> ${CREATED_AT}</p>
    </body>
</html>
EOF
    cd ${path}

    echo -e "\033[1mIndex file created!\033[0m"
}

__upload_release() {
    if [[ "$#" -lt 1 ]]; then
        echo "Usage __upload_release [release_path]"
        exit 1
    fi

    UPLOAD_PATH="${S3_BUCKET}/${1}"

    mc cp -q -r out/ ${S3_ALIAS}/${UPLOAD_PATH}
    echo -e "\033[37mDownload URL: \033[36m${S3_HOST}/${UPLOAD_PATH}/index.html\033[0m"
}



if [[ "$#" -lt 1 ]]; then
    __usage_info
fi

RELEASE="${1}"
case "${RELEASE}" in
    development)
        ;;
    unstable)
        ;;
    stable)
        ;;
    *)
        __usage_info
esac

__check_variables
__prepare_index
__configure_s3_client
__upload_release ${CI_BUILD_REF_NAME}
__upload_release release_${RELEASE}

