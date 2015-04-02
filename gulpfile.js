/**
 * Created by languid on 4/2/15.
 */
var gulp = require('gulp');
var sass = require('gulp-ruby-sass'),
    react = require('gulp-react'),
    mini = require('gulp-minify-css');

gulp.task('sass', function(){
    return sass('resources/scss/base.scss')
        .pipe(mini())
        .pipe(gulp.dest('public/css/'))
});

gulp.task('jsx', function(){
    return gulp.src('./js/jsx/**/*.jsx')
        .pipe(react())
        .pipe(gulp.dest('./js/components/'))
});

gulp.task('watch', function(){
    gulp.watch('./scss/**/*.scss', ['sass']);
   //gulp.watch('./js/jsx/**/*.jsx', ['jsx']);
});

gulp.task('default', ['watch', 'jsx']);