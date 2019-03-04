const {
  override,
  fixBabelImports
} = require("customize-cra");


module.exports = override(
  fixBabelImports("antd-css", {
      libraryName: "antd", style: "css" // import css on demand
    })
);