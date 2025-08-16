# Contributing guidelines

This repository uses [Github Flow](https://docs.github.com/en/get-started/using-github/github-flow) branching strategy; feature branches come off of main, should be short-lived, and then are merged back into main. The configuration of the repository is configured to allow only rebase merging, so any commits you add on your branch will be placed on top of main once it's merged. The idea is to keep the commit history clean and informative. To aid this it is following the [conventional commit standard](https://www.conventionalcommits.org/en/v1.0.0/), which is enforced via a CI check.

## Commit standards

As mentioned above we use conventional commits. In order to ensure this there is a CI workflow that runs on each PR, whichc checks that every commit on the head branch (that isn't on the base) follows this convention. To do this we use [commitlintjs](https://commitlint.js.org/), and you'll find the exact configuration for this check in `commitlint.config.js`.

### Running the commit lint

You can easily run the commit lint script locally for the latest commit with `npx commitlint --latest` from the root of the repository.

## Documentation

The documentation is built with [mkdocs](https://www.mkdocs.org/), and more specifically we use [material](https://squidfunk.github.io/mkdocs-material/) to add few nice features. There are 2 key parts to this, the `mkdocs.yml` config in the root of this repository, and the `docs` folder. The `mkdocs.yml` config defins things like plugins and directories that mkdocs uses to build static html, and the `docs` folder contains all the source markdown files.

The usage of some of the features in mkdocs will result in the docs not rendering particularly well inside github itself, but the github pages site for this repository should present it all nicely. 

### Working locally

If you're writing docs locally and wanna see what the rendered result looks like you can use the mkdocs server in watch mode. This is all hooked up in the `docker-compose.yml` file. Simply run `make dc-up` in the root of this repository, and you'll be able to see the rendered docs on [http://localhost:9898](http://localhost:9898). Every time either the `mkdocs.yml` config or a file in the `docs` directory changes, it will be re-rendered and any browsers will reload the page.