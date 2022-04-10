while true; do
    echo "[run.sh] Starting debugging..."
    dlv debug --headless --log --listen=:2345 --api-version=2 --accept-multiclient --continue &

    PID=$!

    inotifywait -e modify -e move -e create -e delete -e attrib --exclude '(__debug_bin|\.git|media)' -r .

    echo "[run.sh] Stopping process id: $PID"
    
    kill -9 $PID
    pkill -f __debug_bin
done
