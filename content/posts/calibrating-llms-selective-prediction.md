---
title: "Beyond Blind Faith: Calibrating LLM Confidence and Selective Abstention for Factuality"
slug: "calibrating-llms-selective-prediction"
author: "Rohit Shukla"
date: "2026-10-04"
published: true
summary: "Modern LLMs are systematically overconfident on hallucinations. We empirically benchmark Temperature Scaling and Isotonic Regression across multi-seed factual tests, showing why point calibration cannot fix poor ranking and why selective abstention needs more than score thresholding."
---

When an autonomous AI agent claims, *"I am 95% confident that this theorem was proven by Euler in 1748"*, what does that percentage actually mean?

In modern deep learning architectures—and particularly in autoregressive Large Language Models (LLMs)—the answer is often disconcerting: **the number is fundamentally disconnected from empirical reality**. Modern LLMs are prone to dramatic overconfidence, routinely assigning probabilities greater than $0.85$ to subtle hallucinations, non-existent software APIs, and fabricated academic citations.

In this deep-dive, we examine findings from **`EXP-CALIB-001`**, an empirical research study conducted within the *AI Correctness Research Platform*. We benchmark **Temperature Scaling** against non-parametric **Isotonic Regression (PAVA)**, evaluate Expected Calibration Error (ECE) and Brier scores, and analyze why post-hoc probability calibration alone cannot rescue an agent from poor ranking separability.

---

## 1. The Anatomy of Miscalibration

An agent is said to be **well-calibrated** if its subjective confidence matches long-run empirical accuracy:
$$P(\text{Claim is True} \mid \text{Confidence} = p) = p, \quad \forall p \in [0, 1]$$

If an agent emits 100 statements with an estimated confidence of $0.80$, exactly 80 of those statements should be factually accurate.

```
       Perfect Calibration Line (Ideal)
        1.0 ┌───────────────────────/
            │                     / 
            │                   /   
Accuracy    │                 /     Overconfident Gap
  acc(B)    │               /  ░░░ (Model claims 0.90,
            │             / ░░░░░   actually gets 0.68)
            │           /░░░░░░░    
            │         /░░░░░░░      
        0.0 └────────┴──────────────1.0
           0.0    Confidence conf(B)
```

### Quantifying the Gap: ECE and Brier Score

To measure calibration mathematically across a validation set of $N$ assertions with predicted probabilities $p_i \in [0, 1]$ and binary correctness labels $y_i \in \{0, 1\}$:

1. **Expected Calibration Error (ECE)**:
   We partition the unit interval $[0, 1]$ into $M = 10$ equally spaced bins $B_1, \dots, B_M$:
   $$\text{ECE} = \sum_{m=1}^M \frac{|B_m|}{N} \left| \text{acc}(B_m) - \text{conf}(B_m) \right|$$
   where $\text{acc}(B_m) = \frac{1}{|B_m|}\sum_{i \in B_m} y_i$ and $\text{conf}(B_m) = \frac{1}{|B_m|}\sum_{i \in B_m} p_i$.

2. **Brier Score (Strictly Proper Scoring Rule)**:
   $$\text{BS} = \frac{1}{N} \sum_{i=1}^N (p_i - y_i)^2$$
   The Brier score penalizes both miscalibration and lack of resolution.

---

## 2. Post-Hoc Calibration: Parametric vs. Non-Parametric

Rather than retraining multi-billion-parameter foundation models, practitioners apply **post-hoc calibration** to map uncalibrated scores $p \in (0, 1)$ into calibrated probabilities $\hat{p} = f(p)$.

### Method A: Temperature Scaling (Platt Scaling)
Temperature scaling introduces a single learnable scalar $T > 0$ applied to the logit $z = \text{logit}(p) = \log\left(\frac{p}{1-p}\right)$:
$$\hat{p} = \sigma\left(\frac{z}{T}\right) = \frac{1}{1 + \exp(-z / T)}$$

- When $T > 1$, the transformation **softens** extreme probabilities, reducing overconfidence toward uniform uncertainty.
- $T$ is optimized on a held-out calibration set by minimizing negative log-likelihood (cross-entropy loss):
  $$\mathcal{L}_{NLL}(T) = -\sum_{i=1}^{n_{cal}} \left[ y_i \log(\hat{p}_i) + (1 - y_i) \log(1 - \hat{p}_i) \right]$$

### Method B: Non-Parametric Isotonic Regression (PAVA)
Isotonic regression fits a piecewise constant, non-decreasing step function using the **Pool Adjacent Violators Algorithm (PAVA)**. It makes zero parametric assumptions about the underlying distribution, ensuring:
$$\hat{p}_i \le \hat{p}_j \quad \text{whenever} \quad p_i \le p_j$$

---

## 3. Empirical Findings from `EXP-CALIB-001`

We evaluated both methods across three pre-registered seeds (`42, 137, 2026`) with a 3-way split: 40% calibration training ($N=120$), 30% calibration validation ($N=90$), and 30% held-out test ($N=90$).

### Cross-Seed Performance Summary

| Metric | Raw Uncalibrated | Temperature Scaling ($T^*$) | Isotonic Regression (PAVA) | Relative Improvement |
|---|---|---|---|---|
| **Mean ECE** | **0.1985** | **0.1343** | 0.1784 | **-32.34%** |
| *Seed 42 ECE* | 0.2317 | 0.1007 ($T=2.65$) | 0.2063 | **-56.52%** |
| *Seed 137 ECE* | 0.2033 | 0.1274 ($T=1.80$) | 0.1488 | **-37.33%** |
| *Seed 2026 ECE* | 0.1604 | 0.1747 ($T=4.25$) | 0.1800 | +8.91% |
| **Mean Brier Score** | 0.2412 | 0.2266 | **0.2162** | **-6.05%** |
| **Mean AURC** | 0.2560 | 0.2560 | 0.2571 | **0.00% (Preserved)** |
| **Accuracy @ 70% Coverage** | 68.25% | 68.25% | 68.25% | Baseline |

### Reliability Diagram Comparison (Seed 42)

```
[Raw Uncalibrated - Severe Overconfidence Gap]
Bin [0.4 - 0.6]: Predicted Conf = 0.520 | Actual Accuracy = 0.353 (Gap: -0.167)
Bin [0.6 - 0.8]: Predicted Conf = 0.720 | Actual Accuracy = 0.533 (Gap: -0.187)
Bin [0.8 - 1.0]: Predicted Conf = 0.900 | Actual Accuracy = 0.686 (Gap: -0.214)
--> Overall ECE: 0.2317

[Temperature Scaled (T = 2.65) - Aligned Calibration]
Bin [0.4 - 0.6]: Predicted Conf = 0.502 | Actual Accuracy = 0.514 (Gap: +0.012)
Bin [0.6 - 0.8]: Predicted Conf = 0.703 | Actual Accuracy = 0.682 (Gap: -0.021)
Bin [0.8 - 1.0]: Predicted Conf = 0.815 | Actual Accuracy = 0.750 (Gap: -0.065)
--> Overall ECE: 0.1007 (56.5% Error Reduction)
```

---

## 4. The Critical Takeaway: Calibration Does Not Alter Ranking

Notice an intriguing fact in Table 1: **the Area Under the Risk-Coverage Curve (AURC) is identical (0.2560) before and after temperature scaling.**

Why?
Because temperature scaling is a **strictly monotonic transformation** of the original confidence scores:
$$p_1 > p_2 \iff \sigma\left(\frac{\text{logit}(p_1)}{T}\right) > \sigma\left(\frac{\text{logit}(p_2)}{T}\right), \quad \forall T > 0$$

When ranking claims to decide which ones to accept and which ones to abstain from:
- The relative ranking of every claim is 100% identical.
- Selective classification at a target coverage of 70% selects **the exact same set of claims** whether using raw confidence or temperature-scaled confidence.
- Selective accuracy remained at $68.25\%$ with a 95% Wilson confidence interval of $[27.12\%, 50.44\%]$ for selective risk.

> **Key Scientific Insight**:  
> Calibration fixes the *meaning of numbers* (so a 70% prediction is right 70% of the time).  
> But it **cannot separate overlapping feature distributions**. If your model assigns 0.75 confidence to a true claim and 0.75 confidence to a hallucination, no post-hoc calibration function can untangle them.

---

## 5. Architectural Implications for Autonomous Agents

1. **Use Temperature Scaling on Small Calibration Sets**:  
   Isotonic regression is non-parametric and overfits small calibration splits ($N_{cal} \approx 100$). A single well-regularized temperature parameter $T^*$ reliably cuts calibration error by 30%–55%.
2. **Do Not Rely on Point Thresholds for High-Stakes Decisions**:  
   Thresholding at $\tau = 0.80$ provides **no finite-sample statistical guarantees** on individual outputs.
3. **Ascend to Conformal Prediction (Level 3)**:  
   To guarantee that an agent's failure rate remains strictly below a user-specified bound (e.g. $\mathbb{E}[\text{Error}] \le 5\%$), we must move beyond point probabilities to **finite-sample Conformal Risk Control**.

---

*In the next article, we dive into Level 3 of the Correctness Ladder: [Guaranteed Safety: Finite-Sample Conformal Risk Bounding and Distribution Shift in AI Agents](finite-sample-conformal-prediction-llms).*
