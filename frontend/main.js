/**
 * The MIT License (MIT)
 * 
 * Copyright (c) 2015 Christian Speckner
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * 
 */

angular.module('toggles', ['ui.bootstrap'])

    .factory("api-socket", function() {

        return io('', {
            path: '/api/socket.io'
        });

    })

    .factory('utility', function() {
        return {
            parseDurationToMsec: function(duration) {
                var msec = 0,
                    multiplicators = {
                        ms: 1,
                        s: 1000,
                        m: 60000,
                        h: 3600000
                    },
                    matches;

                while (matches = duration.match(/^\s*(\d+)\s*(s|m|h|ms|us|ns)/)) {
                    msec += (parseInt(matches[1], 10) * (multiplicators[matches[2]] || 0));

                    duration = duration.replace(/^\s*\d+\s*(s|m|h|ms|us|ns)/, '');
                }

                return msec;
            }
        };
    })

    .factory('loader', ['$modal', function($modal) {
        var loader,
            locks = 0;

        return {
            show: function() {
                if (locks++ === 0) {
                    loader = $modal.open({
                        template: " ",
                        backdrop: "static",
                        windowClass: "loading-modal",
                        backdropClass: "loading-modal-backdrop",
                        keyboard: false
                    });
                }
            },

            hide: function() {
                if (--locks === 0) {
                        loader.close();
                }
            }
        };
    }])

    .controller('switchController', ['$http', '$scope', 'api-socket', 'utility', 'loader',
        function($http, $scope, socket, utility, loader)
    {
        var switchesById = {},
            intervalHandle = null,
            queuedUpdates = [];

        function createSwitch(dataset) {
            dataset.remaining = function() {
                var duration = moment.duration(this.millisecondsRemaining > 1000 ? this.millisecondsRemaining : 1000),
                    seconds = duration.seconds(),
                    minutes = duration.minutes(),
                    hours = duration.hours(),
                    days = duration.days();

                var result = '';

                if (days > 0) result += (days + 'd ');
                if (hours > 0) result += (hours + 'h ');
                if (minutes > 0) result += (minutes + 'm ');
                if (seconds > 0) result += (seconds + 's ');

                return result.replace(/\s* $/, '');
            };

            if (dataset.type === 'transient') dataset.parsedTimeout = utility.parseDurationToMsec(dataset.timeout);

            dataset.updateTimeoutTimestamp = Date.now();

            return dataset;
        }

        function updateSwitch(dataset) {
            swtch = switchesById[dataset.id];

            if (!swtch || swtch.generation > dataset.generation) return;

            for (var key in dataset) {
                if (!dataset.hasOwnProperty(key)) continue;

                swtch[key] = dataset[key];
            }

            if (swtch.type === 'transient') swtch.parsedTimeout = utility.parseDurationToMsec(swtch.timeout);

            swtch.updateTimeoutTimestamp = Date.now();
        }

        function intervalHandler() {
            $scope.$apply(function() {
                $scope.switches.forEach(function(swtch) {
                    if (swtch.type === 'transient' && swtch.state !== swtch.groundState) {
                        var now = Date.now();

                        swtch.millisecondsRemaining -= (now - swtch.updateTimeoutTimestamp);
                        swtch.updateTimeoutTimestamp = now;
                    }
                });
            });
        }

        function stopInterval() {
            if (intervalHandle !== null) {
                clearInterval(intervalHandle);
                intervalHandle = null;
            }
        }

        function startOrStopInterval() {
            var required = false;

            if ($scope.switches) {
                $scope.switches.forEach(function(swtch) {
                    if (swtch.type === 'transient' && swtch.state !== swtch.groundState) required = true;
                });
            }

            if (required) {
                if (intervalHandle === null) intervalHandle = setInterval(intervalHandler, 1000);
            } else {
                stopInterval();
            }
        }

        function initializeSwitches() {
            loader.show();

            $http.get('/api/structure').then(function(response) {
                if (response.status !== 200) setTimeout(function() {
                    loader.hide();
                    initializeSwitches();
                }, 1000);

                $scope.switches = [];

                response.data.switches.forEach(function(dataset) {
                    var swtch = createSwitch(dataset);

                    $scope.switches.push(swtch);

                    switchesById[swtch.id] = swtch;
                });

                queuedUpdates.forEach(updateSwitch);
                queuedUpdates = [];

                startOrStopInterval();

                loader.hide();
            });
        }

        function updateSwitches() {
            loader.show();

            $http.get('/api/structure').then(function(response) {
                if (response.status !== 200) setTimeout(function() {
                    loader.hide();
                    updateSwitches();
                }, 1000);

                stopInterval();

                response.data.switches.forEach(updateSwitch);

                startOrStopInterval();

                loader.hide();
            });
        }

        this.toggle = function(id, state) {
            loader.show();

            $http.post('/api/switch/' + id + '/' + (state ? 'on' : 'off'), '')
                ['finally'](function() {
                    loader.hide();
                });
        };

        loader.show();

        socket.on('connect', function() {
            if ($scope.switches) return;

            loader.hide();
            initializeSwitches();
        });

        socket.on('switchUpdate', function(dataset) {
            if ($scope.switches) {
                $scope.$apply(function() {
                    updateSwitch(dataset);
                });

                startOrStopInterval();
            } else {
                queuedUpdates.push(dataset);
            }
        });

        socket.on('disconnect', function() {
            loader.show();
        });

        socket.on('reconnect', function() {
            loader.hide();
            updateSwitches();
        });
    }]);
