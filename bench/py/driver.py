#!/usr/bin/env python3
"""Hyperfine workload for the cross-language benchmark in issue #8.

Iterates the full sequence of one ``itertools`` operation, sum-accumulating
every index it sees (masked to 64 bits to match the Go and Rust drivers' wrapping
``u64``), and prints the accumulator once at the end. The data-dependent sum
defeats dead-code elimination and gives a cheap, observable per-iteration cost.

Usage: driver.py <op> <n> [k]
    op  in combinations | cwr | permutations | product | powerset
    n   size of the input (elements are the ints 0..n-1)
    k   required for every op except powerset
"""
import sys
from itertools import (
    chain,
    combinations,
    combinations_with_replacement,
    permutations,
    product,
)


def make_iter(op, n, k):
    data = range(n)
    if op == "combinations":
        return combinations(data, k)
    if op == "cwr":
        return combinations_with_replacement(data, k)
    if op == "permutations":
        return permutations(data, k)
    if op == "product":
        return product(data, repeat=k)
    if op == "powerset":
        # itertools powerset recipe
        return chain.from_iterable(combinations(data, r) for r in range(n + 1))
    raise SystemExit(f"unknown op {op!r}")


def main():
    args = sys.argv[1:]
    if len(args) < 2:
        raise SystemExit("usage: driver.py <op> <n> [k]")
    op = args[0]
    n = int(args[1])
    k = int(args[2]) if op != "powerset" else 0

    acc = 0
    for tup in make_iter(op, n, k):
        for idx in tup:
            acc += idx
    print(acc & 0xFFFFFFFFFFFFFFFF)


if __name__ == "__main__":
    main()
