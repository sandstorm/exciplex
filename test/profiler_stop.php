<?php

// waits 50ms in steps of 10ms
function wait() {
    for($i = 0; $i < 5; $i++) {
        usleep(10000);
    }
}

$profiler = start_profiler(0.01, 0.02);
wait();
$profiler->stop();

// wait and verify no new samples are added
$logBefore = $profiler->getLog();
wait();
$logAfter = $profiler->getLog();

echo $logBefore;
if ($logBefore === $logAfter) {
    echo "stop_works";
} else {
    echo "FAIL: log changed after stop";
}
