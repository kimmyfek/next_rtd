'use strict'

var deps = require('./dependencies');

var gulp = require('gulp'),
    browserify = require('browserify'),
	babelify = require('babelify'),
    buffer = require('gulp-buffer'),
	source = require('vinyl-source-stream'),
    bust = require('gulp-buster'),
    clean = require('gulp-clean'),
    changed = require('gulp-changed'),
    csslint = require('gulp-csslint'),
    eslint = require('gulp-eslint'),
    size = require('gulp-size'),
    concatCss = require('gulp-concat-css');

// Input file.
var bundler = browserify('src/jsx/app.jsx', {
    extensions: ['.js', '.jsx'],
    debug: true
});

// Babel transform
bundler.transform(babelify.configure({
    sourceMapRelative: 'src',
    presets: ["es2015", "react"]
}));

// On updates recompile
bundler.on('update', bundle);

function bundle() {
    return bundler.bundle()
        .on('error', function (err) {
            console.log("=====");
            console.error(err.toString());
            console.log("=====");
            this.emit("end");
        })
    ;
}


gulp.task('concatCss', function () {
  return gulp.src('next/static/css/*.css')
    .pipe(concatCss("bundle.css"))
    .pipe(gulp.dest('next/static/css/'));
});

gulp.task('transformMain', function() {
    return browserify({entries: './next/static/scripts/jsx/app.js', extensions: ['.js'], debug: true})
        .transform(babelify.configure({
			presets: ["env", "react"]
		 }))
        .bundle()
        .pipe(source('./app.js'))
        .pipe(gulp.dest('./next/static/scripts/js'));
});

gulp.task('clean', function() {
  return gulp.src(['./next/static/scripts/js'], {read: false}).pipe(clean());
});

gulp.task('default', ['clean'], function() {
  gulp.start('concat');
  gulp.start('transformMain');
  gulp.start('eslint');
  gulp.start('concatCss');
  gulp.watch('./next/static/css/*.css', ['concat', 'csslint']);
  gulp.watch(['./next/static/scripts/jsx/*.js', './next/static/scripts/jsx/**/*.js'], ['transformMain', 'eslint']);
});

gulp.task('eslint', function () {
//  return gulp.src(['./next/static/scripts/jsx/*.js'])
//    .pipe(eslint())
//    .pipe(eslint.format())
//    .pipe(eslint.failAfterError());
});

gulp.task('csslint', function() {
  gulp.src('./next/static/css/main.css')
    .pipe(csslint('csslintrc.json'))
    .pipe(csslint.formatter('fail'));
});

gulp.task('concat', function() {
  var concat = require('gulp-concat');

  gulp.src(deps.js)
    .pipe(concat('scripts.js'))
    .pipe(gulp.dest('./next/static/scripts/js'))
    .pipe(bust())
    .pipe(gulp.dest('./next/static/scripts/js'));
});
