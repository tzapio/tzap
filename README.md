# Tzap: Think it and it's there. A toolset for GPT workflows, Prompt as Code Functions and simple CLI automations

## What is Tzap?
Tzap is a library that simplifies all things GPT and code. It provides both a CLI tool with pre-selected workflows and a toolkit to build, customize, and extend chatbot prompts in a streamlined and extensible manner. The library is designed to make it easy for developers to create reusable Tzap instances and combinations of Tzaps to quickly and effectively implement desired outcomes in their GPT-based applications.

![prompt demo](/public/promptdemo.gif)


### Quick install (NPM):
```bash
npm install -g tzap
export OPENAI_APIKEY=<apikey>

tzap init

# Do a git add <file> then do:
tzap commit

# Adapt below to your project!
tzap prompt outputfile.txt "can you add a tzap cli command that enables users \
to generate code based on the users code without requiring them to manage prompts themselves? \
Utilize Tzap embedding search."
```

## Notes:
Tzap is in a beta phase.
Tzap has the power to overwrite existing files, so commit local changes first. 
Using embeddings will upload most files to OpenAIs servers. (Alternative solutions are being looked into)
Using external APIs incurs small costs, read [Cost Estimation](#cost-estimation).

## Key Features

- Simple CLI tool
- Built-in local embedding vector database and cache.
- Easily create prompts with domain specific contexts using Tzap functions, workflows, loops and control flows.
- Build apps on top of Tzap and GPT
- Automate GPT copy-pasting tasks
- Integrate magic functions that evaluate GPT prompts instead of code
- Effortlessly manipulate file paths and directories
- Generate multi-modal content 

## How It Works

Tzap allows you to create reusable instances and apply workflows and functions to them, making it convenient to adapt to new use cases, such as automating GPT copy-pasting, creating magic functions that evaluate GPT prompts, and crafting magic CLI tools. In addition, Tzap makes it simple to apply workflows and functions to existing Tzaps, enhancing the library's flexibility.

By using Tzap, you can effortlessly manage file paths and directories, fetch chat responses, and generate content using OpenAI's GPT-4 model. Furthermore, the library provides chat message context in Golang, ensuring a smooth integration process.

With Tzap's intuitive design and powerful capabilities, you can quickly and efficiently implement desired outcomes in your GPT-based applications. So go ahead and give it a try, and let Tzap work its magic for you!

## Getting Started

### Installation and Usage

To install and run Tzap, simply use the following commands:

```bash
npm install -g tzap
```

or 

```bash
npx tzap
```

or

```bash
curl -sSL https://raw.githubusercontent.com/tzapio/tzap/main/cli/install.sh | bash 
```

Once you have installed Tzap, you can start using it immediately by typing `tzap` in your terminal.

## Tzap Embedding Prompt CLI

The `tzap embeddingprompt` command allows you to generate code or content based on embedding search from your existing code. This enables GPT to incorporate your code knowledge and provide more relevant and accurate outputs. By providing additional context, GPT can generate even more accurate code snippets or content based on the provided files.

1. Use the `tzap embeddingprompt` command with the following syntax:
 ```
 tzap embeddingprompt <output_file> "<prompt>" [--inspiration=filepath1,filepath2,...]
 ```

 - `<output_file>`: The name of the file you want to store the generated content.
 - `<prompt>`: The text prompt you want to provide to GPT. Make your prompt as clear as possible and enclose it in double quotes.
 - `--inspiration`: (Optional) A flag followed by a comma-separated list of file paths, which are used as inspiration files to enhance GPT's general understanding.

For example, this text was generated using:

```
npx tzap embeddingprompt -m gpt4 README.md2 "Add to the README.md an explanation on how to use the cli command tzap embeddingprompt" -i README.md 
``` 

Feel free to customize the `<output_file>` and `<prompt>` according to your desired outcomes.

You can also use the `--inspiration` flag to specify relevant inspiration files to provide GPT with better context:

### Tzap Semantic Git Commit CLI

Never write a git commit message again with Tzap! To try this feature, simply run:

```bash
export OPENAI_APIKEY=<openai_key>
tzap semantic:gitcommit
```

This command will automatically generate a meaningful git commit message based on your recent code changes.

## Cost Estimation

Using embeddings and external APIs comes with certain costs. Here's a breakdown of those costs:

### Embeddings

Initialization of embeddings for this project cost: $0.04 (Run once). (200 files, 6,500 lines of code)
Re-fetch when generating with new code changes: $0.0004 per 1,000 "words/tokens"

### GPT Model

You can choose to use GPT-4 or GPT-3.5 (default usually).

The maximum cost per file for GPT-4 is $0.1 (8,000 "word/token" limit).
The maximum cost per file for GPT-3.5 is $0.008 (4,000 "word/token" limit).

It is important to understand and manage these costs while using Tzap.

## Example Code

To start coding with Tzap, you can refer to the provided examples in the `examples` directory of the Tzap repository. These examples demonstrate various use cases and will help you understand how to utilize Tzap effectively.

Feel free to explore the Tzap library, experiment with its features, and create powerful GPT-based applications tailored to your requirements. Enjoy the convenience and efficiency that Tzap brings to your projects!
