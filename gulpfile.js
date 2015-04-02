/**
 * Created by languid on 4/2/15.
 */
var gulp = require('gulp');
var sass = require('gulp-ruby-sass');
var react = require('gulp-react');
var uglify = require('gulp-uglify');
var mini = require('gulp-minify-css');

gulp.task('sass', function(){
    return sass('./resources/scss/base.scss')
        .pipe(mini())
        .pipe(gulp.dest('./public/css/'))
});

gulp.task('jsx', function(){
    return gulp.src('resources/js/components/**/*.jsx')
        .pipe(react())
        .pipe(gulp.dest('public/js/components/'))
});

gulp.task('compressJs', function(){
    return gulp.src('./resources/js/**/*.js')
    .pipe(uglify())
    .pipe(gulp.dest('./public/js/'));
});

gulp.task('watch', function(){
    gulp.watch('./resources/scss/**/*.scss', ['sass']);
    gulp.watch('./resources/js/**/*.js', ['compressJs']);
    gulp.watch('./resources/js/**/*.jsx', ['jsx']);
});

gulp.task('default', ['watch','sass', 'jsx']);