<?php

$timer = exciplex_set_interval(
    function () {
        echo "SHOULD NOT APPEAR";
    },
    0.05,
    0.01
);

// Stop immediately, before the 50ms initial delay fires
$timer->stop();

for($i = 0; $i < 10; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}
echo "done";
