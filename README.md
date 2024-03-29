# Tzap: Find what you need, prompt with it! Embeddings and prompts combined!
[![Discord](https://img.shields.io/badge/Discord-Join%20Our%20Community-blue?logo=discord&logoColor=white)](https://discord.gg/Y8WpMmFH)
[![Official Page](https://img.shields.io/badge/Official%20Page-Visit%20Our%20Site-blue?logo=internet-explorer)](https://tzap.io)
[![NPM](https://img.shields.io/badge/NPM-tzap-blue?logo=npm)](https://npmjs.com/package/tzap)
[![Twitter Follow](https://img.shields.io/twitter/follow/tzap_io?style=social)](https://twitter.com/tzap_io)
## What is Tzap?
Tzap is tool to make it possible for you to search just about anything from "endpoints" and "color" to full code snippets or use cases to find the most relevant files.

When you run the `tzap prompt` command, Tzap combines your prompt with the extracted context and generates a suitable prompt for the GPT model. This allows GPT to generate both very complex and highly specific code.

# Demo (prompt: How would you expose tzap through a golang echo backend?)
![prompt demo](https://raw.githubusercontent.com/tzapio/tzap/main/docs/promptdemo.gif)
# Comparing generation to existing code:
![prompt demo](https://raw.githubusercontent.com/tzapio/tzap/main/docs/comparison.png)
### Quick install
```bash
# Choose how to install. curl, npm, npx
curl https://tzap.io/sh | bash
# npm install -g tzap
# npx tzap

# Provide the apikey. env variable or .env file
export OPENAI_APIKEY=<apikey>
# echo "OPENAI_APIKEY=<apikey>" > .env

tzap init
# Adapt below to your project!
tzap search "where is the background color defined?"
tzap search "using endpoints"

tzap prompt "add a email verification regexp to the user field"
```

## Resources
- Join our [Discord community](https://discord.gg/88xDVYbPVB) to get help, discuss features, and share your projects.
- Visit our official page at [https://tzap.io](https://tzap.io).
- Check out our NPM package: [tzap](https://npmjs.org/package/tzap).
- Follow us on Twitter: [@tzap_io](https://twitter.com/tzap_io)

## Notes:
Tzap is in a beta phase.
Tzap has the power to overwrite existing files, so commit local changes first. 
Using embeddings will upload most files to OpenAIs servers. (Alternative solutions are being looked into)
Using external APIs incurs small costs, read [Cost Estimation](#cost-estimation).

# Tzap provides a few commands:

- `tzap init`: Initializes Tzap in your project directory, creating necessary configuration files and folders.
- `tzap prompt`: Generates code suggestions based on the provided prompt and your project's context.
- other fun commands 
    - `tzap commit`: Makes a git commit suggestion based on diff. It does not use wide context.
    - `tzap ghrelease`: A bit hardcoded and undocumented, makes release notes for github. 

## Key Features

- Simple CLI tool
- Automate GPT copy-pasting tasks
- Effortlessly build context and pull in files
- Local vector database

## How It Works

**Tzap works in the following steps:**
- Init: Initializing a project is done with `tzap init`. In order to limit costs, Tzap requires a specification of both what to ignore and is allowed to include. `.gitignore` and `.tzapignore` is FIRST applied and removes all file matches. THEN `.tzapinclude` further filters out all NON-MATCHES.
- Indexing: When you run `tzap prompt`, Tzap builds a cache, indexes your project directory and builds a vector database of your code files. This allows Tzap to efficiently search for relevant code snippets during the code generation process. Note: This process uploads all file matches to OpenAI.
- Prompt Generation: Tzap takes the prompt string that describes the code you want to generate. Tzap combines your prompt with the extracted context information, such as interfaces, types, ORM, and libraries, to build a specific prompt for the GPT model.
- Code Generation: Tzap sends the generated prompt to the GPT model, which produces code suggestions based on the provided context and the prompt. These suggestions are then presented to you for further evaluation and integration into your codebase.

By automating code generation tasks and leveraging GPT's language capabilities, Tzap simplifies the process of writing code and helps you push features  consistent code styles within your project.

## Tips
The typical use case for `tzap prompt` is when:
1. there is some existing code (does not have to be related code to the prompt, like existing datamodels or endpoints)
2. there is some idea of a starting point (How do I create an endpoint)
3. there is some idea of the goal (that enables customers to change subscription)
4. run `tzap prompt "How do I create an endpoint that enables customers to change subscription"`
5. based on the prompt, GPT might provide a general non-code answer. In such case, instruct it with feedback.

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
curl https://tzap.io/sh | bash 
```

Once you have installed Tzap, you can start using it immediately by typing `tzap` in your terminal.

## Tzap Embedding Prompt CLI

The `tzap prompt` command allows you to generate code or content based on embedding search from your existing code. This enables GPT to incorporate your code knowledge and provide more relevant and accurate outputs. By providing additional context, GPT can generate even more accurate code snippets or content based on the provided files.

1. Use the `tzap embeddingprompt` command with the following syntax:
```
tzap prompt "<prompt>" [--inspiration=filepath1,filepath2,...]
```

- `<prompt>`: The text prompt you want to provide to GPT. Make your prompt as clear as possible and enclose it in double quotes.
- `-i`: (Optional) A flag followed by a comma-separated list of file paths, which are used as inspiration files to enhance GPT's general understanding.
- `--temperature`: Fine-tune the temperature.
- `-s searchQuery`: Split the search and the prompt.
- `-f promptFile`: Use a file as a prompt.

For example, you can run the following command to generate code based on a prompt:

```
tzap prompt "Generate a function that calculates the sum of two numbers" --inspiration=myFile.js,anotherFile.js
```

Feel free to customize `<prompt>` and provide your own inspiration files to get more accurate code snippets or content.

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

