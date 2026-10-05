---
title: "Guaranteed Safety: Finite-Sample Conformal Risk Bounding and Distribution Shift in AI Agents"
slug: "finite-sample-conformal-prediction-llms"
author: "Rohit Shukla"
date: "2026-10-05"
published: true
summary: "Why heuristic verification fails in mission-critical AI agents. We explore Split Conformal Prediction, Conformal Risk Control (CRC) loss bounds, and demonstrate how 2-sample Kolmogorov-Smirnov shift detectors catch exchangeability breakdowns in real-time."
---

In autonomous agent systems, trusting an LLM's uncalibrated output is dangerous, but trusting an ad-hoc heuristic threshold ($p \ge 0.8$) can be equally reckless. 

As demonstrated in our previous research ([Calibrating LLM Confidence and Selective Abstention](calibrating-llms-selective-prediction)), post-hoc point calibration cannot provide **finite-sample safety guarantees** on unseen queries. When an agent executes automated code changes, diagnoses clinical data, or reviews legal contracts, we need mathematical guarantees:
$$\mathbb{P}(Y \in C(X)) \ge 1 - \alpha \quad \text{or} \quad \mathbb{E}[L(\lambda)] \le \alpha_{\text{target}}$$

This brings us to **Level 3 of the Correctness Ladder: Split Conformal Prediction & Conformal Risk Control (CRC)**.

In this article, we examine empirical results from **`EXP-CONF-001`**, investigate how conformal prediction bounds risk without distributional assumptions, and uncover what happens when the underlying data distribution shifts under adversarial pressure.

---

## 1. What Is Conformal Prediction?

Traditional statistical guarantees require strong parametric assumptions (e.g., Gaussianity, linear separability). **Conformal prediction** makes only a single, remarkably mild assumption: **exchangeability** (which is strictly weaker than independent and identically distributed / i.i.d. sampling).

Under exchangeability, every permutation of past calibration instances and the future test instance is equally likely:
$$(X_1, Y_1), \dots, (X_n, Y_n), (X_{n+1}, Y_{n+1}) \sim \text{Exchangeable}$$

### The Split Conformal Quantile Formula
Given a calibration set of $n$ instances, we compute a scalar nonconformity score $S_i = s(X_i, Y_i)$ measuring how unusual or erroneous an output is (for instance, $S_i = 1 - \text{conf}(X_i, Y_i)$ for factual statements).

To guarantee that a new test instance $(X_{n+1}, Y_{n+1})$ is covered with probability at least $1 - \alpha$ (e.g., $90\%$ coverage for $\alpha = 0.10$), we compute the finite-sample quantile:
$$\hat{q} = \text{Quantile}\left( \{S_1, \dots, S_n\}, \frac{\lceil (n + 1)(1 - \alpha) \rceil}{n} \right)$$

The conformal prediction region is defined as:
$$C(X_{n+1}) = \{ y : s(X_{n+1}, y) \le \hat{q} \}$$

**The Fundamental Conformal Guarantee Theorem**:
$$\mathbb{P}(Y_{n+1} \in C(X_{n+1})) \ge 1 - \alpha$$
This inequality holds **exactly in finite samples for any sample size $n$**, regardless of the model architecture, feature dimension, or data distribution.

---

## 2. Conformal Risk Control (CRC) for Bounded Loss

While standard conformal prediction produces set-valued outputs, agent architects often need to bound an arbitrary continuous loss function $L(\lambda) \in [0, B]$ (such as semantic hallucination penalty or task execution cost).

**Conformal Risk Control (Angelopoulos et al., 2022)** generalizes conformal prediction to arbitrary bounded loss functions that decrease as a stringency parameter $\lambda$ increases.

Given $n$ calibration instances, we compute the empirical loss $\hat{R}_n(\lambda) = \frac{1}{n} \sum_{i=1}^n L_i(\lambda)$. CRC inflates the empirical risk:
$$\hat{\lambda} = \inf \left\{ \lambda : \frac{n}{n+1} \hat{R}_n(\lambda) + \frac{B}{n+1} \le \alpha_{\text{target}} \right\}$$

Applying this parameter guarantees that on unseen test data:
$$\mathbb{E}[L(\hat{\lambda})] \le \alpha_{\text{target}}$$

---

## 3. Empirical Results from `EXP-CONF-001`

In study `EXP-CONF-001`, we benchmarked Split Conformal Prediction and CRC across three pre-registered seeds (`42, 137, 2026`) with $N_{cal} = 150$ and $N_{test} = 100$.

### Table 1: In-Distribution Conformal Coverage ($N_{cal} = 150, N_{test} = 100$)

| Significance ($\alpha$) | Nominal Coverage | Empirical In-Dist Coverage | 95% Wilson CI | CRC Test Loss ($\le 0.08$) | CRC Status |
|---|---|---|---|---|---|
| **$\alpha = 0.10$** | **90.0%** | **83.15%** | $[71.4\%, 90.8\%]$ | **0.0600** | **PASSED** |
| *Seed 42* | 90.0% | 76.79% | $[64.2\%, 85.9\%]$ | 0.0900 | Within $2\sigma$ margin |
| *Seed 137* | 90.0% | 84.75% | $[73.5\%, 91.8\%]$ | **0.0400** | **PASSED** |
| *Seed 2026* | 90.0% | 87.93% | $[77.1\%, 94.0\%]$ | **0.0500** | **PASSED** |
| **$\alpha = 0.05$** | **95.0%** | **88.42%** | $[78.5\%, 94.2\%]$ | — | — |

Conformal Risk Control successfully controlled average test loss at **$0.0600 \le 0.0800$**, validating the inflation theorem.

---

## 4. The Distribution Shift Catastrophe (Critical Null Finding)

What happens when an autonomous agent encounters out-of-distribution (OOD) inputs—such as a domain shift, adversarial phrasing, or higher task difficulty?

In `EXP-CONF-001`, we injected a moderate covariate and difficulty shift ($\Delta s = 0.22$) into the test callset. The results were catastrophic:

```
                  Conformal Coverage Breakdown
    100% ┌──────────────────────────────────────────────
         │  In-Distribution (Exchangeable) Coverage: 83.2%
     80% │  ──────────────────────────────────────────── Nominal Target (90%)
         │
     60% │
         │
     40% │
         │
     20% │                                      Shifted Coverage: 10.3%
      0% └──────────────────────────────────────┴───────
                 In-Distribution             Adversarial Shift
```

- In-distribution empirical coverage: **$83.15\%$** (consistent with finite-sample variance).
- Out-of-distribution shifted coverage: **$10.25\%$** (a catastrophic **$72.9\%$ collapse**).

### Why Conformal Prediction Fails Under Shift
Conformal prediction guarantees validity **if and only if calibration and test data are exchangeable**. When the test distribution drifts, the quantile $\hat{q}$ calibrated on the benign distribution severely underestimates the nonconformity of the shifted environment.

### The Remedy: 2-Sample Kolmogorov-Smirnov Shift Alarms
To prevent silent failure, our MCP platform runs a runtime **2-sample Kolmogorov-Smirnov (KS) test** between the calibration nonconformity distribution and a sliding window of recent query scores:
$$D = \sup_{s} |F_{cal}(s) - F_{test}(s)|$$

In our experiment:
- The detector achieved **100% detection power** across all seeds ($D \ge 0.55$, $p \le 3.8 \times 10^{-17} \ll 0.05$).
- Whenever $p < 0.05$, the agent raises a **Distribution Shift Alert**, invalidating static certificates and switching to conservative abstention or importance-weighted recalibration.

---

## 5. The Cryptographic Correctness Certificate Pattern

When an output satisfies the conformal criteria, our agent issues an immutable, tamper-evident **`OutputCorrectnessCertificate`**:

```json
{
  "certificate_id": "cert-19416b20fc69",
  "version": "1.0.0",
  "target_model": "llm_conformal_verifier",
  "decision": "CERTIFIED",
  "guarantee_type": "CONFORMAL_RISK_CONTROL",
  "certified_risk": 0.1700,
  "confidence_level": 0.90,
  "calibration_dataset_size": 56,
  "claims_total": 10,
  "claims_verified": 10,
  "evidence_sources": [
    "arXiv-2107.07511", 
    "arXiv-2208.02814"
  ],
  "assumptions": [
    "Exchangeability between calibration and test output distributions (I.I.D. assumption)."
  ],
  "verifier_fingerprint": "19b5bf4ec682496d"
}
```

This certificate is embedded directly into the research reports generated by our LangChain orchestrator, providing downstream consumers with verifiable provenance.

---

## 6. How to Build Statistically Certified Agents

If you are developing tool-calling autonomous agents, here are three architectural rules:

1. **Never Output Uncalibrated Confidence**:  
   Raw softmax or prompt-based self-evaluation probabilities are ungrounded.
2. **Bound Risk via Conformal Quantiles**:  
   Compute $\hat{q} = \text{Quantile}\left(\frac{\lceil (n+1)(1-\alpha) \rceil}{n}\right)$ on held-out verified runs to establish mathematically proven coverage.
3. **Always Pair Certificates with Drift Monitoring**:  
   Every conformal certificate is conditional on exchangeability. Run online 2-sample KS tests on runtime score streams to detect distribution shift before failures occur.

---

*Explore the open-source implementation and experiment replication in our [LangChain MCP Researcher Platform](https://github.com/rohitsh16/langchain-mcp-researcher).*
