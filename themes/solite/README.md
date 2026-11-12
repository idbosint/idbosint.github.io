# Solite

Based on the [blank theme](https://themes.gohugo.io/blank/)

## Installation

In your Hugo site `themes` directory, run:

```
git clone https://git.sr.ht/~mlthlschr/solite
```

Next, open `config.yaml`/`config.toml` in the base of the Hugo site and ensure the theme option is set to `solite`.

```
theme: "solite"
```

For more information read the official [quick start guide](https://gohugo.io/getting-started/quick-start/) of Hugo.

### JS Dependencies

To compile the `.scss`-files, the following packages are required

```
npm i postcss postcss-cli autoprefixer
```

## Usage

The About page and Privacy Policy are being generated in case you add `contents/about/_index.md` and `contents/privacy/_index.md`. Write the text into that file.

Same procedure with the front page - `contents/_index.md` contains the stuff which should be displayed there.

## Contributing

In case you want to contribute, simply write me a mail.

## License

This theme is released under the [MIT license](https://github.com/Vimux/blank/blob/master/LICENSE).
