module.exports = function(grunt) {

    grunt.loadNpmTasks('grunt-usemin');
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-clean');

    grunt.initConfig({
        useminPrepare: {
            html: "index.html",
            options: {
                dest: '.'
            }
        },
        usemin: {
            html: "index.html"
        },
        copy: {
            main: {
                files: [{
                    src: 'index_unbuilt.html',
                    dest: 'index.html'
                }]
            }
        },
        clean: ['all.min.js', 'all.min.css', '.tmp', 'index.html']
    });

    grunt.registerTask('default', [
        'copy',
        'useminPrepare',
        'concat:generated',
        'cssmin:generated',
        'uglify:generated',
        'usemin'
    ]);
};
