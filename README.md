# 📘 Go Course Modernization (Unofficial)

This project modernizes the [Go programming course by Matt Holiday](https://www.youtube.com/playlist?list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6) — originally recorded using **Go v1.15** — by reviewing each video and updating the content to reflect the latest Go versions up to **Go v1.24**.

The project uses a custom-built Go application and a two-stage **LLM (Large Language Model)** pipeline to:

- **Summarize** each course video into key concepts and Go code examples.
- **Fact-check** those concepts using official Go release notes (v1.16–v1.24).
- **Identify outdated features or practices** and provide updated code when needed.
- **Generate well-structured Markdown documentation** that summarizes each video and highlights relevant updates.

> ⚠️ **Disclaimer:**  
> This is an independent, community-driven effort. I **do not own** the original course material, nor do I intend to monetize any part of it.  
> This project is purely educational — a lot has changed in Go since the course was recorded, and the goal is to help learners get the most out of Matt’s excellent Go tutorials by aligning them with the most up-to-date Go language features and best practices.

## 🎓 Credits & Original Course

This project is based on the excellent Go programming course by Matt Holiday on YouTube:  
👉 [Go Programming Course (YouTube)](https://www.youtube.com/playlist?list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6)

Matt has done a fantastic job making Go approachable and practical. If you haven’t already, please:

- ✅ **Subscribe to Matt’s YouTube channel**
- ✅ **Watch the full course**
- ✅ **Share with other Go learners**

## 🔍 What This Project Contains

Each markdown file in this repository corresponds to one video from the original playlist. Inside each file, you’ll find:

- `# Title` — The original video title  
- `## Summary` — A high-level overview of what the video covers  
- `## Key Points` — Summarized technical topics and original code examples  
- `## What’s New` — Changes since Go 1.15 based on release notes  
- `## Updated Code Snippets` — Revised examples, if needed  
- `## Citations` — Go version references from official release notes

## 🛠 How It Works

Behind the scenes, this repo is powered by a Go application that:

1. Retrieves the URL of each video from the YouTube playlist.
2. Uses an LLM (Google Gemini) to summarize the video as well as cross-references each technical claim against official Go release notes (v1.16 to v1.24).
3. Outputs everything into structured Markdown files.

All changes are carefully scoped to **factual updates only** — no opinions, no re-interpretation.

## 🧪 Local Execution

> Needs Go v1.24 or newer

To run the application locally (e.g., to regenerate the markdown files):

1. Clone the repository:

   ```bash
   git clone https://github.com/Ar4v1nd/go-course-modernizer.git
   cd go-course-modernizer
   ```

2. Set up environment variables:

    Update the existing `.env` file with your API keys:

    ```ini
    YOUTUBE_API_KEY=<INSERT_API_KEY_HERE>
    GEMINI_API_KEY=<INSERT_API_KEY_HERE>
    ```

    - [Get a YouTube API key](https://developers.google.com/youtube/v3)
    - [Get a Gemini API key](https://aistudio.google.com/apikey)

3. Run the application:

    ```bash
    go run .
    ```

## 🙌 Contributing

If you notice something that could be improved or a change in the Go language not yet reflected here, feel free to open a PR or issue.

> This project exists to **support** and **extend the educational value** of Matt Holiday’s course — not to replace it.

## 📜 License

This project is MIT licensed.

> All video content and original teaching materials are © Matt Holiday and hosted on YouTube.
