mkdir -p logs
nohup ./go-cqhttp > logs/runtime.log 2>&1 &
tail -f logs/runtime.log
