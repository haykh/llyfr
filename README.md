# `Llyfr`

Fast rofi-like launcher for book/article collections specifically targeting scientific literature. Written in `Go`/`Wails` + `React`/`TypeScript`.

## Usage

```
llyfr <PATH_TO_BIB> <PATH_TO_PDFs>

# e.g.
llyfr ~/Documents/papers/refs.bib ~/Documents/papers

# or equivalently (assumes PDFs are located in the same directory)
llyfr ~/Documents/papers/refs.bib
```

The `.bib` file must contain the following mandatory keys:

- `year`
- `author`: author list in standard latex bib form (separated with "and"-s)
- `title`
- `file`: PDF file names (will search them in `PATH_TO_PDF`)

Type is determined from `@Article`, `@Book`, etc prefixes. When type is `@Article` it must also contain `journal` entry. 

## Features

- Parsing of ".bib" file
- Fuzzy search by year, author, title
- Opening PDF in external reader
- Keyboard navigation
- Standardization of journal names, author list and titles (removal of latex components, hyphens, etc.)



![img](./demo/demo.png)

## Building from source

You'll need `wails`, `npm`, and `go`. From the root simply run:

```
wails build -clean -o llyfr
```

The executable is then produced in `build/bin`.

# To do

- [ ] Binaries for major platforms
- [ ] nixpkg derivation
