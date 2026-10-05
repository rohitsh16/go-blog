---
title: "Saving 80% Compute: Sequential Verification and Anytime-Valid Early Stopping in AI Agents"
slug: "sequential-verification-early-stopping-llms"
author: "Rohit Shukla"
date: "2026-10-05"
published: true
summary: "Verifying every generated claim or tool execution in an autonomous agent is computationally crippling. We evaluate Wald's SPRT, Empirical Bernstein confidence sequences, and Ville's martingale e-values, demonstrating how to save 74%-83% of verification compute while mathematically bounding error rates."
---

When autonomous agents generate multi-step research reports, complex software diffs, or analytical summaries, how do you verify their correctness without exploding API bills and latency?

If an agent produces 50 claims, invoking a heavyweight verifier (such as a full LLM-as-a-judge call, a sandboxed Python test execution, or a PDF section extractor) on every assertion requires 50 sequential checks. For reliable generations, this is massively wasteful. For deeply flawed generations, waiting for all 50 checks before aborting delays critical failure detection.

In this article, we examine empirical results from **`EXP-SEQ-001`**, investigating **Level 4 of the Correctness Ladder: Sequential & Adaptive Verification**. We show how **Wald's Sequential Probability Ratio Test (SPRT)**, **Empirical Bernstein Confidence Sequences**, and **Ville's Martingale E-Values** cut verification compute by **74% to 83%** with provable Type I and Type II error bounds.

---

## 1. The Fixed-Horizon Verification Problem

Traditional software verification tests a fixed sample of size $T_{\text{max}}$. But in streaming agent verification, fixed-sample testing introduces two painful failure modes:

1. **Over-verification on High-Quality Responses**: If the first 12 assertions are cleanly verified with high confidence, running the remaining 38 verifier checks adds pure latency with diminishing marginal information.
2. **Delayed Abort on Hallucination Spirals**: If an agent starts producing errors early in a reasoning chain, continuing downstream verification wastes compute on an already corrupted premise.

```
       Fixed-Horizon Verification (Wasteful)
       [Claim 1] -> [Claim 2] -> ... -> [Claim 50] (Always runs all 50 checks)

       Sequential SPRT Verification (Adaptive)
       [Claim 1] -> [Claim 2] -> ... -> [Claim 14] ──► STOP & ACCEPT (86% Savings!)
```

---

## 2. Wald's Sequential Probability Ratio Test (SPRT)

Pioneered by Abraham Wald (1945), SPRT formulates verification as a streaming binary hypothesis test:
$$H_0: p \le p_0 = 0.05 \quad (\text{System is reliable}) \quad \text{vs.} \quad H_1: p \ge p_1 = 0.20 \quad (\text{System is flawed})$$

With user-specified safety parameters:
- **$\alpha$ (Type I error)**: Probability of rejecting a reliable system ($0.05$).
- **$\beta$ (Type II error)**: Probability of accepting a flawed system ($0.10$).

At each verification step $t$, we observe whether claim $t$ is verified ($X_t = 0$) or flawed ($X_t = 1$). We update the cumulative log-likelihood ratio:
$$\Lambda_t = \Lambda_{t-1} + \begin{cases} \log(p_1 / p_0) & \text{if error } (X_t = 1) \\ \log\left(\frac{1 - p_1}{1 - p_0}\right) & \text{if verified } (X_t = 0) \end{cases}$$

### The Adaptive Decision Boundaries
Rather than testing fixed $T$, SPRT evaluates Wald's boundaries at every step:
$$\begin{cases}
\Lambda_t \ge \log\left(\frac{1 - \beta}{\alpha}\right) \approx +2.89 & \implies \text{Stop immediately and REJECT (Abort)} \\
\Lambda_t \le \log\left(\frac{\beta}{1 - \alpha}\right) \approx -2.25 & \implies \text{Stop immediately and ACCEPT (Certified)} \\
\text{otherwise} & \implies \text{Continue to next claim } (t \leftarrow t + 1)
\end{cases}$$

---

## 3. Empirical Results from `EXP-SEQ-001`

We benchmarked Wald's SPRT across three pre-registered seeds (`42, 137, 2026`) against a fixed horizon of $T_{\text{max}} = 100$ verification steps.

### Table 1: Early Stopping & Efficiency Performance

| Evaluation Scenario | Metric | Seed 42 | Seed 137 | Seed 2026 | Aggregate Mean | Compute Savings |
|---|---|---|---|---|---|---|
| **Reliable System ($H_0$: Error $\le 3\%$)** | Steps to Decision | 50 | 14 | 14 | **26.0 steps** | **-74.0% Compute** |
| | Decision | ACCEPT_H0 | ACCEPT_H0 | ACCEPT_H0 | **100% Accurate** | 0% False Alarms |
| **Flawed System ($H_1$: Error $\ge 25\%$)** | Steps to Decision | 25 | 8 | 17 | **16.7 steps** | **-83.3% Compute** |
| | Decision | ACCEPT_H1 | ACCEPT_H1 | ACCEPT_H1 | **100% Accurate** | 0% Missed Errors |

### Key Observations:
- **Massive Compute Reduction**: Reliable outputs were accepted in an average of **26 steps**, and as few as **14 steps** (an **86% savings in LLM verifier calls**).
- **Instant Abort on Errors**: Flawed generations were caught and aborted in an average of **16.7 steps** (as fast as **8 steps**).
- **Exact Error Guarantee**: Across all seeds, the observed Type I error was **0.0%** ($\le 5\%$) and Type II error was **0.0%** ($\le 10\%$).

---

## 4. Anytime-Valid Confidence Sequences (Howard et al., 2021)

Classical confidence intervals suffer from **p-hacking** if you inspect them sequentially: if you check a 95% confidence interval repeatedly as samples arrive, the probability of falsely crossing the boundary approaches 100%.

**Empirical Bernstein Confidence Sequences** resolve this by constructing **time-uniform boundaries**:
$$\mathbb{P}\left( \forall t \ge 1, \; \mu^* \in [L_t, U_t] \right) \ge 1 - \alpha$$

In `EXP-SEQ-001`, the anytime-valid interval maintained mathematical inclusion of the true error rate across **100% of evaluation steps and all stopping times**:

```
Time-Uniform Empirical Bernstein Bounds (Seed 137)
Step t = 10:  Estimated Mean = 0.000 | CI = [0.000, 0.412] (Uncertain)
Step t = 14:  Estimated Mean = 0.000 | CI = [0.000, 0.324] (SPRT Accepts H0)
Step t = 50:  Estimated Mean = 0.020 | CI = [0.000, 0.118] (Tight Bound)
Step t = 100: Estimated Mean = 0.020 | CI = [0.000, 0.076] (Converged)
--> True mean (0.03) remained strictly inside [L_t, U_t] for all t >= 1.
```

---

## 5. Ville's Martingale E-Values for Drift Testing

To track whether an agent's runtime error rate drifts over time without fixing sample sizes, we track a Kelly-betting test martingale:
$$M_t = \prod_{i=1}^t E_i, \quad E_t = 1 + \lambda (X_t - p_0)$$

By **Ville's inequality**, the probability that $M_t$ ever exceeds $1/\alpha = 20.0$ under $H_0$ is at most $\alpha$:
$$\mathbb{P}(\exists t \ge 1 : M_t \ge 20.0) \le 0.05$$

- On reliable systems ($H_0$), the martingale decayed safely to near zero ($M_t \approx 0.04$), producing **zero false alarms**.
- On flawed systems ($H_1$), the martingale escalated rapidly, crossing the $20.0$ threshold in every seed (averaging step $41$) and triggering a definitive alarm.

---

## 6. How to Build Sequential Verification in Agent Frameworks

If you are designing agent workflows, adopting sequential verification provides three immediate architectural advantages:

1. **Dynamic Verification Budgets**: Allocate maximum verification checks ($T_{\text{max}} = 50$) as an upper ceiling, but terminate early via Wald's boundary as soon as $\Lambda_t \le -2.25$.
2. **Early Error Pruning**: Stop running multi-agent tasks the moment an early reasoning step crosses $\Lambda_t \ge +2.89$.
3. **Anytime-Valid Telemetry**: Monitor streaming production quality using test martingales and empirical Bernstein sequences rather than batch p-values.

---

*Explore the open-source implementation and test reproduction in our [LangChain MCP Researcher Platform](https://github.com/rohitsh16/langchain-mcp-researcher).*
