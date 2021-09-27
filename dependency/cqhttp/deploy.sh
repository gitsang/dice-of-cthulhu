
if [[ ! -f "go-cqhttp" ]]; then
    if [[ ! -f "go-cqhttp_linux_amd64.tar.gz" ]]; then
        wget https://github.com/Mrs4s/go-cqhttp/releases/download/v1.0.0-beta7-fix2/go-cqhttp_linux_amd64.tar.gz
    fi
    tar zxvf go-cqhttp_linux_amd64.tar.gz
fi
