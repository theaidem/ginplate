
# Imports

gulp      = require 'gulp'
coffee      = require 'gulp-coffee'
sourcemaps  = require 'gulp-sourcemaps'
plumber     = require 'gulp-plumber'
gutil       = require 'gulp-util'
notify      = require 'gulp-notify'
changed     = require 'gulp-changed'
runSequence   = require 'run-sequence'
del         = require 'del'
bower       = require 'gulp-bower'

# Paths

app = 'src/'
dist = 'dist/'
src   = {
	coffee:'src/coffee/'
}

# gulp.watch bug on new files/deleted files.
# For wather bug:
# echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
gulp.task 'watch', ->
	gulp.watch "#{src.coffee}**/*.coffee", ['coffee']

# Utils

gulp.task 'clean', (cb) -> del ['output'], cb

gulp.task 'cleanDist', (cb) -> del [dist], cb

gulp.task 'bowerInstall', -> do bower

# Compilations

gulp.task 'coffee', ['clean'], ->
	gulp.src "#{src.coffee}**/*.coffee"
		.pipe plumber()
		.pipe( changed("#{app}/js/", {extension: '.js'}))
		.pipe( sourcemaps.init())
		.pipe( coffee( {bare: true, sourcemap: {inline: true}} ) )
			.on( 'error', gutil.log )
			.on( 'error', gutil.beep )
			.on( 'error', notify.onError('Error: <%= error.message %>') )
		.pipe( sourcemaps.write() )
		.pipe( gulp.dest("#{app}/js/" ) )    
		.pipe( notify({ onLast: true, message:'Coffee compile with success!' }) )

# Main tasks

gulp.task 'compile', (callback) ->
	runSequence [ 'bowerInstall', 'coffee' ], callback

gulp.task 'default', (callback) ->
	runSequence ['compile'],  ['watch'], callback
