//! Hyperfine workload for the cross-language benchmark in issue #8.
//!
//! Iterates the full sequence of one operation, sum-accumulating (wrapping
//! `u64`) a data-dependent value per item, and prints the accumulator once at
//! the end. The data-dependent sum defeats dead-code elimination and gives a
//! cheap, observable per-iteration cost matching the Go and Python drivers.
//!
//! Usage: rust-driver <op> <args...>
//!     combinations | cwr | permutations | product  <n> <k>
//!     powerset                                     <n>
//!     productof                                    <size>...  (one axis per size)
//!     intpartitions                                <n>
//!     setpartitions                                <n> [k]
//!
//! The itertools-backed ops sum every index. intpartitions sums every part;
//! setpartitions sums assignment[i] * i over the restricted-growth-string
//! block assignment. The itertools crate has no partition generators, so the
//! two partition ops are small hand-rolled reference implementations mirroring
//! the Go library's algorithms — they benchmark "what a Rust programmer would
//! write", not an ecosystem library.
//!
//! Rust's itertools yields an owned `Vec` per item (no buffer reuse), so this
//! driver is most comparable to the Go driver's allocating `All()` path.

use itertools::Itertools;
use std::env;
use std::process::exit;

fn main() {
    let args: Vec<String> = env::args().skip(1).collect();
    if args.len() < 2 {
        eprintln!("usage: rust-driver <op> <args...>");
        exit(2);
    }
    let op = args[0].as_str();
    let nums: Vec<usize> = args[1..]
        .iter()
        .map(|a| a.parse().expect("invalid numeric argument"))
        .collect();

    let range_data = |n: usize| -> Vec<u64> { (0..n as u64).collect() };

    // Each branch folds the per-item index sum into a single wrapping accumulator.
    let fold = |acc: u64, item: Vec<u64>| item.iter().fold(acc, |a, &x| a.wrapping_add(x));

    let acc: u64 = match op {
        "combinations" => range_data(nums[0])
            .into_iter()
            .combinations(nums[1])
            .fold(0, fold),
        "cwr" => range_data(nums[0])
            .into_iter()
            .combinations_with_replacement(nums[1])
            .fold(0, fold),
        "permutations" => range_data(nums[0])
            .into_iter()
            .permutations(nums[1])
            .fold(0, fold),
        "product" => itertools::repeat_n(range_data(nums[0]), nums[1])
            .multi_cartesian_product()
            .fold(0, fold),
        "powerset" => range_data(nums[0]).into_iter().powerset().fold(0, fold),
        "productof" => nums
            .iter()
            .map(|&size| range_data(size))
            .multi_cartesian_product()
            .fold(0, fold),
        "intpartitions" => integer_partitions_sum(nums[0]),
        "setpartitions" => {
            let k = nums.get(1).copied().unwrap_or(0);
            let mut acc = 0u64;
            let mut assignment = vec![0usize; nums[0]];
            next_assignment(&mut assignment, 1, 0, k, &mut acc);
            acc
        }
        other => {
            eprintln!("unknown op {other:?}");
            exit(2);
        }
    };

    println!("{acc}");
}

/// Sums every part of every integer partition of n, generated in reverse
/// lexicographic order with the same trailing-ones redistribution algorithm as
/// the Go library's `IntegerPartitions.AllBorrowed`.
fn integer_partitions_sum(n: usize) -> u64 {
    let sum_parts = |acc: u64, parts: &[usize]| {
        parts.iter().fold(acc, |a, &p| a.wrapping_add(p as u64))
    };

    let mut parts: Vec<usize> = vec![n];
    let mut acc = sum_parts(0, &parts);

    loop {
        // Find the rightmost part greater than 1, counting the trailing 1s.
        let mut ones = 0;
        let mut k = None;
        for i in (0..parts.len()).rev() {
            if parts[i] == 1 {
                ones += 1;
            } else {
                k = Some(i);
                break;
            }
        }
        let Some(k) = k else {
            // All parts are 1: this was the last partition.
            return acc;
        };

        // Decrement that part, then redistribute it plus the trailing 1s into
        // chunks no larger than the new value.
        parts[k] -= 1;
        let mut rem = ones + 1;
        parts.truncate(k + 1);
        while rem > parts[k] {
            parts.push(parts[k]);
            rem -= parts[k];
        }
        parts.push(rem);

        acc = sum_parts(acc, &parts);
    }
}

/// Recursively fills assignment[i..] with every valid restricted growth string
/// suffix (m is the maximum entry so far), accumulating assignment[i] * i for
/// each complete string. k > 0 restricts to exactly k blocks, with the same
/// pruning as the Go library's `SetPartitions.next_assignment`.
fn next_assignment(assignment: &mut [usize], i: usize, m: usize, k: usize, acc: &mut u64) {
    let n = assignment.len();
    if i == n {
        for (idx, &b) in assignment.iter().enumerate() {
            *acc = acc.wrapping_add((b as u64).wrapping_mul(idx as u64));
        }
        return;
    }
    for v in 0..=m + 1 {
        if k > 0 && v > k - 1 {
            break;
        }
        let new_max = m.max(v);
        // Each remaining position can open at most one new block, so prune
        // branches that can no longer reach k blocks.
        if k > 0 && k - 1 > new_max && k - 1 - new_max > n - 1 - i {
            continue;
        }
        assignment[i] = v;
        next_assignment(assignment, i + 1, new_max, k, acc);
    }
}
