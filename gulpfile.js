'use strict'

var deps = require('./dependencies');

var gulp = require('gulp'),
    browserify = require('gulp-browserify'),
    buffer = require('gulp-buffer'),
    bust = require('gulp-buster'),
    clean = require('gulp-clean'),
    changed = require('gulp-changed'),
    csslint = require('gulp-csslint'),
    eslint = require('gulp-eslint'),
    size = require('gulp-size');

gulp.task('transformMain', function() {
  return gulp.src('./next_rtd/static/scripts/jsx/*.js')
    .pipe(changed('./next_rtd/static/scripts/js'))
    .pipe(browserify({transform: ['reactify']}))
    .pipe(gulp.dest('./next_rtd/static/scripts/js'))
    .pipe(buffer())
    .pipe(bust())
    .pipe(gulp.dest('./next_rtd/static/scripts/js'))
    .pipe(size());
});

gulp.task('clean', function() {
  return gulp.src(['./next_rtd/static/scripts/js'], {read: false}).pipe(clean());
});

gulp.task('default', ['clean'], function() {
  gulp.start('copy');
  gulp.start('concat');
  gulp.start('transformMain');
  gulp.start('eslint');
  gulp.watch('./next_rtd/static/css/*.css', ['concat', 'csslint']);
  gulp.watch(['./next_rtd/static/scripts/jsx/*.js', './next_rtd/static/scripts/jsx/**/*.js'], ['transformMain', 'eslint']);
});

gulp.task('eslint', function () {
//  return gulp.src(['./next_rtd/static/scripts/jsx/*.js'])
//    .pipe(eslint())
//    .pipe(eslint.format())
//    .pipe(eslint.failAfterError());
});

gulp.task('csslint', function() {
  gulp.src('./next_rtd/static/css/main.css')
    .pipe(csslint('csslintrc.json'))
    .pipe(csslint.reporter('fail'));
});

gulp.task('copy', function() {
  //  gulp.src(deps.fonts)
  //    .pipe(gulp.dest('static/fonts'));
});

gulp.task('concat', function() {
  var concat = require('gulp-concat');

  gulp.src(deps.js)
    .pipe(concat('scripts.js'))
    .pipe(gulp.dest('./next_rtd/static/scripts/js'))
    .pipe(bust())
    .pipe(gulp.dest('./next_rtd/static/scripts/js'));

  //gulp.src(deps.css)
  //  .pipe(concat('styles.css'))
  //  .pipe(gulp.dest('./next_rtd/static/css'))
  //  .pipe(bust())
  //  .pipe(gulp.dest('./next_rtd/static/css'))
});