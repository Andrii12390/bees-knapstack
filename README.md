# bees_knapsack

A Go implementation of the Bees Algorithm applied to the 0/1 knapsack problem,
with several parallelization strategies for an empirical study of speedup,
scaling, and load balancing.

## Project structure

```
bees_knapsack/
├── go.mod
├── README.md
├── .gitignore
│
├── cmd/
│   ├── demo/                  # interactive demo (also: verify, smoke subcommands)
│   ├── benchmark-seq/         # sequential scaling sweep
│   ├── benchmark-cmp/         # seq vs. all parallel strategies
│   ├── benchmark-series1/     # Series 1: scaling by problem size
│   ├── benchmark-series2/     # Series 2: scaling by worker count
│   └── benchmark-series3/     # Series 3: strategy comparison
│
└── internal/
    ├── problem/               # Problem, Item, Solution + helpers
    ├── algorithm/             # Sequential BeesAlgorithm + parallel strategies
    ├── verify/                # Brute-force correctness checks
    └── benchmark/             # Shared benchmarking helpers
```

## Running

```
go run ./cmd/demo                  # demo run on a small fixed instance
go run ./cmd/demo verify           # brute-force correctness checks
go run ./cmd/demo smoke            # spread check across all strategies
go run ./cmd/benchmark-seq         # sequential scaling benchmark
go run ./cmd/benchmark-cmp         # seq vs par comparison
go run ./cmd/benchmark-series1     # Series 1: scaling by problem size
go run ./cmd/benchmark-series2     # Series 2: scaling by worker count
go run ./cmd/benchmark-series3     # Series 3: strategy comparison
```
