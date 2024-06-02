# TW-Gulp-Test

This project tests the capabilities of Gulp.js in its ability to minify tailwind classes.

## Installation

1. Clone the repository
2. Install Go 1.22+, [air](https://github.com/cosmtrek/air), and Node.js 14+
3. Run `npm install` or `yarn install` or `pnpm install`
4. Run `npm run build` or `yarn build` or `pnpm build`
5. `cd dist` and `MODE=release ./backend` to start the server

## Development

1. Run `npm run dev` or `yarn dev` or `pnpm dev` or `air`
2. Open `http://localhost:3000` in your browser

Any changes made to Go, CSS, or HTML files will automatically reload the server. This will also rebuild the CSS file to minify classes.