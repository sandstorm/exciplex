<?php

$num = 0;
exciplex_set_interval(
    function () use (&$num) {
        echo "interval $num -- ";
        $num += 1;
    },
    0.03,
    0.03
);

for($i = 0; $i < 10; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}