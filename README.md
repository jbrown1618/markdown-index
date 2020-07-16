# markdown-index

A CLI for creating an index for a directory full of markdown files

## Installation

```
go get github.com/jbrown1618/markdown-index
```

## Usage

```
markdown-index --root /path/to/markdown/directory --browser
```

## Example

Given the file tree:

```
root/
|
+- Ideas
|  |
|  +- good-idea.md
|  +- bad-idea.md
|
+- Recipes
   |
   +- good-recipe.md
   +- bad-recipe.md
```

Running `markdown-index --root /path/to/root` will generate a file `/path/to/root/index.md` with contents:

```markdown
# Root
## Ideas
- [My Good Idea](./Ideas/good-idea.md)
- [My Bad Idea](./Ideas/bad-idea.md)
## Recipes
- [My Good Recipe](./Ideas/good-recipe.md)
- [My Bad Recipe](./Ideas/bad-recipe.md)
```