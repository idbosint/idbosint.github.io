// Configuración de PostCSS que usa Hugo (toCSS | postCSS | minify).
// El tema Solite compila assets/scss/style.scss y Hugo invoca postcss-cli.
module.exports = {
  plugins: [
    require('autoprefixer')
  ]
};
