const fs = require('fs');
const path = require('path');

const stylesDir = path.resolve(__dirname, 'node_modules/highlight.js/styles');
const jsOutputFile = path.resolve(__dirname, 'themeStyles.js');
const goOutputFile = path.resolve(__dirname, 'theme_list.go');
const colorInfoOutputFile = path.resolve(__dirname, 'themeColors.js');

// Read all CSS files from the styles directory
let themes = fs.readdirSync(stylesDir).filter(file => file.endsWith('.css'));

// Prioritize .min.css files
const minifiedFiles = new Set(themes.filter(file => file.endsWith('.min.css')).map(file => file.replace('.min.css', '')));
themes = themes.filter(file => {
  const baseName = file.replace('.css', '').replace('.min', '');
  // Skip unminified versions if corresponding .min.css file exists
  return !minifiedFiles.has(baseName) || file.endsWith('.min.css');
});

// Group themes and classify by naming conventions
const themeMap = {};
let themeList = [];
const themeColors = [];
let defaultDarkTheme = null;
let defaultLightTheme = null;

const cssColorNames = {
  black: '#000000',
  white: '#ffffff',
  navy: '#000080',
  // Add more color names as needed if the background color in css is not defined in standard method 
};

function convertColorNameToHex(colorName) {
  return cssColorNames[colorName.toLowerCase()] || null;
}

function normalizeHexColor(hexColor) {
  hexColor = hexColor.startsWith('#') ? hexColor.slice(1) : hexColor;
  if (hexColor.length === 3) {
    hexColor = hexColor.split('').map(char => char + char).join('');
  }
  return hexColor;
}

function isDarkTheme(color) {
  if (!color.startsWith('#') && !color.startsWith('rgb')) {
    const hexColor = convertColorNameToHex(color);
    if (hexColor) {
      color = hexColor;
    } else {
      return false;
    }
  }
  const hexColor = normalizeHexColor(color);
  const rgb = parseInt(hexColor, 16);
  const r = (rgb >> 16) & 0xff;
  const g = (rgb >> 8) & 0xff;
  const b = (rgb >> 0) & 0xff;
  const brightness = (r * 299 + g * 587 + b * 114) / 1000;
  return brightness < 128;
}

themes.forEach(file => {
  const themeName = file.replace('.css', '').replace('.min', '');
  const [base, ...variantParts] = themeName.split('-');
  const variant = variantParts.join('-');

  if (!themeMap[base]) {
    themeMap[base] = {};
  }

  let isDark = false;
  let backgroundColor = null;

  if (variant.includes('light')) {
    if (!themeMap[base].light || themeMap[base].light.length > file.length) {
      themeMap[base].light = `() => import('highlight.js/styles/${file}?inline')`;
      if (!defaultLightTheme) {
        defaultLightTheme = themeMap[base].light;
      }
    }
  } else if (variant.includes('dark')) {
    if (!themeMap[base].dark || themeMap[base].dark.length > file.length) {
      themeMap[base].dark = `() => import('highlight.js/styles/${file}?inline')`;
      if (!defaultDarkTheme) {
        defaultDarkTheme = themeMap[base].dark;
      }
    }
  } else {
    const cssContent = fs.readFileSync(path.resolve(stylesDir, file), 'utf-8');
    const backgroundMatch = cssContent.match(/\.hljs\s*{[^}]*?\s*background(?:-color)?:\s*(#[0-9a-fA-F]{3,6}|rgb\([^)]+\)|[a-zA-Z]+|url\([^)]+\))/i);
    backgroundColor = backgroundMatch ? backgroundMatch[1].trim() : null;

    if (backgroundColor) {
      if (backgroundColor.startsWith('url')) {
        backgroundColor = null;
      } else if (backgroundColor.startsWith('#')) {
        isDark = isDarkTheme(backgroundColor);
      } else if (backgroundColor.startsWith('rgb')) {
        const rgbValues = backgroundColor.match(/\d+/g).map(Number);
        const brightness = (rgbValues[0] * 299 + rgbValues[1] * 587 + rgbValues[2] * 114) / 1000;
        isDark = brightness < 128;
      } else {
        isDark = isDarkTheme(backgroundColor);
      }
    }

    if (isDark) {
      if (!themeMap[base].dark || themeMap[base].dark.length > file.length) {
        themeMap[base].dark = `() => import('highlight.js/styles/${file}?inline')`;
        if (!defaultDarkTheme) {
          defaultDarkTheme = themeMap[base].dark;
        }
      }
    } else {
      if (!themeMap[base].light || themeMap[base].light.length > file.length) {
        themeMap[base].light = `() => import('highlight.js/styles/${file}?inline')`;
        if (!defaultLightTheme) {
          defaultLightTheme = themeMap[base].light;
        }
      }
    }
  }

  if (!themeList.includes(base)) {
    themeList.push(base);
  }

  if (backgroundColor) {
    themeColors.push({
      theme: base,
      variant: isDark ? 'dark' : 'light',
      color: backgroundColor
    });
  }
});

themeList = themeList.map(base => {
  if (themeMap[base].light && !themeMap[base].dark) {
    return `${base}-light`;
  } else if (!themeMap[base].light && themeMap[base].dark) {
    return `${base}-dark`;
  } else if (themeMap[base].light && themeMap[base].dark) {
    return `${base}-all`;
  } else {
    return base;
  }
});

Object.keys(themeMap).forEach(base => {
  if (!themeMap[base].dark && defaultDarkTheme) {
    themeMap[base].dark = defaultDarkTheme;
  }
  if (!themeMap[base].light && defaultLightTheme) {
    themeMap[base].light = defaultLightTheme;
  }
});

const jsOutput = `export const themeStyles = {\n${Object.entries(themeMap)
  .map(([theme, variants]) => 
    `  ${JSON.stringify(theme)}: {\n    light: ${variants.light},\n    dark: ${variants.dark}\n  }`
  ).join(',\n')}\n};`;

fs.writeFileSync(jsOutputFile, jsOutput);

const goOutput = `
package render_markdown_codehighlight

var ThemeList = []string{
${themeList.map(theme => `"${theme}"`).join(",\n  ")},
}
`;

fs.writeFileSync(goOutputFile, goOutput);

console.log('Theme styles, Go theme list, and color information generated successfully!');
