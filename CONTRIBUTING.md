# Contributing to RISKEN

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to RISKEN. These are just guidelines, not rules, use your best judgment and feel free to propose changes to this document in a pull request.

## [Issues](https://github.com/ca-risken/community/issues)

Issues are created [here](https://github.com/ca-risken/community/issues/new/choose). and choose following issue types.

* `Bug report`
* `Feature request`

### Languages

We accept issues in *Japanese or English* language.
Responses to posted issues may or may not be in the original language.

## Pull Requests

Pull Requests are the way concrete changes are made to the code, documentation,
dependencies, and tools contained in this repository.

* [Setting up your local environment](https://docs.security-hub.jp/admin/infra_local/)
  * Step 1(Fork): Fork the project on GitHub and clone your fork locally.
  * Step 2(Build): Build steps and dependencies differ slightly depending on your operating system. see the [README](README.md).
  * Step 3(Branch): To keep your development environment organized, create local branches to hold your work. These should be branched directly off of the `master` branch.
* **Making Changes**
  * Step 4(Code): Please be sure to run `make lint` from time to time on any code changes to ensure that they follow the project's code style.
  * Step 5(Commit): It is recommended to keep your changes grouped logically within individual commits. Many contributors find it easier to review changes that are split across multiple commits. There is no limit to the number of commits in a pull request.
  * Step 6(Rebase): Once you have committed your changes, it is a good idea to use `git rebase` (not `git merge`) to synchronize your work with the main repository. 
    ```bash
    $ git fetch upstream
    $ git rebase upstream/main
    ```
  * Step 7(Test): Please be sure to run `make go-test`, After fixing a bug or adding a feature, you should always make sure that the test passes.
  * Step 8(Push): Once your commits are ready to go -- with passing tests and linting -- begin the process of opening a pull request by pushing your working branch to your fork on GitHub.
    ```bash
    $ git push origin my-branch
    ```
  * Step 9(Pull Request): Open a Pull-Request on this GitHub repository.
  * Step 10(Discuss and Update): You will probably get feedback or requests for changes to your pull request. This is a necessary part of the process in order to evaluate whether the changes are correct and necessary.
    * Approval and Request Changes Workflow: All pull requests require approval from a Code Owner of the area you modified in order to land. Whenever a maintainer reviews a pull request they may request changes.
  * Step 11(Landing): In order to land, a pull request needs to be reviewed and approved by at least one RISKEN Code Owner and pass CI. After that, if there are no objections from other contributors, the pull request can be merged.
  * Continuous Integration Testing: Every pull request is tested on the Continuous Integration (CI) system to confirm that it works on RISKEN's supported platforms. CI starts automatically when you open a pull request, but only core maintainers can restart a CI run. If you believe CI is giving a false negative, ask a maintainer to restart the tests.


## Code Style

RISKEN enforces Linting to improve the efficiency of code review flow. Please be sure to `make lint` before opening a pull-request.
