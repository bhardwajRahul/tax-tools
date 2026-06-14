# schedule-fa

A Go CLI that turns **Interactive Brokers (IBKR)** holdings into a ready-to-use
**Schedule FA** (Foreign Assets) report for the Indian ITR — handling the calendar-year
basis, SBI TT buying-rate conversion to INR, and peak/closing/initial values per
security, with a full audit trail.

See [`docs/schedule-fa-ibkr-plan.md`](docs/schedule-fa-ibkr-plan.md) for the research,
challenges, locked decisions, architecture, and milestones.

## Status

**M0 — skeleton.** CLI scaffold + domain model + package stubs compile; report generation
is not implemented yet (the `generate` command prints what it *will* do and exits).

## Usage (target)

```sh
schedulefa generate \
  --year 2024 \                 # CALENDAR year (Jan 1 – Dec 31), enforced
  --statement statement.xml \   # IBKR Activity Flex Query, XML output (offline mode)
  --rates data/ttbr/usd.csv \   # optional SBI TTBR override
  --out ./report --format md,csv,json
```

## Build

```sh
go build ./cmd/schedulefa       # from the schedule-fa/ directory
```

> **Disclaimer:** Not tax advice. Output is a working draft to be verified before filing.
</content>
