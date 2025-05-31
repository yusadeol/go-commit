# Commit - Automate your Conventional Commit messages with AI

Commit is a Golang CLI tool that automates the generation of Conventional Commit messages using AI.
By analyzing your `staged changes`, it generates clear and structured commit messages
that strictly follow the Conventional Commits standard.

## Official Documentation

##### Download Binary:

To download the latest version of the CLI, go to the [Releases](https://github.com/yusadeol/go-commit/releases) page and download the appropriate binary for your system.

##### Initialize Configuration

To initialize the necessary configuration files:

```shell
commit init
```

This creates the config file at `~/.config/commit.json`.
All customizable settings, such as default AI provider, language preferences, and API keys, are managed in this file.

#### Main functionality

##### Generate a Commit Message

To automatically generate a Conventional Commit message based on the current `staged changes`:

```shell
commit generate
```

By default, this will generate the commit message and immediately create the commit using Git.
If you only want to preview the generated message without committing anything, use the `--commit=false` option:

```shell
commit generate --commit=false
```

##### Using a Custom Diff

You can provide a custom diff instead of using the automatically detected staged changes:

```shell
commit generate YOUR_CUSTOM_DIFF
```

##### Using options to customize the commit

Use the `--provider` option to select the AI provider:

```shell
commit generate --provider=openai
```

Use the `--language` option to specify the language for the commit message:

```shell
commit generate --language=pt_BR
```

Default languages:

- `en_US` for English (United States)
- `pt_BR` for Portuguese (Brazil)
- `es_ES` for Spanish (Spain)

## License

Commit is open-sourced software licensed under the [MIT license](LICENSE.md).
