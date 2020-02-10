define(['knockout-3.2.0', 'unitModel', 'commandBlockModel'], function(ko, unitModel, commandBlockModel) {
        return function areaModel(data) {
                var self = this;

                self.id = data.id;
                self.key = "area-" + self.id;
                self.newDamage = ko.observable(data.newDamage);
                self.shape = data.shape;
                self.coords = data.coords;
                self.selectable = ko.observable(data.selectable);
                self.selected = ko.observable(data.selected);
                self.artilleryFireMarker = ko.observable(data.artilleryFireMarker);
                self.disrupted = ko.observable(data.disrupted);

                self.commandBlock = ko.observable(data.commandBlock ?  new commandBlockModel(data.commandBlock) : null);
                self.toUnits = function(data) {
                        var units = [];
                        if (!data) {
                                return [];
                        }
                        for (var i = 0; i < data.length; i++) {
                                units.push(ko.observable(new unitModel(data[i])));
                        }
                        return units;
                };
                self.units = ko.observableArray(self.toUnits(data.units));

                self.units.myUpdate = function(data) {
                        if (data.units && self.units().length != data.units.length) {
                                self.units(self.toUnits(data.units))
                                        return;
                        }

                        for (var i = 0; i < self.units().length; i++) {
                                self.units()[i]().myUpdate(data.units[i])
                        }
                }

                self.orderDiscs = ko.observable(data.orderDiscs);


                self.hitUnits = ko.computed(function() {
                        if (self.units) {
                                return ko.utils.arrayFilter(self.units(), function(unit) {
                                        return unit().hit();
                                });
                        }
                });

                self.retreatableCav = ko.computed(function() {
                        return ko.utils.arrayFilter(self.units(), function(unit) {
                                return unit().newDamage() > 0 && unit().kind == "cavalry"
                        });
                });

                self.artilleryFor = function(player) {
                       return ko.utils.arrayFilter(self.friendlyUnitsFor(player), function(unit) {
                                return unit().kind == "artillery"
                       });
                };

                self.friendlyUnitsFor = function(player) {
                       return ko.utils.arrayFilter(self.units(), function(unit) {
                               return unit().side == player.side;
                       });
                };

                self.enemyUnitsFor = function(player) {
                       return ko.utils.arrayFilter(self.units(), function(unit) {
                               return unit().side != player.side;
                       });
                };

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

                self.myUpdate = function(data) {
                        if (!data) {
                                return
                        }
                        if (self.newDamage() != data.newDamage) {
                                self.newDamage(data.newDamage);
                        }
                        if (data.commandBlock) {
                                if (self.commandBlock()) {
                                        self.commandBlock().myUpdate(data.commandBlock);
                                } else {
                                        self.commandBlock(new commandBlockModel(data.commandBlock))
                                }
                        } else {
                                self.commandBlock(null);
                        }

                        if (self.selectable() != data.selectable) {
                                self.selectable(data.selectable);
                        };

                        if (self.selected() != data.selected) {
                                self.selected(data.selected);
                        };

                        if (self.artilleryFireMarker() != data.artilleryFireMarker) {
                                self.artilleryFireMarker(data.artilleryFireMarker)
                        };

                        if (self.disrupted() != data.disrupted) {
                                self.disrupted(data.disrupted)
                        };

                        self.orderDiscs(data.orderDiscs);
                        self.units.myUpdate(data);
                };
        }
})
