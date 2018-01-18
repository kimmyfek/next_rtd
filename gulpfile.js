'use strict'

var gulp = require('gulp'),
    browserify = require('browserify'),
	babelify = require('babelify'),
    buffer = require('vinyl-buffer'),
	source = require('vinyl-source-stream'),
    clean = require('gulp-clean'),
    csslint = require('gulp-csslint'),
    concatCss = require('gulp-concat-css'),
    uglify = require('gulp-uglify');

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
    .pipe(gulp.dest('next/static/dist/'));
});

gulp.task('transformMain', function() {
    return browserify({entries: './next/static/scripts/jsx/app.js', extensions: ['.js'], debug: true})
        .transform(babelify.configure({
			presets: ["env", "react"]
		 }))
        .bundle()
        .pipe(source('./app.js'))
        .pipe(buffer())
		.pipe(uglify())
        .pipe(gulp.dest('./next/static/scripts/js'));
});

gulp.task('clean', function() {
  return gulp.src(['./next/static/scripts/js'], {read: false}).pipe(clean());
});

gulp.task('default', ['clean'], function() {
  gulp.start('transformMain');
  gulp.start('concatCss');
  gulp.watch('./next/static/css/*.css', ['csslint']);
  gulp.watch(['./next/static/scripts/jsx/*.js', './next/static/scripts/jsx/**/*.js'], ['transformMain']);
});

gulp.task('csslint', function() {
  gulp.src('./next/static/css/main.css')
    .pipe(csslint('csslintrc.json'))
    .pipe(csslint.formatter('fail'));
});
