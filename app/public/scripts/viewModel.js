define(['knockout-3.2.0', 'put', 'dialogModel', 'areaModel', 'commandBlockModel', 'playerModel', 'unapplyBindings', 'customHandlers/imagemapster'], function(ko, put, dialogModel, areaModel, commandBlockModel, playerModel, unapplyBindings) {
        String.prototype.ucfirst = function() {
                return this.charAt(0).toUpperCase() + this.substr(1);
        }

        String.prototype.lcfirst = function() {
                return this.charAt(0).toLowerCase() + this.substr(1);
        }

        return function viewModel(data) {
                var self = this;

                self.ui = ko.observable(data.ui);
                self.id = data.id;
                self.title = ko.observable(data.title);
                self.round = ko.observable(data.round);
                self.phase = ko.observable(data.phase);
                self.phaseNames = data.phaseNames;
                self.phaseName = ko.computed(function() {
                        return self.phaseNames ? self.phaseNames[self.phase()] : 'None'
                });
                self.subPhaseNames = data.subPhaseNames;
                self.subPhase = ko.observable(data.subPhase);
                self.subPhaseName = ko.computed(function() {
                        return self.subPhaseNames ? self.subPhaseNames[self.subPhase()] : 'None'
                });
                self.status = ko.observable(data.status);
                self.notices = ko.observableArray(data.notices);
                self.errors = ko.observableArray(data.errors);
                self.turn = ko.observable(data.turn);
                self.turnName = data.turnName;
                self.turnClass = function() {
                        if (self.turnName) {
                                return self.turnName.replace(/ /g,'-').toLowerCase();
                        } else {
                                return "none";
                        }
                }();

                self.artilleryFireMarkers = ko.observable(data.artilleryFireMarkers);
                self.disruptionMarkers = ko.observable(data.disruptionMarkers);
                self.entrenchmentSticks = ko.observable(data.entrenchmentSticks);

                self.toPlayers = function(data) {
                        var players = [];
                        if (data) {
                                for (var i = 0; i < data.length; i++) {
                                        players.push(new playerModel(data[i]));
                                }
                        }
                        return players;
                };
                self.players = self.toPlayers(data.players);
                self.players.myUpdate = function(data) {
                        for (var i = 0; i < self.players.length; i++) {
                                self.players[i].myUpdate(data[i]);
                        }
                }

                self.confederatePlayer = ko.computed(function() {
                        return self.players[0];
                });

                self.unionPlayer = ko.computed(function() {
                        return self.players[1];
                });


                self.period = data.period;
                self.periodName = function() {
                        switch (self.period) {
                                case 1:
                                        return "Morning";
                                case 2:
                                        return "Mid-Day";
                                case 3:
                                        return "Afternoon"; 
                                case 4:
                                        return "Night";
                                default:
                                        return "None";
                        }
                }();

                self.cpUserIndices = data.cpUserIndices ? ko.observable(data.cpUserIndices[0]) : ko.observable(-1);
                self.currentPlayer = ko.computed(function() {
                        return ko.utils.arrayFirst(self.players, function(player) {
                                return player.id == self.cpUserIndices();
                        });
                });

                self.availableAreaIDS = function() {
                        var areas = [];
                        for (var i = 0; i < 47; i++) {
                                areas.push(i);
                        }
                        return areas
                }();

                self.boolOptions = [ false, true ];

                self.enemySide = ko.computed(function() {
                        if (self.currentPlayer().side === "union") {
                                return "confederate";
                        } else {
                                return "union";
                        }
                });

                self.winnerIndices = data.winnerIndices ? ko.observable(data.winnerIndices[0]) : ko.observable(-1);

                self.header = ko.computed(function() {
                        return {
                                title: self.title(),
                        turn: self.turn(),
                        round: self.round(),
                        phase: self.phase(),
                        subPhase: self.subPhase(),
                        cpUserIndices: self.cpUserIndices(),
                        winnerIndices: self.winnerIndices(),
                        }
                });
                self.label = ko.observable('dialog test');
                self.isOpen = ko.observable(false);
                self.open = function() {
                        self.isOpen(true);   
                };
                self.close = function() {
                        self.isOpen(false);   
                };

                self.toAreas = function(data) {
                        var areas = [];
                        if (data) {
                                for (var i = 0; i < data.length; i++) {
                                        areas.push(new areaModel(data[i]));
                                }
                        }
                        return areas;
                };
                self.areas = self.toAreas(data.areas);
                self.areas.myUpdate = function(data) {
                        for (var i = 0; i < self.areas.length; i++) {
                                self.areas[i].myUpdate(data[i]);
                        }
                };

                self.selectedAreaID = ko.observable(data.selectedAreaID);
                self.selectedArea = ko.computed(function() {
                        return self.areas[self.selectedAreaID()];
                });
                self.orderedAreaIDS = ko.observableArray(data.orderedAreaIDS ? data.orderedAreaIDS : [] );
                self.orderedAreas = ko.computed(function() {
                        var areas = [];
                        for (var i = 0; i < self.orderedAreaIDS().length; i++) {
                                var areaID = self.orderedAreaIDS()[i];
                                var area = self.areas[areaID];
                                areas.push(area);
                        }
                        return areas;
                });
                self.orderDiscAreaID = ko.observable(data.orderDiscAreaID);
                self.orderDiscArea = ko.computed(function() {
                        return self.areas[self.orderDiscAreaID()];
                });
                self.orderedArtillery = ko.computed(function() {
                        var units = [];
                        for (var i = 0; i < self.orderedAreas().length; i++) {
                                var found = ko.utils.arrayFilter(self.orderedAreas()[i].units(), function(unit) {
                                        return unit().kind == "artillery";
                                });
                                units = units.concat(found);
                        };
                        return units;
                });
                self.targetAreas = ko.computed(function() {
                        var ids = _.uniq(_.collect(self.orderedAreas(), function(area) {
                                return area.units()[0]().targetAreaID()
                        }));
                        var areas = [];
                        for (var i = 0; i < ids.length; i++) {
                                var areaID = ids[i];
                                var area = self.areas[areaID];
                                areas.push(area);
                        }
                        return areas;
                });
                self.targetTypes = ko.computed(function() {
                        var foundArt = false;
                        var foundInf = false;
                        var kinds = [];
                        _.each(self.targetAreas(), function(area) {
                                _.each(area.units(), function(unit) {
                                        if (unit().kind == "artillery") {
                                                foundArt = true;
                                        }
                                        if (unit().kind == "cavalry" || unit().kind == "infantry") {
                                                foundInf = true;
                                        }
                                });
                        });
                        if (foundInf) {
                                kinds.push("infantry");
                        }
                        if (foundArt) {
                                kinds.push("artillery");
                        }
                        return kinds;
                });
                //        return _.uniq(_.flatten(_.collect(self.targetAreas(), function(area) {
                //                console.log("area: " + ko.toJSON(area));
                //                return _.collect(area.units(), function(unit) {
                //                        if (unit().kind == "cavalry") {
                //                                return "infantry";
                //                        } else {
                //                                return unit().kind;
                //                        }
                //                });
                //        })));
                //});
//                                _.each(area.units(), function(unit) {
//                                        if (unit().kind == "cavalry") {
//                                                return "infantry";
//                                        } else {
//                                                return unit().kind;
//                                        }
//                                })
//                self.targetTypesFor = function(player) {
//                        return _.filter(_.uniq(_.collect(self.enemyUnitsFor(player), function(unit) {
//                                if (unit().kind === "cavalry") {
//                                        return "infantry";
//                                }
//                                return unit().kind;
//                       })), function(kind) {
//                               return kind == "infantry" || (kind == "artillery" && self.artilleryFireMarker())
//                       });
//                };
                self.gState = ko.computed(function() {
                        return {
                                selectedAreaID: self.selectedAreaID,
                        orderedAreaIDS: self.orderedAreaIDS,
                        }
                });
                self.confederateEliteInfantryCubeCount = ko.observable(data.cubes ? data.cubes.data.confederateEliteInfantry : 0);
                self.confederateRegularInfantryCubeCount = ko.observable(data.cubes ? data.cubes.data.confederateRegularInfantry : 0);
                self.unionInferiorInfantryCubeCount = ko.observable(data.cubes ? data.cubes.data.unionInferiorInfantry : 0);
                self.unionRegularInfantryCubeCount = ko.observable(data.cubes ? data.cubes.data.unionRegularInfantry : 0);
                self.unionEliteInfantryCubeCount = ko.observable(data.cubes ? data.cubes.data.unionEliteInfantry : 0);
                self.unionRegularCavalryCubeCount = ko.observable(data.cubes ? data.cubes.data.unionRegularCavalry : 0);

                self.dialog = new dialogModel(data.dialog);
                self.myUpdate = function(data) {
                        if (!data) {
                                return
                        }

                        if (data.cpUserIndices && self.cpUserIndices[0] != data.cpUserIndices[0]) {
                                self.cpUserIndices(data.cpUserIndices[0]);
                        };

                        if (self.title() != data.title) {
                                self.title(data.title)
                        }
                        if (self.turn() != data.turn) {
                                self.turn(data.turn)
                        }
                        if (self.phase() != data.phase) {
                                self.phase(data.phase)
                        }
                        if (self.subPhase() != data.subPhase) {
                                self.subPhase(data.subPhase)
                        }
                        if (self.round() != data.round) {
                                self.round(data.round)
                        }
                        if (self.status() != data.status) {
                                self.status(data.status)
                        }
                        if (self.selectedAreaID() != data.selectedAreaID) {
                                self.selectedAreaID(data.selectedAreaID);
                        }
                        if (self.orderedAreaIDS() && self.orderedAreaIDS() != data.orderedAreaIDS) {
                                self.orderedAreaIDS(data.orderedAreaIDS ? data.orderedAreaIDS : [] )
                        }
                        self.notices(data.notices);
                        self.errors(data.errors);
                        self.areas.myUpdate(data.areas);
                        self.players.myUpdate(data.players);
                        if (self.dialog.title() != data.dialog.title) {
                                self.dialog = new dialogModel(data.dialog)
                        }
                        if (self.dialog.content() != data.dialog.content) {
                                self.dialog = new dialogModel(data.dialog)
                        }
                        self.artilleryFireMarkers(data.artilleryFireMarkers);
                        self.disruptionMarkers(data.disruptionMarkers);
                        self.entrenchmentSticks(data.entrenchmentSticks);
                        self.orderDiscAreaID(data.orderDiscAreaID);
                        self.ui(data.ui);
                        self.updateMap(self.areas);
                }

                self.jUpdate = function(data) {
                        self.myUpdate(data);
                        if ($("#ui-box").length > 0) {
                                ko.unapplyBindings($("#ui-box"), true);
                        }
                        $("#flash-box div").prepend("<div id='ui-box'><div data-bind='visible: errors() && errors().length > 0'><div class='error'><ul data-bind='foreach: errors'><li data-bind='html: $data'></li></ul></div></div>" + self.ui() + "</div>");
                        ko.applyBindings(self, $("#ui-box")[0]);

                        if ($("#dialog").length > 0) {
                                ko.unapplyBindings($("#dialog"), true);
                        }
                        $("#board").append("<div id='dialog' data-bind='jqDialog: dialog, dialogVisible: true'>" + self.dialog.content() + "</div>");
                        ko.applyBindings(self, $("#dialog")[0]);
                };

                self.assignDamage = function(element, valueAccessor, allBindingsAccessor, vm, bindingContext) {
                        var value = valueAccessor(),
                            change = $(element).spinner("value") - value(),
                            unit = bindingContext.$data,
                            key = unit.side + unit.rank.ucfirst() + unit.kind.ucfirst() + "CubeCount";
                        var supply = self[key];

                        value($(element).spinner("value"));
                        if (unit.kind != "artillery") {
                                supply(supply() - change);
                        }
                };

                self.put = function(json) {
                        put(json, self.jUpdate);
                };

                self.defaultMapOptions = function() {
                        var opts = {
                                stroke: true,
                                strokeWidth: 10,
                                strokeOpacity: 0.5,
                                fill: false,
                                render_highlight: {
                                        strokeColor: 'ffff00',
                                },
                                render_select: {
                                        fill: true,
                                        fillColor: '000000',
                                        fillColorMask: 'FFFFFF',
                                        fillOpacity: 0.4,
                                        strokeColor: 'ff0000',
                                },
                                onClick: function(data) {
                                        if (data.selected) {
                                                self.put({action: 'select-area', target: "area-" + data.key});
                                        };
                                }, 
                                showToolTip: true,
                        };
                        return opts;
                };

                self.mapCommon = function(areas) {
                        var as = [];
                        var highlight = false;
                        for (var i = 0; i < areas.length; i++) {
                                var a = {};
                                var area = areas[i];
                                a.key = area.id.toString();
                                a.selected = area.selected();
                                a.isSelectable = area.selectable();
                                a.toolTip = "Area: " + area.id;
                                as.push(a);
                                if (area.selectable()) {
                                        highlight = true;
                                }
                        };
                        var opts = self.defaultMapOptions();
                        opts.areas = as;
                        opts.highlight = highlight;
                        return opts
                }

                self.updateMap = function(areas) {
                        var opts = self.mapCommon(areas);
                        $("area").mapster("deselect");
                        $(".clickmap").mapster("rebind", opts);
                };

                self.initMap = function(areas) {
                        var opts = self.mapCommon(areas);
                        $(".clickmap").mapster(opts);
                };

                self.init = function() {
                        ko.applyBindings(self);
                        self.initMap(self.areas);

                        if ($("#ui-box").length <= 0) {
                                $("#flash-box div").prepend("<div id='ui-box'>" + self.ui() + "</div>");
                        }
                        ko.applyBindings(self, $("#ui-box")[0]);

                        if ($("#dialog").length <= 0) {
                                $("#board").append("<div id='dialog' data-bind='jqDialog: dialog, root: $root, dialogVisible: true'>" + self.dialog.content() + "</div>");
                        }
                        ko.applyBindings(self, $("#dialog")[0]);
                        $("#areas").on("click", "img.unit", function() {
                                var context = ko.contextFor(this);
                                var index = context.$index;
                                var areaID = "area-" + context.$parent.id;
                                var $ajaxTarget = window.location.href;
                                var unitID = context.$data.id;

                                $.ajax({
                                        url: $ajaxTarget,
                                        data: { "target": areaID, "action": "select-unit", "unit": unitID }, 
                                        dataType: "json",
                                        success: function(data) {
                                                self.jUpdate(data);
                                        },
                                        type: "PUT",
                                        error: function(data) {
                                                window.location.href = $ajaxTarget;
                                        },
                                }); 
                        });
                };
        };
});
