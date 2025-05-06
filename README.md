# 🏔 DriftDune

**Catch LLM Drift Before It Hits Production.**

> A CLI tool that detects when your LLM’s output starts to behave differently — even if no one told you the model changed.

---

## ⚡ Quick Start

Install:

```bash
git clone https://github.com/DriftDuneAI/driftdune-cli.git
cd driftdune-cli
pip install -r requirements.txt
```

Run:

```bash
python driftdune.py diff examples/v1_outputs.json examples/v2_outputs.json
```

Output:

```
🚨 Drift Detected: 18.2% semantic change across 100 outputs.
Top offenders: prompt_04.json, prompt_17.json
```

---

## 🔍 What It Does

| Feature                | Description                                           |
|------------------------|-------------------------------------------------------|
| ✅ Embedding Diff      | Compares outputs using `text-embedding-3-small`      |
| ✅ Semantic Drift %    | Calculates average vector shift between versions     |
| ✅ Top Offenders       | Flags highest-drift prompts                          |
| ✅ CLI-First           | One-line install & usage. No dashboard needed.       |
| ⏳ Coming Soon         | YAML policy packs, Vault replays, Slack alerts       |

---

## 🧠 Why This Matters

LLMs silently change:
- Model updates aren’t always announced
- Prompt outputs degrade or hallucinate
- You don’t know it’s broken until your users do

**DriftDune** helps you spot drift early — before it costs you credibility, customers, or compliance.

---

## 📦 Inputs & Examples

Inputs:
- `v1_outputs.json` — outputs from your model on Day 1
- `v2_outputs.json` — outputs from the same prompts on Day N

Try it with files in the [`examples/`](./examples) folder.

---

## 🛠 How It Works

1. Loads two versions of output for each prompt
2. Converts responses to vectors using OpenAI’s `text-embedding-3-small`
3. Compares them via cosine similarity
4. Flags drifted prompts below similarity threshold (default: 0.85)

---

## 🔐 API Keys & Privacy

- You’ll need an OpenAI API key.
- Data is processed locally; no logs are sent anywhere.
- Future versions will include optional Vault + alerting.

---

## 🚀 Roadmap

- [ ] CLI-based YAML policy enforcement
- [ ] Slack/webhook drift alerts
- [ ] Drift Vector Vault with replay
- [ ] Hosted dashboard (opt-in, non-required)

---

## 🤝 Contributing

Found a bug? Want to suggest a feature?  
Open an issue or PR — contributions welcome.

---

## 🪙 License

MIT — free for personal or commercial use.

