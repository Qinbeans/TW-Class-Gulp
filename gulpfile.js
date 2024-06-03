const gulp = require('gulp');
const postcss = require('gulp-postcss');
const tailwindcss = require('tailwindcss');
const autoprefixer = require('autoprefixer');
const cssnano = require('cssnano');
const htmlparser2 = require('htmlparser2');
const fs = require('fs');

const classMap = {};
const jsonMap = {};

function extractClassesFromHtml() {
  return new Promise(async (resolve, reject) => {
    const { nanoid } = await import('nanoid');
    gulp.src('public/views/**/*.html')
      .on('data', file => {
        const parser = new htmlparser2.Parser({
          onopentag(_name, attribs) {
            if (attribs.class) {
              const classes = attribs.class.split(/\s+/).sort().join(' ');
              if (!classMap[classes]) {
                // replace any ., _, #, or : with a -
                const newClassName = "c"+nanoid(8).replace(/[\.\_\#\:\[\]]/g, '-');
                classMap[classes] = newClassName;
                jsonMap[classes.split(' ').join(' ')] = newClassName;  // Reverse key-value pair
              }
            }
          }
        }, { decodeEntities: true });

        parser.write(file.contents.toString());
        parser.end();
      })
      .on('end', resolve)
      .on('error', reject);
  });
}

gulp.task('extract-classes', function (done) {
  extractClassesFromHtml().then(() => {
    fs.writeFileSync('./classMap.json', JSON.stringify(jsonMap, null, 2));
    done();
  }).catch(err => {
    console.error(err);
    done(err);
  });
});

gulp.task('extract-classes:release', function (done) {
  extractClassesFromHtml().then(() => {
    done();
  }).catch(err => {
    console.error(err);
    done(err);
  });
});

gulp.task('css', function () {
  return gulp.src('public/static/style/app.css')
    .pipe(postcss([
      require('./scripts/postcss-class-extractor')({ classMap }),
      tailwindcss,
      autoprefixer,
      cssnano({
        preset: 'default',
      }),
    ]))
    .pipe(gulp.dest('tmp/static/style/'));
});

gulp.task('css:release', function () {
  return gulp.src('public/static/style/app.css')
    .pipe(postcss([
      tailwindcss,
      autoprefixer,
      cssnano({
        preset: 'default',
      }),
    ]))
    .pipe(gulp.dest('dist/static/style/'));
});

gulp.task('template', function () {
  // move html files from public/views to temp/views with class names replaced
  return gulp.src('public/views/**/*.html')
    .pipe(require('./scripts/html-replacer')({ classMap }))
    .pipe(gulp.dest('tmp/views/'));
});

gulp.task('template:release', function () {
  // move html files from public/views to temp/views with class names replaced
  return gulp.src('public/views/**/*.html')
    .pipe(require('./scripts/html-replacer')({ classMap }))
    .pipe(gulp.dest('dist/views/'));
});

gulp.task('statics', function () {
  // move everything except views and styles to temp
  return gulp.src(['public/**/*', '!public/views/**/*', '!public/static/styles/**/*'])
    .pipe(gulp.dest('tmp/'));
});

gulp.task('statics:release', function () {
  // move everything except views and styles to dist
  return gulp.src(['public/**/*', '!public/views/**/*', '!public/static/styles/**/*'])
    .pipe(gulp.dest('dist/'));
});

gulp.task('release', gulp.series('statics:release', 'extract-classes:release', 'css:release', 'template:release'));

gulp.task('default', gulp.series('statics', 'extract-classes', 'css', 'template'));
