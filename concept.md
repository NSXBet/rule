# rule

This rule engine must pass the test/fixtures.go set of tests. It should evaluate rules in the form of "x eq 10 and y gt
20".

There are some constraints:

1. The evaluator will be loaded ahead of time so it can allocate as much memory as it can to optimize for runtime speed.
2. Evaluation of context based rules MUST not allocate memory.
3. Evaluation of context based rules MUST be completed in less than 1000 nanoseconds.
4. The rules that must work are present in the test/fixtures.go file in the form of test tables.

## Coding guidelines

You are a seasoned golang developer. You know what you are doing. You are also a champion of TDD. You will always follow
the RED->GREEN->REFACTOR pattern.

Our code must follow modern golang constructs like any, for range 10, and so on. Please refrain from using fmt.Sprintf
as much as possible. We will only have typed errors in our code. No errorf or error.New. Any code that allocated memory
must be weighted against the performance of the code. We want to optimize for speed. Any code that locks for concurrency
should be replaced with lock-free structures (xsync v4 or uber's atomic).
