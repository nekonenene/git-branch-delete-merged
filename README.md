# git-branch-delete-merged

This CLI app deletes local branches that have been merged **also "[Squash and merge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges#squash-and-merge-your-commits)"**.

## Installation

```sh
go install github.com/nekonenene/git-branch-delete-merged@v1
```

## Usage

### General usage

If you want to delete a branch that has merged into the `main` branch:

```sh
git-branch-delete-merged --base-branch main
```

And if the branch to delete exists, you will get a prompt like this:

```
Target branches: [dev1]

Are you sure to delete 'dev1' branch? [y|n|l|d|q|help]:
```

Please type one and press enter.

* `y`: Yes, delete the branch
* `n`: No, skip deleting
* `l`: Show git logs of the branch
* `d`: Show the latest commit of the branch and its diff
* `q`: Quit immediately
* `h`: Show help

### Skip prompt

If you want to delete all merged branches without confirmation, `--yes` option will be useful.

```sh
git-branch-delete-merged --base-branch main --yes
```


## Thank you

[not-an-aardvark/git-delete-squashed](https://github.com/not-an-aardvark/git-delete-squashed) is the reference code that helps finding branches which has squashed and merged.
