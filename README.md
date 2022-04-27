# snark

`snark` is a project designed to let terminal warriors quickly generate snarky emoticon faces and sarcastic fong/case for strings.

## Usage

### Printing Emoticon Faces

#### To see all available faces, run:

```
snark list
```

``` title="Sample Output"
lenny   =>      ( ͡° ͜ʖ ͡°)
shrug   =>      ¯\_(ツ)_/¯
cat     =>      (•ㅅ•)
shock   =>      ಠ_ಠ
```

#### To print a single face, run the `print` command, and provide it with the name of one of the entries in the list:

```
snark print shrug
```

```
¯\_(ツ)_/¯
```

The specific face will also be immediately copied to your clipboard.

### Sarcastic Casing

Simply provide a quoted string to get a "sarcastic" version of that string printed out.

The "sarcastic" version of the string will also be immediately copied to your clipboard.

## Building from source

With the appropriate Go environment in place, the Makefile provides the necessary commands for build, run, etc.