# Monte-Carlo — SYNTH (sma-cross(3/8))

- Period: 2023-01-02 → 2023-06-16 (120 bars)
- Bootstrap trials: 500 (seed 12345)

| Metric       | Actual |     P5 |    P25 | Median |    P75 |     P95 |
| ------------ | -----: | -----: | -----: | -----: | -----: | ------: |
| Total return | 29.01% | 21.42% | 25.54% | 29.22% | 32.57% |  38.02% |
| CAGR         | 75.74% | 53.67% | 65.46% | 76.37% | 86.67% | 104.07% |
| Max drawdown |  1.05% |  0.45% |  0.49% |  0.64% |  0.82% |   1.16% |
| Sharpe       |   9.21 |   7.24 |   8.44 |   9.25 |  10.02 |   11.30 |

> 100% of 500 bootstrap trials were profitable, and the middle 90% of total returns span 21.4% to 38.0%. A distribution that hugs or crosses zero means the headline result leaned on luck; a tight, clearly-positive spread is what a real edge looks like.
> Bootstrap shuffles days independently, so it measures sampling luck, not regime or autocorrelation effects — read it next to walk-forward, not instead of it.

_NOTE: not investment advice. A backtest is a hypothesis, not a forecast — it is fit to the past, ignores regime change, and flatters strategies that overfit. Costs and slippage are estimates; live results will be worse._
