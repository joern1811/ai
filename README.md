# AI Command Line Tool

The AI command line tool is a utility that helps transcribe audio files and generate concise summaries of the transcribed text. It provides a simple interface to interact with OpenAI's GPT-based language model.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
    - [Ask a Question](#ask-a-question)
    - [Transcribe Audio](#transcribe-audio)
    - [Summarize Text](#summarize-text)
- [How to Install](#how-to-install)
- [Configuration](#configuration)
    - [Getting an OpenAI API Key](#getting-an-openai-api-key)
    - [Storing the API Key](#storing-the-api-key)
- [License](#license)

## Installation

You can easily install the AI tool using Homebrew:

```bash
brew tap domjoe1811/ai
brew install ai
```
## Usage

The AI command line tool supports the following commands:
### Ask a Question

Use the speech ask command to ask a question to ChatGPT and receive a response.
    
```bash
ai speech ask "What is the meaning of life?"
   ```

### Transcribe Audio

The speech transcript command transcribes an audio file and outputs the transcribed text.

```bash 
ai speech transcript -i input_audio.wav -o output_transcript.txt
```

### Summarize Text

The speech summarize command can be used to transcribe an audio file or read text from a file and generate a summary of the content.

```bash
ai speech summarize -i input_audio.wav -o output_summary.txt
```
or
```bash
ai speech summarize -i input_text.txt -o output_summary.txt
```
## How to Install

To install the AI tool via Homebrew, follow these steps:

1. Open your terminal.

2. Tap into the repository containing the AI tool:

```bash
brew tap domjoe1811/ai
```
3. Install the tool:
    
```bash
brew install ai
```
## Configuration

Before you can use the AI tool, you need to sign up for an OpenAI account and generate an API key. 

### Getting an OpenAI API Key:

1. Visit the OpenAI website at [https://platform.openai.com/account/api-keys](https://platform.openai.com/account/api-keys).
2. If you don't have an account, click on "Sign Up" to create one. If you do, click "Log In".
3. Once logged in, navigate to your API key management page.
4. Click on "Create new secret key".
5. Enter a name for your new key, then click "Create secret key".
6. Your new API key will be displayed. Use this key to interact with the OpenAI API.

**Note:** Your API key is sensitive information. Do not share it with anyone.

### Storing the API Key:

Once you have the API key, create a configuration file in JSON format and store the API key in it. The configuration file must be named ".ai" and placed in the home directory, for example in $HOME/.ai.

```json
{
  "openAIConfig": {
    "openAIAuthToken": "[api-key]"
  },
  "promptConfig": {
    "summarizePrompt": "Fasse die folgende Nachricht stichpunktartig zusammen. Versuche dabei, die Themen der Nachricht zu erkennen und diese bei der Zusammenfassung zu gruppieren."
  }
}
```

Replace [api-key] with your actual OpenAI API key.

## Costs

Please note that the usage of the OpenAI API may incur costs after the initial credits have been exhausted. It's important to monitor your usage and review the pricing information provided by OpenAI to understand the associated costs.

## License

This project is licensed under the terms of the Apache License 2.0 license. See the [LICENSE](LICENSE) file for details.

