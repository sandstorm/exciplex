<?php

exciplex_set_timeout(
    function () {
        echo "interrupt was called";
    },
    0.05
);

for($i = 0; $i < 10; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}