
ps -ef | grep cqhttp | grep -v grep | awk '{print $2}' | xargs -i -t kill {}
rm -f session.token
run.sh
