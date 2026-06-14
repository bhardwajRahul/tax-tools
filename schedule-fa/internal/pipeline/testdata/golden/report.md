# Schedule FA — calendar year 2024

> Not tax advice. A working draft to verify before filing. Schedule FA covers the CALENDAR year; see the audit trail for every figure.

## Table A2 — Foreign Custodial Account

- **Institution:** Interactive Brokers LLC
- **Address:** One Pickwick Plaza, Greenwich, CT 06830 · country code 1
- **Account number:** U1234567 · **Status:** Owner · **Opened:** 2021-05-10
- **Peak balance:** ₹3352670.00 · **Closing balance:** ₹2600720.00 · **Gross credited:** ₹3330.75
- ⚠︎ _peak balance is an upper bound (per-security peaks summed, not simultaneous); cash not included_

## Table A3 — Foreign Equity and Debt Interest

| # | Entity | Country (code) | Acquired | Initial (INR) | Peak (INR) | Closing (INR) | Dividend (INR) | Proceeds (INR) | Review |
|---|--------|----------------|----------|--------------:|-----------:|--------------:|---------------:|---------------:|:------:|
| 1 | APPLE INC | United States of America (1) | 2023-01-10 | 498600.00 | 2138750.00 | 2138750.00 | 2077.50 | 0.00 | ⚠︎ |
| 2 | MICROSOFT CORP | United States of America (1) | 2024-02-01 | 664800.00 | 751950.00 | 0.00 | 0.00 | 751950.00 | ⚠︎ |
| 3 | VANGUARD S&P 500 ETF | United States of America (1) | 2024-06-20 | 417750.00 | 461970.00 | 461970.00 | 1253.25 | 0.00 | ⚠︎ |

## Reconciliation

- Securities (A3 rows): **3**
- Rows needing manual review: **3**
- Total closing value: **₹2600720.00**
- Total gross dividend: **₹3330.75**
- Total sale proceeds: **₹751950.00**

## Audit trail

Each INR figure and the SBI TTBR (rate date actually used) behind it.

### 1. APPLE INC

_Review: entity address/ZIP missing (set via --entities); 1 corporate action(s) in year — verify quantity/cost basis; fx: no USD TTBR on or before 2023-01-10 (earliest available is 2024-01-05)_

| Figure | Source | TTBR | Rate date | INR |
|--------|--------|-----:|-----------|----:|
| Initial value | USD 6000.00 | 83.1000 | 2024-01-05 | 498600.00 |
| Peak value | USD 25000.00 | 85.5500 | 2024-12-31 | 2138750.00 |
| Closing value | USD 25000.00 | 85.5500 | 2024-12-31 | 2138750.00 |
| Dividend | USD 25.00 | 83.1000 | 2024-01-05 | 2077.50 |
| Sale proceeds | — | — | — | 0.00 |

### 2. MICROSOFT CORP

_Review: entity address/ZIP missing (set via --entities)_

| Figure | Source | TTBR | Rate date | INR |
|--------|--------|-----:|-----------|----:|
| Initial value | USD 8000.00 | 83.1000 | 2024-01-05 | 664800.00 |
| Peak value | USD 9000.00 | 83.5500 | 2024-06-14 | 751950.00 |
| Closing value | — | — | — | 0.00 |
| Dividend | — | — | — | 0.00 |
| Sale proceeds | USD 9000.00 | 83.5500 | 2024-06-14 | 751950.00 |

### 3. VANGUARD S&P 500 ETF

_Review: entity address/ZIP missing (set via --entities)_

| Figure | Source | TTBR | Rate date | INR |
|--------|--------|-----:|-----------|----:|
| Initial value | USD 3500.00 | 83.5500 | 2024-06-14 | 292425.00 |
|  | USD 1500.00 | 83.5500 | 2024-06-14 | 125325.00 |
| Peak value | USD 5400.00 | 85.5500 | 2024-12-31 | 461970.00 |
| Closing value | USD 5400.00 | 85.5500 | 2024-12-31 | 461970.00 |
| Dividend | USD 15.00 | 83.5500 | 2024-06-14 | 1253.25 |
| Sale proceeds | — | — | — | 0.00 |

