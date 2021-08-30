# Carthage

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-carthage?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-carthage/releases)

Runs the selected Carthage command.

<details>
<summary>Description</summary>

Downloads and builds the dependencies your Cocoa application needs.

### Configuring the Step
1. Add the **Carthage** Step in your Workflow.
2. Select `bootstrap` Carthage command in the **Carthage command to run** input. Make sure you have the **Bitrise.io Cache:Pull** Step before and the **Bitrise.io Cache:Push** Step after the **Carthage** Step in your Workflow to cache files and speed up your Bitrise build.
2. Provide your GitHub credentials in the **GitHub Personal Access Token** input to avoid GitHub rate limit issues. Don't worry, your credentials are safe with us since we store them encrypted and do not print them out in build logs.
3. Optionally, you can provide any extra flag for the Carthage command you wish to run in the **Additional options for Carthage command** input.
5. To get more information printed out, set the **Enable verbose logging** to `yes`.

### Troubleshooting
It is important that you use `bootstrap` Carthage command, as this is the only command that can leverage the cache! If you run, for example, the `update` command, it won't generate the required cache information, because the `update` command will disregard the available files or the cache.

### Useful links
- [Official Carthage documentation](https://github.com/Carthage/Carthage)
- [About Secrets and Env Vars ](https://devcenter.bitrise.io/builds/env-vars-secret-env-vars/)

### Related Steps
- [Bitrise.io Cache Push](https://www.bitrise.io/integrations/steps/cache-push)
- [Bitrise.io Cache Pull](https://www.bitrise.io/integrations/steps/cache-pull)
- [iOS Auto Provision](https://www.bitrise.io/integrations/steps/ios-auto-provision)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `carthage_command` | Select a command to set up your dependencies.  The step will cache your dependencies only when using `bootstrap` in this input and you have `cache-pull` and `cache-push` steps in your workflow.  To see available commands run: `carthage help` on your local machine. | required | `bootstrap` |
| `carthage_options` | Options added to the end of the Carthage call. You can use multiple options, separated by a space character.  To see available command's options, call `carthage help COMMAND`   Format example: `--platform ios` |  |  |
| `github_access_token` | Use this input to avoid Github rate limit issues.  See the github's guide: [Creating an access token for command-line use](https://help.github.com/articles/creating-an-access-token-for-command-line-use/),    how to create Personal Access Token.  __UNCHECK EVERY SCOPE BOX__ when creating this token. There is no reason this token needs access to private information. | sensitive | `$GITHUB_ACCESS_TOKEN` |
| `xcconfig` | Use this input to provide an `xcconfig` file as a workaround for the Xcode 12 issue. For more information, see [the Github issue](https://github.com/Carthage/Carthage/issues/3019).  Can either be a local file provided with the `file://` scheme (like `file://path/to/file.xcconfig`) or an URL (like https://domain.com/file.xconfig). |  |  |
| `verbose_log` | Enable verbose logging? | required | `no` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-carthage/pulls) and [issues](https://github.com/bitrise-steplib/steps-carthage/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
