<?php

$num = 0;
$timer = exciplex_set_interval(
    function () use (&$num, &$timer) {
        echo "interval $num -- ";
        $num += 1;

        // Stop after the first callback
        $timer->stop();
    },
    0.02,
    0.01
);

for($i = 0; $i < 10; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}
echo "done";
