<?php

exciplex_set_timeout(
    function () {
        die("last visible text, no further counting");
    },
    0.05
);

for($i = 0; $i < 10; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}