const postcss = require('postcss');

module.exports = postcss.plugin('postcss-class-extractor', (opts = {}) => {
  return (root) => {
    const classMap = opts.classMap || {};
    root.walkRules(rule => {
      const classes = rule.selector.match(/\.[^ ,]+/g);
      if (classes) {
        const className = classes.sort().join(' ');
        if (classMap[className]) {
          rule.selector = `.${classMap[className]}`;
        }
      }
    });

    // Generate the CSS with the new class mappings
    const outputCss = Object.keys(classMap).map(originalClasses => {
      return `.${classMap[originalClasses]} { @apply ${originalClasses.replace(/\./g, '')}; }`;
    }).join('\n');

    root.append(postcss.parse(outputCss));
  };
});