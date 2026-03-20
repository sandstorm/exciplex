<?php

/**
 * Represents a scheduled callback — either a one-shot timeout or a repeating interval.
 * Obtained from {@see exciplex_set_timeout()} or {@see exciplex_set_interval()}.
 */
final class ExciplexTimer {
    /**
     * Cancels the timer, preventing any further callback invocations.
     *
     * Safe to call even if the timer has already fired (one-shot) or from within
     * the callback itself (repeating interval).
     */
    public function stop(): void {}
}

/**
 * Schedules $callable to be invoked once after $interval seconds.
 *
 * Returns a timer handle that can be used to cancel the callback before it fires.
 *
 * Example: Abort a long running function with a timeout:
 * ```php
 * $timer = exciplex_set_timeout(function (): void {
 *     throw new Exception("timed out after 500ms");
 * }, 0.5);
 * potentiallyLongRunningFunction();
 * // prevent exception from being thrown after the function ended successfully
 * $timer->stop();
 * ```
 *
 * @param callable(): void $callable Callback to invoke after the delay.
 * @param float $interval Delay in seconds before the callback fires.
 * @return ExciplexTimer A handle to cancel the timeout before it fires.
 */
function exciplex_set_timeout(callable $callable, float $interval): ExciplexTimer {}

/**
 * Schedules $callable to be invoked repeatedly at a fixed interval.
 *
 * The callback fires first after $initialDelay seconds, then every $interval seconds
 * until {@see ExciplexTimer::stop()} is called.
 *
 * Example:
 * ```php
 * $counter = 0;
 * $timer = exciplex_set_interval(function () use (&$counter, &$timer): void {
 *     $counter++;
 *     if ($counter >= 3) {
 *         $timer->stop();
 *     }
 * }, 0.1, 0.1);
 * ```
 *
 * @param callable(): void $callable Callback to invoke on each tick.
 * @param float $initialDelay Delay in seconds before the first invocation.
 * @param float $interval Interval in seconds between subsequent invocations.
 * @return ExciplexTimer A handle to cancel the interval.
 */
function exciplex_set_interval(callable $callable, float $initialDelay, float $interval): ExciplexTimer {}

/**
 * A sampling profiler that periodically captures the current PHP call stack.
 * Obtained from {@see start_profiler()}.
 */
final class ExciplexProfiler {
    /**
     * Stops the profiler. No further samples will be captured after this call.
     *
     * The collected data remains accessible via {@see getLog()} after stopping.
     */
    public function stop(): void {}

    /**
     * Returns the collected samples as a string in the collapsed-stack format
     * (compatible with flamegraph.pl and similar flame graph tools).
     *
     * Each line represents one unique call stack observed during profiling:
     *
     *   <frame>;<frame>;...<frame> <count>
     *
     * Frames run from outermost (entry file) to innermost (deepest call), separated
     * by semicolons. Frame formats:
     * - Entry file / file-level code: /absolute/path/to/file.php
     * - Named function:               functionName
     * - Static or instance method:    ClassName::methodName
     * - Closure:                      {closure:/absolute/path/to/file.php(definitionLine)}
     *
     * Returns an empty string if no samples were collected (e.g. profiler was stopped
     * before the first sample interval elapsed).
     *
     * Example output:
     * ```
     * /app/Web/index.php;processRequest;executeQuery 42
     * /app/Web/index.php;processRequest 8
     * /app/Web/index.php;{closure:/app/Web/index.php(15)} 3
     * ```
     *
     * @return string Collapsed-stack profile data, or empty string if no samples were collected.
     */
    public function getLog(): string {}
}

/**
 * Creates and starts a sampling profiler that captures the PHP call stack at a fixed interval.
 *
 * The profiler begins capturing after $initialDelay seconds, then takes a sample every
 * $interval seconds until {@see ExciplexProfiler::stop()} is called.
 *
 * Example usage for sampling pages in production:
 * ```php
 * $initialDelay = rand(min: 0, max: 60000) / 1000.0;
 * $profiler = start_profiler(initialDelay: $initialDelay, interval: 60);
 * register_shutdown_function(function () use ($profiler) {
 *     $profiler->stop();
 *     $data = $profiler->getLog();
 *
 *     // return early if no traces collected
 *     if ($data === '') {
 *         return;
 *     }
 *
 *     // write log to file or send it to an api endpoint
 *     file_put_contents('/tmp/exciplex-traces.log', $data, FILE_APPEND);
 * });
 * ```
 *
 * This will start capturing a single sample after a random delay of at most 60 seconds (or never, if the page is
 * rendered first) and further frames once every 60 seconds after. This will do nothing for most requests but for the
 * remaining fraction one stack trace will be logged. Due to the random initial delay a flame graph can then be
 * generated from multiple (a lot) requests.
 *
 * @param float $initialDelay Delay in seconds before the first sample is taken.
 * @param float $interval Interval in seconds between samples.
 * @return ExciplexProfiler A handle to stop profiling and retrieve results.
 */
function start_profiler(float $initialDelay, float $interval): ExciplexProfiler {}
