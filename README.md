# Tzap
Tzap is a library designed to seamlessly integrate GPT prompts into your code. It simplifies the process of building, customizing, and extending GPT prompts, making it more efficient for developers to incorporate desired outcomes in their GPT-based applications. With Tzap, you can create reusable instances and combine them in various ways to meet your specific needs.

## Key Features
- Easily create reusable Tzap instances and templates
- Apply templates and functions to existing Tzaps
- Automate GPT copy-pasting tasks
- Integrate magic functions that evaluate GPT prompts instead of code
- Develop magic CLI tools
- Effortlessly manipulate file paths and directories
- Generate content using OpenAI's GPT-4 model
- Provide chat message context in Golang

## How It Works
Tzap allows you to create reusable instances and apply templates and functions to them, making it convenient to adapt to new use cases, such as automating GPT copy-pasting, creating magic functions that evaluate GPT prompts, and crafting magic CLI tools. In addition, Tzap makes it simple to apply templates and functions to existing Tzaps, enhancing the library's flexibility.

By using Tzap, you can effortlessly manage file paths and directories, fetch chat responses, and generate content using OpenAI's GPT-4 model. Furthermore, the library provides chat message context in Golang, ensuring a smooth integration process.

With Tzap's intuitive design and powerful capabilities, you can quickly and efficiently implement desired outcomes in your GPT-based applications. So go ahead and give it a try, and let Tzap work its magic for you!

# Getting Started with Tzap Command Line (GPT Git Commits)

Getting started with Tzap is easy! Follow the steps below to begin using Tzap in your projects.

## Installation

1. To install Tzap, simply download the `install.sh` script:

```bash
curl -sSL https://github.com/tzapio/tzap/releases/latest/download/install.sh -o install.sh
```

2. Make the script executable and run it:

```bash
chmod +x install.sh
./install.sh
```

Now, Tzap is installed on your system and ready to use!

## Usage

### Tzap GitCommit

Never write a git commit message again with Tzap! To try this feature, simply run:

```bash
tzap gitcommit
```

This command will automatically generate a meaningful git commit message based on your recent code changes.

### Dev Container

To start developing with Tzap in a dev container, follow these steps:

1. Ensure you have Docker installed on your machine.

2. Clone the Tzap repository:

```bash
git clone https://github.com/tzapio/tzap.git
```

3. Navigate to the Tzap directory:

```bash
cd tzap
```

4. Build the dev container:

```bash
docker-compose build
```

5. Run the dev container:

```bash
docker-compose up -d
```

Now you can access the Tzap development environment in the running container. Edit the source code, run tests, and build new features using Tzap!

## Start Coding

To start coding with Tzap, you can refer to the provided examples in the `examples` directory of the Tzap repository. These examples demonstrate various use cases and will help you understand how to utilize Tzap effectively.

Feel free to explore the Tzap library, experiment with its features, and create powerful GPT-based applications tailored to your requirements. Enjoy the convenience and efficiency that Tzap brings to your projects!