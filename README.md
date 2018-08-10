# jsonf

`jsonf` is a simple tool to format JSON files.

## Installation

Installation requires [golang](https://golang.org/).

	go get github.com/cbednarski/jsonf

Run `jsonf` or `jsonf -h` for usage info.

## Usage Examples

```
$ find . -name *.json
./test-fixtures/dogs/corgi.json
./test-fixtures/dogs/dachshund.json
./test-fixtures/icecream.json
$ cat test-fixtures/icecream.json
{
	"icecream" : [
"chocolate",
"strawberry",
"vanilla"
]
}

```

Use `jsonf filename` to reformat a file (output to stdout with two spaces
indentation by default).

```
$ jsonf test-fixtures/icecream.json
{
  "icecream": [
    "chocolate",
    "strawberry",
    "vanilla"
  ]
}
```

Pass `-i` to change indentation. You can use `t` for tab, `ss` for two spaces, `ssss` for four spaces, etc.

```
$ jsonf -i t test-fixtures/icecream.json
{
	"icecream": [
		"chocolate",
		"strawberry",
		"vanilla"
	]
}
```

Pass `-c` to compact (minify) the JSON file.

```
$ jsonf -c test-fixtures/icecream.json
{"icecream":["chocolate","strawberry","vanilla"]}
```

Pass `-w` to rewrite the file in place.

```
$ jsonf -w -i t test-fixtures/icecream.json
$ cat test-fixtures/icecream.json
{
	"icecream": [
		"chocolate",
		"strawberry",
		"vanilla"
	]
}
```

You can also apply formatting to an entire directory at once. Pass `-r` to
recurse into subdirectories.

```
$ jsonf -r test-fixtures/
test-fixtures/dogs/corgi.json
{
  "height": "short",
  "ears": "pointy",
  "tail": "fluffy"
}
test-fixtures/dogs/dachshund.json
{
  "height": "short",
  "ears": "floppy",
  "tail": "pointy"
}
test-fixtures/icecream.json
{
  "icecream": [
    "chocolate",
    "strawberry",
    "vanilla"
  ]
}
```
