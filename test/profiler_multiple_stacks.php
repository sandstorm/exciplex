<?php

function work_a() {
    for($i = 0; $i < 50; $i++) {
        usleep(1000);
    }
}

function work_b() {
    for($i = 0; $i < 50; $i++) {
        usleep(1000);
    }
}

$profiler = start_profiler(0.01, 0.02);
work_a();
work_b();
$profiler->stop();
$log = $profiler->getLog();
echo $log;
echo "---END---";
