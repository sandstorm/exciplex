<?php

$timer = exciplex_set_timeout(
    function () {
        echo "SHOULD NOT APPEAR";
    },
    0.05
);

// Stop immediately, before the 50ms timeout fires
$timer->stop();

for($i = 0; $i < 10; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}
echo "done";
