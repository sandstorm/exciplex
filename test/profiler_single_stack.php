<?php

function test() {
    // Wait in a loop so the profiler could sample us multiple times (but should only do so once)
    for ($i = 0; $i < 15; $i++) {
        usleep(1000);
    }
}

$profiler = start_profiler(0.005, 0.013);
test();
$profiler->stop();

echo $profiler->getLog();
