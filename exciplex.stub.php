<?php

/** @generate-class-entries */

final class ExciplexTimer {
    public function stop(): void {}
}

function exciplex_set_timeout(callable $callable, float $interval): ExciplexTimer {}

function exciplex_set_interval(callable $callable, float $initialDelay, float $interval): ExciplexTimer {}
