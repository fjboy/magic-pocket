
function logInfo() {
    echo `date "+%F %T" ` "INFO:" $@ 1>&2
}

function main(){
    logInfo "Get version"
    version=$(go run cmd/magic-pocket.go -v |awk '{print $3}')
    if [[ -z $version ]] || [[ "${version}" == "" ]]; then
        exit 1
    fi
    logInfo "Start to build with option main.Version=${version}"
    go build  -ldflags "-X main.Version=${version}" cmd/magic-pocket.go
    logInfo "Build success"
}
main
