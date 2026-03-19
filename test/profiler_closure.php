<?php

$closure = function() {
    for ($i = 0; $i < 15; $i++) {
        usleep(10000);
    }
};

$profiler = start_profiler(0.01, 0.05);
$closure();
$profiler->stop();

echo $profiler->getLog();
