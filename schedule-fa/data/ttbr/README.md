# SBI TTBR rate data

Bundled SBI **TT Buying Rate** history, one CSV per currency (or a combined file).
Used by `internal/fx` (M2) with preceding-working-day fallback for missing dates.

Expected columns (header required):

```csv
date,currency,inr_per_unit
2024-12-31,USD,85.55
2024-12-30,USD,85.40
```

- `date` — `YYYY-MM-DD`, the date the rate was published for.
- `currency` — ISO-4217, e.g. `USD`.
- `inr_per_unit` — INR per 1 unit of the currency (decimal).

Users can extend or override with `--rates <file.csv>`. There is no official free SBI
TTBR API, so these are curated/curatable data files, not fetched at runtime.
