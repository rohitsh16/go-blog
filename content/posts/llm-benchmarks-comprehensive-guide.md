---
title: "The Complete Guide to LLM Benchmarks: Every Evaluation That Matters"
slug: "llm-benchmarks-comprehensive-guide"
author: "Rohit Shukla"
date: "2026-09-20"
published: true
summary: "A comprehensive reference of every major LLM benchmark — from reasoning and coding to safety and multimodal tasks — with direct links to papers, leaderboards, and datasets."
---

Large Language Models are only as trustworthy as the benchmarks used to evaluate them. With dozens of evaluations measuring everything from common-sense reasoning to PhD-level science, knowing **what each benchmark tests** and **where to find it** is essential for anyone building with or evaluating LLMs.

This post is a living reference of every major LLM benchmark in active use today.

## General Knowledge & Reasoning

- **MMLU (Massive Multitask Language Understanding)** — 57 subjects from STEM to humanities, testing broad factual recall and reasoning. [Paper](https://arxiv.org/abs/2009.03300) | [Dataset](https://huggingface.co/datasets/cais/mmlu)
- **MMLU-Pro** — A harder, more discriminating version of MMLU with 10 answer choices and chain-of-thought focus. [Paper](https://arxiv.org/abs/2406.01574) | [Dataset](https://huggingface.co/datasets/TIGER-Lab/MMLU-Pro)
- **HellaSwag** — Sentence completion benchmark testing commonsense natural language inference. [Paper](https://arxiv.org/abs/1905.07830) | [Dataset](https://rowanzellers.com/hellaswag/)
- **ARC (AI2 Reasoning Challenge)** — Grade-school science questions in Easy and Challenge sets. [Paper](https://arxiv.org/abs/1803.05457) | [Dataset](https://allenai.org/data/arc)
- **WinoGrande** — Large-scale Winograd Schema challenge for commonsense coreference resolution. [Paper](https://arxiv.org/abs/1907.10641) | [Dataset](https://huggingface.co/datasets/allenai/winogrande)
- **TruthfulQA** — Tests whether models generate truthful answers vs. common misconceptions. [Paper](https://arxiv.org/abs/2109.07958) | [Dataset](https://github.com/sylinrl/TruthfulQA)
- **BoolQ** — Yes/no question answering from naturally occurring questions. [Paper](https://arxiv.org/abs/1905.10044) | [Dataset](https://huggingface.co/datasets/google/boolq)
- **PIQA (Physical Intuition QA)** — Tests physical commonsense reasoning about everyday objects. [Paper](https://arxiv.org/abs/1911.11641) | [Dataset](https://yonatanbiri.github.io/PIQA/)

## Mathematics

- **GSM8K** — 8,500 grade-school math word problems requiring multi-step arithmetic reasoning. [Paper](https://arxiv.org/abs/2110.14168) | [Dataset](https://huggingface.co/datasets/openai/gsm8k)
- **MATH** — 12,500 competition-level problems from AMC, AIME, and Olympiad with LaTeX solutions. [Paper](https://arxiv.org/abs/2103.03874) | [Dataset](https://huggingface.co/datasets/hendrycks/competition_math)
- **MathVista** — Visual math reasoning combining diagrams, charts, and geometry with mathematical problem-solving. [Paper](https://arxiv.org/abs/2310.02255) | [Leaderboard](https://mathvista.github.io/)
- **MGSM (Multilingual GSM)** — GSM8K translated into 10 languages for multilingual math evaluation. [Paper](https://arxiv.org/abs/2210.03057) | [Dataset](https://huggingface.co/datasets/juletxara/mgsm)

## Coding & Programming

- **HumanEval** — 164 hand-written Python programming problems measuring functional code generation. [Paper](https://arxiv.org/abs/2107.03374) | [Dataset](https://github.com/openai/human-eval)
- **HumanEval+** — Extended HumanEval with 80x more tests per problem to catch false positives. [Paper](https://arxiv.org/abs/2305.01210) | [Dataset](https://github.com/evalplus/evalplus)
- **MBPP (Mostly Basic Python Problems)** — 974 crowd-sourced Python programming tasks for entry-level coding. [Paper](https://arxiv.org/abs/2108.07732) | [Dataset](https://huggingface.co/datasets/google-research-datasets/mbpp)
- **SWE-bench** — Real-world GitHub issues from popular Python repos; models must generate working patches. [Paper](https://arxiv.org/abs/2310.06770) | [Leaderboard](https://www.swebench.com/)
- **SWE-bench Verified** — Human-validated subset of SWE-bench for more reliable evaluation. [Leaderboard](https://www.swebench.com/)
- **LiveCodeBench** — Continuously updated coding benchmark from competitive programming contests. [Paper](https://arxiv.org/abs/2403.07974) | [Leaderboard](https://livecodebench.github.io/)
- **MultiPL-E** — HumanEval and MBPP translated to 18+ programming languages. [Paper](https://arxiv.org/abs/2208.08227) | [Dataset](https://huggingface.co/datasets/nuprl/MultiPL-E)
- **BigCodeBench** — 1,140 practical programming tasks requiring diverse library usage. [Paper](https://arxiv.org/abs/2406.15877) | [Leaderboard](https://bigcode-bench.github.io/)

## Science & Domain Knowledge

- **GPQA (Graduate-Level Google-Proof QA)** — Expert-written questions in physics, chemistry, and biology that domain experts struggle with. [Paper](https://arxiv.org/abs/2311.12022) | [Dataset](https://huggingface.co/datasets/Idavidrein/gpqa)
- **GPQA Diamond** — The hardest, most vetted subset of GPQA. [Dataset](https://huggingface.co/datasets/Idavidrein/gpqa)
- **MedQA** — USMLE-style medical licensing exam questions. [Paper](https://arxiv.org/abs/2009.13081) | [Dataset](https://huggingface.co/datasets/bigbio/med_qa)
- **PubMedQA** — Biomedical research question answering from PubMed abstracts. [Paper](https://arxiv.org/abs/1909.06146) | [Dataset](https://pubmedqa.github.io/)
- **SciQ** — Crowdsourced science exam questions with supporting evidence passages. [Dataset](https://huggingface.co/datasets/allenai/sciq)

## Instruction Following & Chat

- **MT-Bench** — Multi-turn conversation benchmark scored by GPT-4 judges across 8 categories. [Paper](https://arxiv.org/abs/2306.05685) | [Code](https://github.com/lm-sys/FastChat/tree/main/fastchat/llm_judge)
- **AlpacaEval** — Automated evaluation of instruction-following against a reference model. [Leaderboard](https://tatsu-lab.github.io/alpaca_eval/)
- **Chatbot Arena** — Live human preference rankings via blind pairwise comparisons on LMSYS. [Leaderboard](https://chat.lmsys.org/?leaderboard)
- **IFEval (Instruction Following Eval)** — Verifiable instruction-following with format constraints (e.g., "write exactly 3 paragraphs"). [Paper](https://arxiv.org/abs/2311.07911) | [Dataset](https://huggingface.co/datasets/google/IFEval)
- **WildBench** — Challenging real-world user queries with automated scoring and length-controlled win rates. [Paper](https://arxiv.org/abs/2406.04770) | [Leaderboard](https://huggingface.co/spaces/allenai/WildBench)

## Multimodal (Vision + Language)

- **MMMU (Massive Multi-discipline Multimodal Understanding)** — College-level questions with images spanning 30 subjects. [Paper](https://arxiv.org/abs/2311.16502) | [Leaderboard](https://mmmu-benchmark.github.io/)
- **MMLUx** — Multimodal extension of MMLU incorporating visual reasoning. [Dataset](https://huggingface.co/datasets/TIGER-Lab/MMLU-Pro)
- **VQAv2 (Visual Question Answering)** — Open-ended questions about natural images. [Paper](https://arxiv.org/abs/1612.00837) | [Dataset](https://visualqa.org/)
- **TextVQA** — VQA requiring reading and reasoning about text in images. [Paper](https://arxiv.org/abs/1904.08920) | [Dataset](https://textvqa.org/)
- **DocVQA** — Question answering on scanned document images. [Paper](https://arxiv.org/abs/2007.00398) | [Leaderboard](https://rrc.cvc.uab.es/?ch=17)
- **ChartQA** — Reasoning about charts and data visualizations. [Paper](https://arxiv.org/abs/2203.10244) | [Dataset](https://huggingface.co/datasets/ahmed-masry/ChartQA)

## Reasoning & Problem Solving

- **ARC-AGI** — Abstraction and Reasoning Corpus measuring fluid intelligence / novel pattern recognition. [Paper](https://arxiv.org/abs/1911.01547) | [Dataset](https://github.com/fchollet/ARC-AGI)
- **BBH (BIG-Bench Hard)** — 23 hardest tasks from BIG-Bench where prior models failed. [Paper](https://arxiv.org/abs/2210.09261) | [Dataset](https://huggingface.co/datasets/lukaemon/bbh)
- **DROP** — Discrete Reasoning Over Paragraphs requiring numerical extraction and arithmetic. [Paper](https://arxiv.org/abs/1903.00161) | [Dataset](https://allenai.org/data/drop)
- **MUSR** — Multi-step Soft Reasoning with murder mysteries, team assignments, and object placement. [Paper](https://arxiv.org/abs/2310.16049) | [Dataset](https://huggingface.co/datasets/TAUR-Lab/MuSR)
- **ZebraLogic** — Logic grid puzzles testing systematic constraint satisfaction. [Paper](https://arxiv.org/abs/2407.01370) | [Leaderboard](https://huggingface.co/spaces/allenai/ZebraLogic)
- **FrontierMath** — Extremely hard, original math problems created by professional mathematicians. [Announcement](https://epochai.org/frontiermath)

## Long Context & Retrieval

- **RULER** — Synthetic long-context evaluation measuring recall across 4K–128K token windows. [Paper](https://arxiv.org/abs/2404.06654)
- **Needle in a Haystack** — Retrieval of a planted fact within long contexts of varying length. [Blog](https://www.anthropic.com/news/claude-2-1-prompting)
- **LongBench** — Diverse long-context tasks including summarization, QA, and code completion. [Paper](https://arxiv.org/abs/2308.14508) | [Dataset](https://huggingface.co/datasets/THUDM/LongBench)
- **InfiniteBench** — Tasks requiring 100K+ tokens with novel documents, code debug, and dialogue. [Paper](https://arxiv.org/abs/2402.13718) | [Dataset](https://huggingface.co/datasets/xinrongzhang2022/InfiniteBench)
- **HELMET** — Holistic Evaluation of Long-context Models across 7 categories. [Paper](https://arxiv.org/abs/2410.02694)

## Safety, Bias & Alignment

- **ToxiGen** — Large-scale machine-generated toxic and benign statements targeting 13 minority groups. [Paper](https://arxiv.org/abs/2203.09509) | [Dataset](https://huggingface.co/datasets/skg/toxigen-data)
- **BBQ (Bias Benchmark for QA)** — Tests social biases across 11 categories including race, gender, and religion. [Paper](https://arxiv.org/abs/2110.08193) | [Dataset](https://github.com/nyu-mll/BBQ)
- **RealToxicityPrompts** — 100K naturally occurring prompts scored for toxicity of model completions. [Paper](https://arxiv.org/abs/2009.11462) | [Dataset](https://allenai.org/data/real-toxicity-prompts)
- **WMDP (Weapons of Mass Destruction Proxy)** — Measures dangerous knowledge in biosecurity, cybersecurity, and chemical weapons. [Paper](https://arxiv.org/abs/2403.03218) | [Dataset](https://huggingface.co/datasets/cais/wmdp)
- **HarmBench** — Standardized evaluation framework for automated red-teaming of LLMs. [Paper](https://arxiv.org/abs/2402.04249) | [Code](https://github.com/centerforaisafety/HarmBench)

## Agentic & Tool Use

- **BFCL (Berkeley Function Calling Leaderboard)** — Tests structured function/API calling across simple, parallel, and multi-step scenarios. [Leaderboard](https://gorilla.cs.berkeley.edu/leaderboard.html)
- **τ-bench (Tau-Bench)** — Real-world tool-use tasks simulating customer service workflows with database operations. [Paper](https://arxiv.org/abs/2406.12045) | [Code](https://github.com/sierra-research/tau-bench)
- **WebArena** — Autonomous web agent benchmark across e-commerce, forums, and maps. [Paper](https://arxiv.org/abs/2307.13854) | [Leaderboard](https://webarena.dev/)
- **OSWorld** — Desktop computer use benchmark testing GUI interaction across OS environments. [Paper](https://arxiv.org/abs/2404.07972) | [Leaderboard](https://os-world.github.io/)
- **GAIA** — General AI Assistant benchmark requiring multi-step web browsing, file handling, and reasoning. [Paper](https://arxiv.org/abs/2311.12983) | [Leaderboard](https://huggingface.co/spaces/gaia-benchmark/leaderboard)

## Multilingual

- **MMMLU** — MMLU machine-translated into 14 languages for cross-lingual knowledge evaluation. [Dataset](https://huggingface.co/datasets/openai/MMMLU)
- **FLORES-200** — Machine translation benchmark covering 200 languages. [Paper](https://arxiv.org/abs/2207.04672) | [Dataset](https://huggingface.co/datasets/facebook/flores)
- **XL-Sum** — Abstractive summarization in 44 languages from BBC articles. [Paper](https://arxiv.org/abs/2106.13822) | [Dataset](https://huggingface.co/datasets/csebuetnlp/xlsum)

## Aggregate Leaderboards

These meta-leaderboards compile results across many benchmarks:

- **Open LLM Leaderboard (Hugging Face)** — Community-driven evaluation of open models. [Leaderboard](https://huggingface.co/spaces/open-llm-leaderboard/open_llm_leaderboard)
- **Chatbot Arena (LMSYS)** — ELO-based ranking from human blind comparisons. [Leaderboard](https://chat.lmsys.org/?leaderboard)
- **SEAL Leaderboards (Scale AI)** — Expert-curated evaluations for coding, instruction following, and math. [Leaderboard](https://scale.com/leaderboard)
- **LiveBench** — Contamination-free benchmark with monthly-refreshed questions. [Leaderboard](https://livebench.ai/)
- **Artificial Analysis** — Tracks quality, speed, and pricing of API providers side by side. [Leaderboard](https://artificialanalysis.ai/)

> This post is a living document. As new benchmarks emerge and existing ones evolve, it will be updated to reflect the latest landscape of LLM evaluation.
