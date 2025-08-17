# Github Actions

This repository uses github actions for CI, and publishing artifacts.

## Conventions

The main conventions followed for github actions are outlined in this section.

The name of each workflow is typically prefixed with some kind of group like `GROUP: ...`. This is just a nicety to help quickly scan through the actions when using the Github UI, as it's not always easy at a glance to know the underlying point of a workflow. For each workflow, there is usually a run-name that is almost identical the naem of the workflow but typically adds a little more info around the context in which the job is running. Typically anything that is running in the context of a PR will have `#{pr number}` to help identify it as a job for a given PR.

The groups are:

| Group | Purpose |
|-------|---------|
| **CI** | Anything thing happens to ensure that the code is in a mergeable state. |
| **RELEASE** | Processes that take place to publish a new version of the application. |
| **DOCS** | Processes responsible for deploying & validating documentation. |
| **_INC** | In some cases there are workflows that are never directly triggered, but are there only to be used by other workflows. Within the repository any file that is intended to be an "includable workflow" and not triggered by any typical event, is prefixed with an `_` e.g. `_publish_release_for_tag.yaml`. |


## Workflows

This section gives a little more info on each workflow.

### `commit_lint.yaml`

**Must pass in order to merge**

This workflow uses [commitlintjs](https://commitlint.js.org/) to verify that all the commits in the PR are following the conventional commit standards. As detailed in [Release process](./release_process.md) it is crucial to the release process that this standard is followed, and this job will fail if any commit does not conform. The `commitlint.config.js`
file defines exactly what rules are in place.

*Trigger: [pull_request](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#pull_request)*

### `ci_build.yaml`

**Must pass in order to merge**

This runs a snapshot build of the application using goreleaser. A snaopshot is used to test that the application can be built correctly, without publishing a release anywhere. This is a basic test to ensure that code compiles.

*Trigger: [pull_request](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#pull_request)*

### `create_next_release.yaml`

This workflow is the entrypoint of the release process, and will firstly use semantic-release to calculate what next version of the app should be. Once it has calculated this it will create a tag on the repository accordingly (although we use a `v` prefix for the tags). Once the tag has been created, it uses the [`_publish_release_for_tag.yaml`](#_publish_release_for_tagyaml) workflow which generates the release.

Depending on the included commits, we may not always need a release. This logic is handled in this workflow, mostly by semantic-release, if it is determined that the commits don't require a new release, then there will be no new tag or release created.

Their is a concurrency rule on this workflow that prevents it running in parallel for more than one commit on the branch; this is to ensure that the versions are incremented correctly. Due to the way semantic-release works (it fetches tags to calculate the next), it could be very unpredictable if multiple instances were to run at the same time.

*Trigger: [push](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#push) on main*

### `_publish_release_for_tag.yaml`

| Input | Required | Detail |
|-------|----------|--------|
| `tag` | Yes | This specifies the tag that a release should be created for. |

This workflow is used by [`create_next_release.yaml`](#create_next_releaseyaml), once it knows the version to create. It will pass the new tag into this workflow, and run goreleaser to generate a new release in Github. Goreleaser attaches the executables for various architectures to the release for distribution. 

*Trigger: [workflow_call](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#workflow_call)*

### `deploy_github_pages.yaml`

This workflow renders and publishes the documentation using mkdocs. It pushes the rendered files to the `gh-pages` branch in this repo, and from there github pages will take over and deploy the contents of `gh-pages` branch.

*Trigger: [push](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#push) on main, but only if it includes changes to documentation*

### `static_analysis.yaml`

**Must pass in order to merge**

This workflow runs a series of static checks across the full codebase. As it stands this just checks that the code conforms to the go formatting standard. More checks will be added in the future.

*Trigger: [pull_request](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#pull_request)*

### [Github managed] `pages-build-deployment`

This workflow isn't contained within the repository, but is managed by Github when we enable the github pages deployment. It will run whenever the `gh-pages` branch is updated.

