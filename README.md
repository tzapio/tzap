# Tzap: Think it and it's there. A toolset for GPT templating, prompt as Code and simple CLI automations

### Quick install (NPM):
```bash
npm install -g tzap
export OPENAI_APIKEY=<apikey>

# Do a git add <file> then do:
tzap semantic:gitcommit
# Adapt below to your project!
tzap embeddingprompt outputfile.txt "can you add a tzap cli command that enables users \
to generate code based on the users code without requiring them to manage prompts themselves? \
Utilize Tzap embedding search."
```

## What is Tzap?
Tzap is a library that simplifies all things GPT and code. It provides both a CLI tool with pre-selected workflows and a toolkit to build, customize, and extend chatbot prompts in a streamlined and extensible manner. The library is designed to make it easy for developers to create reusable Tzap instances and combinations of Tzaps to quickly and effectively implement desired outcomes in their GPT-based applications.

## Notes:
Tzap is in a beta phase.
Using embeddings will upload most files to OpenAIs servers. (Alternative solutions are being looked into)
Using external APIs incurs small costs, read [Cost Estimation](#cost-estimation).

## Key Features

- Simple CLI tool
- Built-in local embedding vector database and cache.
- Easily create prompts with domain specific contexts using Tzap functions, templates, loops and control flows.
- Build apps on top of Tzap and GPT
- Automate GPT copy-pasting tasks
- Integrate magic functions that evaluate GPT prompts instead of code
- Effortlessly manipulate file paths and directories
- Generate multi-modal content 


## How It Works

Tzap allows you to create reusable instances and apply templates and functions to them, making it convenient to adapt to new use cases, such as automating GPT copy-pasting, creating magic functions that evaluate GPT prompts, and crafting magic CLI tools. In addition, Tzap makes it simple to apply templates and functions to existing Tzaps, enhancing the library's flexibility.

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

### Tzap Semantic Git Commit

Never write a git commit message again with Tzap! To try this feature, simply run:

```bash
export OPENAI_APIKEY=<openai_key>
tzap semantic:gitcommit
```

This command will automatically generate a meaningful git commit message based on your recent code changes.

## Cost estimations:
In order to search for relevant code Tzap index all files. This is done with openais    
Initialization of embeddings for this project: 0.04$ (Run once). (200 files, 6500 lines of code)
Re-fetch when generating with new code changes: 0.0004$ per 1000 "words/tokens"

You can chose to use GPT4 or GPT3.5 (default usually). 
The ceiling per file for gpt4, is 0.1$. (8k "word/token" limit) 
The ceiling per file for gpt3.5, is 0.008$. (4k "word/token" limit) 

## Example Code

To start coding with Tzap, you can refer to the provided examples in the `examples` directory of the Tzap repository. These examples demonstrate various use cases and will help you understand how to utilize Tzap effectively.

Feel free to explore the Tzap library, experiment with its features, and create powerful GPT-based applications tailored to your requirements. Enjoy the convenience and efficiency that Tzap brings to your projects!