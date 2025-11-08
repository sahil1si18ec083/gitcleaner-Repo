# 🧹 gitcleaner — Clean up your GitHub repositories easily

`gitcleaner` is a command-line tool written in **Golang** that helps you **list** and **delete** your GitHub repositories in bulk — safely, interactively, and efficiently.

---

## ✨ Features

- 🔍 **List all repositories** of the authenticated user  
- 🗑️ **Delete repositories**:
  - `--all` → delete all repos
  - `--only` → delete specific repos
  - `--except` → delete all except a few
- 🧪 **Dry run mode** → preview what will be deleted before executing  
- ⚠️ **Confirmation prompt** before actual deletion  
- 🧠 Designed with **Cobra CLI** for intuitive command usage  
- 🔐 Uses your **GitHub personal access token** securely via environment variable  

---

## 🚀 Getting Started

### 1️⃣ Clone the repo

```bash
git clone https://github.com/sahil1si18ec083/gitcleaner-Repo.git
cd gitcleaner-Repo
